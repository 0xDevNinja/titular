package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/handlers"
	"github.com/0xDevNinja/titular/services/gateway-go/internal/middleware"
)

// New builds and returns the HTTP router.
// agentHandlers must be non-nil; it is injected so tests can provide a
// fixture-backed instance without the singleton loader.
func New(agentHandlers *handlers.AgentHandlers) http.Handler {
	r := chi.NewRouter()

	// Middleware chain — order matters:
	//   Recovery must wrap everything so panics surface as 500.
	//   RequestID before Logger so the log line includes the id.
	r.Use(middleware.Recovery)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger("gateway"))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/agents", agentHandlers.ListAgents)
		r.Get("/agents/{id}", agentHandlers.GetAgent)
		r.Get("/agents/{id}/trades", agentHandlers.ListAgentTrades)
	})

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	return r
}
