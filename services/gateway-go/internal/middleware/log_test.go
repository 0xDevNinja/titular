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

// captureLogger returns a zerolog logger wired to a buffer the test can
// inspect. Using zerolog's JSON output keeps assertions structural.
func captureLogger() (zerolog.Logger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	return zerolog.New(buf), buf
}

// firstLogLine decodes the first log line in buf into a generic map so tests
// can assert on individual fields without coupling to the full schema.
func firstLogLine(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	line, _, _ := strings.Cut(buf.String(), "\n")
	if line == "" {
		t.Fatalf("no log line emitted; buffer = %q", buf.String())
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(line), &out); err != nil {
		t.Fatalf("decode log line %q: %v", line, err)
	}
	return out
}

func TestLog_EmitsRequiredFields(t *testing.T) {
	logger, buf := captureLogger()

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Log(middleware.LogConfig{Logger: logger, Service: "gateway-test"}))
	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	rec2 := firstLogLine(t, buf)
	for _, field := range []string{
		"service", "request_id", "method", "path",
		"status", "latency_ms", "client_ip", "bytes",
	} {
		if _, ok := rec2[field]; !ok {
			t.Errorf("missing field %q in log line: %v", field, rec2)
		}
	}
	if rec2["service"] != "gateway-test" {
		t.Errorf("service: got %v, want gateway-test", rec2["service"])
	}
	if rec2["method"] != http.MethodGet {
		t.Errorf("method: got %v", rec2["method"])
	}
	if rec2["path"] != "/ping" {
		t.Errorf("path: got %v", rec2["path"])
	}
	if rec2["status"] != float64(http.StatusOK) {
		t.Errorf("status: got %v", rec2["status"])
	}
	if id, ok := rec2["request_id"].(string); !ok || id == "" {
		t.Errorf("request_id missing or empty in log line: %v", rec2["request_id"])
	}
}

func TestLog_LevelFollowsStatus(t *testing.T) {
	cases := []struct {
		name      string
		status    int
		wantLevel string
	}{
		{"2xx is info", http.StatusOK, "info"},
		{"4xx is warn", http.StatusBadRequest, "warn"},
		{"5xx is error", http.StatusInternalServerError, "error"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			logger, buf := captureLogger()
			r := gin.New()
			r.Use(middleware.RequestID())
			r.Use(middleware.Log(middleware.LogConfig{Logger: logger, Service: "svc"}))
			r.GET("/x", func(c *gin.Context) { c.Status(tc.status) })

			req := httptest.NewRequest(http.MethodGet, "/x", nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			line := firstLogLine(t, buf)
			if line["level"] != tc.wantLevel {
				t.Fatalf("level: got %v, want %s", line["level"], tc.wantLevel)
			}
		})
	}
}

func TestLog_SkipsConfiguredPaths(t *testing.T) {
	logger, buf := captureLogger()
	r := gin.New()
	r.Use(middleware.Log(middleware.LogConfig{
		Logger:    logger,
		Service:   "svc",
		SkipPaths: []string{"/healthz"},
	}))
	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if buf.Len() != 0 {
		t.Fatalf("expected no log output for skipped path; got %q", buf.String())
	}
}

func TestLog_DefaultsServiceWhenEmpty(t *testing.T) {
	logger, buf := captureLogger()
	r := gin.New()
	r.Use(middleware.Log(middleware.LogConfig{Logger: logger}))
	r.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if got := firstLogLine(t, buf)["service"]; got != "gateway" {
		t.Fatalf("default service: got %v, want gateway", got)
	}
}

func TestLog_LevelOverrideSuppressesQuieterEvents(t *testing.T) {
	// Parent logger inherits the package default (no level filter) so the
	// 200 INFO line below would normally render. With LogConfig.Level set to
	// Warn, the middleware's logger is filtered to >=warn — INFO must drop.
	logger, buf := captureLogger()
	r := gin.New()
	r.Use(middleware.Log(middleware.LogConfig{
		Logger: logger,
		Level:  zerolog.WarnLevel,
	}))
	r.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if buf.Len() != 0 {
		t.Fatalf("expected no log line at warn-level filter for 200; got %q", buf.String())
	}

	// 5xx still escalates above the floor and should render.
	r2 := gin.New()
	logger2, buf2 := captureLogger()
	r2.Use(middleware.Log(middleware.LogConfig{
		Logger: logger2,
		Level:  zerolog.WarnLevel,
	}))
	r2.GET("/x", func(c *gin.Context) { c.Status(http.StatusInternalServerError) })

	req2 := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec2 := httptest.NewRecorder()
	r2.ServeHTTP(rec2, req2)

	if buf2.Len() == 0 {
		t.Fatal("expected a log line for 500 at warn-level filter; got none")
	}
}
