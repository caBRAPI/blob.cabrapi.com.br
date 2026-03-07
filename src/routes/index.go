package routes

import (
	"blob/src/controllers"
	"blob/src/middleware"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, limiter *middleware.RateLimiter) {
	// GET /health (public)
	mux.Handle(
		"GET /health",
		limiter.Middleware(http.HandlerFunc(HealthHandler)),
	)

	// GET / (public)
	mux.Handle(
		"GET /",
		limiter.Middleware(http.HandlerFunc(GETHandler)),
	)

	// PUT /blob (private)
	mux.Handle(
		"PUT /blob",
		limiter.Middleware(middleware.AuthMiddleware(http.HandlerFunc(controllers.UploadBlobController))),
	)
}
