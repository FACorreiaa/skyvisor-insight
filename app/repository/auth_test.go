package repository

import (
	"testing"

	"github.com/FACorreiaa/Aviation-tracker/app/auth"
)

func TestUsernameCandidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		principal auth.Principal
		want      string
	}{
		{
			name:      "from display name",
			principal: auth.Principal{Name: "Ada Lovelace", Email: "ada@example.com"},
			want:      "ada-lovelace",
		},
		{
			name:      "falls back to email local part",
			principal: auth.Principal{Name: "★★★", Email: "ada.l@example.com"},
			want:      "ada-l",
		},
		{
			name:      "falls back to default",
			principal: auth.Principal{Name: "★★★", Email: ""},
			want:      "traveler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := usernameCandidate(tt.principal); got != tt.want {
				t.Fatalf("usernameCandidate() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSlugifyTruncates(t *testing.T) {
	t.Parallel()
	long := slugify("abcdefghij abcdefghij abcdefghij abcdefghij")
	if len(long) > 32 {
		t.Fatalf("slugify() length = %d, want <= 32", len(long))
	}
}
