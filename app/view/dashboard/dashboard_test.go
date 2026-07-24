package dashboard

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
)

func TestPageRendersOperationalEvidence(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()
	dashboard := apiclient.OperationsDashboard{
		Summary: apiclient.OperationsSummary{ActiveWatches: 1, FlightsAtRisk: 1},
		Attention: []apiclient.OperationsAttention{{
			ID: "watch:1", Severity: "high", Title: "FX100 is delayed 42 minutes",
			FlightNumber: "FX100", Route: "MEM → CDG", Reason: "Check the downstream handoff.",
		}},
		Watches: []apiclient.OperationsWatch{{
			WatchID: "1", FlightNumber: "FX100", Route: "MEM → CDG", Status: "active", FreshnessStatus: "fresh", UpdatedAt: now,
		}},
		Freshness:   apiclient.OperationsFreshness{Provider: "aviationstack", Status: "fresh", FreshRecords: 1, StaleAfterSeconds: 120, LatestObservationAt: now},
		GeneratedAt: now,
	}
	var output bytes.Buffer
	if err := Page(dashboard, "").Render(context.Background(), &output); err != nil {
		t.Fatalf("render dashboard: %v", err)
	}
	html := output.String()
	for _, expected := range []string{"What needs attention now", "FX100 is delayed 42 minutes", "MEM → CDG", "aviationstack", "Live data current"} {
		if !strings.Contains(html, expected) {
			t.Fatalf("dashboard HTML missing %q", expected)
		}
	}
}

func TestPageRendersHonestEmptyState(t *testing.T) {
	t.Parallel()
	var output bytes.Buffer
	if err := Page(apiclient.OperationsDashboard{}, "").Render(context.Background(), &output); err != nil {
		t.Fatalf("render dashboard: %v", err)
	}
	html := output.String()
	if !strings.Contains(html, "No operational risks detected") || !strings.Contains(html, "Missing evidence is shown as unknown") {
		t.Fatalf("empty dashboard did not explain missing data: %s", html)
	}
}
