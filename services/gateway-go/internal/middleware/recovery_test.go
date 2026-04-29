package middleware_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/middleware"
)

func TestRecovery_ConvertsPanicTo500(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := zerolog.New(buf)

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Recovery(logger))
	r.GET("/boom", func(_ *gin.Context) {
		panic("kaboom")
	})

	req := httptest.NewRequest(http.MethodGet, "/boom", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want 500", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `"internal_server_error"`) {
		t.Fatalf("body missing canonical error code: %q", body)
	}
	// MUST NOT leak the panic value or stack to the client.
	if strings.Contains(body, "kaboom") {
		t.Fatalf("response leaked panic value: %q", body)
	}
}

func TestRecovery_LogsStackAndRequestID(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := zerolog.New(buf)

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Recovery(logger))
	r.GET("/boom", func(_ *gin.Context) { panic("explode") })

	req := httptest.NewRequest(http.MethodGet, "/boom", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	line := buf.String()
	if line == "" {
		t.Fatal("expected a log line, got none")
	}
	var rec2 map[string]any
	if err := json.Unmarshal([]byte(strings.TrimRight(line, "\n")), &rec2); err != nil {
		t.Fatalf("decode log line %q: %v", line, err)
	}

	if rec2["level"] != "error" {
		t.Fatalf("level: got %v, want error", rec2["level"])
	}
	if rec2["panic"] != "explode" {
		t.Fatalf("panic: got %v, want explode", rec2["panic"])
	}
	if id, ok := rec2["request_id"].(string); !ok || id == "" {
		t.Fatalf("expected request_id in log line, got %v", rec2["request_id"])
	}
	// The stack is base64-encoded by zerolog.Bytes; we only assert it is
	// present and non-empty so the test is robust to encoding tweaks.
	if rec2["stack"] == nil || rec2["stack"] == "" {
		t.Fatal("expected non-empty stack in log line")
	}
}

func TestRecovery_NoPanicLetsResponsePassThrough(t *testing.T) {
	r := gin.New()
	r.Use(middleware.Recovery(zerolog.Nop()))
	r.GET("/ok", func(c *gin.Context) { c.String(http.StatusOK, "fine") })

	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}
	if rec.Body.String() != "fine" {
		t.Fatalf("body: got %q, want fine", rec.Body.String())
	}
}
