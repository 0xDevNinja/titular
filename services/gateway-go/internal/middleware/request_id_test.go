package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/middleware"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// newRequestIDEngine builds a gin engine with only the RequestID middleware
// attached and a single test endpoint that echoes the request id from the
// context as the response body. This isolates the unit-under-test from the
// rest of the chain.
func newRequestIDEngine() *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID())
	r.GET("/echo", func(c *gin.Context) {
		c.String(http.StatusOK, middleware.GetRequestIDFromGin(c))
	})
	return r
}

func TestRequestID_GeneratesUUIDWhenAbsent(t *testing.T) {
	r := newRequestIDEngine()

	req := httptest.NewRequest(http.MethodGet, "/echo", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}

	got := rec.Header().Get(middleware.RequestIDHeader)
	if got == "" {
		t.Fatalf("expected %s header, got empty", middleware.RequestIDHeader)
	}
	if _, err := uuid.Parse(got); err != nil {
		t.Fatalf("expected UUID, got %q: %v", got, err)
	}

	if body := rec.Body.String(); body != got {
		t.Fatalf("context id %q != header id %q", body, got)
	}
}

func TestRequestID_ReusesValidInboundUUID(t *testing.T) {
	r := newRequestIDEngine()

	want := uuid.NewString()
	req := httptest.NewRequest(http.MethodGet, "/echo", nil)
	req.Header.Set(middleware.RequestIDHeader, want)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	got := rec.Header().Get(middleware.RequestIDHeader)
	if got != want {
		t.Fatalf("expected echoed id %q, got %q", want, got)
	}
}

func TestRequestID_RejectsMalformedInbound(t *testing.T) {
	r := newRequestIDEngine()

	cases := []string{
		"not-a-uuid",
		"",
		// Newline injection — a hostile client should not be able to splice
		// content into the response header.
		"abcd\nX-Evil: 1",
	}

	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/echo", nil)
			if in != "" {
				req.Header.Set(middleware.RequestIDHeader, in)
			}
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			got := rec.Header().Get(middleware.RequestIDHeader)
			if got == in && in != "" {
				t.Fatalf("expected fresh UUID, got echoed %q", got)
			}
			if _, err := uuid.Parse(got); err != nil {
				t.Fatalf("expected UUID replacement, got %q", got)
			}
		})
	}
}

func TestRequestID_AvailableViaRequestContext(t *testing.T) {
	// Verify that handlers that only see the underlying *http.Request (e.g.
	// chi handlers behind gin.WrapH) can still resolve the request id.
	r := gin.New()
	r.Use(middleware.RequestID())
	r.GET("/echo", func(c *gin.Context) {
		id := middleware.GetRequestID(c.Request.Context())
		c.String(http.StatusOK, id)
	})

	req := httptest.NewRequest(http.MethodGet, "/echo", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if got := rec.Body.String(); got == "" {
		t.Fatal("expected request id on request.Context, got empty")
	}
}

func TestGetRequestID_NilContext(t *testing.T) {
	if got := middleware.GetRequestID(nil); got != "" {
		t.Fatalf("expected empty string for nil context, got %q", got)
	}
}
