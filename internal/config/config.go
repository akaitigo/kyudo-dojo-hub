// Package config loads application configuration from environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all application settings.
type Config struct {
	Port        int
	DatabaseURL string
	// CORSOrigins is the list of allowed CORS origins. Set via CORS_ORIGIN as a
	// comma-separated list (e.g. "http://localhost:5173,https://app.example.com").
	CORSOrigins []string
}

// defaultCORSOrigin is used when CORS_ORIGIN is not set.
const defaultCORSOrigin = "http://localhost:5173"

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	port := 8080
	if v := os.Getenv("PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid PORT value %q: %w", v, err)
		}
		port = p
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		CORSOrigins: parseCORSOrigins(os.Getenv("CORS_ORIGIN")),
	}, nil
}

// parseCORSOrigins splits a comma-separated CORS_ORIGIN value into a list of
// trimmed, non-empty origins. It falls back to the default single origin when
// the value is empty or contains only blanks.
func parseCORSOrigins(raw string) []string {
	origins := make([]string, 0, 1)
	for _, part := range strings.Split(raw, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	if len(origins) == 0 {
		return []string{defaultCORSOrigin}
	}
	return origins
}
