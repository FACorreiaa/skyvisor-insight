package pricing

import (
	"context"
	"strings"
	"testing"
)

func TestPricingPageRenders(t *testing.T) {
	var sb strings.Builder
	if err := PricingPage().Render(context.Background(), &sb); err != nil {
		t.Fatalf("render: %v", err)
	}
	html := sb.String()
	for _, want := range []string{
		"Watch one flight free",
		"$19", "$180", "$49", "$490", "Custom",
		"Most popular",
		"x-data=\"{ yearly: true }\"",
		"/register?plan=pro",
		"mailto:fernandocorreia316@gmail.com",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q in rendered pricing page", want)
		}
	}
	if strings.Contains(html, "—") || strings.Contains(html, "–") {
		t.Error("em/en dash found in pricing page copy")
	}
}
