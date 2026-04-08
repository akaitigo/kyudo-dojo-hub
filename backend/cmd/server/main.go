// Package main starts the kyudo-dojo-hub REST API server.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ryusei/kyudo-dojo-hub/backend/internal/handler"
	"github.com/ryusei/kyudo-dojo-hub/backend/internal/middleware"
	"github.com/ryusei/kyudo-dojo-hub/backend/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	allowOrigin := os.Getenv("CORS_ALLOW_ORIGIN")
	if allowOrigin == "" {
		allowOrigin = "http://localhost:5173"
	}

	s := store.New()
	h := handler.New(s)
	mux := http.NewServeMux()

	// Users
	mux.HandleFunc("GET /api/users", h.GetUsers)
	mux.HandleFunc("GET /api/users/{id}", h.GetUser)

	// Dojos
	mux.HandleFunc("GET /api/dojos", h.GetDojos)
	mux.HandleFunc("GET /api/dojos/{id}", h.GetDojo)

	// Practices
	mux.HandleFunc("GET /api/practices", h.GetPractices)
	mux.HandleFunc("GET /api/practices/{id}", h.GetPractice)
	mux.HandleFunc("POST /api/practices", h.CreatePractice)

	// Videos
	mux.HandleFunc("GET /api/videos", h.GetVideos)
	mux.HandleFunc("GET /api/videos/{id}", h.GetVideo)
	mux.HandleFunc("POST /api/videos", h.CreateVideo)

	// Analyses
	mux.HandleFunc("GET /api/analyses", h.GetAnalyses)
	mux.HandleFunc("GET /api/analyses/{id}", h.GetAnalysis)
	mux.HandleFunc("GET /api/analyses/video/{videoId}", h.GetAnalysisByVideo)
	mux.HandleFunc("POST /api/analyses/analyze", h.AnalyzeVideo)

	// Reservations
	mux.HandleFunc("GET /api/reservations", h.GetReservations)
	mux.HandleFunc("GET /api/reservations/{id}", h.GetReservation)
	mux.HandleFunc("POST /api/reservations", h.CreateReservation)
	mux.HandleFunc("DELETE /api/reservations/{id}", h.DeleteReservation)

	// Exam Checklists
	mux.HandleFunc("GET /api/exam-checklists", h.GetExamChecklists)
	mux.HandleFunc("GET /api/exam-checklists/{id}", h.GetExamChecklist)
	mux.HandleFunc("PATCH /api/exam-checklists/{checklistId}/items/{itemId}/toggle", h.ToggleChecklistItem)

	// Dashboard
	mux.HandleFunc("GET /api/dashboard/{dojoId}", h.GetDashboardSummary)

	// Health check
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
			log.Printf("error writing health check response: %v", err)
		}
	})

	corsHandler := middleware.CORS(allowOrigin)(mux)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           corsHandler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Printf("Received signal %v, shutting down...", sig)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Starting kyudo-dojo-hub API server on :%s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
	log.Println("Server stopped gracefully")
}
