package sse

import (
	"path/filepath"
	"runtime"
	"testing"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/server"
	srvtest "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
)

// runEmbeddedNATS boots an in-process nats-server with JetStream so
// the SSE multiplexer tests exercise the same wire path the indexer's
// publisher uses. Mirrors graph.runEmbeddedNATS — duplicated here
// because the helper is unexported in that package and cross-package
// test helpers are messier than a copy.
func runEmbeddedNATS(t *testing.T) *nats.Conn {
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
	})
	if !srv.ReadyForConnections(2 * time.Second) {
		t.Fatal("nats-server not ready")
	}
	nc, err := nats.Connect(srv.ClientURL())
	if err != nil {
		t.Fatalf("nats.Connect: %v", err)
	}
	t.Cleanup(nc.Close)
	return nc
}

// silence the natsserver import so future maintainers see it as
// intentional even when no test directly references the type.
var _ = natsserver.RANDOM_PORT

// Test_Multiplexer_DispatchAllSubjects asserts the hub fans every
// subject the indexer publisher emits to subscribers without a filter.
// This is the smoke test for the package — it covers the full surface
// (Start/Subscribe/dispatch/Close/Stop) end-to-end with a live NATS
// server.
func Test_Multiplexer_DispatchAllSubjects(t *testing.T) {
	nc := runEmbeddedNATS(t)
	mux, err := NewMultiplexer(Config{NATS: nc})
	if err != nil {
		t.Fatalf("NewMultiplexer: %v", err)
	}
	if err := mux.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer mux.Stop()

	sub, err := mux.Subscribe(SubscribeOptions{})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	defer sub.Close()

	// Publish one envelope on each top-level prefix and confirm we
	// receive all four. We use plain bytes here — the multiplexer is
	// envelope-agnostic.
	cases := []struct{ subject, body string }{
		{"agents.registered", `{"name":"AgentRegistry.Registered","block_number":1,"log_index":0}`},
		{"jobs.funded", `{"name":"Escrow.Funded","block_number":2,"log_index":1}`},
		{"tokens.fee_split", `{"name":"FeeSplitter.FeeSplit","block_number":3,"log_index":2}`},
		{"trades.bought", `{"name":"BondingCurve.Bought","block_number":4,"log_index":3}`},
	}

	for _, c := range cases {
		if err := nc.Publish(c.subject, []byte(c.body)); err != nil {
			t.Fatalf("Publish %q: %v", c.subject, err)
		}
	}

	got := make(map[string]bool)
	deadline := time.After(3 * time.Second)
	for len(got) < len(cases) {
		select {
		case ev, ok := <-sub.C:
			if !ok {
				t.Fatal("subscriber channel closed early")
			}
			got[ev.Subject] = true
		case <-deadline:
			t.Fatalf("only saw %d of %d subjects: %v", len(got), len(cases), got)
		}
	}
	for _, c := range cases {
		if !got[c.subject] {
			t.Errorf("missing subject %q", c.subject)
		}
	}
}

// Test_Multiplexer_FilterAppliesPattern pins per-subscriber subject
// filtering. A subscriber asking for `trades.*` must NOT see jobs
// events; a separate subscriber asking for `jobs.>` must see only
// jobs events.
func Test_Multiplexer_FilterAppliesPattern(t *testing.T) {
	nc := runEmbeddedNATS(t)
	mux, err := NewMultiplexer(Config{NATS: nc})
	if err != nil {
		t.Fatalf("NewMultiplexer: %v", err)
	}
	if err := mux.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer mux.Stop()

	tradesSub, err := mux.Subscribe(SubscribeOptions{Filter: Filter{"trades.*"}})
	if err != nil {
		t.Fatalf("Subscribe trades: %v", err)
	}
	defer tradesSub.Close()

	jobsSub, err := mux.Subscribe(SubscribeOptions{Filter: Filter{"jobs.>"}})
	if err != nil {
		t.Fatalf("Subscribe jobs: %v", err)
	}
	defer jobsSub.Close()

	_ = nc.Publish("trades.bought", []byte(`{"x":1}`))
	_ = nc.Publish("jobs.funded", []byte(`{"x":2}`))
	_ = nc.Publish("agents.registered", []byte(`{"x":3}`))

	// trades subscriber: expect exactly trades.bought.
	select {
	case ev := <-tradesSub.C:
		if ev.Subject != "trades.bought" {
			t.Errorf("trades sub: got %q, want trades.bought", ev.Subject)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("trades sub: timed out")
	}

	// jobs subscriber: expect exactly jobs.funded.
	select {
	case ev := <-jobsSub.C:
		if ev.Subject != "jobs.funded" {
			t.Errorf("jobs sub: got %q, want jobs.funded", ev.Subject)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("jobs sub: timed out")
	}

	// Drain settle: wait briefly to be sure the rejected `agents.*`
	// publish did NOT slip through the trades filter. We rely on a
	// short polling window; a longer wait would slow the suite for no
	// new signal.
	select {
	case ev := <-tradesSub.C:
		t.Fatalf("trades sub leaked: got %q", ev.Subject)
	case <-time.After(200 * time.Millisecond):
		// expected — nothing else
	}
}

// Test_Multiplexer_DropsOnOverflow asserts the bounded-buffer behaviour:
// a slow subscriber that does NOT drain its channel must see Drops
// increment past the buffer size, while a fast subscriber on the same
// hub keeps receiving every message.
//
// This is the guard against the classic "one slow client takes down the
// fanout" failure mode.
func Test_Multiplexer_DropsOnOverflow(t *testing.T) {
	nc := runEmbeddedNATS(t)
	mux, err := NewMultiplexer(Config{NATS: nc})
	if err != nil {
		t.Fatalf("NewMultiplexer: %v", err)
	}
	if err := mux.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer mux.Stop()

	// Slow subscriber: tiny buffer, never drained.
	const bufSize = 2
	slow, err := mux.Subscribe(SubscribeOptions{Buffer: bufSize})
	if err != nil {
		t.Fatalf("Subscribe slow: %v", err)
	}
	defer slow.Close()

	// Publish far more than the buffer can hold. NATS publishes are
	// async — the client batches them onto the wire, the server fans
	// them back to the dispatch callback on its own read loop. We
	// poll until the total observed (drops + buffered) settles at
	// the expected count, then assert drops dominate.
	const n = 50
	for i := 0; i < n; i++ {
		_ = nc.Publish("trades.bought", []byte(`{"i":1}`))
	}
	if err := nc.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	// Wait until the hub has dispatched every published message. The
	// total dispatched is "drops + items currently in buffer"; the
	// items in buffer are at most `bufSize`. The deadline is generous
	// because the race-detector slows scheduling materially.
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		dispatched := int(slow.Drops.Load()) + len(slow.C)
		if dispatched >= n {
			break
		}
		runtime.Gosched()
		time.Sleep(20 * time.Millisecond)
	}

	if slow.Drops.Load() == 0 {
		t.Fatal("expected drops on saturated subscriber, got 0")
	}
	// Once everything has dispatched, drops MUST equal n - bufSize
	// exactly (the buffer absorbed bufSize messages, the rest were
	// dropped). We allow a small ±1 tolerance for the rare case
	// where the dispatcher has not finished the last message at the
	// moment we sample — running with -count=N catches that.
	if got := slow.Drops.Load(); got < uint64(n-bufSize-1) || got > uint64(n) {
		t.Errorf("drops out of range: got %d, want in [%d, %d]",
			got, n-bufSize-1, n)
	}
}

// Test_Multiplexer_StopClosesSubscribers verifies graceful shutdown:
// after Stop, every subscriber's channel is closed (so the SSE handler
// loop exits) and Subscribe returns ErrStopped.
func Test_Multiplexer_StopClosesSubscribers(t *testing.T) {
	nc := runEmbeddedNATS(t)
	mux, err := NewMultiplexer(Config{NATS: nc})
	if err != nil {
		t.Fatalf("NewMultiplexer: %v", err)
	}
	if err := mux.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}

	sub, err := mux.Subscribe(SubscribeOptions{})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	mux.Stop()

	// Channel must close.
	select {
	case _, ok := <-sub.C:
		if ok {
			t.Fatal("expected closed channel")
		}
	case <-time.After(time.Second):
		t.Fatal("subscriber channel was not closed by Stop")
	}

	// Subscribe must refuse new connections.
	if _, err := mux.Subscribe(SubscribeOptions{}); err == nil {
		t.Fatal("Subscribe after Stop: expected error")
	}

	// Stop is idempotent.
	mux.Stop()
}

// Test_Multiplexer_CloseDetachesNoLeak confirms a Close on one
// subscriber does NOT affect other subscribers' channels and removes
// the closed subscriber from the registry so subsequent dispatches
// don't keep targeting it.
//
// Catches the regression where the fanout map is keyed by stale ids.
func Test_Multiplexer_CloseDetachesNoLeak(t *testing.T) {
	nc := runEmbeddedNATS(t)
	mux, err := NewMultiplexer(Config{NATS: nc})
	if err != nil {
		t.Fatalf("NewMultiplexer: %v", err)
	}
	if err := mux.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer mux.Stop()

	a, _ := mux.Subscribe(SubscribeOptions{})
	b, _ := mux.Subscribe(SubscribeOptions{})

	a.Close()

	// a's channel must close.
	select {
	case _, ok := <-a.C:
		if ok {
			t.Fatal("subscriber a channel: expected closed")
		}
	case <-time.After(time.Second):
		t.Fatal("subscriber a channel: not closed by Close")
	}

	// b must still receive events.
	_ = nc.Publish("trades.bought", []byte(`{"x":1}`))
	select {
	case ev := <-b.C:
		if ev.Subject != "trades.bought" {
			t.Errorf("subscriber b: got %q", ev.Subject)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("subscriber b: timed out (close on a may have leaked)")
	}
	b.Close()
}

// Test_Multiplexer_NilConn pins the constructor's nil guard.
func Test_Multiplexer_NilConn(t *testing.T) {
	if _, err := NewMultiplexer(Config{}); err == nil {
		t.Fatal("expected error on nil conn")
	}
}

// Test_subjectMatches drives the pattern matcher with the cases the
// SDK is documented to accept. Exhaustive enough to catch regressions
// in the token-walk loop without becoming a property test.
func Test_subjectMatches(t *testing.T) {
	cases := []struct {
		pattern, subject string
		want             bool
	}{
		// exact
		{"trades.bought", "trades.bought", true},
		{"trades.bought", "trades.sold", false},

		// `*` single token
		{"trades.*", "trades.bought", true},
		{"trades.*", "trades.bonding.bought", false},
		{"*.bought", "trades.bought", true},
		{"*.bought", "agents.bought", true},

		// `>` rest of subject
		{"trades.>", "trades.bought", true},
		{"trades.>", "trades.bonding.bought", true},
		{">", "anything.at.all", true},
		{"trades.>", "trades", false}, // > requires >= 1 trailing token

		// length mismatches
		{"a.b", "a.b.c", false},
		{"a.b.c", "a.b", false},

		// edge: empty pattern never matches
		{"", "trades.bought", false},
	}
	for _, tc := range cases {
		if got := subjectMatches(tc.pattern, tc.subject); got != tc.want {
			t.Errorf("subjectMatches(%q,%q) = %v, want %v",
				tc.pattern, tc.subject, got, tc.want)
		}
	}
}

// Test_SubjectsFromCSV pins the query-string parser. The handler relies
// on this for `?subjects=...`; an empty / whitespace-only input must
// produce a nil filter so the multiplexer treats it as "all subjects".
func Test_SubjectsFromCSV(t *testing.T) {
	cases := []struct {
		in   string
		want Filter
	}{
		{"", nil},
		{"   ", nil},
		{",,,", nil},
		{"trades.*", Filter{"trades.*"}},
		{"trades.*,jobs.>", Filter{"trades.*", "jobs.>"}},
		{"  trades.* ,  jobs.>  ", Filter{"trades.*", "jobs.>"}},
	}
	for _, tc := range cases {
		got := SubjectsFromCSV(tc.in)
		if len(got) != len(tc.want) {
			t.Errorf("SubjectsFromCSV(%q) = %v, want %v", tc.in, got, tc.want)
			continue
		}
		for i, p := range got {
			if p != tc.want[i] {
				t.Errorf("SubjectsFromCSV(%q)[%d] = %q, want %q",
					tc.in, i, p, tc.want[i])
			}
		}
	}
}

// TestSubjectsFromCSV_RejectsGreaterThan asserts the parser refuses a
// bare `>` pattern. A leading `>` would let a client subscribe across
// every prefix the multiplexer holds, which defeats the closed-set
// posture documented on SubjectPrefixes.
func TestSubjectsFromCSV_RejectsGreaterThan(t *testing.T) {
	if got := SubjectsFromCSV(">"); got != nil {
		t.Errorf("SubjectsFromCSV(\">\") = %v, want nil (rejected)", got)
	}
	if _, err := ParseSubjectsCSV(">"); err == nil {
		t.Error("ParseSubjectsCSV(\">\"): expected error, got nil")
	}
	// The reject must propagate even when mixed with a valid pattern:
	// silently dropping the bad entry would let a client think the
	// remaining narrow filter is being honoured while the broad one
	// is quietly tossed.
	if got := SubjectsFromCSV(">,trades.*"); got != nil {
		t.Errorf("SubjectsFromCSV(\">,trades.*\") = %v, want nil", got)
	}
}

// TestSubjectsFromCSV_RejectsBareStar asserts the parser refuses a
// leading `*`. `*` matches any single token, so `*` alone or
// `*.bought` would let a client cross the registered-prefix boundary
// (e.g. fan an `attacker.>` test subject if one were ever published).
func TestSubjectsFromCSV_RejectsBareStar(t *testing.T) {
	if got := SubjectsFromCSV("*"); got != nil {
		t.Errorf("SubjectsFromCSV(\"*\") = %v, want nil", got)
	}
	if got := SubjectsFromCSV("*.bought"); got != nil {
		t.Errorf("SubjectsFromCSV(\"*.bought\") = %v, want nil", got)
	}
	if _, err := ParseSubjectsCSV("*.bought"); err == nil {
		t.Error("ParseSubjectsCSV(\"*.bought\"): expected error, got nil")
	}
}

// TestSubjectsFromCSV_AllowsPrefixWildcard locks in that `*` and `>`
// remain legal within a registered prefix; the validation only
// rejects patterns whose FIRST token is the wildcard.
func TestSubjectsFromCSV_AllowsPrefixWildcard(t *testing.T) {
	cases := []struct {
		in   string
		want Filter
	}{
		{"agents.*", Filter{"agents.*"}},
		{"trades.>", Filter{"trades.>"}},
		{"jobs.*,tokens.>", Filter{"jobs.*", "tokens.>"}},
		{"agents.registered", Filter{"agents.registered"}},
	}
	for _, tc := range cases {
		got := SubjectsFromCSV(tc.in)
		if len(got) != len(tc.want) {
			t.Errorf("SubjectsFromCSV(%q) = %v, want %v", tc.in, got, tc.want)
			continue
		}
		for i, p := range got {
			if p != tc.want[i] {
				t.Errorf("SubjectsFromCSV(%q)[%d] = %q, want %q",
					tc.in, i, p, tc.want[i])
			}
		}
	}
}

// Test_parseJetStreamID covers the Last-Event-ID parser for the canonical
// js:<seq> shape and the rejected forms.
func Test_parseJetStreamID(t *testing.T) {
	cases := []struct {
		in     string
		seq    uint64
		wantOK bool
	}{
		{"js:1", 1, true},
		{"js:0", 0, true},
		{"js:18446744073709551615", 1<<64 - 1, true}, // uint64 max
		{"js:", 0, false},
		{"hub:1", 0, false},
		{"", 0, false},
		{"js:abc", 0, false},
		{"js:1.5", 0, false},
	}
	for _, tc := range cases {
		seq, ok := parseJetStreamID(tc.in)
		if ok != tc.wantOK || seq != tc.seq {
			t.Errorf("parseJetStreamID(%q) = (%d,%v), want (%d,%v)",
				tc.in, seq, ok, tc.seq, tc.wantOK)
		}
	}
}
