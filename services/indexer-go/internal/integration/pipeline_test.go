//go:build integration

package integration

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/0xDevNinja/titular/services/indexer-go/internal/decoder"
	"github.com/0xDevNinja/titular/services/indexer-go/internal/subscriber"
)

// pipelineSink is the test-only glue that ties the three production
// packages together exactly as the eventual indexer main loop will.
//
// Flow per ConfirmedBlock:
//
//  1. Pin to the canonical block hash carried on the ConfirmedBlock so a
//     mid-block reorg cannot make us pick up logs from a stale fork.
//     `eth_getLogs` with a `BlockHash` filter is the primary defence here;
//     anvil supports it the same way Alchemy / geth do.
//  2. Dispatch each log through the production decoder. Logs whose topic0
//     is not registered (Emitter only emits CapabilitiesUpdated, but a
//     deployed chain may interleave other events) are silently skipped via
//     [decoder.DispatchAndHandle].
//  3. Forward every decoded event to the configured Handler — the real
//     [publisher.Publisher] in tests, so dedup is exercised end-to-end.
//
// The sink also keeps an in-memory log of every observation so a test
// can assert "block N at hash H emitted M events with these names" without
// snooping inside JetStream.
type pipelineSink struct {
	client  *ethclient.Client
	addr    common.Address
	handler decoder.Handler

	mu       sync.Mutex
	observed []observedEmit // append-only event log for assertions
}

// observedEmit records one (block, log, decoded-name, reorged-flag) tuple.
// The flag is propagated from the ConfirmedBlock — if the same Number
// appears twice the second entry will have reorged=true, which is the
// invariant the reorg test asserts.
type observedEmit struct {
	blockNumber uint64
	blockHash   common.Hash
	txHash      common.Hash
	logIndex    uint
	eventName   string
	reorged     bool
}

func newPipelineSink(t *testing.T, c *ethclient.Client, addr common.Address, h decoder.Handler) *pipelineSink {
	t.Helper()
	return &pipelineSink{client: c, addr: addr, handler: h}
}

// Emit satisfies subscriber.Sink. Returning an error is fatal for the
// subscriber, so we surface only the unrecoverable cases (RPC down,
// handler refused). Decoder errors for a single log are logged via t.Logf
// and skipped — the indexer's contract is "drop the bad log, keep the
// block" since a malformed log is contract-bug territory and fatally
// stalling the chain follower would punish operators for someone else's
// regression.
func (s *pipelineSink) Emit(ctx context.Context, b subscriber.ConfirmedBlock) error {
	logs, err := s.client.FilterLogs(ctx, ethereum.FilterQuery{
		BlockHash: &b.Hash,
		Addresses: []common.Address{s.addr},
	})
	if err != nil {
		return fmt.Errorf("FilterLogs(block=%s): %w", b.Hash.Hex(), err)
	}
	for _, log := range logs {
		// Take a copy: DispatchAndHandle stores the log inside the
		// DecodedEvent and the slice header is reused on the next
		// FilterLogs call.
		log := log
		ok, err := decoder.DispatchAndHandle(ctx, log, s.handler)
		if err != nil {
			return fmt.Errorf("dispatch log idx=%d: %w", log.Index, err)
		}
		if !ok {
			// Unknown topic0 / no topics: not an error, just not part
			// of the tracked surface (e.g. a constructor-emitted
			// access-control event from a real ACP deploy).
			continue
		}
		// Re-derive the qualified name via Dispatch so the assertion
		// log records the exact string the publisher saw. This is
		// cheap (O(1) topic0 lookup) and avoids exposing the entry
		// from the handler back through the Handler interface.
		ev, derr := decoder.Dispatch(log)
		if derr != nil {
			// Should be impossible — DispatchAndHandle just succeeded
			// for this log. Surface as a fatal error so a future
			// contract-state-change in Dispatch does not silently lose
			// observations from the assertion log.
			return fmt.Errorf("re-dispatch log idx=%d: %w", log.Index, derr)
		}
		s.mu.Lock()
		s.observed = append(s.observed, observedEmit{
			blockNumber: b.Number,
			blockHash:   b.Hash,
			txHash:      log.TxHash,
			logIndex:    log.Index,
			eventName:   ev.Name,
			reorged:     b.Reorged,
		})
		s.mu.Unlock()
	}
	return nil
}

// snapshot returns a defensive copy of the observation log so tests can
// iterate it without holding the sink's lock.
func (s *pipelineSink) snapshot() []observedEmit {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]observedEmit, len(s.observed))
	copy(out, s.observed)
	return out
}

// expectClean stops a running subscriber and asserts the loop exited
// because of context cancellation, not a fatal error. Any other terminal
// state (sink-rejected, reorg-deeper-than-window) is a real bug; failing
// loudly here keeps the test from ending with a confusing "no events"
// message when the subscriber actually died ten lines earlier.
func expectClean(t *testing.T, cancel func(), done <-chan error) {
	t.Helper()
	cancel()
	err := <-done
	if err == nil || errors.Is(err, context.Canceled) {
		return
	}
	t.Fatalf("subscriber.Run exited with non-cancel error: %v", err)
}
