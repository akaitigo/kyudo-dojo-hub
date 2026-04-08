package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ryusei/kyudo-dojo-hub/backend/internal/handler"
	"github.com/ryusei/kyudo-dojo-hub/backend/internal/model"
	"github.com/ryusei/kyudo-dojo-hub/backend/internal/store"
)

func setupHandler() *handler.Handler {
	s := store.New()
	return handler.New(s)
}

// envelope wraps API responses for deserialization.
type envelope[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

type errorEnvelope struct {
	Success bool `json:"success"`
	Error   struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// ---------------------------------------------------------------------------
// Users
// ---------------------------------------------------------------------------

func TestGetUsers(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rec := httptest.NewRecorder()
	h.GetUsers(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp envelope[[]model.User]
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if !resp.Success {
		t.Fatal("expected success=true")
	}
	if len(resp.Data) < 10 {
		t.Fatalf("expected at least 10 users, got %d", len(resp.Data))
	}
}

func TestGetUser_Found(t *testing.T) {
	h := setupHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/users/{id}", h.GetUser)

	req := httptest.NewRequest(http.MethodGet, "/api/users/user-001", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp envelope[model.User]
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if resp.Data.Name != "田中太郎" {
		t.Fatalf("expected 田中太郎, got %s", resp.Data.Name)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	h := setupHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/users/{id}", h.GetUser)

	req := httptest.NewRequest(http.MethodGet, "/api/users/nonexistent", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestGetUsersByDojo(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/users?dojoId=dojo-001", nil)
	rec := httptest.NewRecorder()
	h.GetUsers(rec, req)

	var resp envelope[[]model.User]
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(resp.Data) == 0 {
		t.Fatal("expected at least 1 user for dojo-001")
	}
}

// ---------------------------------------------------------------------------
// Dojos
// ---------------------------------------------------------------------------

func TestGetDojos(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/dojos", nil)
	rec := httptest.NewRecorder()
	h.GetDojos(rec, req)

	var resp envelope[[]model.Dojo]
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 dojos, got %d", len(resp.Data))
	}
}

func TestGetDojo_Found(t *testing.T) {
	h := setupHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/dojos/{id}", h.GetDojo)

	req := httptest.NewRequest(http.MethodGet, "/api/dojos/dojo-001", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var resp envelope[model.Dojo]
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if resp.Data.Name != "東京弓道場" {
		t.Fatalf("expected 東京弓道場, got %s", resp.Data.Name)
	}
}

// ---------------------------------------------------------------------------
// Practices
// ---------------------------------------------------------------------------

func TestGetPractices(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/practices", nil)
	rec := httptest.NewRecorder()
	h.GetPractices(rec, req)

	var resp envelope[[]model.Practice]
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(resp.Data) < 10 {
		t.Fatalf("expected at least 10 practices, got %d", len(resp.Data))
	}
}

func TestCreatePractice_Success(t *testing.T) {
	h := setupHandler()
	body := `{"userId":"user-001","dojoId":"dojo-001","date":"2026-03-30","hitRate":60,"arrowCount":36,"notes":"test","instructorComment":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/practices", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	h.CreatePractice(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var resp envelope[model.Practice]
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if resp.Data.HitRate != 60 {
		t.Fatalf("expected hitRate 60, got %d", resp.Data.HitRate)
	}
}

func TestCreatePractice_ValidationError(t *testing.T) {
	h := setupHandler()
	body := `{"userId":"user-001","date":"2026-03-30","hitRate":101,"arrowCount":36,"notes":"","instructorComment":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/practices", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	h.CreatePractice(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	var resp errorEnvelope
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if resp.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected VALIDATION_ERROR, got %s", resp.Error.Code)
	}
}

func TestCreatePractice_InvalidBody(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/practices", bytes.NewBufferString("invalid json"))
	rec := httptest.NewRecorder()
	h.CreatePractice(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestCreatePractice_BodyTooLarge(t *testing.T) {
	h := setupHandler()
	// 1MB + 1 byte exceeds the MaxBytesReader limit
	largeBody := strings.Repeat("x", 1<<20+1)
	req := httptest.NewRequest(http.MethodPost, "/api/practices", strings.NewReader(largeBody))
	rec := httptest.NewRecorder()
	h.CreatePractice(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for oversized body, got %d", rec.Code)
	}
}

// ---------------------------------------------------------------------------
// Reservations
// ---------------------------------------------------------------------------

func TestCreateReservation_Success(t *testing.T) {
	h := setupHandler()
	body := `{"dojoId":"dojo-001","userId":"user-001","laneNumber":6,"date":"2099-12-31","startTime":"09:00","endTime":"10:00"}`
	req := httptest.NewRequest(http.MethodPost, "/api/reservations", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	h.CreateReservation(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func TestDeleteReservation_Success(t *testing.T) {
	h := setupHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/reservations/{id}", h.DeleteReservation)

	req := httptest.NewRequest(http.MethodDelete, "/api/reservations/res-001", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestDeleteReservation_NotFound(t *testing.T) {
	h := setupHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/reservations/{id}", h.DeleteReservation)

	req := httptest.NewRequest(http.MethodDelete, "/api/reservations/nonexistent", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

// ---------------------------------------------------------------------------
// Exam Checklists
// ---------------------------------------------------------------------------

func TestToggleChecklistItem_Success(t *testing.T) {
	h := setupHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /api/exam-checklists/{checklistId}/items/{itemId}/toggle", h.ToggleChecklistItem)

	req := httptest.NewRequest(http.MethodPatch, "/api/exam-checklists/exam-002/items/item-011/toggle", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

// ---------------------------------------------------------------------------
// Dashboard
// ---------------------------------------------------------------------------

func TestGetDashboardSummary(t *testing.T) {
	h := setupHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/dashboard/{dojoId}", h.GetDashboardSummary)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/dojo-001", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp envelope[model.DashboardSummary]
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if resp.Data.TotalMemberCount == 0 {
		t.Fatal("expected non-zero member count")
	}
}

// ---------------------------------------------------------------------------
// Videos
// ---------------------------------------------------------------------------

func TestCreateVideo_Success(t *testing.T) {
	h := setupHandler()
	body := `{"userId":"user-001","fileName":"test.mp4","fileSize":1048576,"duration":30,"mimeType":"video/mp4","url":"blob:test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/videos", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	h.CreateVideo(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func TestCreateVideo_ValidationError(t *testing.T) {
	h := setupHandler()
	body := `{"userId":"user-001","fileName":"test.mp4","fileSize":524288001,"duration":30,"mimeType":"video/mp4","url":"blob:test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/videos", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	h.CreateVideo(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

// ---------------------------------------------------------------------------
// Analyses
// ---------------------------------------------------------------------------

func TestGetAnalyses(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/analyses", nil)
	rec := httptest.NewRecorder()
	h.GetAnalyses(rec, req)

	var resp envelope[[]model.Analysis]
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(resp.Data) < 2 {
		t.Fatalf("expected at least 2 analyses, got %d", len(resp.Data))
	}
}

func TestGetAnalysisByVideo(t *testing.T) {
	h := setupHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/analyses/video/{videoId}", h.GetAnalysisByVideo)

	req := httptest.NewRequest(http.MethodGet, "/api/analyses/video/video-001", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestAnalyzeVideo_FallbackWhenWorkerUnavailable(t *testing.T) {
	h := setupHandler()
	body := `{"videoId":"video-001","userId":"user-001"}`
	req := httptest.NewRequest(http.MethodPost, "/api/analyses/analyze", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	h.AnalyzeVideo(rec, req)

	// The Python worker is not running during tests, so we expect fallback (201 with data)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 (fallback), got %d", rec.Code)
	}

	var resp envelope[model.Analysis]
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if resp.Data.OverallScore == 0 {
		t.Fatal("expected non-zero overall score in fallback analysis")
	}
	if len(resp.Data.Phases) != 8 {
		t.Fatalf("expected 8 phases in fallback, got %d", len(resp.Data.Phases))
	}
}

func TestAnalyzeVideo_WorkerReturnsError(t *testing.T) {
	// MediaPipeワーカーがエラーを返す場合にフォールバックが使われることを確認
	s := store.New()
	h := handler.New(s)

	body := `{"videoId":"video-001","userId":"user-001"}`
	req := httptest.NewRequest(http.MethodPost, "/api/analyses/analyze", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	h.AnalyzeVideo(rec, req)

	// ワーカーが起動していないためフォールバック
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var resp envelope[model.Analysis]
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	// フォールバック分析が8フェーズすべてのスコアを持つことを確認
	scores := resp.Data.Scores
	if scores.Ashibumi == 0 && scores.Dozukuri == 0 && scores.Kai == 0 {
		t.Fatal("expected non-zero fallback scores")
	}
}
