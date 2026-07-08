package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/ryusei/kyudo-dojo-hub/backend/internal/middleware"
)

func TestParseOrigins(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want []string
	}{
		{name: "empty falls back to default", raw: "", want: []string{middleware.DefaultCORSOrigin}},
		{name: "blanks fall back to default", raw: " , , ", want: []string{middleware.DefaultCORSOrigin}},
		{name: "single", raw: "https://a.example.com", want: []string{"https://a.example.com"}},
		{
			name: "multiple with spaces",
			raw:  " https://a.example.com , https://b.example.com ",
			want: []string{"https://a.example.com", "https://b.example.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := middleware.ParseOrigins(tt.raw)
			if !slices.Equal(got, tt.want) {
				t.Errorf("ParseOrigins(%q) = %v, want %v", tt.raw, got, tt.want)
			}
		})
	}
}

func newCORSHandler(origins []string) http.Handler {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return middleware.CORS(origins)(next)
}

func TestCORS_AllowedOriginIsEchoed(t *testing.T) {
	h := newCORSHandler([]string{"http://localhost:5173", "https://app.example.com"})

	for _, origin := range []string{"http://localhost:5173", "https://app.example.com"} {
		req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		req.Header.Set("Origin", origin)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != origin {
			t.Errorf("Allow-Origin = %q, want %q", got, origin)
		}
		if got := rec.Header().Get("Vary"); got != "Origin" {
			t.Errorf("Vary = %q, want Origin", got)
		}
	}
}

func TestCORS_DisallowedOriginNotEchoed(t *testing.T) {
	h := newCORSHandler([]string{"http://localhost:5173"})
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	req.Header.Set("Origin", "https://evil.example.com")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("Allow-Origin = %q, want empty for disallowed origin", got)
	}
}

func TestCORS_PreflightReturnsNoContent(t *testing.T) {
	h := newCORSHandler([]string{"http://localhost:5173"})
	req := httptest.NewRequest(http.MethodOptions, "/api/users", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want 204", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Error("expected Access-Control-Allow-Methods to be set for preflight")
	}
}
