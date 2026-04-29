package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthZ is the chassis liveness probe extracted from the router so it
// has a stable, named symbol the OpenAPI generator can attach
// annotations to. Anonymous handlers in router setup cannot carry swag
// directives.
//
// @Summary      Liveness probe
// @Description  Returns {"status":"ok"} as soon as the HTTP server is up.
// @Description  This endpoint does NOT depend on Postgres, Redis, or NATS;
// @Description  it is intentionally a chassis-only check so orchestrators
// @Description  can distinguish "binary running" from "fully wired".
// @Tags         health
// @Produce      json
// @Success      200  {object}  openapi.HealthResponse
// @Router       /healthz [get]
func HealthZ(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(http.StatusOK, `{"status":"ok"}`)
}
