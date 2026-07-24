package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLegalRoutesArePublic(t *testing.T) {
	t.Parallel()

	router := Router(
		nil,
		[]byte("test-session-key-with-at-least-32-characters"),
		false,
		nil,
		nil,
		nil,
	)
	tests := []struct {
		path    string
		content string
	}{
		{path: "/terms", content: "Terms of Service"},
		{path: "/privacy", content: "Privacy Policy"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			if response.Code != http.StatusOK {
				t.Fatalf("GET %s status = %d, want %d", tt.path, response.Code, http.StatusOK)
			}
			if !strings.Contains(response.Body.String(), tt.content) {
				t.Errorf("GET %s response missing %q", tt.path, tt.content)
			}
		})
	}
}
