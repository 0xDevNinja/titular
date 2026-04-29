package sse

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

// newTestRouter wires a Multiplexer + Handler onto a gin engine for an
// in-process httptest server. Returns the engine, the live NATS
// connection (so individual tests can publish on it), and the
// multiplexer (so tests can introspect subscriber state and trigger
// shutdown). Per-IP cap is disabled so the test surface mirrors the
// pre-#90 baseline; the cap has its own dedicated test below.
func newTestRouter(t *testing.T, hbInterval time.Duration) (*gin.Engine, *nats.Conn, *Multiplexer) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	nc := runEmbeddedNATS(t)

	mux, err := NewMultiplexer(Config{NATS: nc})
	if err != nil {
		t.Fatalf("NewMultiplexer: %v", err)
	}
	if err := mux.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	t.Cleanup(mux.Stop)

	js, _ := nc.JetStream()
	h, err := NewHandler(HandlerConfig{
		Multiplexer:       mux,
		JetStream:         js,
		HeartbeatInterval: hbInterval,
		// Disable the per-IP cap for unrelated tests; every connection
		// in the in-process httptest server resolves to 127.0.0.1, so
		// otherwise the cap would skew the disconnect/replay assertions.
		MaxPerIP: -1,
	})
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}

	engine := gin.New()
	h.Register(engine.Group(""))
	return engine, nc, mux
}

// readSSEFrame consumes one full SSE frame (terminated by a blank
// line) from r and returns the raw frame bytes. Returns io.EOF when
// the stream closes.
func readSSEFrame(r *bufio.Reader) (string, error) {
	var b bytes.Buffer
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			// Return what we've buffered so callers can inspect a
			// partial frame on EOF.
			if b.Len() > 0 {
				return b.String(), err
			}
			return "", err
		}
		b.WriteString(line)
		if line == "\n" || line == "\r\n" {
			return b.String(), nil
		}
	}
}

// Test_Handler_HeadersAndContentType pins the SSE response headers
// required by #90: text/event-stream, no-cache, X-Accel-Buffering: no.
// Browsers and L7 proxies key on these to disable buffering and avoid
// caching the never-ending response.
func Test_Handler_HeadersAndContentType(t *testing.T) {
	engine, _, _ := newTestRouter(t, time.Hour)
	srv := httptest.NewServer(engine)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/events")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status: got %d, want 200", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/event-stream") {
		t.Errorf("content-type: got %q, want text/event-stream*", ct)
	}
	if cc := resp.Header.Get("Cache-Control"); cc != "no-cache" {
		t.Errorf("cache-control: got %q, want no-cache", cc)
	}
	if x := resp.Header.Get("X-Accel-Buffering"); x != "no" {
		t.Errorf("x-accel-buffering: got %q, want no", x)
	}
}

// Test_Handler_DispatchEvent publishes a single envelope and asserts
// the SSE handler renders it with id, event and data fields.
func Test_Handler_DispatchEvent(t *testing.T) {
	engine, nc, mux := newTestRouter(t, time.Hour)
	srv := httptest.NewServer(engine)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/events")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()

	r := bufio.NewReader(resp.Body)
	// First frame is the initial heartbeat (`: ping`).
	if _, err := readSSEFrame(r); err != nil {
		t.Fatalf("read initial heartbeat: %v", err)
	}

	// Wait for the handler to register on the hub before publishing —
	// otherwise the publish can race ahead of Subscribe and be lost
	// (live tail is at-most-once from subscribe time).
	waitForSubscribers(t, mux, 1)

	payload := `{"name":"BondingCurve.Bought","block_number":7,"log_index":2,"payload":{}}`
	if err := nc.Publish("trades.bought", []byte(payload)); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	frame, err := readSSEFrame(r)
	if err != nil {
		t.Fatalf("read event frame: %v", err)
	}
	if !strings.Contains(frame, "event: trades.bought") {
		t.Errorf("frame missing event line: %q", frame)
	}
	if !strings.Contains(frame, "id: ") {
		t.Errorf("frame missing id line: %q", frame)
	}
	if !strings.Contains(frame, "data: "+payload) {
		t.Errorf("frame missing/incorrect data line: %q", frame)
	}
}

// Test_Handler_FilterDropsOther asserts ?subjects= is honoured: a
// subscriber asking for trades.* must NOT see jobs.* events.
func Test_Handler_FilterDropsOther(t *testing.T) {
	engine, nc, mux := newTestRouter(t, time.Hour)
	srv := httptest.NewServer(engine)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/events?subjects=trades.*")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()

	r := bufio.NewReader(resp.Body)
	_, _ = readSSEFrame(r) // drain initial heartbeat
	waitForSubscribers(t, mux, 1)

	// Wrong-subject publish first, right-subject publish second. The
	// handler MUST skip the first entirely and emit the second.
	_ = nc.Publish("jobs.funded", []byte(`{"name":"Escrow.Funded"}`))
	_ = nc.Publish("trades.bought", []byte(`{"name":"BondingCurve.Bought"}`))

	frame, err := readSSEFrame(r)
	if err != nil {
		t.Fatalf("read frame: %v", err)
	}
	if !strings.Contains(frame, "event: trades.bought") {
		t.Errorf("got non-trades frame: %q", frame)
	}
	if strings.Contains(frame, "jobs") {
		t.Errorf("filter leak — jobs subject in trades-only stream: %q", frame)
	}
}

// Test_Handler_HeartbeatPing asserts that the heartbeat ticker emits a
// `: ping` comment frame at the configured interval. We use a 50ms
// interval so the test stays fast.
func Test_Handler_HeartbeatPing(t *testing.T) {
	engine, _, _ := newTestRouter(t, 50*time.Millisecond)
	srv := httptest.NewServer(engine)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/events")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()

	r := bufio.NewReader(resp.Body)
	// First frame — initial heartbeat fired before the loop enters
	// select.
	first, err := readSSEFrame(r)
	if err != nil {
		t.Fatalf("read first frame: %v", err)
	}
	if !strings.HasPrefix(first, ": ping") {
		t.Errorf("first frame: got %q, want comment ping", first)
	}

	// Second frame — periodic heartbeat from the ticker, arriving
	// after one tick of the configured interval.
	second, err := readSSEFrame(r)
	if err != nil {
		t.Fatalf("read second frame: %v", err)
	}
	if !strings.HasPrefix(second, ": ping") {
		t.Errorf("second frame: got %q, want comment ping", second)
	}
}

// Test_Handler_DisconnectCleanup asserts the handler unblocks promptly
// when the client closes the connection, and that the multiplexer's
// subscriber registry returns to its pre-connection size.
//
// This is the goroutine-leak guard required by #90's test plan.
func Test_Handler_DisconnectCleanup(t *testing.T) {
	engine, _, mux := newTestRouter(t, time.Hour)
	srv := httptest.NewServer(engine)
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL+"/events", nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}

	// Read the initial heartbeat so we know the handler is past
	// header-write and has registered its subscriber on the hub.
	r := bufio.NewReader(resp.Body)
	if _, err := readSSEFrame(r); err != nil {
		t.Fatalf("read initial heartbeat: %v", err)
	}
	waitForSubscribers(t, mux, 1)

	// Cancel the request context: the server-side handler must
	// observe ctx.Done and unwind, removing its subscriber.
	cancel()
	_ = resp.Body.Close()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if subscriberCount(mux) == 0 {
			return
		}
		runtime.Gosched()
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("subscriber not removed within deadline (count = %d)", subscriberCount(mux))
}

// Test_Handler_LastEventIDReplay drives the JetStream replay path:
// publish a sequence of events to a JetStream-backed stream, then open
// an SSE connection with Last-Event-ID set to the second message's id
// and verify only messages 3..N are replayed.
//
// This is the canonical resumption test required by #90. We bootstrap
// the stream manually rather than calling into the indexer publisher
// because the gateway test package can't import that without a
// dependency cycle; the wire shape is the contract.
func Test_Handler_LastEventIDReplay(t *testing.T) {
	engine, nc, _ := newTestRouter(t, time.Hour)
	srv := httptest.NewServer(engine)
	defer srv.Close()

	js, err := nc.JetStream()
	if err != nil {
		t.Fatalf("JetStream: %v", err)
	}

	// Create the canonical stream covering all four prefixes.
	if _, err := js.AddStream(&nats.StreamConfig{
		Name:     "TITULAR_EVENTS",
		Subjects: SubjectPrefixes,
		Storage:  nats.MemoryStorage,
	}); err != nil {
		t.Fatalf("AddStream: %v", err)
	}

	// Publish 5 envelopes onto trades.bought via JetStream so each
	// one carries a stream sequence we can resume from.
	pubAcks := make([]uint64, 0, 5)
	for i := 0; i < 5; i++ {
		ack, err := js.Publish("trades.bought", []byte(`{"i":`+strconv.Itoa(i)+`}`))
		if err != nil {
			t.Fatalf("js.Publish %d: %v", i, err)
		}
		pubAcks = append(pubAcks, ack.Sequence)
	}
	// Resume from after the second message: client claims to have
	// seen js:<seq2>, expects 3..5 replayed.
	resumeFrom := pubAcks[1]

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/events?subjects=trades.bought", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	req.Header.Set("Last-Event-ID", "js:"+strconv.FormatUint(resumeFrom, 10))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	defer resp.Body.Close()

	r := bufio.NewReader(resp.Body)
	// Drain frames until we have collected the three replay events.
	// Skip comment frames (heartbeats) and any non-event frames.
	seen := make([]uint64, 0, 3)
	deadline := time.Now().Add(5 * time.Second)
	for len(seen) < 3 && time.Now().Before(deadline) {
		frame, err := readSSEFrame(r)
		if err != nil && err != io.EOF {
			t.Fatalf("read frame: %v", err)
		}
		if strings.HasPrefix(frame, ": ") {
			continue // heartbeat / status comment
		}
		id := extractID(frame)
		if id == "" {
			continue
		}
		seq, ok := parseJetStreamID(id)
		if !ok {
			t.Errorf("non-js id in replay: %q", id)
			continue
		}
		seen = append(seen, seq)
	}
	if len(seen) < 3 {
		t.Fatalf("replay incomplete: got %d, want 3 (seen=%v)", len(seen), seen)
	}
	// Every replayed sequence must be > resumeFrom.
	for _, s := range seen {
		if s <= resumeFrom {
			t.Errorf("replay leaked older message: seq %d <= resume %d", s, resumeFrom)
		}
	}
}

// Test_Handler_ServerShutdown asserts the handler emits a final
// comment frame and exits promptly when the multiplexer is stopped.
// This is the gateway-graceful-shutdown counterpart of the disconnect
// test above.
func Test_Handler_ServerShutdown(t *testing.T) {
	engine, _, mux := newTestRouter(t, time.Hour)
	srv := httptest.NewServer(engine)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/events")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()

	r := bufio.NewReader(resp.Body)
	_, _ = readSSEFrame(r) // drain initial heartbeat
	waitForSubscribers(t, mux, 1)

	mux.Stop()

	// Read until EOF; we should see at most a "server shutdown"
	// comment then the body closes.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		_, err := readSSEFrame(r)
		if err != nil {
			return // EOF — handler exited
		}
	}
	t.Fatal("handler did not exit after multiplexer stop")
}

// extractID pulls the value of the `id:` line from an SSE frame, or ""
// if there is no id line. Used by the replay test.
func extractID(frame string) string {
	for _, line := range strings.Split(frame, "\n") {
		if strings.HasPrefix(line, "id: ") {
			return strings.TrimSpace(line[len("id: "):])
		}
	}
	return ""
}

// subscriberCount returns the multiplexer's current subscriber count.
// Defined as a test-only helper so the test can introspect the hub
// without exposing the internals on the public API.
func subscriberCount(m *Multiplexer) int {
	m.subsMu.RLock()
	defer m.subsMu.RUnlock()
	return len(m.subs)
}

// waitForSubscribers polls the multiplexer until at least n
// subscribers are registered. Used to defeat the inherent race
// between an HTTP request reaching the handler goroutine and a test
// publishing onto NATS — without this wait, a fast publish can race
// ahead of Subscribe and be missed by the live tail.
func waitForSubscribers(t *testing.T, m *Multiplexer, n int) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if subscriberCount(m) >= n {
			return
		}
		runtime.Gosched()
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("waitForSubscribers: only saw %d (want %d)", subscriberCount(m), n)
}

// TestSSE_RejectsExcessConcurrentConnectionsPerIP opens MaxPerIP
// concurrent /events streams from the same IP and asserts the
// (MaxPerIP+1)th request gets a 429 with Retry-After: 5. Defends the
// memory-pinning attack where a single client opens hundreds of
// long-lived streams to chew up gateway buffers.
func TestSSE_RejectsExcessConcurrentConnectionsPerIP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	nc := runEmbeddedNATS(t)

	mux, err := NewMultiplexer(Config{NATS: nc})
	if err != nil {
		t.Fatalf("NewMultiplexer: %v", err)
	}
	if err := mux.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	t.Cleanup(mux.Stop)

	const maxPerIP = 10
	js, _ := nc.JetStream()
	h, err := NewHandler(HandlerConfig{
		Multiplexer:       mux,
		JetStream:         js,
		HeartbeatInterval: time.Hour,
		MaxPerIP:          maxPerIP,
	})
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}

	engine := gin.New()
	// Trust the loopback proxy so c.ClientIP() resolves to the test
	// connection's source ip (127.0.0.1) rather than any forwarded
	// value. Without trusted proxies set, gin still returns the remote
	// address — fine here too — but we pin it explicitly so the test
	// is independent of gin's default trust posture.
	_ = engine.SetTrustedProxies(nil)
	h.Register(engine.Group(""))

	srv := httptest.NewServer(engine)
	defer srv.Close()

	// Open `maxPerIP` connections in parallel; each must succeed. We
	// hold a context per connection so the goroutines stay parked on
	// the stream until the test's defer cancels them — releasing the
	// per-IP slot would invalidate the (maxPerIP+1)th-must-fail
	// assertion.
	holds := make([]context.CancelFunc, 0, maxPerIP)
	defer func() {
		for _, c := range holds {
			c()
		}
	}()

	for i := 0; i < maxPerIP; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		holds = append(holds, cancel)

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL+"/events", nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("connection %d: Do: %v", i, err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("connection %d: status %d, want 200", i, resp.StatusCode)
		}
		// Drain the initial heartbeat so we KNOW the handler reached
		// the per-IP register step before we open the next connection;
		// otherwise the (maxPerIP+1)th request can race ahead of
		// acquireConn and erroneously succeed.
		r := bufio.NewReader(resp.Body)
		if _, err := readSSEFrame(r); err != nil {
			t.Fatalf("connection %d: read initial heartbeat: %v", i, err)
		}
		// Don't close the body — we want the connection to stay open.
		// resp.Body is closed indirectly when ctx is cancelled above.
		_ = resp
	}

	// Wait until the multiplexer sees maxPerIP subscribers so we know
	// every prior connection has fully registered, not just acquired
	// the per-IP slot.
	waitForSubscribers(t, mux, maxPerIP)

	// (maxPerIP+1)th request from the same source must be refused.
	resp, err := http.Get(srv.URL + "/events")
	if err != nil {
		t.Fatalf("(maxPerIP+1)th GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("(maxPerIP+1)th status: got %d, want 429", resp.StatusCode)
	}
	if ra := resp.Header.Get("Retry-After"); ra != "5" {
		t.Errorf("(maxPerIP+1)th Retry-After: got %q, want \"5\"", ra)
	}
}
