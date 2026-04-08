// Package handler provides HTTP request handlers for the kyudo-dojo-hub REST API.
package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/ryusei/kyudo-dojo-hub/backend/internal/model"
	"github.com/ryusei/kyudo-dojo-hub/backend/internal/store"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	store *store.Store
}

// New creates a new Handler.
func New(s *store.Store) *Handler {
	return &Handler{store: s}
}

// ---------------------------------------------------------------------------
// Response helpers
// ---------------------------------------------------------------------------

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("error encoding JSON response: %v", err)
	}
}

func writeSuccess(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusOK, model.APIResponse[interface{}]{
		Success: true,
		Data:    data,
	})
}

func writeCreated(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusCreated, model.APIResponse[interface{}]{
		Success: true,
		Data:    data,
	})
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	resp := model.APIError{
		Success: false,
	}
	resp.Error.Code = code
	resp.Error.Message = message
	writeJSON(w, status, resp)
}

func writeNotFound(w http.ResponseWriter, resource string) {
	writeError(w, http.StatusNotFound, "NOT_FOUND", resource+" が見つかりません")
}

func writeValidationError(w http.ResponseWriter, message string) {
	writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", message)
}

func decodeBody(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := r.Body.Close(); cerr != nil {
			log.Printf("error closing request body: %v", cerr)
		}
	}()
	return json.Unmarshal(body, v)
}

// ---------------------------------------------------------------------------
// Users
// ---------------------------------------------------------------------------

// GetUsers handles GET /api/users and GET /api/users?dojoId=xxx
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	dojoID := r.URL.Query().Get("dojoId")
	if dojoID != "" {
		users := h.store.GetUsersByDojo(dojoID)
		if users == nil {
			users = []model.User{}
		}
		writeSuccess(w, users)
		return
	}
	users := h.store.GetUsers()
	writeSuccess(w, users)
}

// GetUser handles GET /api/users/{id}
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	user, ok := h.store.GetUser(id)
	if !ok {
		writeNotFound(w, "ユーザー")
		return
	}
	writeSuccess(w, user)
}

// ---------------------------------------------------------------------------
// Dojos
// ---------------------------------------------------------------------------

// GetDojos handles GET /api/dojos
func (h *Handler) GetDojos(w http.ResponseWriter, r *http.Request) {
	dojos := h.store.GetDojos()
	writeSuccess(w, dojos)
}

// GetDojo handles GET /api/dojos/{id}
func (h *Handler) GetDojo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	dojo, ok := h.store.GetDojo(id)
	if !ok {
		writeNotFound(w, "道場")
		return
	}
	writeSuccess(w, dojo)
}

// ---------------------------------------------------------------------------
// Practices
// ---------------------------------------------------------------------------

// GetPractices handles GET /api/practices and GET /api/practices?userId=xxx
func (h *Handler) GetPractices(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	practices := h.store.GetPractices(userID)
	if practices == nil {
		practices = []model.Practice{}
	}
	writeSuccess(w, practices)
}

// GetPractice handles GET /api/practices/{id}
func (h *Handler) GetPractice(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	practice, ok := h.store.GetPractice(id)
	if !ok {
		writeNotFound(w, "稽古記録")
		return
	}
	writeSuccess(w, practice)
}

// CreatePractice handles POST /api/practices
func (h *Handler) CreatePractice(w http.ResponseWriter, r *http.Request) {
	var input model.CreatePracticeInput
	if err := decodeBody(r, &input); err != nil {
		writeValidationError(w, "リクエストボディが不正です")
		return
	}

	practice, err := h.store.CreatePractice(input)
	if err != nil {
		writeValidationError(w, err.Error())
		return
	}
	writeCreated(w, practice)
}

// ---------------------------------------------------------------------------
// Videos
// ---------------------------------------------------------------------------

// GetVideos handles GET /api/videos and GET /api/videos?userId=xxx
func (h *Handler) GetVideos(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	videos := h.store.GetVideos(userID)
	if videos == nil {
		videos = []model.Video{}
	}
	writeSuccess(w, videos)
}

// GetVideo handles GET /api/videos/{id}
func (h *Handler) GetVideo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	video, ok := h.store.GetVideo(id)
	if !ok {
		writeNotFound(w, "動画")
		return
	}
	writeSuccess(w, video)
}

// CreateVideo handles POST /api/videos
func (h *Handler) CreateVideo(w http.ResponseWriter, r *http.Request) {
	var input model.CreateVideoInput
	if err := decodeBody(r, &input); err != nil {
		writeValidationError(w, "リクエストボディが不正です")
		return
	}

	video, err := h.store.CreateVideo(input)
	if err != nil {
		writeValidationError(w, err.Error())
		return
	}
	writeCreated(w, video)
}

// ---------------------------------------------------------------------------
// Analyses
// ---------------------------------------------------------------------------

// GetAnalyses handles GET /api/analyses and GET /api/analyses?userId=xxx
func (h *Handler) GetAnalyses(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	analyses := h.store.GetAnalyses(userID)
	if analyses == nil {
		analyses = []model.Analysis{}
	}
	writeSuccess(w, analyses)
}

// GetAnalysis handles GET /api/analyses/{id}
func (h *Handler) GetAnalysis(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	analysis, ok := h.store.GetAnalysis(id)
	if !ok {
		writeNotFound(w, "分析結果")
		return
	}
	writeSuccess(w, analysis)
}

// GetAnalysisByVideo handles GET /api/analyses/video/{videoId}
func (h *Handler) GetAnalysisByVideo(w http.ResponseWriter, r *http.Request) {
	videoID := r.PathValue("videoId")
	analysis, ok := h.store.GetAnalysisByVideo(videoID)
	if !ok {
		writeNotFound(w, "分析結果")
		return
	}
	writeSuccess(w, analysis)
}

// AnalyzeVideo handles POST /api/analyses/analyze
// This endpoint calls the Python MediaPipe worker to analyze a video.
func (h *Handler) AnalyzeVideo(w http.ResponseWriter, r *http.Request) {
	var req model.AnalyzeVideoRequest
	if err := decodeBody(r, &req); err != nil {
		writeValidationError(w, "リクエストボディが不正です")
		return
	}

	if req.VideoID == "" {
		writeValidationError(w, "動画IDは必須です")
		return
	}
	if req.UserID == "" {
		writeValidationError(w, "ユーザーIDは必須です")
		return
	}

	video, ok := h.store.GetVideo(req.VideoID)
	if !ok {
		writeNotFound(w, "動画")
		return
	}

	// Call Python MediaPipe worker
	analysis, err := callMediaPipeWorker(r.Context(), video, req.UserID)
	if err != nil {
		log.Printf("MediaPipe worker unavailable, using fallback: %v", err)

		// Fallback: generate simulated analysis when the worker is unavailable.
		fallback := generateFallbackAnalysis(video, req.UserID)
		h.store.AddAnalysis(fallback)
		writeCreated(w, fallback)
		return
	}

	h.store.AddAnalysis(analysis)
	writeCreated(w, analysis)
}

// ---------------------------------------------------------------------------
// Reservations
// ---------------------------------------------------------------------------

// GetReservations handles GET /api/reservations?dojoId=xxx&date=xxx
func (h *Handler) GetReservations(w http.ResponseWriter, r *http.Request) {
	dojoID := r.URL.Query().Get("dojoId")
	date := r.URL.Query().Get("date")
	reservations := h.store.GetReservations(dojoID, date)
	if reservations == nil {
		reservations = []model.Reservation{}
	}
	writeSuccess(w, reservations)
}

// GetReservation handles GET /api/reservations/{id}
func (h *Handler) GetReservation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	reservation, ok := h.store.GetReservation(id)
	if !ok {
		writeNotFound(w, "予約")
		return
	}
	writeSuccess(w, reservation)
}

// CreateReservation handles POST /api/reservations
func (h *Handler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var input model.CreateReservationInput
	if err := decodeBody(r, &input); err != nil {
		writeValidationError(w, "リクエストボディが不正です")
		return
	}

	reservation, err := h.store.CreateReservation(input)
	if err != nil {
		writeValidationError(w, err.Error())
		return
	}
	writeCreated(w, reservation)
}

// DeleteReservation handles DELETE /api/reservations/{id}
func (h *Handler) DeleteReservation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.store.DeleteReservation(id); err != nil {
		writeNotFound(w, "予約")
		return
	}
	writeSuccess(w, map[string]bool{"deleted": true})
}

// ---------------------------------------------------------------------------
// Exam Checklists
// ---------------------------------------------------------------------------

// GetExamChecklists handles GET /api/exam-checklists and GET /api/exam-checklists?userId=xxx
func (h *Handler) GetExamChecklists(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	checklists := h.store.GetExamChecklists(userID)
	if checklists == nil {
		checklists = []model.ExamChecklist{}
	}
	writeSuccess(w, checklists)
}

// GetExamChecklist handles GET /api/exam-checklists/{id}
func (h *Handler) GetExamChecklist(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	checklist, ok := h.store.GetExamChecklist(id)
	if !ok {
		writeNotFound(w, "審査チェックリスト")
		return
	}
	writeSuccess(w, checklist)
}

// ToggleChecklistItem handles PATCH /api/exam-checklists/{checklistId}/items/{itemId}/toggle
func (h *Handler) ToggleChecklistItem(w http.ResponseWriter, r *http.Request) {
	checklistID := r.PathValue("checklistId")
	itemID := r.PathValue("itemId")

	checklist, err := h.store.ToggleChecklistItem(checklistID, itemID)
	if err != nil {
		writeNotFound(w, "審査チェックリスト")
		return
	}
	writeSuccess(w, checklist)
}

// ---------------------------------------------------------------------------
// Dashboard
// ---------------------------------------------------------------------------

// GetDashboardSummary handles GET /api/dashboard/{dojoId}
func (h *Handler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	dojoID := r.PathValue("dojoId")
	summary := h.store.GetDashboardSummary(dojoID)
	writeSuccess(w, summary)
}
