package config

import (
	"strings"
	"testing"
)

func TestNewServerConfig(t *testing.T) {
	tests := []struct {
		name         string
		sessionKey   string
		cookieSecure string
		wantErr      string
		wantSecure   bool
	}{
		{name: "requires session key", wantErr: "SESSION_KEY"},
		{name: "rejects invalid secure value", sessionKey: strings.Repeat("a", 32), cookieSecure: "sometimes", wantErr: "COOKIE_SECURE"},
		{name: "accepts secure cookies", sessionKey: strings.Repeat("b", 32), cookieSecure: "true", wantSecure: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SESSION_KEY", tt.sessionKey)
			t.Setenv("COOKIE_SECURE", tt.cookieSecure)
			cfg, err := NewServerConfig()
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("NewServerConfig() error = %v, want containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewServerConfig() error = %v", err)
			}
			if cfg.CookieSecure != tt.wantSecure {
				t.Fatalf("CookieSecure = %v, want %v", cfg.CookieSecure, tt.wantSecure)
			}
		})
	}
}

func TestNewDatabaseConfigUsesDatabaseURL(t *testing.T) {
	const databaseURL = "postgres://example.invalid/skyvisor"
	t.Setenv("DATABASE_URL", databaseURL)
	t.Setenv("DB_PASS", "")

	cfg, err := NewDatabaseConfig()
	if err != nil {
		t.Fatalf("NewDatabaseConfig() error = %v", err)
	}
	if cfg.ConnectionURL != databaseURL {
		t.Fatalf("ConnectionURL = %q, want %q", cfg.ConnectionURL, databaseURL)
	}
}
