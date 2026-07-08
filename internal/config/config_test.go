package config

import (
	"slices"
	"testing"
)

func TestLoad_MissingDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing DATABASE_URL")
	}
}

func TestLoad_InvalidPort(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/test")
	t.Setenv("PORT", "abc")
	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid PORT")
	}
}

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/test")
	t.Setenv("PORT", "")
	t.Setenv("CORS_ORIGIN", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Port)
	}
	if !slices.Equal(cfg.CORSOrigins, []string{"http://localhost:5173"}) {
		t.Errorf("expected default CORS origins, got %v", cfg.CORSOrigins)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://custom/db")
	t.Setenv("PORT", "9090")
	t.Setenv("CORS_ORIGIN", "https://example.com")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Port)
	}
	if cfg.DatabaseURL != "postgres://custom/db" {
		t.Errorf("expected custom DATABASE_URL, got %s", cfg.DatabaseURL)
	}
	if !slices.Equal(cfg.CORSOrigins, []string{"https://example.com"}) {
		t.Errorf("expected single custom CORS origin, got %v", cfg.CORSOrigins)
	}
}

func TestLoad_MultipleCORSOrigins(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/test")
	t.Setenv("CORS_ORIGIN", "http://localhost:5173, https://app.example.com ,https://admin.example.com")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"http://localhost:5173", "https://app.example.com", "https://admin.example.com"}
	if !slices.Equal(cfg.CORSOrigins, want) {
		t.Errorf("expected %v, got %v", want, cfg.CORSOrigins)
	}
}

func TestParseCORSOrigins(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want []string
	}{
		{name: "empty falls back to default", raw: "", want: []string{defaultCORSOrigin}},
		{name: "blanks fall back to default", raw: " ,  , ", want: []string{defaultCORSOrigin}},
		{name: "single", raw: "https://a.example.com", want: []string{"https://a.example.com"}},
		{
			name: "multiple with spaces",
			raw:  " https://a.example.com , https://b.example.com ",
			want: []string{"https://a.example.com", "https://b.example.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCORSOrigins(tt.raw)
			if !slices.Equal(got, tt.want) {
				t.Errorf("parseCORSOrigins(%q) = %v, want %v", tt.raw, got, tt.want)
			}
		})
	}
}
