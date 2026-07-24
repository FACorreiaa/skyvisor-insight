package mcpdocs

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
)

func TestPlaygroundPageRendersUsageAndGovernance(t *testing.T) {
	t.Parallel()

	usage := apiclient.UsageSnapshot{
		MCPReads:            17,
		MCPDailyReadLimit:   100,
		MCPActions:          0,
		MCPDailyActionLimit: 0,
		AssistantCalls:      3,
		AssistantDailyLimit: 10,
		Entitlements: apiclient.Entitlements{
			Plan: "free",
		},
	}

	var output bytes.Buffer
	if err := PlaygroundPage(usage, true, "").Render(context.Background(), &output); err != nil {
		t.Fatalf("render MCP playground: %v", err)
	}

	html := output.String()
	for _, expected := range []string{"API connected", "FREE", "17 / 100", "Read only", "get_operations_dashboard"} {
		if !strings.Contains(html, expected) {
			t.Fatalf("MCP playground HTML missing %q", expected)
		}
	}
}

func TestPlaygroundPageRendersUnavailableUsageState(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	if err := PlaygroundPage(apiclient.UsageSnapshot{}, false, "Usage could not be loaded.").Render(context.Background(), &output); err != nil {
		t.Fatalf("render MCP playground: %v", err)
	}

	html := output.String()
	for _, expected := range []string{"Usage unavailable", "Usage could not be loaded."} {
		if !strings.Contains(html, expected) {
			t.Fatalf("unavailable MCP playground HTML missing %q", expected)
		}
	}
	if strings.Contains(html, "MCP reads") {
		t.Fatal("unavailable MCP playground rendered quota metrics")
	}
}
