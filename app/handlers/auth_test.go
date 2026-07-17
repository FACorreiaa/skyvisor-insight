package handlers

import "testing"

func TestSafeReturnPath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "empty", path: "", want: ""},
		{name: "absolute path", path: "/settings", want: "/settings"},
		{name: "nested path", path: "/flights/flight/TP1363", want: "/flights/flight/TP1363"},
		{name: "protocol-relative", path: "//evil.example", want: ""},
		{name: "absolute URL", path: "https://evil.example/", want: ""},
		{name: "backslash trick", path: "/\\evil.example", want: ""},
		{name: "relative path", path: "settings", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := safeReturnPath(tt.path); got != tt.want {
				t.Fatalf("safeReturnPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
