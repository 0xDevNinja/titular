package sse

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Handler serves GET /events: a Server-Sent Events stream that
// multiplexes the NATS subjects published by the indexer into a
// per-client text/event-stream response.
//
// Wiring lives in cmd/gateway: the handler takes a *Multiplexer (live
// tail) and an optional nats.JetStreamContext (history replay when the
// client sends Last-Event-ID). When JetStream is nil, replay is
// silently skipped — the stream still works for live tail, which is
// the common case.
//
// # Response shape
//
//	Content-Type: text/event-stream; charset=utf-8
//	Cache-Control: no-cache
//	Connection: keep-alive
//	X-Accel-Buffering: no
//
// Each event is rendered as:
//
//	id: js:<seq>
//	event: <subject>
//	data: <json envelope>
//
// Heartbeat frames (`: ping\n\n`) are emitted on HeartbeatInterval so
// the connection stays alive through proxies that idle-close TCP. The
// SSE spec defines lines beginning with `:` as comments; clients ignore
// them, which is exactly the desired behaviour.
//
// # Last-Event-ID replay
//
// When the client reconnects with a `Last-Event-ID` header in the form
// `js:<seq>`, the handler opens an ephemeral JetStream consumer with
// nats.StartSequence(seq+1) and drains every backlog message before
// attaching the connection to the live multiplexer. Backlog messages
// keep the same `id:` shape as live messages (js:<seq>) so the client
// can reconnect again from any point. When the header is missing /
// malformed / not in `js:` form, the handler skips replay and starts at
// the live tail; this is the safe fallback because the canonical
// SSE-spec semantic for Last-Event-ID is "best effort".
//
// # Buffering / backpressure
//
// Each connection owns a bounded sse.Subscriber buffer (DefaultBufferSize).
// The hub drops on overflow; the handler logs the running drop count
// when the connection terminates so a slow client surfaces in metrics.
//
// # Per-IP connection cap
//
// A single client opening N concurrent SSE connections is the cheapest
// way to pin gateway memory: each connection holds a 64-event buffer,
// a goroutine, and a NATS subscriber slot. The handler refuses new
// /events connections from an IP once that IP already has MaxPerIP
// connections open; the rejection is a 429 with Retry-After: 5 so the
// client backs off rather than hot-looping. Defaults to 10 per IP and
// is configurable via HandlerConfig.MaxPerIP (cmd/gateway reads
// GATEWAY_SSE_MAX_PER_IP and forwards). Setting MaxPerIP <= 0 disables
// the cap, which is appropriate for tests and tightly-trusted internal
// deployments only.
type Handler struct {
	mux              *Multiplexer
	js               nats.JetStreamContext
	heartbeatInterval time.Duration
	logger           zerolog.Logger
	now              func() time.Time

	maxPerIP int
	connsMu  sync.Mutex
	conns    map[string]int
}

// HandlerConfig wires the handler.
type HandlerConfig struct {
	// Multiplexer is the live-tail hub. Required.
	Multiplexer *Multiplexer

	// JetStream is used for Last-Event-ID replay. Optional; when nil,
	// the Last-Event-ID header is ignored and the connection starts at
	// the live tail.
	JetStream nats.JetStreamContext

	// HeartbeatInterval is how often `: ping` frames are emitted. Zero
	// uses DefaultHeartbeatInterval.
	HeartbeatInterval time.Duration

	// Logger receives drop / error log lines. Zero uses zerolog.Nop().
	Logger zerolog.Logger

	// MaxPerIP caps concurrent /events connections from a single
	// client IP. Zero falls back to DefaultMaxPerIP; explicitly setting
	// a negative value disables the cap (intended for tests and
	// trusted-network deployments only).
	MaxPerIP int
}

// DefaultHeartbeatInterval is the heartbeat cadence specified by #90 —
// every 30 seconds. Long enough to be cheap, short enough to keep most
// idle-killing intermediaries happy (cloudflare's default is 100s,
// nginx is 60s). Operators who terminate at a more aggressive proxy
// can lower this via HandlerConfig.
const DefaultHeartbeatInterval = 30 * time.Second

// DefaultMaxPerIP caps concurrent SSE connections per source IP at 10.
// Sized to allow a normal browsing session (a small handful of tabs)
// while refusing the cheap memory-pinning attack of a single IP
// holding hundreds of streams. Operators behind a known L7 proxy with
// trusted-IP forwarding may raise this via GATEWAY_SSE_MAX_PER_IP.
const DefaultMaxPerIP = 10

// NewHandler validates cfg and returns a configured Handler. Returns an
// error when required fields are missing — silent zero-config would let
// a misconfigured deploy serve clients but never deliver events.
func NewHandler(cfg HandlerConfig) (*Handler, error) {
	if cfg.Multiplexer == nil {
		return nil, errors.New("sse: nil multiplexer")
	}
	hb := cfg.HeartbeatInterval
	if hb <= 0 {
		hb = DefaultHeartbeatInterval
	}
	logger := cfg.Logger
	if reflect.DeepEqual(logger, zerolog.Logger{}) {
		logger = zerolog.Nop()
	}
	logger = logger.With().Str("component", "sse_handler").Logger()

	// Zero falls back to the default; a negative value is the explicit
	// "disable cap" signal so tests can opt out without depending on
	// magic numbers in env.
	maxPerIP := cfg.MaxPerIP
	if maxPerIP == 0 {
		maxPerIP = DefaultMaxPerIP
	}

	return &Handler{
		mux:               cfg.Multiplexer,
		js:                cfg.JetStream,
		heartbeatInterval: hb,
		logger:            logger,
		now:               time.Now,
		maxPerIP:          maxPerIP,
		conns:             make(map[string]int),
	}, nil
}

// Register mounts GET /events on the provided router group. The group
// inherits the gateway's outer middleware chain (rate-limit, CORS,
// recovery), so per-IP connection limiting is applied uniformly with
// the rest of the surface.
func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.GET("/events", h.Stream)
}

// acquireConn reserves a per-IP slot, returning false when the cap is
// already at the configured ceiling. ip == "" or maxPerIP <= 0 short-
// circuits to true: an empty IP usually signals unkeyable transport
// (test recorder, unix socket) where the cap is meaningless, and a
// non-positive cap is the explicit opt-out.
func (h *Handler) acquireConn(ip string) bool {
	if ip == "" || h.maxPerIP <= 0 {
		return true
	}
	h.connsMu.Lock()
	defer h.connsMu.Unlock()
	if h.conns[ip] >= h.maxPerIP {
		return false
	}
	h.conns[ip]++
	return true
}

// releaseConn drops a per-IP slot. Mirrors the acquireConn guards so
// release-without-acquire is a no-op rather than a counter underflow.
// Removes the map entry once the count returns to zero so an IP that
// connects once and disconnects doesn't permanently inflate the map.
func (h *Handler) releaseConn(ip string) {
	if ip == "" || h.maxPerIP <= 0 {
		return
	}
	h.connsMu.Lock()
	defer h.connsMu.Unlock()
	n := h.conns[ip]
	if n <= 1 {
		delete(h.conns, ip)
		return
	}
	h.conns[ip] = n - 1
}

// Stream is the GET /events handler. It runs for the lifetime of the
// SSE connection: registers a subscriber on the multiplexer, optionally
// replays history from JetStream, then loops forwarding multiplexer
// events and heartbeats until the client disconnects or the
// multiplexer stops.
func (h *Handler) Stream(c *gin.Context) {
	w := c.Writer
	r := c.Request

	// SSE requires a flushable writer. Gin's ResponseWriter implements
	// http.Flusher; if the underlying chain ever swaps it for one that
	// doesn't (e.g. a buffering middleware), refuse rather than silently
	// dropping every event.
	flusher, ok := w.(http.Flusher)
	if !ok {
		c.String(http.StatusInternalServerError, "streaming not supported")
		return
	}

	// Validate the subject filter BEFORE writing any response headers
	// so a bad pattern can be rejected with a 400 rather than swallowed
	// into the stream. The rule is narrow: a leading `>` or `*` would
	// fan every subject the multiplexer could ever see (or worse, every
	// subject NATS knows about) past the closed-prefix posture, so we
	// refuse those. Patterns that narrow within a registered prefix
	// (e.g. `agents.*`, `trades.>`) are accepted as-is.
	filter, err := ParseSubjectsCSV(c.Query("subjects"))
	if err != nil {
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.String(http.StatusBadRequest, `{"error":{"code":"invalid_subjects","message":"subjects parameter contains a disallowed wildcard"}}`)
		return
	}

	// Per-IP concurrent connection cap. A 429 with Retry-After: 5
	// tells the EventSource client to back off rather than reconnect
	// immediately. We register BEFORE any response bytes are written
	// so the body of a refused request stays a plain 429; the
	// corresponding deferred release must always fire even on early
	// returns below.
	clientIP := c.ClientIP()
	if !h.acquireConn(clientIP) {
		c.Header("Retry-After", "5")
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.String(http.StatusTooManyRequests, `{"error":{"code":"too_many_connections","message":"per-ip sse connection cap exceeded"}}`)
		return
	}
	defer h.releaseConn(clientIP)

	// Per-connection lifecycle span. Lives for the entire lifetime of
	// the SSE connection (could be hours), so this span is the parent
	// of any future per-event child spans. SpanKindServer is used
	// because the gin inbound span (SpanKindServer) is what
	// otel.GetTracerProvider hands us via the request context, and a
	// nested SpanKindServer would be incorrect — we are NOT a separate
	// server. Using a plain Internal kind because this span tracks the
	// internal handler's lifecycle, not a transport-level event.
	tracer := otel.Tracer("titular/gateway/sse")
	streamCtx, streamSpan := tracer.Start(c.Request.Context(), "sse.stream",
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(
			attribute.StringSlice("sse.subjects", []string(filter)),
		),
	)
	defer streamSpan.End()
	// Replace the request's context so child spans (replay, fanout)
	// hang off the lifecycle span.
	c.Request = c.Request.WithContext(streamCtx)

	// Clear the per-response write deadline. The shared *http.Server
	// sets WriteTimeout=30s for REST/GraphQL; that's the right posture
	// for short responses but it would kill every long-lived SSE
	// connection at the 30s heartbeat. http.NewResponseController is
	// the supported escape hatch (Go 1.20+) — it surfaces the underlying
	// transport's SetWriteDeadline, and a zero time disables it for
	// just this response without changing the server-wide setting.
	//
	// The read deadline is left alone: ReadTimeout (10s) only bounds
	// the request-line + header read on GET /events (no body), which
	// has already completed by the time we reach this code.
	rc := http.NewResponseController(w)
	if err := rc.SetWriteDeadline(time.Time{}); err != nil {
		// Some test transports (httptest's default ResponseRecorder
		// path) don't implement SetWriteDeadline — log at debug and
		// continue. In production the gin/net.http chain always does.
		h.logger.Debug().Err(err).Msg("sse: clear write deadline unsupported; continuing")
	}

	lastEventID := strings.TrimSpace(r.Header.Get("Last-Event-ID"))

	// Headers MUST be written before the first event. WriteHeader is
	// implicit on the first Write, so set the headers explicitly here
	// and then call Flush so curl / EventSource clients see the
	// connection as established.
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// X-Accel-Buffering disables nginx response buffering. Without
	// this, nginx will hold each event in a buffer until either the
	// buffer is full or the connection idles, which defeats the
	// real-time purpose of SSE.
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	sub, err := h.mux.Subscribe(SubscribeOptions{Filter: filter})
	if err != nil {
		// Subscribe only fails on Stop; the connection has already had
		// its 200 written, so emit a single error event and return.
		// The client EventSource will treat the close as a transient
		// error and reconnect, at which point we'll respond cleanly.
		writeComment(w, "service unavailable")
		flusher.Flush()
		return
	}
	defer sub.Close()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Optional history replay. We deliberately do replay BEFORE
	// attaching the live tail to avoid the "live event arrives during
	// replay" interleaving problem; the handler's subscriber is already
	// registered, so any live event arriving during replay is buffered
	// in the subscriber's channel and consumed naturally after the
	// replay loop exits. Buffer overflow during a long replay is
	// possible — the drop log line surfaces that to operators.
	if lastEventID != "" && h.js != nil {
		if err := h.replay(ctx, w, flusher, lastEventID, filter); err != nil {
			// Replay errors are non-fatal: log, send a comment, and
			// fall through to the live tail. The client will see
			// missing history but the connection stays alive.
			h.logger.Warn().Err(err).Str("last_event_id", lastEventID).Msg("sse replay failed")
			writeComment(w, "replay unavailable")
			flusher.Flush()
		}
	}

	// Live tail loop. Three sources of work:
	//
	//   1. ctx.Done    — client disconnected or server shutting down.
	//   2. sub.C       — multiplexer fanout.
	//   3. heartbeat   — periodic ping to keep proxies open.
	//
	// We use a single ticker for heartbeat rather than time.After per
	// loop so a high-throughput stream does not allocate a fresh
	// timer on every iteration.
	hbTicker := time.NewTicker(h.heartbeatInterval)
	defer hbTicker.Stop()

	// Initial heartbeat so the EventSource onopen fires even if the
	// stream is currently silent. Cheap, and removes a class of
	// "is the connection dead?" debugging questions.
	writeComment(w, "ping")
	flusher.Flush()

	for {
		select {
		case <-ctx.Done():
			streamSpan.SetAttributes(attribute.String("sse.close_reason", "client disconnect"))
			h.logSummary(sub, "client disconnect")
			return
		case ev, ok := <-sub.C:
			if !ok {
				// Multiplexer stopped; emit a comment and exit. The
				// client will treat connection close as transient and
				// reconnect, by which time the gateway will either be
				// back up or returning 503.
				writeComment(w, "server shutdown")
				flusher.Flush()
				streamSpan.SetAttributes(attribute.String("sse.close_reason", "multiplexer stopped"))
				h.logSummary(sub, "multiplexer stopped")
				return
			}
			if err := writeEvent(w, ev); err != nil {
				streamSpan.SetStatus(codes.Error, "write error")
				streamSpan.RecordError(err)
				h.logSummary(sub, "write error: "+err.Error())
				return
			}
			flusher.Flush()
		case <-hbTicker.C:
			writeComment(w, "ping")
			flusher.Flush()
		}
	}
}

// replay opens an ephemeral JetStream consumer starting at the
// sequence following lastEventID and drains every available message
// onto the response writer before returning. Filter is applied
// client-side because the JetStream consumer subscribes to the
// stream's full subject set; pre-filtering on the server requires a
// per-pattern consumer per call, which is wasteful for the common case
// of a small replay window.
//
// Returns nil if replay completes (or there's nothing to replay).
// Returns a non-nil error only on JetStream / network failure; a
// malformed Last-Event-ID is treated as "no replay" and returns nil.
func (h *Handler) replay(
	ctx context.Context,
	w gin.ResponseWriter,
	flusher http.Flusher,
	lastEventID string,
	filter Filter,
) error {
	seq, ok := parseJetStreamID(lastEventID)
	if !ok {
		// Non-JetStream id (e.g. hub:*) — no canonical replay anchor,
		// silently skip. The SSE spec lets the server treat
		// Last-Event-ID best-effort.
		return nil
	}

	// Bound the replay scan: a misconfigured client that sets
	// Last-Event-ID to "js:1" on a stream with millions of messages
	// would otherwise pin a goroutine for minutes. The cap is
	// enforced by counting drained messages; once the cap is hit we
	// emit a comment and stop replay, letting the live tail take
	// over. Operators who need bigger replays can paginate via the
	// REST API (#88), which is what the GraphQL transport advises.
	const maxReplay = 1000

	// Per-pattern subjects: when the client sent ?subjects=, replay
	// only those; otherwise replay every published prefix. Limiting
	// the JS consumer to the requested patterns is cheaper than
	// subscribing to the whole stream and discarding non-matches.
	subjects := filter
	if len(subjects) == 0 {
		subjects = Filter(append([]string(nil), SubjectPrefixes...))
	}

	for _, pattern := range subjects {
		drained, err := h.replayOne(ctx, w, flusher, pattern, seq+1, maxReplay)
		if err != nil {
			return err
		}
		if drained >= maxReplay {
			writeComment(w, "replay truncated")
			flusher.Flush()
			return nil
		}
	}
	return nil
}

// replayOne pulls historical messages on `pattern` starting at
// startSeq, writing each as an SSE event. Returns the number of
// messages drained or an error. Stops at:
//
//   - ctx.Done
//   - max messages drained
//   - JetStream "no more messages" / batch timeout (treated as
//     end-of-history; not an error)
func (h *Handler) replayOne(
	ctx context.Context,
	w gin.ResponseWriter,
	flusher http.Flusher,
	pattern string,
	startSeq uint64,
	max int,
) (int, error) {
	sub, err := h.js.PullSubscribe(
		pattern,
		"", // empty durable name → ephemeral consumer
		nats.StartSequence(startSeq),
	)
	if err != nil {
		return 0, fmt.Errorf("pull subscribe %q: %w", pattern, err)
	}
	// Always tear the consumer down; ephemeral consumers idle out on
	// the server, but Unsubscribe is the supported synchronous
	// release.
	defer func() { _ = sub.Unsubscribe() }()

	const batch = 32
	drained := 0
	for drained < max {
		select {
		case <-ctx.Done():
			return drained, nil
		default:
		}

		msgs, err := sub.Fetch(batch, nats.MaxWait(500*time.Millisecond))
		if err != nil {
			if errors.Is(err, nats.ErrTimeout) || errors.Is(err, context.DeadlineExceeded) {
				// No more historical messages on this pattern.
				return drained, nil
			}
			return drained, fmt.Errorf("fetch %q: %w", pattern, err)
		}
		if len(msgs) == 0 {
			return drained, nil
		}
		for _, m := range msgs {
			id := ""
			if md, mderr := m.Metadata(); mderr == nil && md != nil {
				id = formatJetStreamID(md.Sequence.Stream)
			}
			ev := Event{
				ID:      id,
				Subject: m.Subject,
				Data:    m.Data,
			}
			if werr := writeEvent(w, ev); werr != nil {
				return drained, werr
			}
			drained++
			// Ack so the ephemeral consumer's pending count drops; an
			// un-acked ephemeral that we abandon would still be torn
			// down by Unsubscribe, but explicit Ack is cheap and
			// keeps the JetStream server logs clean.
			_ = m.Ack()
			if drained >= max {
				flusher.Flush()
				return drained, nil
			}
		}
		flusher.Flush()
	}
	return drained, nil
}

// logSummary emits a single info line per terminated connection so an
// operator can see how many events were dropped for that client. The
// drop count is the most useful per-connection signal we have without
// adding metrics infrastructure.
func (h *Handler) logSummary(sub *Subscriber, reason string) {
	drops := uint64(0)
	if sub != nil && sub.Drops != nil {
		drops = sub.Drops.Load()
	}
	h.logger.Info().
		Uint64("drops", drops).
		Str("reason", reason).
		Msg("sse connection closed")
}

// writeEvent renders an SSE message frame onto w. The order of fields
// follows the WHATWG SSE spec: id, event, data, blank line. data is
// emitted as a single line because the publisher envelope is
// single-line JSON; if the publisher ever switches to multi-line JSON
// each newline-delimited chunk would need its own `data:` line.
func writeEvent(w gin.ResponseWriter, ev Event) error {
	var b strings.Builder
	b.Grow(len(ev.ID) + len(ev.Subject) + len(ev.Data) + 32)
	if ev.ID != "" {
		b.WriteString("id: ")
		b.WriteString(ev.ID)
		b.WriteByte('\n')
	}
	if ev.Subject != "" {
		b.WriteString("event: ")
		b.WriteString(ev.Subject)
		b.WriteByte('\n')
	}
	b.WriteString("data: ")
	// Defence in depth: a stray newline in the envelope would split
	// the SSE frame and confuse the client. The publisher emits
	// json.Marshal output which never contains a literal newline, but
	// we strip CR/LF as a belt-and-braces guard. Cost is one byte
	// scan per event.
	if hasLineBreak(ev.Data) {
		b.Write(stripLineBreaks(ev.Data))
	} else {
		b.Write(ev.Data)
	}
	b.WriteString("\n\n")
	if _, err := w.Write([]byte(b.String())); err != nil {
		return err
	}
	return nil
}

// writeComment emits an SSE comment frame (a heartbeat or human-
// readable status line). Per the spec, lines starting with `:` are
// ignored by the client.
func writeComment(w gin.ResponseWriter, msg string) {
	_, _ = w.Write([]byte(": " + msg + "\n\n"))
}

// hasLineBreak / stripLineBreaks: cheap CR/LF scrubbers for the data
// line. Kept as helpers so the writeEvent body stays readable.
func hasLineBreak(b []byte) bool {
	for _, c := range b {
		if c == '\n' || c == '\r' {
			return true
		}
	}
	return false
}

func stripLineBreaks(b []byte) []byte {
	out := make([]byte, 0, len(b))
	for _, c := range b {
		if c == '\n' || c == '\r' {
			continue
		}
		out = append(out, c)
	}
	return out
}
