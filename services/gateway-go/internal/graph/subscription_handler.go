package graph

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/store"
)

// Subscribe handles GET /graphql, performing the websocket upgrade for
// the `graphql-transport-ws` sub-protocol (the modern replacement for
// the older Apollo `graphql-ws`). The protocol is documented at
// https://github.com/enisdenjo/graphql-ws/blob/master/PROTOCOL.md.
//
// Only the subset that subscriptions need is implemented:
//
//   - connection_init / connection_ack
//   - subscribe (operation type must be `subscription`)
//   - next      (server -> client value frame)
//   - complete  (server -> client end-of-stream)
//   - error     (server -> client validation error)
//
// Mutations and queries over the websocket transport are intentionally
// unsupported here — the POST /graphql path covers those. A client that
// sends `query { ... }` over this socket gets an error frame back.
//
// The handler runs entirely in this goroutine; per-subscription
// fan-out goroutines are spawned by RunNewTrade etc. and shut down
// when the parent context is cancelled (which happens when the
// websocket is closed).
func (h *Handler) Subscribe(c *gin.Context) {
	// Only upgrade if the client actually requested a websocket. A bare
	// GET /graphql without the upgrade header gets a 400 — we do not
	// want to silently confuse browsers that accidentally hit this
	// path.
	if !websocket.IsWebSocketUpgrade(c.Request) {
		writeGraphQLError(c, http.StatusBadRequest, "websocket upgrade required for subscriptions")
		return
	}

	upgrader := websocket.Upgrader{
		Subprotocols:    []string{"graphql-transport-ws"},
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(_ *http.Request) bool {
			// Origin enforcement happens in the CORS middleware that
			// wraps this handler. Returning true here means the
			// websocket upgrade itself does not double-validate; the
			// invariant is that a Browser request that fails CORS
			// never reaches us, so any request that did reach us is
			// pre-authorised.
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, http.Header{
		"Sec-WebSocket-Protocol": []string{"graphql-transport-ws"},
	})
	if err != nil {
		// upgrader has already written an error response; nothing
		// more we can do.
		return
	}

	// Each connection runs in its own context derived from the
	// request. Cancelling ctx propagates to every active subscription
	// so the fan-out goroutines exit promptly when the client closes
	// the socket.
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	(&wsSession{
		conn:     conn,
		resolver: h.resolver,
		ctx:      ctx,
	}).run()
}

// ---------------------------------------------------------------------------
// graphql-transport-ws frame shapes
// ---------------------------------------------------------------------------

type wsFrame struct {
	ID      string          `json:"id,omitempty"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// subscribePayload is the inner shape of a subscribe-frame's payload.
type subscribePayload struct {
	Query         string         `json:"query"`
	OperationName string         `json:"operationName,omitempty"`
	Variables     map[string]any `json:"variables,omitempty"`
}

// errPayload is the inner shape of a server-side error frame. We keep
// the same shape as the HTTP response error envelope so client
// libraries can deserialise both with one struct.
type errPayload struct {
	Errors []graphqlErrorEnvelope `json:"errors"`
}

// ---------------------------------------------------------------------------
// wsSession — per-connection state machine
// ---------------------------------------------------------------------------

type wsSession struct {
	conn     *websocket.Conn
	resolver *Resolver
	ctx      context.Context

	mu           sync.Mutex
	wmu          sync.Mutex // serialises websocket Writes
	initialised  bool
	subs         map[string]context.CancelFunc
}

// run is the read loop. It reads frames sequentially and dispatches
// them to handler methods. The connection lives until the client
// disconnects or the server context is cancelled.
func (s *wsSession) run() {
	defer s.conn.Close()
	if s.subs == nil {
		s.subs = make(map[string]context.CancelFunc)
	}

	// Liveness: bound how long we wait for the first init frame so a
	// silent client cannot hold a connection indefinitely.
	_ = s.conn.SetReadDeadline(time.Now().Add(15 * time.Second))

	for {
		_, raw, err := s.conn.ReadMessage()
		if err != nil {
			s.cancelAll()
			return
		}
		// Reset the deadline after each successful read; the protocol
		// is bidirectional and the client may legitimately stay quiet
		// between subscribes.
		_ = s.conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		var f wsFrame
		if err := json.Unmarshal(raw, &f); err != nil {
			s.writeError("", "invalid frame")
			continue
		}
		switch f.Type {
		case "connection_init":
			s.handleInit()
		case "subscribe":
			s.handleSubscribe(f)
		case "complete":
			s.handleComplete(f.ID)
		case "ping":
			// graphql-transport-ws keepalive; respond with pong.
			s.writeFrame(wsFrame{Type: "pong"})
		case "pong":
			// Client reply to our ping; ignore.
		default:
			s.writeError(f.ID, "unsupported frame type")
		}
	}
}

func (s *wsSession) handleInit() {
	s.mu.Lock()
	if s.initialised {
		s.mu.Unlock()
		s.writeError("", "connection already initialised")
		return
	}
	s.initialised = true
	s.mu.Unlock()
	s.writeFrame(wsFrame{Type: "connection_ack"})
}

func (s *wsSession) handleSubscribe(f wsFrame) {
	s.mu.Lock()
	if !s.initialised {
		s.mu.Unlock()
		s.writeError(f.ID, "connection not initialised")
		return
	}
	if _, dup := s.subs[f.ID]; dup {
		s.mu.Unlock()
		s.writeError(f.ID, "duplicate subscription id")
		return
	}
	s.mu.Unlock()

	if f.ID == "" {
		s.writeError("", "subscription id is required")
		return
	}

	var p subscribePayload
	if err := json.Unmarshal(f.Payload, &p); err != nil {
		s.writeError(f.ID, "invalid subscribe payload")
		return
	}

	op, args, err := parseSubscriptionOp(p.Query, p.OperationName, p.Variables)
	if err != nil {
		s.writeError(f.ID, err.Error())
		return
	}

	subCtx, cancel := context.WithCancel(s.ctx)
	s.mu.Lock()
	s.subs[f.ID] = cancel
	s.mu.Unlock()

	go s.runSubscription(subCtx, f.ID, op, args)
}

// handleComplete cancels an active subscription. Per protocol, the
// server then emits no further frames for that id.
func (s *wsSession) handleComplete(id string) {
	s.mu.Lock()
	cancel, ok := s.subs[id]
	delete(s.subs, id)
	s.mu.Unlock()
	if ok {
		cancel()
	}
}

// cancelAll tears down every active subscription. Called when the
// websocket read loop exits.
func (s *wsSession) cancelAll() {
	s.mu.Lock()
	subs := s.subs
	s.subs = nil
	s.mu.Unlock()
	for _, c := range subs {
		c()
	}
}

// runSubscription drives the per-subscription channel onto the websocket.
// Each value is wrapped in a `next` frame whose payload is the
// canonical {data:{<field>: …}} shape. End-of-stream emits `complete`;
// errors emit an `error` frame and then stop without `complete`
// (per-protocol).
func (s *wsSession) runSubscription(ctx context.Context, id, op string, args map[string]any) {
	defer s.handleComplete(id) // remove from map even if upstream closed

	switch op {
	case "newTrade":
		token, _ := args["agentToken"].(string)
		ch, err := s.resolver.RunNewTrade(ctx, token)
		if err != nil {
			s.writeError(id, sanitiseGraphQLErr(err.Error()))
			return
		}
		s.pumpTrades(ctx, id, ch)
	case "jobUpdated":
		jobID, _ := args["id"].(string)
		ch, err := s.resolver.RunJobUpdated(ctx, jobID)
		if err != nil {
			s.writeError(id, sanitiseGraphQLErr(err.Error()))
			return
		}
		s.pumpJobs(ctx, id, ch)
	case "agentRegistered":
		ch, err := s.resolver.RunAgentRegistered(ctx)
		if err != nil {
			s.writeError(id, sanitiseGraphQLErr(err.Error()))
			return
		}
		s.pumpAgents(ctx, id, ch)
	default:
		s.writeError(id, "unknown subscription field")
	}
}

func (s *wsSession) pumpTrades(ctx context.Context, id string, ch <-chan store.Trade) {
	for {
		select {
		case <-ctx.Done():
			return
		case t, ok := <-ch:
			if !ok {
				s.writeFrame(wsFrame{ID: id, Type: "complete"})
				return
			}
			s.emitNext(id, "newTrade", marshalTrade(t))
		}
	}
}

func (s *wsSession) pumpJobs(ctx context.Context, id string, ch <-chan store.Job) {
	for {
		select {
		case <-ctx.Done():
			return
		case j, ok := <-ch:
			if !ok {
				s.writeFrame(wsFrame{ID: id, Type: "complete"})
				return
			}
			s.emitNext(id, "jobUpdated", marshalJob(j))
		}
	}
}

func (s *wsSession) pumpAgents(ctx context.Context, id string, ch <-chan store.Agent) {
	for {
		select {
		case <-ctx.Done():
			return
		case a, ok := <-ch:
			if !ok {
				s.writeFrame(wsFrame{ID: id, Type: "complete"})
				return
			}
			s.emitNext(id, "agentRegistered", marshalAgent(a))
		}
	}
}

// emitNext sends a `next` frame with the canonical {data: {<field>: v}}
// payload. The marshal* helpers produce JSON in the same shape the
// query path returns so a client can use one parser for both
// transports.
func (s *wsSession) emitNext(id, field string, payload json.RawMessage) {
	wrapped, err := json.Marshal(map[string]any{
		"data": map[string]json.RawMessage{field: payload},
	})
	if err != nil {
		s.writeError(id, "encode failed")
		return
	}
	s.writeFrame(wsFrame{ID: id, Type: "next", Payload: wrapped})
}

func (s *wsSession) writeFrame(f wsFrame) {
	s.wmu.Lock()
	defer s.wmu.Unlock()
	_ = s.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_ = s.conn.WriteJSON(f)
}

func (s *wsSession) writeError(id, msg string) {
	body, _ := json.Marshal(errPayload{
		Errors: []graphqlErrorEnvelope{{Message: msg}},
	})
	s.writeFrame(wsFrame{ID: id, Type: "error", Payload: body})
}

// parseSubscriptionOp is a minimal extractor for the subscription field
// name + arguments. We deliberately do not run the full graphql-go
// executor for subscriptions — that would require a stream-aware
// runtime which graphql-go does not provide. Instead we parse the
// query, find the `subscription { field(args) }` shape, and dispatch
// on the literal field name.
//
// This keeps the dependency surface small and matches how the
// Subscription resolvers are designed (one Run* helper per field).
// The trade-off is that we do not validate the operation against the
// schema's subscription type at this level — a typo in the field name
// surfaces as `unknown subscription field` rather than a typed
// validation error. That is acceptable for the alpha and is covered
// by the websocket integration test.
func parseSubscriptionOp(query, operationName string, vars map[string]any) (string, map[string]any, error) {
	const subKeyword = "subscription"
	q := stripWhitespace(query)
	if !startsWith(q, subKeyword) {
		return "", nil, errors.New("only subscription operations are supported on this transport")
	}
	// Find the first `{`, then the field name.
	body := q[len(subKeyword):]
	for i := 0; i < len(body); i++ {
		if body[i] == '{' {
			body = body[i+1:]
			break
		}
	}
	// field name is everything up to '(' or '{' or whitespace.
	end := len(body)
	for i, c := range body {
		if c == '(' || c == '{' || c == ' ' || c == '\n' || c == '\t' {
			end = i
			break
		}
	}
	field := body[:end]
	if field == "" {
		return "", nil, errors.New("subscription field name missing")
	}

	// Variables are passed verbatim. The schema's resolver path will
	// trim/lowercase as needed; we only need the field-name dispatch
	// here. operationName is informational on the websocket transport.
	_ = operationName
	return field, vars, nil
}

// stripWhitespace removes leading whitespace so the simple subscription
// parser can match the keyword without a regex. We only trim the
// front; whitespace inside the body is handled by the field-name
// extractor.
func stripWhitespace(s string) string {
	for i, c := range s {
		if c != ' ' && c != '\n' && c != '\t' && c != '\r' {
			return s[i:]
		}
	}
	return ""
}

// ---------------------------------------------------------------------------
// Per-type marshallers — produce the JSON shape the query path returns
// so the websocket payload schema matches the HTTP one.
// ---------------------------------------------------------------------------

func marshalTrade(t store.Trade) json.RawMessage {
	out := map[string]any{
		"id":          intToString(t.ID),
		"side":        t.Side,
		"trader":      lower(t.Trader),
		"curve":       lower(t.Curve),
		"fee":         t.Fee,
		"blockNumber": int(t.BlockNumber),
		"txHash":      lower(t.TxHash),
		"logIndex":    int(t.LogIndex),
		"createdAt":   t.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
	if t.AgentTokenPK != nil {
		out["agentTokenPK"] = int(*t.AgentTokenPK)
	}
	addOptString(out, "quoteIn", t.QuoteIn)
	addOptString(out, "agentOut", t.AgentOut)
	addOptString(out, "agentIn", t.AgentIn)
	addOptString(out, "quoteOut", t.QuoteOut)
	b, _ := json.Marshal(out)
	return b
}

func marshalJob(j store.Job) json.RawMessage {
	out := map[string]any{
		"id":          intToString(j.ID),
		"onChainId":   j.OnChainID,
		"principal":   lower(j.Principal),
		"jobType":     int(j.JobType),
		"budget":      j.Budget,
		"phase":       string(j.Phase),
		"blockNumber": int(j.BlockNumber),
		"txHash":      lower(j.TxHash),
		"logIndex":    int(j.LogIndex),
		"createdAt":   j.CreatedAt.UTC().Format(time.RFC3339Nano),
		"updatedAt":   j.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
	addOptString(out, "cloneAddress", j.CloneAddress)
	if j.AgentPK != nil {
		out["agentPK"] = int(*j.AgentPK)
	}
	addOptString(out, "agentId", j.AgentID)
	addOptString(out, "token", j.Token)
	if j.Deadline != nil {
		out["deadline"] = j.Deadline.UTC().Format(time.RFC3339Nano)
	}
	addOptString(out, "resultUri", j.ResultURI)
	b, _ := json.Marshal(out)
	return b
}

func marshalAgent(a store.Agent) json.RawMessage {
	out := map[string]any{
		"id":          intToString(a.ID),
		"kind":        string(a.Kind),
		"blockNumber": int(a.BlockNumber),
		"txHash":      lower(a.TxHash),
		"logIndex":    int(a.LogIndex),
		"createdAt":   a.CreatedAt.UTC().Format(time.RFC3339Nano),
		"updatedAt":   a.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
	addOptString(out, "agentId", a.AgentID)
	addOptString(out, "controller", a.Controller)
	addOptString(out, "creator", a.Creator)
	addOptString(out, "metadataUri", a.MetadataURI)
	addOptString(out, "capabilities", a.Capabilities)
	b, _ := json.Marshal(out)
	return b
}

func addOptString(m map[string]any, key, val string) {
	if val != "" {
		m[key] = val
	}
}

func intToString(n int64) string {
	// One use site, kept small to avoid pulling in strconv just for
	// the marshallers above.
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

func lower(s string) string {
	out := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		out[i] = c
	}
	return string(out)
}
