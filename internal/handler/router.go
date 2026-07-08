package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/ryusei/kyudo-dojo-hub/internal/config"
)

// NewRouter creates a chi router with all API routes registered.
func NewRouter(h *Handler, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", h.HealthCheck)

		// Users
		r.Route("/users", func(r chi.Router) {
			r.Get("/", h.ListUsers)
			r.Get("/{id}", h.GetUser)
		})

		// Dojos
		r.Route("/dojos", func(r chi.Router) {
			r.Get("/", h.ListDojos)
			r.Get("/{id}", h.GetDojo)
			r.Get("/{dojoId}/users", h.ListUsersByDojo)
			r.Get("/{dojoId}/dashboard", h.GetDashboardSummary)
		})

		// Practices
		r.Route("/practices", func(r chi.Router) {
			r.Get("/", h.ListPractices)
			r.Post("/", h.CreatePractice)
			r.Get("/{id}", h.GetPractice)
		})

		// Videos
		r.Route("/videos", func(r chi.Router) {
			r.Get("/", h.ListVideos)
			r.Post("/", h.CreateVideo)
			r.Get("/{id}", h.GetVideo)
		})

		// Analyses
		r.Route("/analyses", func(r chi.Router) {
			r.Get("/", h.ListAnalyses)
			r.Get("/{id}", h.GetAnalysis)
			r.Get("/by-video/{videoId}", h.GetAnalysisByVideo)
		})

		// Reservations
		r.Route("/reservations", func(r chi.Router) {
			r.Get("/", h.ListReservations)
			r.Post("/", h.CreateReservation)
			r.Get("/{id}", h.GetReservation)
			r.Delete("/{id}", h.DeleteReservation)
		})

		// Exam Checklists
		r.Route("/exam-checklists", func(r chi.Router) {
			r.Get("/", h.ListExamChecklists)
			r.Get("/{id}", h.GetExamChecklist)
			r.Patch("/{id}/items/{itemId}/toggle", h.ToggleChecklistItem)
		})
	})

	return r
}
