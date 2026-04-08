// Package handler provides HTTP handlers for the REST API.
package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ryusei/kyudo-dojo-hub/internal/service"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	svc    *service.Service
	logger *slog.Logger
}

// New creates a new Handler.
func New(svc *service.Service, logger *slog.Logger) *Handler {
	return &Handler{svc: svc, logger: logger}
}

// apiResponse is a success envelope.
type apiResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

// apiError is an error envelope.
type apiError struct {
	Success bool          `json:"success"`
	Error   apiErrorInner `json:"error"`
}

type apiErrorInner struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}

func writeSuccess(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: data})
}

func writeCreated(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusCreated, apiResponse{Success: true, Data: data})
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, apiError{
		Success: false,
		Error:   apiErrorInner{Code: code, Message: message},
	})
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	var valErr *service.ValidationError
	if errors.As(err, &valErr) {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", valErr.Message)
		return
	}

	var notFoundErr *service.NotFoundError
	if errors.As(err, &notFoundErr) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", notFoundErr.Error())
		return
	}

	var conflictErr *service.ConflictError
	if errors.As(err, &conflictErr) {
		writeError(w, http.StatusConflict, "CONFLICT", conflictErr.Message)
		return
	}

	h.logger.Error("internal server error", "error", err)
	writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "内部サーバーエラーが発生しました")
}

func decodeJSON(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// ---------------------------------------------------------------------------
// Users
// ---------------------------------------------------------------------------

// ListUsers handles GET /api/v1/users
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.ListUsers(r.Context())
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeSuccess(w, ensureSlice(users))
}

// GetUser handles GET /api/v1/users/{id}
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := h.svc.GetUser(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeSuccess(w, user)
}

// ListUsersByDojo handles GET /api/v1/dojos/{dojoId}/users
func (h *Handler) ListUsersByDojo(w http.ResponseWriter, r *http.Request) {
	dojoID := chi.URLParam(r, "dojoId")
	users, err := h.svc.ListUsersByDojo(r.Context(), dojoID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeSuccess(w, ensureSlice(users))
}

// ---------------------------------------------------------------------------
// Dojos
// ---------------------------------------------------------------------------

// ListDojos handles GET /api/v1/dojos
func (h *Handler) ListDojos(w http.ResponseWriter, r *http.Request) {
	dojos, err := h.svc.ListDojos(r.Context())
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeSuccess(w, ensureSlice(dojos))
}

// GetDojo handles GET /api/v1/dojos/{id}
func (h *Handler) GetDojo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	dojo, err := h.svc.GetDojo(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeSuccess(w, dojo)
}

// ---------------------------------------------------------------------------
// Practices
// ---------------------------------------------------------------------------

// ListPractices handles GET /api/v1/practices
func (h *Handler) ListPractices(w http.ResponseWriter, r *http.Request) {
	var userID *string
	if uid := r.URL.Query().Get("userId"); uid != "" {
		userID = &uid
	}
	practices, err := h.svc.ListPractices(r.Context(), userID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeSuccess(w, ensureSlice(practices))
}

// GetPractice handles GET /api/v1/practices/{id}
func (h *Handler) GetPractice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	practice, err := h.svc.GetPractice(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeSuccess(w, practice)
}

// CreatePractice handles POST /api/v1/practices
func (h *Handler) CreatePractice(w http.ResponseWriter, r *http.Request) {
	var input service.CreatePracticeInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "リクエストボディの形式が不正です")
		return
	}
	practice, err := h.svc.CreatePractice(r.Context(), input)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeCreated(w, practice)
}

// ---------------------------------------------------------------------------
// Videos
// ---------------------------------------------------------------------------

// ListVideos handles GET /api/v1/videos
func (h *Handler) ListVideos(w http.ResponseWriter, r *http.Request) {
	var userID *string
	if uid := r.URL.Query().Get("userId"); uid != "" {
		userID = &uid
	}
	videos, err := h.svc.ListVideos(r.Context(), userID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeSuccess(w, ensureSlice(videos))
}

// GetVideo handles GET /api/v1/videos/{id}
func (h *Handler) GetVideo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	video, err := h.svc.GetVideo(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeSuccess(w, video)
}

// CreateVideo handles POST /api/v1/videos
func (h *Handler) CreateVideo(w http.ResponseWriter, r *http.Request) {
	var input service.CreateVideoInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "リクエストボディの形式が不正です")
		return
	}
	video, err := h.svc.CreateVideo(r.Context(), input)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeCreated(w, video)
}

// ---------------------------------------------------------------------------
// Analyses
// ---------------------------------------------------------------------------

// ListAnalyses handles GET /api/v1/analyses
func (h *Handler) ListAnalyses(w http.ResponseWriter, r *http.Request) {
	var userID *string
	if uid := r.URL.Query().Get("userId"); uid != "" {
		userID = &uid
	}
	analyses, err := h.svc.ListAnalyses(r.Context(), userID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeSuccess(w, ensureSlice(analyses))
}

// GetAnalysis handles GET /api/v1/analyses/{id}
func (h *Handler) GetAnalysis(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	analysis, err := h.svc.GetAnalysis(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeSuccess(w, analysis)
}

// GetAnalysisByVideo handles GET /api/v1/analyses/by-video/{videoId}
func (h *Handler) GetAnalysisByVideo(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "videoId")
	analysis, err := h.svc.GetAnalysisByVideo(r.Context(), videoID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeSuccess(w, analysis)
}

// ---------------------------------------------------------------------------
// Reservations
// ---------------------------------------------------------------------------

// ListReservations handles GET /api/v1/reservations
func (h *Handler) ListReservations(w http.ResponseWriter, r *http.Request) {
	var dojoID, date *string
	if did := r.URL.Query().Get("dojoId"); did != "" {
		dojoID = &did
	}
	if d := r.URL.Query().Get("date"); d != "" {
		date = &d
	}
	reservations, err := h.svc.ListReservations(r.Context(), dojoID, date)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeSuccess(w, ensureSlice(reservations))
}

// GetReservation handles GET /api/v1/reservations/{id}
func (h *Handler) GetReservation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	reservation, err := h.svc.GetReservation(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeSuccess(w, reservation)
}

// CreateReservation handles POST /api/v1/reservations
func (h *Handler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var input service.CreateReservationInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "リクエストボディの形式が不正です")
		return
	}
	reservation, err := h.svc.CreateReservation(r.Context(), input)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeCreated(w, reservation)
}

// DeleteReservation handles DELETE /api/v1/reservations/{id}
func (h *Handler) DeleteReservation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeleteReservation(r.Context(), id); err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeSuccess(w, map[string]bool{"deleted": true})
}

// ---------------------------------------------------------------------------
// Exam Checklists
// ---------------------------------------------------------------------------

// ListExamChecklists handles GET /api/v1/exam-checklists
func (h *Handler) ListExamChecklists(w http.ResponseWriter, r *http.Request) {
	var userID *string
	if uid := r.URL.Query().Get("userId"); uid != "" {
		userID = &uid
	}
	checklists, err := h.svc.ListExamChecklists(r.Context(), userID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeSuccess(w, ensureSlice(checklists))
}

// GetExamChecklist handles GET /api/v1/exam-checklists/{id}
func (h *Handler) GetExamChecklist(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	checklist, err := h.svc.GetExamChecklist(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeSuccess(w, checklist)
}

// ToggleChecklistItem handles PATCH /api/v1/exam-checklists/{id}/items/{itemId}/toggle
func (h *Handler) ToggleChecklistItem(w http.ResponseWriter, r *http.Request) {
	checklistID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "itemId")
	checklist, err := h.svc.ToggleChecklistItem(r.Context(), checklistID, itemID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeSuccess(w, checklist)
}

// ---------------------------------------------------------------------------
// Dashboard
// ---------------------------------------------------------------------------

// GetDashboardSummary handles GET /api/v1/dojos/{dojoId}/dashboard
func (h *Handler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	dojoID := chi.URLParam(r, "dojoId")
	summary, err := h.svc.GetDashboardSummary(r.Context(), dojoID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeSuccess(w, summary)
}

// ---------------------------------------------------------------------------
// Health
// ---------------------------------------------------------------------------

// HealthCheck handles GET /api/v1/health
func (h *Handler) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	writeSuccess(w, map[string]string{"status": "ok"})
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// ensureSlice returns an empty JSON array instead of null for nil slices.
func ensureSlice[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}
