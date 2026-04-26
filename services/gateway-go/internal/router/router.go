package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/handlers"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/middleware"
)

// New builds and returns the HTTP router.
// agentHandlers must be non-nil; it is injected so tests can provide a
// fixture-backed instance without the singleton loader.
func New(agentHandlers *handlers.AgentHandlers) http.Handler {
	jh, err := handlers.NewJobHandlers()
	if err != nil {
		// Fatal startup: fixture files missing or malformed.
		panic(fmt.Sprintf("load job handlers: %v", err))
	}
	return NewWithHandlers(agentHandlers, jh)
}

// NewWithHandlers builds and returns the HTTP router with explicit handler
// injection. Preferred in tests so fixtures are controlled by the caller.
func NewWithHandlers(agentHandlers *handlers.AgentHandlers, jobHandlers *handlers.JobHandlers) http.Handler {
	r := chi.NewRouter()

	// Middleware chain — order matters:
	//   Recovery must wrap everything so panics surface as 500.
	//   RequestID before Logger so the log line includes the id.
	r.Use(middleware.Recovery)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger("gateway"))

	r.Route("/v1", func(r chi.Router) {
		// M2 — agents (do not modify)
		r.Get("/agents", agentHandlers.ListAgents)
		r.Get("/agents/{id}", agentHandlers.GetAgent)
		r.Get("/agents/{id}/trades", agentHandlers.ListAgentTrades)

		// M3 — ACP jobs
		r.Get("/jobs", jobHandlers.ListJobs)
		r.Post("/jobs/prepare", jobHandlers.PrepareJob)
		r.Get("/jobs/{id}", jobHandlers.GetJob)
	})

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	return r
}
