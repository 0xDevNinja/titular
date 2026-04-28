package publisher

import (
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

// StreamManager is the subset of [nats.JetStreamContext] needed to
// bootstrap the events stream. We narrow it to keep tests fast (the
// embedded server we use covers exactly these calls) and to make the
// boundary obvious to reviewers — `EnsureStream` is the only place in
// the package that mutates server state.
type StreamManager interface {
	StreamInfo(name string, opts ...nats.JSOpt) (*nats.StreamInfo, error)
	AddStream(cfg *nats.StreamConfig, opts ...nats.JSOpt) (*nats.StreamInfo, error)
}

// StreamOptions tunes the stream config produced by [EnsureStream]. Zero
// value uses the alpha defaults documented on each field. Callers that
// need to override anything should construct a [StreamOptions] explicitly
// rather than mutating one returned from a helper, because the zero value
// is the only API contract we publish.
type StreamOptions struct {
	// MaxAge is how long messages stay in the stream. Defaults to 7
	// days, which is long enough for an SDK consumer to recover from a
	// week-long outage but short enough that file storage stays
	// bounded for the alpha. Production deployments will likely raise
	// this; the value lives in env config and not in code.
	MaxAge time.Duration

	// MaxMsgsPerSubject caps per-subject message count, mostly as a
	// belt-and-suspenders against a runaway producer. Defaults to 0
	// (unlimited) — the publisher only publishes one message per
	// (tx, log_index) pair, so the natural ceiling is the chain's own
	// throughput.
	MaxMsgsPerSubject int64

	// DuplicateWindow is the JetStream dedup window. Messages with the
	// same `Nats-Msg-Id` published within this window are dropped
	// server-side. Defaults to 2 minutes — long enough to absorb a
	// reorg-and-replay of confirmed blocks but short enough that the
	// dedup table stays tiny.
	DuplicateWindow time.Duration

	// Storage selects in-memory vs file storage. Defaults to
	// [nats.FileStorage] so a server restart preserves the stream;
	// tests use [nats.MemoryStorage] to avoid touching disk.
	Storage nats.StorageType
}

// EnsureStream creates the [StreamName] stream if it does not already
// exist, bound to the four top-level subject prefixes in
// [SubjectPrefixes]. If the stream exists with the same name, the call is
// a no-op — we deliberately do not reconcile config on every boot
// because that would let a bad rollout silently flip retention or
// storage on a populated stream. Operators run a one-shot reconcile job
// when retention needs to change.
//
// The returned error is wrapped with the stream name so log lines from
// multiple indexers are distinguishable in shared environments.
func EnsureStream(mgr StreamManager, opts StreamOptions) error {
	if mgr == nil {
		return errors.New("publisher: nil StreamManager")
	}

	if opts.MaxAge == 0 {
		opts.MaxAge = 7 * 24 * time.Hour
	}
	if opts.DuplicateWindow == 0 {
		opts.DuplicateWindow = 2 * time.Minute
	}
	if opts.Storage == 0 {
		// nats.FileStorage is the zero value's neighbour, but explicit
		// is better than implicit and it makes the test override read
		// naturally.
		opts.Storage = nats.FileStorage
	}

	// Probe first. nats.go returns nats.ErrStreamNotFound when the
	// stream is missing; any other error is a transport / auth issue
	// that AddStream would also hit, so we surface it now.
	info, err := mgr.StreamInfo(StreamName)
	switch {
	case errors.Is(err, nats.ErrStreamNotFound):
		// Fall through to AddStream below.
	case err != nil:
		return fmt.Errorf("publisher: stream %s probe: %w", StreamName, err)
	default:
		// Stream already exists. Sanity check that it covers our
		// prefixes — if a previous deploy left a stream with a
		// narrower subject set, the publisher would silently drop
		// messages, which is exactly the failure mode operators have
		// hit elsewhere with shared NATS clusters.
		if err := checkSubjects(info); err != nil {
			return err
		}
		return nil
	}

	cfg := &nats.StreamConfig{
		Name:              StreamName,
		Subjects:          append([]string(nil), SubjectPrefixes...),
		Retention:         nats.LimitsPolicy,
		Storage:           opts.Storage,
		MaxAge:            opts.MaxAge,
		MaxMsgsPerSubject: opts.MaxMsgsPerSubject,
		Duplicates:        opts.DuplicateWindow,
		Discard:           nats.DiscardOld,
	}

	if _, err := mgr.AddStream(cfg); err != nil {
		return fmt.Errorf("publisher: stream %s create: %w", StreamName, err)
	}
	return nil
}

// checkSubjects asserts that an existing stream covers every prefix in
// [SubjectPrefixes]. Surfaces the missing prefix in the error so the
// operator does not have to diff configs by hand.
func checkSubjects(info *nats.StreamInfo) error {
	have := make(map[string]bool, len(info.Config.Subjects))
	for _, s := range info.Config.Subjects {
		have[s] = true
	}
	for _, p := range SubjectPrefixes {
		if !have[p] {
			return fmt.Errorf(
				"publisher: existing stream %s missing subject %q (have %v)",
				StreamName, p, info.Config.Subjects,
			)
		}
	}
	return nil
}
