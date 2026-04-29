package openapi_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"

	gatewaydocs "github.com/0xDevNinja/titular/services/gateway-go/docs"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/openapi"
)

func init() { gin.SetMode(gin.TestMode) }

// TestSpecValidatesAgainstOpenAPI30 runs the embedded spec through the
// kin-openapi loader, which applies the OpenAPI 3.0 JSON Schema and the
// library's own structural checks. A failure here means a handler
// annotation produced an invalid path / parameter / schema and the spec
// would be rejected by SDK generators.
func TestSpecValidatesAgainstOpenAPI30(t *testing.T) {
	t.Parallel()

	if len(gatewaydocs.SpecJSON) == 0 {
		t.Fatal("SpecJSON is empty — run `go run ./scripts/openapi-gen`")
	}

	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(gatewaydocs.SpecJSON)
	if err != nil {
		t.Fatalf("load spec: %v", err)
	}
	if err := doc.Validate(context.Background()); err != nil {
		t.Fatalf("openapi 3.0 validation: %v", err)
	}

	if doc.OpenAPI == "" || !strings.HasPrefix(doc.OpenAPI, "3.") {
		t.Errorf("openapi field = %q; want a 3.x version", doc.OpenAPI)
	}
	if doc.Info == nil || doc.Info.Title == "" {
		t.Error("info.title is required and must be non-empty")
	}
}

// TestSpecCoversAllRoutes asserts every gateway route is documented.
// The list is the canonical set the milestone PR closes — keep this in
// lockstep with router/router.go so a newly added route can't silently
// ship without docs.
func TestSpecCoversAllRoutes(t *testing.T) {
	t.Parallel()

	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(gatewaydocs.SpecJSON)
	if err != nil {
		t.Fatalf("load spec: %v", err)
	}

	required := []struct {
		path   string
		method string
	}{
		{"/healthz", "get"},
		{"/api/v1/agents", "get"},
		{"/api/v1/agents/{id}", "get"},
		{"/api/v1/trades", "get"},
		{"/api/v1/jobs", "get"},
		{"/api/v1/stats", "get"},
		{"/auth/siwe/nonce", "post"},
		{"/auth/siwe/verify", "post"},
		{"/auth/siws/nonce", "post"},
		{"/auth/siws/verify", "post"},
		{"/auth/logout", "post"},
	}

	if doc.Paths == nil {
		t.Fatal("spec has no paths map")
	}
	for _, r := range required {
		item := doc.Paths.Find(r.path)
		if item == nil {
			t.Errorf("missing path %q", r.path)
			continue
		}
		op := item.GetOperation(strings.ToUpper(r.method))
		if op == nil {
			t.Errorf("missing %s on %q", r.method, r.path)
			continue
		}
		if op.Summary == "" {
			t.Errorf("%s %s: empty summary; the swag annotation is incomplete", r.method, r.path)
		}
	}
}

// TestCustomScalarFormats checks that the openapi-gen post-processor
// applied the protocol-specific format hints. If a schema name
// changes, this test also flags drift between the wire types and the
// patcher's case list.
func TestCustomScalarFormats(t *testing.T) {
	t.Parallel()

	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(gatewaydocs.SpecJSON)
	if err != nil {
		t.Fatalf("load spec: %v", err)
	}

	mustField := func(schemaName, field string) *openapi3.Schema {
		s, ok := doc.Components.Schemas[schemaName]
		if !ok || s.Value == nil {
			t.Fatalf("schema %q missing", schemaName)
		}
		p, ok := s.Value.Properties[field]
		if !ok || p.Value == nil {
			t.Fatalf("%s.%s missing", schemaName, field)
		}
		return p.Value
	}

	// store.Trade.tx_hash → bytes32 / 0x-hex 64
	tx := mustField("store.Trade", "tx_hash")
	if tx.Format != "bytes32" {
		t.Errorf("tx_hash format = %q; want bytes32", tx.Format)
	}
	if tx.Pattern != `^0x[a-fA-F0-9]{64}$` {
		t.Errorf("tx_hash pattern = %q; want 32-byte hex", tx.Pattern)
	}

	// store.Trade.trader → address / 0x-hex 40
	trader := mustField("store.Trade", "trader")
	if trader.Format != "address" {
		t.Errorf("trader format = %q; want address", trader.Format)
	}

	// store.Job.budget → bigint
	budget := mustField("store.Job", "budget")
	if budget.Format != "bigint" {
		t.Errorf("budget format = %q; want bigint", budget.Format)
	}
	if budget.Pattern != `^[0-9]+$` {
		t.Errorf("budget pattern = %q; want decimal-only", budget.Pattern)
	}

	// store.Agent.created_at → date-time (round-trip from time.Time)
	created := mustField("store.Agent", "created_at")
	if created.Format != "date-time" {
		t.Errorf("created_at format = %q; want date-time", created.Format)
	}
}

// TestServeJSONReturnsEmbedded asserts the served bytes are exactly the
// in-binary spec, byte-for-byte. SDK generators that diff against a
// pinned snapshot rely on this stability.
func TestServeJSONReturnsEmbedded(t *testing.T) {
	t.Parallel()

	spec, err := openapi.NewSpec(gatewaydocs.SpecJSON, gatewaydocs.SpecYAML, false)
	if err != nil {
		t.Fatalf("new spec: %v", err)
	}

	r := newEngineWithSpec(spec)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d; want 200", rr.Code)
	}
	if got, want := rr.Header().Get("Content-Type"), "application/json"; !strings.HasPrefix(got, want) {
		t.Errorf("content-type = %q; want prefix %q", got, want)
	}

	body, _ := io.ReadAll(rr.Body)
	if string(body) != string(gatewaydocs.SpecJSON) {
		t.Error("served body does not match embedded SpecJSON")
	}

	// The served bytes must still parse as a valid JSON object.
	var generic map[string]any
	if err := json.Unmarshal(body, &generic); err != nil {
		t.Fatalf("invalid json on the wire: %v", err)
	}
}

// TestServeYAMLReturnsEmbedded mirrors the JSON test for the YAML form.
func TestServeYAMLReturnsEmbedded(t *testing.T) {
	t.Parallel()

	spec, err := openapi.NewSpec(gatewaydocs.SpecJSON, gatewaydocs.SpecYAML, false)
	if err != nil {
		t.Fatalf("new spec: %v", err)
	}

	r := newEngineWithSpec(spec)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/swagger/doc.yaml", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d; want 200", rr.Code)
	}
	if got := rr.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/yaml") {
		t.Errorf("content-type = %q; want application/yaml", got)
	}

	body, _ := io.ReadAll(rr.Body)
	if string(body) != string(gatewaydocs.SpecYAML) {
		t.Error("served body does not match embedded SpecYAML")
	}
}

// TestSwaggerUIDisabledByDefault is the security gate: a default-config
// gateway must NOT serve the HTML viewer. We assert 404 on /swagger/
// when EnableUI is false.
func TestSwaggerUIDisabledByDefault(t *testing.T) {
	t.Parallel()

	spec, err := openapi.NewSpec(gatewaydocs.SpecJSON, gatewaydocs.SpecYAML, false)
	if err != nil {
		t.Fatalf("new spec: %v", err)
	}

	r := newEngineWithSpec(spec)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/swagger/", nil))
	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d; want 404 when UI is gated off", rr.Code)
	}
}

// TestSwaggerUIEnabledServesHTML asserts the operator opt-in path. With
// EnableUI true, GET /swagger/ returns the HTML viewer page.
func TestSwaggerUIEnabledServesHTML(t *testing.T) {
	t.Parallel()

	spec, err := openapi.NewSpec(gatewaydocs.SpecJSON, gatewaydocs.SpecYAML, true)
	if err != nil {
		t.Fatalf("new spec: %v", err)
	}

	r := newEngineWithSpec(spec)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/swagger/", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d; want 200", rr.Code)
	}
	if got := rr.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/html") {
		t.Errorf("content-type = %q; want text/html", got)
	}
	body, _ := io.ReadAll(rr.Body)
	if !strings.Contains(string(body), `id="swagger-ui"`) {
		t.Error("served HTML missing swagger-ui mount node")
	}
	if !strings.Contains(string(body), `/swagger/doc.json`) {
		t.Error("served HTML does not point at the embedded spec")
	}
}

// TestSwaggerUIRedirect ensures the no-slash form 301s to the canonical
// trailing-slash form so relative URLs in the page resolve correctly.
func TestSwaggerUIRedirect(t *testing.T) {
	t.Parallel()

	spec, err := openapi.NewSpec(gatewaydocs.SpecJSON, gatewaydocs.SpecYAML, true)
	if err != nil {
		t.Fatalf("new spec: %v", err)
	}

	r := newEngineWithSpec(spec)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/swagger", nil))
	if rr.Code != http.StatusMovedPermanently {
		t.Errorf("status = %d; want 301", rr.Code)
	}
	if got := rr.Header().Get("Location"); got != "/swagger/" {
		t.Errorf("location = %q; want /swagger/", got)
	}
}

// TestEnabledFromEnvParsing nails down the case-insensitive "true"
// contract so a deploy that exports GATEWAY_SWAGGER_UI=TRUE is
// equivalent to GATEWAY_SWAGGER_UI=true (and anything else is off).
func TestEnabledFromEnvParsing(t *testing.T) {
	t.Parallel()

	cases := map[string]bool{
		"":         false,
		"false":    false,
		"FALSE":    false,
		"0":        false,
		"yes":      false,
		"  true  ": true,
		"true":     true,
		"TRUE":     true,
		"True":     true,
	}
	for in, want := range cases {
		if got := openapi.EnabledFromEnv(in); got != want {
			t.Errorf("EnabledFromEnv(%q) = %v; want %v", in, got, want)
		}
	}
}

// TestNewSpecRejectsEmpty asserts startup-time validation — shipping a
// binary built without the embedded spec must fail, not silently serve
// an empty document.
func TestNewSpecRejectsEmpty(t *testing.T) {
	t.Parallel()

	if _, err := openapi.NewSpec(nil, []byte("y"), false); err == nil {
		t.Error("nil json spec accepted")
	}
	if _, err := openapi.NewSpec([]byte("{}"), nil, false); err == nil {
		t.Error("nil yaml spec accepted")
	}
}

// newEngineWithSpec wires a minimal Gin engine with just the openapi
// routes — keeps the tests independent of router.NewWithConfig and its
// dependencies (auth, store, NATS).
func newEngineWithSpec(spec *openapi.Spec) *gin.Engine {
	engine := gin.New()
	spec.Register(engine.Group(""))
	return engine
}
