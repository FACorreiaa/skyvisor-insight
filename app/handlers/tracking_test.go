package handlers

import "testing"

func TestNormalizeFlightNumber(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
		valid bool
	}{
		{name: "normalizes whitespace and case", input: " tp 1363 ", want: "TP1363", valid: true},
		{name: "accepts three character airline code", input: "TAP1363", want: "TAP1363", valid: true},
		{name: "rejects missing number", input: "TP", want: "TP", valid: false},
		{name: "rejects punctuation", input: "TP-1363", want: "TP-1363", valid: false},
		{name: "rejects empty value", input: "  ", want: "", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, valid := normalizeFlightNumber(tt.input)
			if got != tt.want || valid != tt.valid {
				t.Fatalf("normalizeFlightNumber(%q) = (%q, %v), want (%q, %v)", tt.input, got, valid, tt.want, tt.valid)
			}
		})
	}
}
