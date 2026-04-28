// Package publisher fans decoded EVM events out to NATS JetStream so the
// rest of the M4 surface (gateway, SDK, frontends) can consume Titular
// activity over a durable pub/sub bus.
//
// The publisher is the only [decoder.Handler] used by the indexer's main
// loop in M4: it accepts a [decoder.DecodedEvent], maps the
// contract-qualified event name to a JetStream subject, and publishes the
// abigen payload as JSON. Persistence (Postgres writes, materialised
// views) is layered on top in later phases by adding additional
// JetStream consumers, not by tapping the indexer directly. That keeps
// the indexer's responsibility narrow — confirmed log to durable
// message — and means a slow consumer cannot back-pressure the chain
// follower.
//
// # Subjects
//
// Subjects use four top-level prefixes that mirror the SDK's surface:
//
//	agents.*  — AgentRegistry lifecycle (registration, score, controller)
//	jobs.*    — Job + JobFactory lifecycle (creation, execution, dispute,
//	            and the Escrow side of fund movement)
//	tokens.*  — FeeSplitter, BuybackBurner, HookRegistry — anything that
//	            moves protocol-level token flow or upgrades the hook set
//	trades.*  — reserved for the bonding-curve / DEX flow tracked by the
//	            M2 handlers; populated by a later phase but pinned here so
//	            the SDK's subject grammar is stable now.
//
// Subject naming rules:
//
//   - Lowercase, dot-separated, version-less. The wire format is the
//     evolution boundary, not the subject string.
//   - Each subject corresponds to exactly one contract-qualified event
//     name (one decoded event -> one publish), so consumers can subscribe
//     at any granularity without server-side fan-in.
//
// # Source of truth
//
// [SubjectFor] is the only mapping consulted at runtime. The package-level
// table [eventSubjects] is enumerated in tests against
// [decoder.KnownTopics] so any future regen that adds a tracked event will
// fail CI until a subject is assigned, and any subject that no longer has
// a corresponding decoder will likewise surface. This is the same pattern
// that keeps `decoder/decoders.go` in sync with `abi/abi_test.go`'s
// pinned event surface.
package publisher

import "fmt"

// StreamName is the JetStream stream that holds every subject this package
// publishes. A single stream simplifies retention/storage configuration
// and makes per-subject consumers cheap; the trade-off is that all four
// top-level prefixes share the same MaxAge / Storage policy. That is the
// right call for the alpha — we want every event to live the same length
// of time so replay is uniform — and can be split later by binding the
// existing subjects to additional streams without changing consumers.
const StreamName = "TITULAR_EVENTS"

// SubjectPrefixes is the closed set of top-level subject prefixes the
// publisher emits to. JetStream stream creation uses these directly:
// `<prefix>.>` lets a single stream capture every event without listing
// individual subjects, while still letting consumers narrow down with a
// concrete subject like `jobs.funded`.
//
// Order is alphabetical and stable; the test suite asserts the exact
// slice so any addition forces an explicit decision (and a stream config
// migration in deployed environments).
var SubjectPrefixes = []string{
	"agents.>",
	"jobs.>",
	"tokens.>",
	"trades.>",
}

// eventSubjects maps a [decoder.DecodedEvent.Name] (i.e. the
// "<contract>.<event>" key produced by `decoder/decoders.go`) to its
// JetStream subject.
//
// Mapping rationale:
//
//   - AgentRegistry.* -> agents.*  — the entire agent identity surface.
//   - Job.*           -> jobs.*    — per-job state machine.
//   - JobFactory.*    -> jobs.*    — clone deployments and admin updates;
//     SDK consumers want all "job lifecycle" events on one prefix.
//   - Escrow.*        -> jobs.funded / jobs.released / jobs.refunded —
//     funds custody is conceptually part of the job pipeline. Splitting
//     this onto a separate `escrow.*` prefix would force the SDK to
//     subscribe to two prefixes for a single user-visible flow ("how is
//     job N progressing?"), which is the opposite of what we want.
//   - FeeSplitter.*   -> tokens.*  — protocol-level token routing.
//   - BuybackBurner.* -> tokens.*  — same.
//   - HookRegistry.*  -> tokens.*  — hook governance affects token flow
//     since hooks are how routing is composed.
//
// The `trades.*` prefix is intentionally absent here: the M2 BondingCurve
// `Bought` / `Sold` events are not part of the M3 ACP set tracked by the
// decoder yet, but the prefix is reserved in [SubjectPrefixes] so the
// stream is provisioned for it from day one.
var eventSubjects = map[string]string{
	// AgentRegistry — `agents.*`.
	"AgentRegistry.AgentRegistered":             "agents.registered",
	"AgentRegistry.MetadataUpdated":              "agents.metadata_updated",
	"AgentRegistry.CapabilitiesUpdated":          "agents.capabilities_updated",
	"AgentRegistry.ActiveStatusChanged":          "agents.active_status_changed",
	"AgentRegistry.ScorePosted":                  "agents.score_posted",
	"AgentRegistry.ControllerTransferProposed":   "agents.controller_transfer_proposed",
	"AgentRegistry.ControllerTransferAccepted":   "agents.controller_transfer_accepted",
	"AgentRegistry.ControllerTransferCancelled":  "agents.controller_transfer_cancelled",

	// Job — `jobs.*`. State machine for a single Job clone.
	"Job.JobInitialised":    "jobs.initialised",
	"Job.AgentAccepted":     "jobs.agent_accepted",
	"Job.ResultSubmitted":   "jobs.result_submitted",
	"Job.ResultApproved":    "jobs.result_approved",
	"Job.ResultRejected":    "jobs.result_rejected",
	"Job.JobCompleted":      "jobs.completed",
	"Job.JobCancelled":      "jobs.cancelled",
	"Job.DisputeRaised":     "jobs.dispute_raised",
	"Job.DisputeResolved":   "jobs.dispute_resolved",
	"Job.EvaluatorAssigned": "jobs.evaluator_assigned",

	// JobFactory — `jobs.*`. Clone deployments / admin.
	"JobFactory.JobCreated":             "jobs.created",
	"JobFactory.ImplementationUpdated":  "jobs.implementation_updated",
	"JobFactory.DefaultArbiterUpdated":  "jobs.default_arbiter_updated",

	// Escrow — `jobs.*` (funds movement is part of the job lifecycle).
	"Escrow.Funded":   "jobs.funded",
	"Escrow.Released": "jobs.released",
	"Escrow.Refunded": "jobs.refunded",

	// FeeSplitter — `tokens.*`. Protocol-level fee routing.
	"FeeSplitter.FeeSplit":             "tokens.fee_split",
	"FeeSplitter.TreasuryUpdated":      "tokens.treasury_updated",
	"FeeSplitter.BuybackBurnerUpdated": "tokens.buyback_burner_updated",

	// BuybackBurner — `tokens.*`. TITU buyback execution / config.
	"BuybackBurner.BuybackAndBurn":    "tokens.buyback_and_burn",
	"BuybackBurner.RouterUpdated":     "tokens.router_updated",
	"BuybackBurner.SwapPathUpdated":   "tokens.swap_path_updated",
	"BuybackBurner.MinTituOutUpdated": "tokens.min_titu_out_updated",

	// HookRegistry — `tokens.*`. Approved-hook governance.
	"HookRegistry.HookRegistered":   "tokens.hook_registered",
	"HookRegistry.HookDeregistered": "tokens.hook_deregistered",
}

// SubjectFor returns the JetStream subject for the given decoded event
// name. The returned `ok` is false when the event has no mapping; callers
// (the [Publisher]) treat that as a soft skip rather than a fatal error,
// because the decoder is allowed to surface events ahead of the SDK's
// subject grammar (e.g. during a regen that adds a binding before the
// matching subject is decided).
//
// SubjectFor is read-only and safe for concurrent use after package init.
func SubjectFor(eventName string) (string, bool) {
	subj, ok := eventSubjects[eventName]
	return subj, ok
}

// AllSubjects returns the slice of subjects currently mapped, sorted by
// event name. Used by tests to assert the table is complete and by the
// stream-bootstrap helper for sanity logging — not by the per-event hot
// path, which goes through [SubjectFor].
func AllSubjects() map[string]string {
	out := make(map[string]string, len(eventSubjects))
	for k, v := range eventSubjects {
		out[k] = v
	}
	return out
}

// validatePrefix sanity-checks that a configured subject starts with one
// of the [SubjectPrefixes]. It is invoked at package init via the test
// suite (not at runtime) — any drift here is a programmer error caught
// long before deploy.
func validatePrefix(subject string) error {
	for _, p := range SubjectPrefixes {
		// `agents.>` matches `agents.foo`, `agents.foo.bar`, etc.
		// We only need a prefix-with-dot check because subjects in
		// [eventSubjects] are concrete (no wildcards).
		dot := len(p) - 1 // strip trailing '>' to get e.g. "agents."
		if len(subject) >= dot && subject[:dot] == p[:dot] {
			return nil
		}
	}
	return fmt.Errorf("publisher: subject %q does not match any registered prefix", subject)
}
