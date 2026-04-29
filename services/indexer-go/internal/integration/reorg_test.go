//go:build integration

package integration

import (
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/json"
	"log/slog"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	natsserver "github.com/nats-io/nats-server/v2/server"
	srvtest "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"

	"github.com/0xDevNinja/titular/services/indexer-go/internal/publisher"
	"github.com/0xDevNinja/titular/services/indexer-go/internal/subscriber"
)

// TestReorgPipeline drives a full subscriber → decoder → publisher loop
// against an out-of-process anvil and asserts the four contract bullets
// from the M4 reorg-sim phase issue:
//
//  1. Events emit at confirmed depth (head − Confirmations + 1) and not
//     before.
//  2. A reorg at the confirmed-depth window is detected by the
//     subscriber.
//  3. The same height is re-emitted with the new canonical hash, with
//     [subscriber.ConfirmedBlock.Reorged] flagged true.
//  4. The publisher's JetStream dedup window correctly admits the
//     re-emit because the new tx has a different (tx_hash, log_index)
//     pair, so the dedup ID is distinct and both messages land on
//     the wire.
//
// We use a small confirmation depth (3 instead of the production 12) so
// the test does not have to mine 24+ blocks per phase. The subscriber's
// behaviour is independent of the depth value, and the unit-test suite
// in `internal/subscriber` already covers the boundary cases at
// production depth.
func TestReorgPipeline(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test (use -tags integration without -short)")
	}

	a := startAnvil(t)
	if a == nil {
		return // startAnvil already called t.Skip
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Install our minimal Emitter at a deterministic address. Using
	// `anvil_setCode` skips the deploy tx and keeps the test focused on
	// the read path under test.
	if err := a.setCode(ctx, emitterAddress, emitterRuntimeCode); err != nil {
		t.Fatalf("anvil_setCode: %v", err)
	}

	pk, devAddr := devKey(t)

	// --------------------------------------------------------------------
	// Embedded NATS for the publisher boundary. Memory-only storage so
	// the test does not touch disk.
	// --------------------------------------------------------------------
	js := startJetStream(t)
	if err := publisher.EnsureStream(js, publisher.StreamOptions{
		Storage:         nats.MemoryStorage,
		DuplicateWindow: 5 * time.Minute,
	}); err != nil {
		t.Fatalf("EnsureStream: %v", err)
	}
	pub, err := publisher.New(js, publisher.Options{SkipUnmapped: true})
	if err != nil {
		t.Fatalf("publisher.New: %v", err)
	}
	// SkipUnmapped is true because anvil's genesis often emits no
	// events, but we want to keep the door open for future variants of
	// this test that interleave non-tracked events without breaking
	// the assertion log.

	// --------------------------------------------------------------------
	// Wire the production subscriber against the anvil ws endpoint.
	// The sink's `eth_getLogs` calls go through a SEPARATE http client
	// so a slow log-fetch cannot back-pressure the WebSocket
	// subscription that drives the subscriber loop. In production the
	// indexer would use the same separation.
	// --------------------------------------------------------------------
	wsCli := a.dialEth(t)
	httpCli := a.dialEthHTTP(t)
	sink := newPipelineSink(t, httpCli, emitterAddress, pub)

	const depth = uint64(3)
	cfg := subscriber.Config{
		Confirmations:      depth,
		ReorgWindow:        16,
		HeaderFetchTimeout: 5 * time.Second,
		// Aggressive backoff because the test deliberately triggers a
		// transient `not found` error during the post-revert window
		// where the new fork has not caught up to the old confirmed
		// tip. Production uses the package defaults (1s → 30s).
		InitialBackoff: 50 * time.Millisecond,
		MaxBackoff:     200 * time.Millisecond,
	}
	if testing.Verbose() {
		cfg.Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	sub, err := subscriber.New(wsCli, sink, cfg)
	if err != nil {
		t.Fatalf("subscriber.New: %v", err)
	}
	loopCtx, stop := context.WithCancel(ctx)
	errCh := make(chan error, 1)
	go func() { errCh <- sub.Run(loopCtx) }()
	defer expectClean(t, stop, errCh)

	// Drive a few blocks so the subscriber's WS subscription receives
	// at least one head and anchors `nextEmit`. Without an anchor the
	// loop is dormant: subscriber.Run only walks the confirmation
	// window after the first head arrives, and a fresh anvil emits
	// no `newHeads` until the first mine call after subscribe.
	for i := 0; i < 5; i++ {
		if sub.LastSeenBlock() > 0 {
			break
		}
		if err := a.mine(ctx, 1); err != nil {
			t.Fatalf("anchor mine #%d: %v", i, err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	if sub.LastSeenBlock() == 0 {
		// Drain the error channel non-blockingly to surface any early
		// failure from sub.Run.
		select {
		case err := <-errCh:
			t.Fatalf("subscriber.Run exited early: %v", err)
		default:
		}
		t.Fatalf("subscriber never observed a head; anvil head=%d", mustBlockNumber(t, ctx, a))
	}

	// --------------------------------------------------------------------
	// Phase 1 — pre-reorg: emit two events (agentId=1 and agentId=2)
	// then bury them past the confirmation depth.
	//
	// We snapshot AFTER warm-up but BEFORE the first emit so the revert
	// undoes the emits without rolling back the genesis-and-warm-up
	// chain that anchored the subscriber's internal counters.
	// --------------------------------------------------------------------

	// Mine a couple of warm-up blocks so the subscriber's first head is
	// at a height where confirmedTip is well-defined (>= 1). Without
	// these, the very first head from anvil's genesis would race the
	// subscription open.
	if err := a.mine(ctx, 2); err != nil {
		t.Fatalf("mine warm-up: %v", err)
	}

	snapID, err := a.snapshot(ctx)
	if err != nil {
		t.Fatalf("evm_snapshot: %v", err)
	}

	// Pre-reorg tx #1: emitOne(1, 0xAA).
	r1, err := emitTx(ctx, a, pk, devAddr, 1, 0xAA)
	if err != nil {
		t.Fatalf("emit pre-reorg #1: %v", err)
	}
	// Pre-reorg tx #2: emitOne(2, 0xBB), one block later.
	r2, err := emitTx(ctx, a, pk, devAddr, 2, 0xBB)
	if err != nil {
		t.Fatalf("emit pre-reorg #2: %v", err)
	}

	// Bury both events past the confirmation depth so the subscriber
	// will emit them. depth=3 means we need head >= block(r2) + 3 - 1
	// for r2's height to be confirmed; mine `depth+2` to give a margin.
	if err := a.mine(ctx, int(depth)+2); err != nil {
		t.Fatalf("mine to confirm pre-reorg: %v", err)
	}

	// Check whether the subscriber loop terminated early. A clean
	// failure mode here keeps "no events" diagnostics from racing the
	// real cause (sink-rejected, reorg-deeper-than-window).
	select {
	case err := <-errCh:
		t.Fatalf("subscriber.Run terminated before pre-reorg waits: %v", err)
	default:
	}

	preR1 := waitForEmit(t, sink, r1.BlockNumber.Uint64(), 15*time.Second)
	preR2 := waitForEmit(t, sink, r2.BlockNumber.Uint64(), 15*time.Second)
	if preR1.eventName != "AgentRegistry.CapabilitiesUpdated" {
		t.Fatalf("pre #1 event name = %q", preR1.eventName)
	}
	if preR2.eventName != "AgentRegistry.CapabilitiesUpdated" {
		t.Fatalf("pre #2 event name = %q", preR2.eventName)
	}
	if preR1.reorged || preR2.reorged {
		t.Fatalf("initial emits should not have Reorged=true: %+v %+v", preR1, preR2)
	}

	// --------------------------------------------------------------------
	// Phase 2 — force a reorg via evm_revert and emit alternate events
	// at the same heights. We expect the subscriber to spot the hash
	// change at the previously-confirmed heights and re-emit.
	// --------------------------------------------------------------------

	ok, err := a.revert(ctx, snapID)
	if err != nil {
		t.Fatalf("evm_revert: %v", err)
	}
	if !ok {
		t.Fatal("evm_revert returned false (snapshot id consumed?)")
	}

	// After revert, anvil's chain head is back at the warm-up tip. We
	// now build an alternate fork: emit different events with the same
	// heights so the canonical hashes change at the same Numbers the
	// subscriber already remembers.
	r1b, err := emitTx(ctx, a, pk, devAddr, 1, 0xCC)
	if err != nil {
		t.Fatalf("emit post-reorg #1: %v", err)
	}
	r2b, err := emitTx(ctx, a, pk, devAddr, 2, 0xDD)
	if err != nil {
		t.Fatalf("emit post-reorg #2: %v", err)
	}

	// Assert the heights line up with the pre-reorg block heights
	// (otherwise the test is not actually exercising the reorg path —
	// it would just be an extension of the chain).
	if r1b.BlockNumber.Uint64() != r1.BlockNumber.Uint64() {
		t.Fatalf("post-reorg r1 height %d != pre %d", r1b.BlockNumber.Uint64(), r1.BlockNumber.Uint64())
	}
	if r2b.BlockNumber.Uint64() != r2.BlockNumber.Uint64() {
		t.Fatalf("post-reorg r2 height %d != pre %d", r2b.BlockNumber.Uint64(), r2.BlockNumber.Uint64())
	}
	if r1b.BlockHash == r1.BlockHash {
		t.Fatalf("post-reorg r1 has the same block hash as pre — anvil snapshot/revert is broken: %s", r1b.BlockHash.Hex())
	}
	if r2b.BlockHash == r2.BlockHash {
		t.Fatalf("post-reorg r2 has the same block hash as pre — anvil snapshot/revert is broken: %s", r2b.BlockHash.Hex())
	}

	// Bury again past the confirmation depth on the alternate fork.
	if err := a.mine(ctx, int(depth)+2); err != nil {
		t.Fatalf("mine to confirm post-reorg: %v", err)
	}

	// Wait for the subscriber to (a) detect the reorg at the confirmed
	// heights and (b) re-emit with the new hash and Reorged=true.
	//
	// We "tick" the chain by mining one extra block per second while
	// waiting. This compensates for a subscriber-package quirk: when
	// the WS reconnect window overlaps the burst of post-revert mines,
	// the freshly opened subscription has no head queued because anvil
	// only pushes newHeads going forward. Without an extra head the
	// subscriber's onHead never re-runs detectReorgs and the test
	// would deadlock until ctx expires. The production indexer hits
	// the same edge but on a real chain the next block arrives within
	// 2s on Base Sepolia, so this is faithful to the operator
	// behaviour we expect.
	postR1 := waitForReorgEmitWithTick(t, ctx, a, sink, r1.BlockNumber.Uint64(), preR1.blockHash, 15*time.Second)
	postR2 := waitForReorgEmitWithTick(t, ctx, a, sink, r2.BlockNumber.Uint64(), preR2.blockHash, 15*time.Second)

	if !postR1.reorged || !postR2.reorged {
		t.Fatalf("post-reorg emits must have Reorged=true: r1=%+v r2=%+v", postR1, postR2)
	}
	if postR1.blockHash == preR1.blockHash {
		t.Fatalf("post-reorg r1 hash unchanged (%s) — subscriber missed the reorg", postR1.blockHash.Hex())
	}
	if postR1.txHash == preR1.txHash {
		t.Fatalf("post-reorg r1 tx hash unchanged — anvil reused the pre-reorg tx, this test cannot validate dedup")
	}

	// --------------------------------------------------------------------
	// Phase 3 — drain the JetStream subject and verify the publisher
	// admitted four distinct messages (two pre + two post) with four
	// distinct dedup IDs. This is the M4 dedup-window contract under
	// reorg: same Number, different (tx_hash, log_index) ⇒ admitted.
	// --------------------------------------------------------------------

	got := drainSubject(t, js, "agents.capabilities_updated", 4, 5*time.Second)
	if len(got) != 4 {
		t.Fatalf("publisher delivered %d messages, want 4", len(got))
	}

	dedupIDs := make(map[string]struct{}, len(got))
	for _, m := range got {
		id := m.Header.Get(nats.MsgIdHdr)
		if id == "" {
			t.Errorf("message missing %q header: %s", nats.MsgIdHdr, string(m.Data))
			continue
		}
		if _, dup := dedupIDs[id]; dup {
			t.Errorf("duplicate dedup id %q in delivered set", id)
		}
		dedupIDs[id] = struct{}{}
	}
	if len(dedupIDs) != 4 {
		t.Errorf("dedup ids = %d, want 4: %v", len(dedupIDs), dedupIDs)
	}

	// Sanity-decode each envelope to make sure the bigint capabilities
	// field round-trips. We only check the value count, not the order
	// (NATS does not guarantee delivery order between PullSubscribe
	// fetches across reorgs).
	caps := map[string]int{}
	for _, m := range got {
		var env struct {
			Payload struct {
				Capabilities string `json:"Capabilities"`
			} `json:"payload"`
		}
		if err := json.Unmarshal(m.Data, &env); err != nil {
			t.Fatalf("envelope decode: %v", err)
		}
		caps[env.Payload.Capabilities]++
	}
	for _, want := range []string{"170", "187", "204", "221"} {
		if caps[want] != 1 {
			t.Errorf("capabilities=%s appears %d times, want 1 (full distribution: %v)", want, caps[want], caps)
		}
	}
}

// emitTx builds, signs, sends, and includes a tx that calls
// `Emitter.emitOne(agentId, capabilities)`. Returns the receipt so the
// caller can pin block-level invariants (Number, Hash, Status).
func emitTx(
	ctx context.Context,
	a *anvil,
	pk *ecdsa.PrivateKey,
	from common.Address,
	agentID, capabilities uint64,
) (*receiptLike, error) {
	data := encodeEmitOne(agentID, capabilities)
	tx, err := signedCall(ctx, a, pk, from, emitterAddress, data)
	if err != nil {
		return nil, err
	}
	r, err := sendAndMine(ctx, a, tx)
	if err != nil {
		return nil, err
	}
	if r.Status != 1 {
		return nil, errReverted{Tx: r.TxHash, Block: r.BlockNumber}
	}
	return &receiptLike{
		TxHash:      r.TxHash,
		BlockNumber: r.BlockNumber,
		BlockHash:   r.BlockHash,
		LogIndex:    func() uint {
			if len(r.Logs) == 0 {
				return 0
			}
			return r.Logs[0].Index
		}(),
	}, nil
}

// receiptLike is a narrowed view of *types.Receipt carrying only the
// fields the test asserts on. Defining a smaller struct keeps the test
// cases readable without `r.BlockNumber.Uint64()` ceremony at every
// reference.
type receiptLike struct {
	TxHash      common.Hash
	BlockNumber *big.Int
	BlockHash   common.Hash
	LogIndex    uint
}

// errReverted is returned by emitTx when anvil includes the tx but the
// EVM reverted. We surface this distinctly from a network error so the
// test sees "contract didn't emit" rather than "RPC down".
type errReverted struct {
	Tx    common.Hash
	Block *big.Int
}

func (e errReverted) Error() string {
	return "tx reverted: " + e.Tx.Hex() + " in block " + e.Block.String()
}

// encodeEmitOne assembles the calldata for `emitOne(uint256,uint256)`.
// Solidity ABI-encoding for two uint256 arguments is just the 4-byte
// selector followed by two 32-byte big-endian words.
func encodeEmitOne(agentID, capabilities uint64) []byte {
	buf := make([]byte, 4+32+32)
	copy(buf[0:4], emitOneSelector[:])
	binary.BigEndian.PutUint64(buf[4+24:4+32], agentID)
	binary.BigEndian.PutUint64(buf[4+32+24:4+32+32], capabilities)
	return buf
}

// mustBlockNumber returns the current anvil head, failing the test on
// RPC error. Used in failure-mode messages to give the operator a
// concrete starting point for triage.
func mustBlockNumber(t *testing.T, ctx context.Context, a *anvil) uint64 {
	t.Helper()
	n, err := a.blockNumber(ctx)
	if err != nil {
		t.Fatalf("anvil block number: %v", err)
	}
	return n
}

// waitForEmit blocks until the pipelineSink records an event at the
// given block number, or fails the test on timeout.
func waitForEmit(t *testing.T, sink *pipelineSink, height uint64, timeout time.Duration) observedEmit {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		for _, e := range sink.snapshot() {
			if e.blockNumber == height {
				return e
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("waitForEmit: no observation at height %d within %s", height, timeout)
	return observedEmit{}
}

// waitForReorgEmitWithTick is like waitForReorgEmit but mines one extra
// anvil block every second until the condition is met. See the call
// site for the rationale.
func waitForReorgEmitWithTick(
	t *testing.T,
	ctx context.Context,
	a *anvil,
	sink *pipelineSink,
	height uint64,
	oldHash common.Hash,
	timeout time.Duration,
) observedEmit {
	t.Helper()
	deadline := time.Now().Add(timeout)
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	for time.Now().Before(deadline) {
		for _, e := range sink.snapshot() {
			if e.blockNumber == height && e.blockHash != oldHash {
				return e
			}
		}
		select {
		case <-tick.C:
			if err := a.mine(ctx, 1); err != nil {
				t.Logf("tick mine: %v (continuing)", err)
			}
		case <-time.After(50 * time.Millisecond):
		}
	}
	t.Fatalf("waitForReorgEmitWithTick: no re-emit at height %d (old=%s) within %s",
		height, oldHash.Hex(), timeout)
	return observedEmit{}
}

// drainSubject pulls messages off the given JetStream subject until
// `want` are accumulated or the deadline passes. Used to assert on the
// publisher's output without flaky sleeps.
func drainSubject(t *testing.T, js nats.JetStreamContext, subject string, want int, timeout time.Duration) []*nats.Msg {
	t.Helper()
	sub, err := js.PullSubscribe(subject, "reorg-test", nats.AckExplicit())
	if err != nil {
		t.Fatalf("PullSubscribe(%s): %v", subject, err)
	}
	defer sub.Unsubscribe()

	got := make([]*nats.Msg, 0, want)
	deadline := time.Now().Add(timeout)
	for len(got) < want && time.Now().Before(deadline) {
		batch, err := sub.Fetch(want-len(got), nats.MaxWait(500*time.Millisecond))
		if err != nil && err != nats.ErrTimeout {
			t.Fatalf("Fetch(%s): %v", subject, err)
		}
		for _, m := range batch {
			_ = m.Ack()
			got = append(got, m)
		}
	}
	return got
}

// startJetStream boots an embedded nats-server with JetStream backed by
// a temp dir. Mirrors the helper in publisher_test.go but kept local so
// the integration package does not import _test.go from a sibling.
func startJetStream(t *testing.T) nats.JetStreamContext {
	t.Helper()
	dir := t.TempDir()

	opts := srvtest.DefaultTestOptions
	opts.Port = -1
	opts.JetStream = true
	opts.StoreDir = filepath.Join(dir, "jetstream")

	srv := srvtest.RunServer(&opts)
	t.Cleanup(func() {
		srv.Shutdown()
		srv.WaitForShutdown()
		_ = os.RemoveAll(opts.StoreDir)
	})

	if !srv.ReadyForConnections(2 * time.Second) {
		t.Fatal("nats-server not ready")
	}

	nc, err := nats.Connect(srv.ClientURL())
	if err != nil {
		t.Fatalf("nats.Connect: %v", err)
	}
	t.Cleanup(nc.Close)

	js, err := nc.JetStream()
	if err != nil {
		t.Fatalf("nc.JetStream: %v", err)
	}
	// Reference the natsserver package so the import is preserved if a
	// future change wants to inject a logger.
	_ = natsserver.RANDOM_PORT
	return js
}
