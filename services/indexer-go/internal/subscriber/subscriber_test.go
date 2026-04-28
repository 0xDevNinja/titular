package subscriber_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/0xDevNinja/titular/services/indexer-go/internal/subscriber"
)

// ----------------------------------------------------------------------------
// Test doubles
// ----------------------------------------------------------------------------

// mockSub satisfies ethereum.Subscription for tests. The err channel lets a
// test simulate a websocket disconnect by sending a non-nil error.
type mockSub struct {
	err  chan error
	once sync.Once
}

func newMockSub() *mockSub           { return &mockSub{err: make(chan error, 1)} }
func (m *mockSub) Err() <-chan error { return m.err }
func (m *mockSub) Unsubscribe()      { m.once.Do(func() { close(m.err) }) }
func (m *mockSub) fail(e error)      { m.err <- e }

// mockClient implements subscriber.EthClient. Tests push headers onto the
// subscription channel via pushHead and seed canonical hashes for
// HeaderByNumber via setHeader. Concurrency-safe so a test goroutine can
// drive the subscriber from outside the Run loop.
type mockClient struct {
	mu sync.Mutex

	headers   map[uint64]*types.Header
	subCh     chan<- *types.Header
	sub       *mockSub
	subscribe func() (*mockSub, error) // override for SubscribeNewHead
}

func newMockClient() *mockClient {
	return &mockClient{headers: make(map[uint64]*types.Header)}
}

func (m *mockClient) SubscribeNewHead(_ context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.subscribe != nil {
		s, err := m.subscribe()
		if err != nil {
			return nil, err
		}
		m.subCh = ch
		m.sub = s
		return s, nil
	}
	s := newMockSub()
	m.subCh = ch
	m.sub = s
	return s, nil
}

func (m *mockClient) HeaderByNumber(_ context.Context, number *big.Int) (*types.Header, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if number == nil {
		return nil, errors.New("mock: HeaderByNumber(nil) not supported in tests")
	}
	h, ok := m.headers[number.Uint64()]
	if !ok {
		return nil, fmt.Errorf("mock: no header at %d", number.Uint64())
	}
	// Return a copy so a later mutation in the test (reorg) does not
	// retroactively change a previously returned header.
	cp := *h
	return &cp, nil
}

// setHeader registers the canonical header at a given height. To simulate
// a reorg, call setHeader again with a different ParentHash / nonce.
func (m *mockClient) setHeader(number uint64, parentHash common.Hash, nonce uint64) common.Hash {
	m.mu.Lock()
	defer m.mu.Unlock()
	h := &types.Header{
		Number:     new(big.Int).SetUint64(number),
		ParentHash: parentHash,
		Time:       1_700_000_000 + number,
		// Nonce is only used here to perturb the hash so test reorgs
		// produce a deterministic-but-different Hash().
		Nonce: types.EncodeNonce(nonce),
	}
	m.headers[number] = h
	return h.Hash()
}

// hasSub returns true once the subscriber has registered a subscription.
// Lock-protected so concurrent waitFor probes don't race with
// SubscribeNewHead.
func (m *mockClient) hasSub() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sub != nil
}

// pushHead delivers a head header to the subscriber. Blocks until the
// subscriber drains it (channel buffer = 16).
func (m *mockClient) pushHead(t *testing.T, number uint64) {
	t.Helper()
	m.mu.Lock()
	h, ok := m.headers[number]
	ch := m.subCh
	m.mu.Unlock()
	if !ok {
		t.Fatalf("pushHead: no header registered at %d", number)
	}
	cp := *h
	select {
	case ch <- &cp:
	case <-time.After(time.Second):
		t.Fatalf("pushHead: timeout pushing head %d", number)
	}
}

// chanSink captures every emitted ConfirmedBlock so tests can assert on
// the exact sequence and reorg flag.
type chanSink struct {
	mu     sync.Mutex
	blocks []subscriber.ConfirmedBlock
	emitFn func(subscriber.ConfirmedBlock) error // optional override
}

func newChanSink() *chanSink { return &chanSink{} }

func (s *chanSink) Emit(_ context.Context, b subscriber.ConfirmedBlock) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.emitFn != nil {
		if err := s.emitFn(b); err != nil {
			return err
		}
	}
	s.blocks = append(s.blocks, b)
	return nil
}

func (s *chanSink) snapshot() []subscriber.ConfirmedBlock {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]subscriber.ConfirmedBlock, len(s.blocks))
	copy(out, s.blocks)
	return out
}

// waitFor polls until pred returns true or the deadline passes. Used to
// synchronise against the subscriber's internal goroutine without sleeps.
func waitFor(t *testing.T, timeout time.Duration, pred func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if pred() {
			return
		}
		time.Sleep(2 * time.Millisecond)
	}
	t.Fatalf("waitFor: condition not met within %s", timeout)
}

// ----------------------------------------------------------------------------
// New() validates inputs
// ----------------------------------------------------------------------------

func TestNewRejectsNilDependencies(t *testing.T) {
	t.Parallel()
	if _, err := subscriber.New(nil, newChanSink(), subscriber.Config{}); err == nil {
		t.Fatal("expected error for nil client")
	}
	if _, err := subscriber.New(newMockClient(), nil, subscriber.Config{}); err == nil {
		t.Fatal("expected error for nil sink")
	}
}

// ----------------------------------------------------------------------------
// Confirmation depth: only emit at head - Confirmations + 1
// ----------------------------------------------------------------------------

func TestEmitsOnlyAfterConfirmationDepth(t *testing.T) {
	t.Parallel()

	const depth = 3
	mc := newMockClient()
	sink := newChanSink()

	// Build a 10-block linear chain.
	var prev common.Hash
	for n := uint64(1); n <= 10; n++ {
		prev = mc.setHeader(n, prev, n)
	}

	s, err := subscriber.New(mc, sink, subscriber.Config{
		Confirmations:      depth,
		HeaderFetchTimeout: time.Second,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()

	// Wait for subscription to open.
	waitFor(t, time.Second, func() bool { return mc.hasSub() })

	// Push head=5: anchor latches at confirmedTip = 5 - 3 + 1 = 3, so on
	// the first head only block 3 emits (no backfill from genesis — that
	// is a different module's job).
	mc.pushHead(t, 5)
	waitFor(t, time.Second, func() bool { return s.LastConfirmedBlock() == 3 })

	got := sink.snapshot()
	if len(got) != 1 || got[0].Number != 3 {
		t.Fatalf("first emit: got %+v, want [#3]", emitSummary(got))
	}
	if got[0].Reorged {
		t.Errorf("block 3 emitted with Reorged=true on initial emit")
	}
	if s.LastSeenBlock() != 5 {
		t.Errorf("LastSeenBlock = %d, want 5", s.LastSeenBlock())
	}

	// Head=6 → confirmedTip=4, emit just block 4.
	mc.pushHead(t, 6)
	waitFor(t, time.Second, func() bool { return s.LastConfirmedBlock() == 4 })
	got = sink.snapshot()
	if len(got) != 2 || got[1].Number != 4 {
		t.Fatalf("after head=6: got %+v, want [#3, #4]", emitSummary(got))
	}

	// Head=8 → confirmedTip=6, emit blocks 5 and 6 (forward catch-up
	// when newHeads jumped a number).
	mc.pushHead(t, 8)
	waitFor(t, time.Second, func() bool { return s.LastConfirmedBlock() == 6 })
	got = sink.snapshot()
	if len(got) != 4 {
		t.Fatalf("after head=8: emitted %d, want 4", len(got))
	}

	cancel()
	if err := <-done; !errors.Is(err, context.Canceled) {
		t.Errorf("Run returned %v, want context.Canceled", err)
	}
}

// TestSkipsBeforeDepthReached ensures we never emit when the head is
// shallower than the confirmation depth.
func TestSkipsBeforeDepthReached(t *testing.T) {
	t.Parallel()

	mc := newMockClient()
	sink := newChanSink()

	var prev common.Hash
	for n := uint64(1); n <= 5; n++ {
		prev = mc.setHeader(n, prev, n)
	}

	s, err := subscriber.New(mc, sink, subscriber.Config{
		Confirmations:      12,
		HeaderFetchTimeout: time.Second,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()

	waitFor(t, time.Second, func() bool { return mc.hasSub() })

	// Head=5, depth=12: confirmedTip would be -7. Nothing should emit.
	mc.pushHead(t, 5)

	// Give the loop a tick to process. We expect zero emissions and the
	// liveness gauge to have moved.
	waitFor(t, time.Second, func() bool { return s.LastSeenBlock() == 5 })
	if got := sink.snapshot(); len(got) != 0 {
		t.Errorf("emitted %d blocks before depth, want 0", len(got))
	}
	if s.LastConfirmedBlock() != 0 {
		t.Errorf("LastConfirmedBlock = %d, want 0", s.LastConfirmedBlock())
	}

	cancel()
	<-done
}

// ----------------------------------------------------------------------------
// Reorg detection: same height, different hash
// ----------------------------------------------------------------------------

func TestDetectsReorgWithinWindow(t *testing.T) {
	t.Parallel()

	const depth = 2
	mc := newMockClient()
	sink := newChanSink()

	// Initial chain: 1..6, all confirmed at head=6 (confirmedTip=5).
	var prev common.Hash
	for n := uint64(1); n <= 6; n++ {
		prev = mc.setHeader(n, prev, n)
	}
	hash5Old := mc.headers[5].Hash()

	s, err := subscriber.New(mc, sink, subscriber.Config{
		Confirmations:      depth,
		ReorgWindow:        4,
		HeaderFetchTimeout: time.Second,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()

	waitFor(t, time.Second, func() bool { return mc.hasSub() })

	mc.pushHead(t, 6)
	waitFor(t, time.Second, func() bool { return s.LastConfirmedBlock() == 5 })

	// Anchor latches at the first confirmable height, so only block 5
	// emits initially (not 1..5).
	if got := sink.snapshot(); len(got) != 1 || got[0].Number != 5 {
		t.Fatalf("initial emit: got %+v, want [#5]", emitSummary(sink.snapshot()))
	}

	// Reorg: replace block 5 with a different nonce so its hash changes.
	// Parent (block 4) stays the same so the "first match from tip
	// downward" walk lands on height 5.
	parent4 := mc.headers[4].Hash()
	hash5New := mc.setHeader(5, parent4, 999)
	if hash5New == hash5Old {
		t.Fatal("test setup: forced reorg produced identical hash")
	}

	// Push head=7 (also a fresh canonical block) to trigger another
	// drain. ConfirmedTip = 6, which forces a reorg recheck of [3..6].
	mc.setHeader(7, hash5New, 7) // 7 is now a child of the new 5
	// Re-bind block 6 onto the new 5 so the canonical chain at 6 reflects
	// the reorg, otherwise our "is the tip stable" loop would not see a
	// mismatch at 6 either.
	mc.setHeader(6, hash5New, 6)
	mc.pushHead(t, 7)

	waitFor(t, 2*time.Second, func() bool {
		got := sink.snapshot()
		// Expect 1 original (#5) + at least 1 reorg emission for #5 +
		// 1 fresh confirmation at #6.
		return len(got) >= 3
	})

	got := sink.snapshot()
	// Find the first reorg emission and assert it is for height 5.
	var sawReorgAt5 bool
	var sawFreshAt6 bool
	for _, b := range got {
		if b.Reorged && b.Number == 5 {
			sawReorgAt5 = true
			if b.Hash != hash5New {
				t.Errorf("reorg emit hash = %x, want %x", b.Hash, hash5New)
			}
		}
		if !b.Reorged && b.Number == 6 {
			sawFreshAt6 = true
		}
	}
	if !sawReorgAt5 {
		t.Errorf("no Reorged=true emission at height 5; got %+v", emitSummary(got))
	}
	if !sawFreshAt6 {
		t.Errorf("no fresh emit at height 6; got %+v", emitSummary(got))
	}
	if s.LastReorgAt() == 0 {
		t.Errorf("LastReorgAt not updated after reorg")
	}

	cancel()
	<-done
}

// emitSummary renders a sequence of ConfirmedBlocks compactly for failure
// messages.
func emitSummary(bs []subscriber.ConfirmedBlock) []string {
	out := make([]string, len(bs))
	for i, b := range bs {
		flag := ""
		if b.Reorged {
			flag = " (reorg)"
		}
		out[i] = fmt.Sprintf("#%d %s%s", b.Number, b.Hash.Hex()[:10], flag)
	}
	return out
}

// TestReorgDeeperThanWindowIsFatal: if the reorg edge of the window
// disagrees with the canonical chain, Run must return rather than rewind
// silently.
//
// Setup: depth=1, window=3. We confirm a long sequence so the in-memory
// hash cache is fully populated across the window. Then we replace the
// hashes at every height in the window AND at the edge boundary so that
// the recheck walk hits a mismatch at its oldest scanned height.
func TestReorgDeeperThanWindowIsFatal(t *testing.T) {
	t.Parallel()

	const depth = 1
	const window = 3
	mc := newMockClient()
	sink := newChanSink()

	// Build chain 1..15, then walk the head forward one block at a time
	// so the subscriber populates its hash cache for every height.
	var prev common.Hash
	for n := uint64(1); n <= 15; n++ {
		prev = mc.setHeader(n, prev, n)
	}

	s, err := subscriber.New(mc, sink, subscriber.Config{
		Confirmations:      depth,
		ReorgWindow:        window,
		HeaderFetchTimeout: time.Second,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()

	waitFor(t, time.Second, func() bool { return mc.hasSub() })

	// Walk the head from 5 → 12, one push at a time. After each push the
	// subscriber emits one new block, so by head=12 the cache holds
	// hashes for #5..#12 (window-trim culls older entries).
	mc.pushHead(t, 5)
	waitFor(t, time.Second, func() bool { return s.LastConfirmedBlock() == 5 })
	for h := uint64(6); h <= 12; h++ {
		mc.pushHead(t, h)
		hh := h
		waitFor(t, time.Second, func() bool { return s.LastConfirmedBlock() == hh })
	}

	// At this point: nextEmit=13, last=12, oldest=12-3+1=10. The cache
	// has {10,11,12,...} (window-trim keeps 4 entries: 9,10,11,12).
	// Replace every height from 10 upward so the entire scanned window
	// mismatches; the loop emits reorgs at 12, 11, 10, then hits oldest
	// (10) and returns fatal because oldest > 0.
	parent9 := mc.headers[9].Hash()
	mc.setHeader(10, parent9, 9999)
	h10 := mc.headers[10].Hash()
	mc.setHeader(11, h10, 9998)
	h11 := mc.headers[11].Hash()
	mc.setHeader(12, h11, 9997)
	h12 := mc.headers[12].Hash()
	mc.setHeader(13, h12, 13)

	mc.pushHead(t, 13)

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("Run returned nil; expected fatal reorg error")
		}
		if errors.Is(err, context.Canceled) {
			t.Fatalf("Run returned context.Canceled; expected fatal reorg")
		}
		// Good: any non-nil non-cancel error from Run on a deep reorg.
	case <-time.After(3 * time.Second):
		t.Fatal("Run did not return after deep reorg")
	}
}

// ----------------------------------------------------------------------------
// Reconnect on subscription failure
// ----------------------------------------------------------------------------

func TestReconnectsAfterSubscriptionDrop(t *testing.T) {
	t.Parallel()

	const depth = 1
	mc := newMockClient()
	sink := newChanSink()

	var prev common.Hash
	for n := uint64(1); n <= 10; n++ {
		prev = mc.setHeader(n, prev, n)
	}

	// First Subscribe call returns a subscription that fails immediately;
	// second Subscribe call returns a healthy one. After this the test
	// pushes a head and verifies confirmation continued past the reconnect.
	calls := 0
	mc.subscribe = func() (*mockSub, error) {
		calls++
		s := newMockSub()
		if calls == 1 {
			// Fail asynchronously so SubscribeNewHead returns first.
			go func() {
				time.Sleep(5 * time.Millisecond)
				s.fail(errors.New("simulated ws disconnect"))
			}()
		}
		return s, nil
	}

	s, err := subscriber.New(mc, sink, subscriber.Config{
		Confirmations:      depth,
		InitialBackoff:     5 * time.Millisecond,
		MaxBackoff:         20 * time.Millisecond,
		HeaderFetchTimeout: time.Second,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()

	// Wait for the second Subscribe call (i.e. the reconnect).
	waitFor(t, 2*time.Second, func() bool {
		mc.mu.Lock()
		defer mc.mu.Unlock()
		return calls >= 2
	})

	// Push head=5 → confirmedTip=5 (depth=1). First head after reconnect
	// anchors at #5 and emits just that block.
	mc.pushHead(t, 5)
	waitFor(t, 2*time.Second, func() bool { return s.LastConfirmedBlock() == 5 })

	if got := sink.snapshot(); len(got) != 1 || got[0].Number != 5 {
		t.Errorf("after reconnect: got %+v, want [#5]", emitSummary(sink.snapshot()))
	}

	cancel()
	<-done
}

// ----------------------------------------------------------------------------
// Sink errors are fatal — Run returns without infinite reconnect.
// ----------------------------------------------------------------------------

func TestSinkErrorIsFatal(t *testing.T) {
	t.Parallel()

	mc := newMockClient()
	sink := newChanSink()
	sink.emitFn = func(_ subscriber.ConfirmedBlock) error {
		return errors.New("boom: db down")
	}

	var prev common.Hash
	for n := uint64(1); n <= 5; n++ {
		prev = mc.setHeader(n, prev, n)
	}

	s, err := subscriber.New(mc, sink, subscriber.Config{
		Confirmations:      1,
		InitialBackoff:     time.Millisecond,
		HeaderFetchTimeout: time.Second,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()

	waitFor(t, time.Second, func() bool { return mc.hasSub() })
	mc.pushHead(t, 3)

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("Run returned nil; expected fatal sink error")
		}
		if errors.Is(err, context.Canceled) {
			t.Fatalf("Run returned context.Canceled; expected fatal sink error")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after sink error")
	}
}

// ----------------------------------------------------------------------------
// Cancellation
// ----------------------------------------------------------------------------

func TestRunReturnsOnContextCancel(t *testing.T) {
	t.Parallel()

	mc := newMockClient()
	sink := newChanSink()
	s, err := subscriber.New(mc, sink, subscriber.Config{Confirmations: 1})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()

	waitFor(t, time.Second, func() bool { return mc.hasSub() })

	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Run returned %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Run did not return after cancel")
	}
}

// TestNilHeaderTolerated: a malformed payload (header with nil Number)
// must not crash the loop.
func TestNilHeaderTolerated(t *testing.T) {
	t.Parallel()

	mc := newMockClient()
	sink := newChanSink()
	s, err := subscriber.New(mc, sink, subscriber.Config{Confirmations: 1})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()

	waitFor(t, time.Second, func() bool { return mc.hasSub() })

	// Push a header with a nil Number — should be silently dropped.
	mc.mu.Lock()
	ch := mc.subCh
	mc.mu.Unlock()
	select {
	case ch <- &types.Header{}:
	case <-time.After(time.Second):
		t.Fatal("could not push nil-number header")
	}

	// Subscriber should still be alive: register a real header and push it.
	var prev common.Hash
	for n := uint64(1); n <= 3; n++ {
		prev = mc.setHeader(n, prev, n)
	}
	mc.pushHead(t, 3)
	waitFor(t, time.Second, func() bool { return s.LastConfirmedBlock() == 3 })

	cancel()
	<-done
}
