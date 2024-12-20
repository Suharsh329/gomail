package routes

import (
	"gomail/internal/handlers"
	"gomail/internal/middleware"
	"gomail/internal/services"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	mailService := services.NewMailService(false)
	mailHandler := handlers.NewMailHandler(mailService)

	mux.Handle("POST /mail/games", middleware.PerClientRateLimiter(mailHandler.PostMail))

	healthHandler := handlers.NewHealthHandler()
	mux.Handle("GET /health", middleware.PerClientRateLimiter(healthHandler.HealthCheck))
}
