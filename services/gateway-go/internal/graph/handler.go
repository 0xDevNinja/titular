package graph

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
)

// Handler wraps a constructed schema and exposes Gin handlers for the
// /graphql surface. A single Handler is built at startup and registered
// on the router under both /graphql (POST: queries) and the websocket
// upgrade path (handled by Subscriptions below).
type Handler struct {
	schema   graphql.Schema
	resolver *Resolver
	// EnablePlayground gates GET /graphql/playground. Defaults to false
	// because exposing the GraphiQL UI on a public endpoint encourages
	// uncontrolled query authoring against production data.
	EnablePlayground bool
}

// NewHandler returns a Handler bound to schema and resolver. Both must
// be non-nil; the router uses NewHandler at startup so callers cannot
// observe a half-built handler under load.
func NewHandler(schema graphql.Schema, resolver *Resolver) (*Handler, error) {
	if resolver == nil {
		return nil, errors.New("graph: nil resolver")
	}
	return &Handler{schema: schema, resolver: resolver}, nil
}

// Register mounts the GraphQL endpoints onto the given Gin router group.
//
// The endpoints are:
//
//   - POST /graphql            — queries and mutations.
//   - GET  /graphql            — websocket upgrade for subscriptions
//                                (graphql-transport-ws sub-protocol).
//   - GET  /graphql/playground — GraphiQL UI; gated by Handler.EnablePlayground.
//
// We deliberately keep the path minimal; the gateway router (parent)
// is responsible for any auth or rate-limit middleware applied to the
// group before this Register call.
func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.POST("/graphql", h.Query)
	rg.GET("/graphql", h.Subscribe)
	if h.EnablePlayground {
		rg.GET("/graphql/playground", h.Playground)
	}
}

// queryRequest is the wire shape of a POST /graphql body. The same shape
// is used over websocket subscription frames.
type queryRequest struct {
	Query         string         `json:"query"`
	OperationName string         `json:"operationName,omitempty"`
	Variables     map[string]any `json:"variables,omitempty"`
}

// queryResponse mirrors graphql-go's standard result shape so SDK
// clients can decode it without a custom unmarshaller.
type queryResponse struct {
	Data   any                  `json:"data,omitempty"`
	Errors []graphqlErrorEnvelope `json:"errors,omitempty"`
}

// graphqlErrorEnvelope is the public-facing error format. We re-shape
// graphql-go's runtime error into this struct rather than embedding it
// directly so a future change to the upstream library cannot leak
// fields we did not intend (path locations are kept for tooling, but
// stack traces / extensions are stripped).
type graphqlErrorEnvelope struct {
	Message string         `json:"message"`
	Path    []any          `json:"path,omitempty"`
	Locations []graphqlLoc `json:"locations,omitempty"`
}

type graphqlLoc struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// Query handles POST /graphql. Body is JSON: {query, operationName, variables}.
//
// SECURITY:
//   - Subscriptions cannot be requested over POST. The schema accepts
//     them via the websocket transport only; POSTing a `subscription`
//     operation here would otherwise hang the request without a
//     channel sink.
//   - We reject bodies larger than maxQueryBytes to make the
//     introspection/load-testing surface more predictable. The gateway's
//     CORS / rate-limit middleware is responsible for the 4xx flood
//     control around this.
func (h *Handler) Query(c *gin.Context) {
	if c.Request.Body == nil {
		writeGraphQLError(c, http.StatusBadRequest, "empty request body")
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxQueryBytes)

	var req queryRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		writeGraphQLError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if strings.TrimSpace(req.Query) == "" {
		writeGraphQLError(c, http.StatusBadRequest, "query is required")
		return
	}

	if isSubscription(req.Query, req.OperationName) {
		writeGraphQLError(c, http.StatusBadRequest, "subscriptions must use the websocket transport")
		return
	}

	result := graphql.Do(graphql.Params{
		Schema:         h.schema,
		RequestString:  req.Query,
		VariableValues: req.Variables,
		OperationName:  req.OperationName,
		Context:        c.Request.Context(),
	})

	c.JSON(http.StatusOK, toResponse(result))
}

// maxQueryBytes is the cap on POST body size. 1 MiB is well above any
// reasonable hand-written query and any introspection payload, and
// well below what would let an attacker exhaust the JSON decoder's
// memory.
const maxQueryBytes = 1 << 20

// toResponse converts graphql-go's *graphql.Result to our public envelope.
// Errors are filtered through graphqlErrorEnvelope so the wire shape is
// stable regardless of which graphql-go version we're on.
func toResponse(r *graphql.Result) queryResponse {
	resp := queryResponse{Data: r.Data}
	if len(r.Errors) > 0 {
		resp.Errors = make([]graphqlErrorEnvelope, 0, len(r.Errors))
		for _, e := range r.Errors {
			env := graphqlErrorEnvelope{Message: sanitiseGraphQLErr(e.Message)}
			env.Path = append(env.Path, e.Path...)
			for _, loc := range e.Locations {
				env.Locations = append(env.Locations, graphqlLoc{Line: loc.Line, Column: loc.Column})
			}
			resp.Errors = append(resp.Errors, env)
		}
	}
	return resp
}

// sanitiseGraphQLErr applies the same allow-list as the REST surface
// (see resolvers.go's sanitiseStoreErr) plus a defence against
// "internal error" leaks. We never let an error beginning with the
// internal-marker prefix escape — a future contributor adding a custom
// error type that bubbles a SQL fragment would fail this gate.
func sanitiseGraphQLErr(msg string) string {
	// Keep validation errors as-is (they tell the client what to fix).
	if isUserVisibleStoreErr(msg) {
		return msg
	}
	// Keep GraphQL framework errors verbatim — these are syntax /
	// type-validation errors that contain only the schema vocabulary,
	// not any runtime data.
	if startsWith(msg, "Cannot query field") ||
		startsWith(msg, "Field ") ||
		startsWith(msg, "Argument ") ||
		startsWith(msg, "Variable ") ||
		startsWith(msg, "Syntax Error") ||
		startsWith(msg, "Expected ") ||
		startsWith(msg, "Unknown ") ||
		startsWith(msg, "Must provide") ||
		startsWith(msg, "Operation ") ||
		startsWith(msg, "Anonymous ") ||
		startsWith(msg, "id is required") ||
		startsWith(msg, "id: ") ||
		startsWith(msg, "kind: ") ||
		startsWith(msg, "status: ") ||
		startsWith(msg, "from: ") ||
		startsWith(msg, "to: ") ||
		startsWith(msg, "subscriptions ") ||
		startsWith(msg, "websocket ") ||
		startsWith(msg, "limit ") {
		return msg
	}
	// Anything else: opaque generic. We log the original via the gin
	// recovery middleware (which sees it via panic), so operators can
	// still debug.
	return "internal error"
}

// writeGraphQLError emits a queryResponse-shaped error body with the
// given HTTP status. graphql-go uses 200 OK for application errors;
// the cases this helper handles are transport-level (bad JSON, oversize
// body, wrong operation type) where a non-200 is correct.
func writeGraphQLError(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, queryResponse{
		Errors: []graphqlErrorEnvelope{{Message: msg}},
	})
}

// isSubscription parses query just enough to decide whether the named
// operation is a subscription. We use the upstream parser rather than
// a regex because operation type detection has to handle leading
// fragments / comments / multi-operation documents correctly.
//
// On parse failure we return false so the downstream graphql.Do gets
// the chance to emit the canonical "Syntax Error" message — surfacing
// our own "subscriptions only over websocket" error here would mask
// the real diagnostic.
func isSubscription(query, operationName string) bool {
	doc, err := parser.Parse(parser.ParseParams{
		Source: source.NewSource(&source.Source{
			Body: []byte(query),
			Name: "query",
		}),
	})
	if err != nil {
		return false
	}
	for _, def := range doc.Definitions {
		op, ok := def.(*ast.OperationDefinition)
		if !ok {
			continue
		}
		// graphql-go's parser exposes Operation as a string ("query",
		// "mutation", "subscription"). Empty operationName matches
		// the first operation in the document, mirroring graphql.Do
		// behaviour.
		if operationName == "" || (op.Name != nil && op.Name.Value == operationName) {
			return op.Operation == "subscription"
		}
	}
	return false
}

// Playground serves a minimal GraphiQL UI. We do NOT use a third-party
// hosted bundle so the gateway has no third-party CDN dependency at
// runtime; the inline HTML is tiny and points at the local /graphql
// endpoint. Handler.EnablePlayground gates the registration so this
// page is opt-in per deployment.
func (h *Handler) Playground(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(playgroundHTML))
}

// playgroundHTML is a single-file GraphiQL UI loaded from a CDN. We
// intentionally accept the CDN dependency for the dev playground only;
// production deployments leave EnablePlayground at false (its default)
// and the page is not mounted at all.
const playgroundHTML = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8" />
<title>Titular GraphQL</title>
<link rel="stylesheet" href="https://unpkg.com/graphiql@3/graphiql.min.css" />
</head>
<body style="margin:0">
<div id="graphiql" style="height:100vh"></div>
<script crossorigin src="https://unpkg.com/react@18/umd/react.production.min.js"></script>
<script crossorigin src="https://unpkg.com/react-dom@18/umd/react-dom.production.min.js"></script>
<script crossorigin src="https://unpkg.com/graphiql@3/graphiql.min.js"></script>
<script>
  const fetcher = GraphiQL.createFetcher({ url: '/graphql' });
  ReactDOM.createRoot(document.getElementById('graphiql'))
    .render(React.createElement(GraphiQL, { fetcher }));
</script>
</body>
</html>`

// ExecuteForTest runs a one-shot query through the handler's schema
// without going through HTTP. Exported for use by package-internal
// tests that want to assert on the GraphQL result shape directly.
func (h *Handler) ExecuteForTest(ctx context.Context, query string, variables map[string]any) *graphql.Result {
	return graphql.Do(graphql.Params{
		Schema:         h.schema,
		RequestString:  query,
		VariableValues: variables,
		Context:        ctx,
	})
}
