package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLegalPages(t *testing.T) {
	t.Parallel()

	h := &Handler{}
	tests := []struct {
		name    string
		path    string
		render  func(*httptest.ResponseRecorder, string) error
		content []string
	}{
		{
			name: "terms",
			path: "/terms",
			render: func(recorder *httptest.ResponseRecorder, path string) error {
				return h.TermsPage(recorder, httptest.NewRequest("GET", path, nil))
			},
			content: []string{"Terms of Service", "Fernando Correia", `href="/privacy"`},
		},
		{
			name: "privacy",
			path: "/privacy",
			render: func(recorder *httptest.ResponseRecorder, path string) error {
				return h.PrivacyPage(recorder, httptest.NewRequest("GET", path, nil))
			},
			content: []string{"Privacy Policy", "Google Gemini", "Stripe", "PostHog", "CNPD"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			recorder := httptest.NewRecorder()
			if err := tt.render(recorder, tt.path); err != nil {
				t.Fatalf("render %s: %v", tt.path, err)
			}
			for _, content := range tt.content {
				if !strings.Contains(recorder.Body.String(), content) {
					t.Errorf("%s response missing %q", tt.path, content)
				}
			}
		})
	}
}
