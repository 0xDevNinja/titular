// Package subscriber tails an EVM chain head over a WebSocket connection,
// confirms blocks at a configurable depth, detects reorgs at the confirmed
// tip, and emits canonical headers downstream in monotonic order.
//
// Design goals (from M4 indexer requirements):
//
//   - WebSocket subscription via ethclient.SubscribeNewHead (typically an
//     Alchemy / Base Sepolia endpoint).
//   - Auto-reconnect with exponential backoff on transient failures
//     (subscription drop, ws hangup, RPC error). The loop only exits when
//     the caller cancels the parent context.
//   - 12-block confirmation: only emit a header at height H once the chain
//     tip has reached H + Confirmations. Default is 12.
//   - Reorg handling: when the canonical block at a confirmed height changes
//     hash, emit it again with the new hash. Re-emission is bounded by the
//     ReorgWindow so a malicious peer cannot force unbounded rewinds.
//   - Health gauge: LastSeenBlock returns the most recent head observed
//     (NOT yet confirmed). LastConfirmedBlock returns the last height that
//     was emitted to the sink. Both are atomically readable from any
//     goroutine and intended to back a Prometheus gauge or /healthz check.
//
// The package exposes only an EthClient interface (a minimal subset of
// ethclient.Client) so unit tests can drive the loop with a mock. The
// production wiring lives at the call site (e.g. cmd/indexer/main.go).
package subscriber

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"sync/atomic"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// DefaultConfirmations is the canonical confirmation depth used by the
// indexer when Config.Confirmations is left at zero. 12 blocks is the
// historical Ethereum/Base default for "safe" finality on PoS chains; the
// numbers we publish downstream are documented against this value.
const DefaultConfirmations = 12

// DefaultReorgWindow caps how far back the subscriber will rescan when it
// observes a hash change at a previously confirmed height. A reorg deeper
// than this is treated as operator-level corruption and surfaced as a
// fatal error rather than silently rewound.
const DefaultReorgWindow = 64

// EthClient is the minimal slice of *ethclient.Client used by Subscriber.
// Defining it as an interface lets unit tests inject a deterministic mock
// without touching network code.
//
// The signatures match go-ethereum's *ethclient.Client exactly so the
// production binding is simply:
//
//	c, err := ethclient.DialContext(ctx, wsURL)
//	sub, err := subscriber.New(c, sink, cfg)
type EthClient interface {
	// SubscribeNewHead opens a "newHeads" WebSocket subscription. The
	// returned ethereum.Subscription closes its Err() channel on
	// disconnect; the subscriber treats any non-nil error as a transient
	// disconnect and reconnects with backoff.
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)

	// HeaderByNumber fetches the canonical header at a given height. Used
	// to (a) backfill confirmed blocks when a new head jumps the queue
	// and (b) re-fetch blocks whose hash may have changed under a reorg.
	// Pass nil to fetch the latest head.
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

// ConfirmedBlock is the canonical, confirmation-aged block emitted to the
// sink. It carries enough information for downstream consumers to detect a
// re-emission (Reorged == true) without re-deriving canonical state on
// their own.
type ConfirmedBlock struct {
	// Number is the block height.
	Number uint64

	// Hash is the canonical block hash at Number after the configured
	// number of confirmations.
	Hash common.Hash

	// ParentHash is the parent block hash, useful for downstream consumers
	// that maintain their own canonical chain pointer.
	ParentHash common.Hash

	// Timestamp is the block timestamp (seconds since the unix epoch).
	Timestamp uint64

	// Reorged is true when this block has been emitted before with a
	// different Hash. The caller must treat the new (Number, Hash) pair as
	// canonical and unwind any state derived from the prior hash.
	Reorged bool
}

// Sink is the downstream consumer of confirmed blocks. Implementations MUST
// be safe for synchronous calls from a single goroutine and MUST be fast:
// the subscriber's main loop blocks on Emit, so a slow sink directly delays
// reorg detection.
type Sink interface {
	// Emit is called once per confirmed block in monotonic Number order,
	// except across a reorg window where the same Number may be emitted
	// twice (the second time with Reorged == true and a different Hash).
	// A non-nil error from Emit terminates the subscriber loop.
	Emit(ctx context.Context, block ConfirmedBlock) error
}

// Config controls subscriber behaviour. The zero value is valid and uses
// the documented defaults; callers typically set only Logger and the
// reconnect knobs.
type Config struct {
	// Confirmations is the number of blocks the subscriber waits before
	// emitting a header. Defaults to DefaultConfirmations (12).
	Confirmations uint64

	// ReorgWindow caps how far back the subscriber will rescan when a
	// hash change is observed at a previously confirmed height. Defaults
	// to DefaultReorgWindow (64).
	ReorgWindow uint64

	// InitialBackoff is the first sleep on reconnect failure. Doubles up
	// to MaxBackoff. Defaults to 1s.
	InitialBackoff time.Duration

	// MaxBackoff caps the reconnect sleep. Defaults to 30s.
	MaxBackoff time.Duration

	// HeaderFetchTimeout caps each HeaderByNumber call. Defaults to 10s.
	// Tests that drive the mock with synchronous returns may set this to
	// something tiny.
	HeaderFetchTimeout time.Duration

	// Logger receives structured progress events. nil disables logging.
	Logger *slog.Logger

	// now is overridable for tests. Production code leaves this nil; the
	// subscriber falls back to time.Now.
	now func() time.Time
}

func (c *Config) defaults() {
	if c.Confirmations == 0 {
		c.Confirmations = DefaultConfirmations
	}
	if c.ReorgWindow == 0 {
		c.ReorgWindow = DefaultReorgWindow
	}
	if c.InitialBackoff <= 0 {
		c.InitialBackoff = time.Second
	}
	if c.MaxBackoff <= 0 {
		c.MaxBackoff = 30 * time.Second
	}
	if c.HeaderFetchTimeout <= 0 {
		c.HeaderFetchTimeout = 10 * time.Second
	}
	if c.now == nil {
		c.now = time.Now
	}
}

// Subscriber tails an EVM chain head and emits confirmed blocks to a sink.
// Construct with New, drive with Run, observe with the LastSeen* / Last
// Confirmed* accessors.
type Subscriber struct {
	client EthClient
	sink   Sink
	cfg    Config

	// confirmedHashes tracks the hash we last emitted at each confirmed
	// height inside [tip-ReorgWindow, tip]. Used to detect reorgs.
	confirmedHashes map[uint64]common.Hash

	// nextEmit is the next height the subscriber will try to confirm.
	// Initialised on the first observed head as max(0, head-Confirmations).
	nextEmit uint64

	// initialised is set to true after the first head is processed; until
	// then nextEmit is unset and we cannot start emitting.
	initialised bool

	// Health gauges (atomic so callers can read without holding any lock).
	lastSeenBlock      atomic.Uint64
	lastConfirmedBlock atomic.Uint64
	lastReorgAt        atomic.Int64 // unix seconds; 0 means "never"
}

// New constructs a Subscriber. It does not open any network connection;
// call Run to start the loop.
//
// client and sink MUST be non-nil. Returns an error otherwise so the
// failure surfaces at boot rather than as a nil-deref deep in the loop.
func New(client EthClient, sink Sink, cfg Config) (*Subscriber, error) {
	if client == nil {
		return nil, errors.New("subscriber.New: client is nil")
	}
	if sink == nil {
		return nil, errors.New("subscriber.New: sink is nil")
	}
	cfg.defaults()
	return &Subscriber{
		client:          client,
		sink:            sink,
		cfg:             cfg,
		confirmedHashes: make(map[uint64]common.Hash),
	}, nil
}

// LastSeenBlock returns the height of the most recent head observed on the
// subscription, regardless of whether it has been confirmed yet. Suitable
// for a "is the websocket alive" liveness probe.
func (s *Subscriber) LastSeenBlock() uint64 {
	return s.lastSeenBlock.Load()
}

// LastConfirmedBlock returns the height of the most recent block emitted
// to the sink. A widening gap between this and LastSeenBlock indicates the
// sink is falling behind.
func (s *Subscriber) LastConfirmedBlock() uint64 {
	return s.lastConfirmedBlock.Load()
}

// LastReorgAt returns the unix timestamp (seconds) of the most recent
// reorg observed, or zero if none has been seen. Suitable for a reorg
// counter / alert in operator dashboards.
func (s *Subscriber) LastReorgAt() int64 {
	return s.lastReorgAt.Load()
}

// Run drives the subscription loop until ctx is cancelled or a fatal error
// occurs (sink failure, reorg deeper than the window). Transient errors —
// subscription drops, RPC hiccups — trigger an exponential-backoff
// reconnect and never return from Run.
//
// On clean ctx cancellation Run returns ctx.Err().
func (s *Subscriber) Run(ctx context.Context) error {
	backoff := s.cfg.InitialBackoff
	for {
		err := s.runOnce(ctx)
		switch {
		case err == nil:
			// runOnce only returns nil when ctx is done; surface that.
			return ctx.Err()
		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			return err
		case errors.Is(err, errFatal):
			return err
		default:
			s.logf("subscriber: transient error, will reconnect", "err", err, "backoff", backoff)
		}

		// Sleep with cancellation support before retrying.
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}

		backoff *= 2
		if backoff > s.cfg.MaxBackoff {
			backoff = s.cfg.MaxBackoff
		}
	}
}

// errFatal wraps an error that should NOT trigger a reconnect; the loop
// exits and the caller is expected to crash or alert.
var errFatal = errors.New("subscriber: fatal")

// fatal annotates an error so Run treats it as terminal rather than
// transient. Used for sink errors (downstream is broken) and reorgs deeper
// than the window (we cannot reconcile).
func fatal(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w: %v", errFatal, err)
}

// runOnce opens a single subscription and processes headers until the
// connection drops, the context is cancelled, or a fatal error occurs.
// Returns nil only when ctx is done; otherwise returns the error that
// terminated the inner loop so Run can decide to reconnect or exit.
func (s *Subscriber) runOnce(ctx context.Context) error {
	headers := make(chan *types.Header, 16)

	sub, err := s.client.SubscribeNewHead(ctx, headers)
	if err != nil {
		return fmt.Errorf("subscribe newHeads: %w", err)
	}
	defer sub.Unsubscribe()

	s.logf("subscriber: subscription open")

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-sub.Err():
			if err == nil {
				// nil on Err() means the subscription closed cleanly
				// (server shut down). Treat as transient.
				return errors.New("subscription closed by server")
			}
			return fmt.Errorf("subscription error: %w", err)
		case h, ok := <-headers:
			if !ok {
				return errors.New("headers channel closed")
			}
			if err := s.onHead(ctx, h); err != nil {
				return err
			}
		}
	}
}

// onHead processes a single head header: updates the liveness gauge and
// drains the confirmed-block backlog up to head-Confirmations. Returns a
// fatal error if the sink fails or a reorg exceeds the window.
func (s *Subscriber) onHead(ctx context.Context, h *types.Header) error {
	if h == nil || h.Number == nil {
		// Defensive: a malformed payload is a transient RPC issue, not
		// fatal. Drop and wait for the next head.
		return nil
	}
	headNum := h.Number.Uint64()
	s.lastSeenBlock.Store(headNum)

	// First head ever: anchor nextEmit at the deepest still-confirmable
	// block (head - Confirmations). For very fresh chains where the head
	// is shallower than Confirmations, start at 0.
	if !s.initialised {
		var start uint64
		if headNum >= s.cfg.Confirmations {
			start = headNum - s.cfg.Confirmations + 1
		}
		s.nextEmit = start
		s.initialised = true
		s.logf("subscriber: anchored", "head", headNum, "next_emit", s.nextEmit)
	}

	// Drain everything that has now aged past the confirmation depth.
	if headNum < s.cfg.Confirmations {
		return nil
	}
	confirmedTip := headNum - s.cfg.Confirmations + 1

	// First, walk forward through any never-emitted heights up to
	// confirmedTip. WebSocket newHeads can skip numbers when the client
	// is briefly disconnected, so we always backfill via HeaderByNumber.
	for s.nextEmit <= confirmedTip {
		hdr, err := s.fetchHeader(ctx, s.nextEmit)
		if err != nil {
			return fmt.Errorf("fetch confirmed header %d: %w", s.nextEmit, err)
		}
		blk := ConfirmedBlock{
			Number:     hdr.Number.Uint64(),
			Hash:       hdr.Hash(),
			ParentHash: hdr.ParentHash,
			Timestamp:  hdr.Time,
			Reorged:    false,
		}
		if err := s.sink.Emit(ctx, blk); err != nil {
			return fatal(fmt.Errorf("sink emit %d: %w", blk.Number, err))
		}
		s.recordConfirmed(blk.Number, blk.Hash)
		s.nextEmit++
	}

	// Then re-check the trailing window for reorgs. We compare the hash
	// we previously emitted against the canonical hash at the same
	// height; any mismatch is a reorg and we re-emit with Reorged=true.
	if err := s.detectReorgs(ctx, confirmedTip); err != nil {
		return err
	}

	return nil
}

// fetchHeader wraps HeaderByNumber with the configured per-call timeout so
// a hung RPC does not stall the subscription loop.
func (s *Subscriber) fetchHeader(ctx context.Context, number uint64) (*types.Header, error) {
	cctx, cancel := context.WithTimeout(ctx, s.cfg.HeaderFetchTimeout)
	defer cancel()
	return s.client.HeaderByNumber(cctx, new(big.Int).SetUint64(number))
}

// recordConfirmed updates the liveness gauge and the in-memory hash cache
// used for reorg detection. The cache is trimmed to ReorgWindow entries
// behind the current tip; older confirmations are considered final and a
// reorg deeper than the window is fatal.
func (s *Subscriber) recordConfirmed(num uint64, hash common.Hash) {
	s.confirmedHashes[num] = hash
	s.lastConfirmedBlock.Store(num)

	// Trim entries older than tip - ReorgWindow. We bound the map so a
	// long-running indexer does not leak memory.
	if num <= s.cfg.ReorgWindow {
		return
	}
	cutoff := num - s.cfg.ReorgWindow
	for k := range s.confirmedHashes {
		if k < cutoff {
			delete(s.confirmedHashes, k)
		}
	}
}

// detectReorgs walks the trailing ReorgWindow and re-emits any block whose
// canonical hash now differs from the one we stored. Returns a fatal error
// if a mismatch is observed at the very edge of the window (depth ==
// ReorgWindow), because the subscriber cannot prove the reorg did not go
// deeper than its memory.
func (s *Subscriber) detectReorgs(ctx context.Context, confirmedTip uint64) error {
	if confirmedTip == 0 {
		return nil
	}

	// We have already emitted heights [..., nextEmit-1]. The most recent
	// emission is at nextEmit-1, which equals confirmedTip after the
	// forward drain above. Walk backwards across the window.
	last := s.nextEmit - 1
	var oldest uint64
	if last >= s.cfg.ReorgWindow {
		oldest = last - s.cfg.ReorgWindow + 1
	}

	// Walk every cached height in [oldest..last]. We do not early-exit on
	// a match because a reorg can leave the canonical hash at a recent
	// height untouched while rewriting an earlier block (e.g. an
	// uncle-aunt swap that the new tip happens to graft over). Window is
	// bounded (default 64) so the O(window) scan is cheap.
	//
	// Loop is written with an explicit terminator so `oldest==0` does not
	// underflow uint64 when we decrement past zero.
	for height := last; ; height-- {
		want, known := s.confirmedHashes[height]
		if !known {
			// Either pre-anchor or already trimmed. Stop scanning.
			break
		}
		hdr, err := s.fetchHeader(ctx, height)
		if err != nil {
			return fmt.Errorf("reorg recheck %d: %w", height, err)
		}
		got := hdr.Hash()
		if got != want {
			// Hash changed at this height: a reorg.
			if height == oldest && oldest > 0 {
				// Mismatch at the very edge of what we remember, so the
				// reorg may be deeper. Refuse to silently rewind state
				// the downstream sink cannot reconstruct.
				return fatal(fmt.Errorf("reorg deeper than window at height %d", height))
			}

			blk := ConfirmedBlock{
				Number:     hdr.Number.Uint64(),
				Hash:       got,
				ParentHash: hdr.ParentHash,
				Timestamp:  hdr.Time,
				Reorged:    true,
			}
			if err := s.sink.Emit(ctx, blk); err != nil {
				return fatal(fmt.Errorf("sink reorg-emit %d: %w", height, err))
			}
			s.confirmedHashes[height] = got
			s.lastReorgAt.Store(s.cfg.now().Unix())
			s.logf("subscriber: reorg observed", "height", height, "old", want, "new", got)
		}

		if height == oldest {
			break
		}
	}
	return nil
}

// logf is a nil-safe wrapper over the optional structured logger.
func (s *Subscriber) logf(msg string, kv ...any) {
	if s.cfg.Logger == nil {
		return
	}
	s.cfg.Logger.Info(msg, kv...)
}
