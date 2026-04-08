package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEnsureSlice(t *testing.T) {
	t.Run("nil returns empty slice", func(t *testing.T) {
		var s []string
		result := ensureSlice(s)
		if result == nil {
			t.Fatal("expected non-nil slice")
		}
		if len(result) != 0 {
			t.Fatalf("expected empty slice, got len=%d", len(result))
		}
	})

	t.Run("non-nil returns same slice", func(t *testing.T) {
		s := []string{"a", "b"}
		result := ensureSlice(s)
		if len(result) != 2 {
			t.Fatalf("expected len=2, got len=%d", len(result))
		}
	})
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	writeJSON(w, http.StatusOK, map[string]string{"key": "value"})

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}

	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["key"] != "value" {
		t.Errorf("expected key=value, got key=%s", body["key"])
	}
}

func TestWriteSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	writeSuccess(w, map[string]string{"status": "ok"})

	var body apiResponse
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if !body.Success {
		t.Error("expected success=true")
	}
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "bad input")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var body apiError
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Success {
		t.Error("expected success=false")
	}
	if body.Error.Code != "VALIDATION_ERROR" {
		t.Errorf("expected code VALIDATION_ERROR, got %s", body.Error.Code)
	}
	if body.Error.Message != "bad input" {
		t.Errorf("expected message 'bad input', got %s", body.Error.Message)
	}
}
