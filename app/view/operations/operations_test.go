package operations

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
)

func TestCasesPageRendersOperationalEvidence(t *testing.T) {
	t.Parallel()
	items := []apiclient.OperationalCase{{ID: "case-1", Kind: "shipment", Reference: "AWB-1", Title: "Cold chain", FlightNumbers: []string{"TP1363"}, Status: "action_required", Severity: "high", Owner: "Duty manager", SLACurrency: "EUR", FailureCostMinor: 250000}}
	trust := apiclient.DecisionTrustMetrics{Precision: .8, FalsePositiveRate: .2, MedianLeadTimeMinutes: 47, ExecutedActions: 5, ActionSuccessRate: .6}
	var output bytes.Buffer
	if err := CasesPage(items, trust, "", "", "csrf").Render(context.Background(), &output); err != nil {
		t.Fatalf("render cases: %v", err)
	}
	for _, expected := range []string{"Cases, decisions, outcomes", "AWB-1", "TP1363", "EUR 2500.00", "80.0%", "47 min"} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("cases HTML missing %q", expected)
		}
	}
}

func TestCaseDetailRendersDecisionAndAudit(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()
	probability := .82
	detail := apiclient.OperationalCaseDetail{
		Case:      apiclient.OperationalCase{ID: "case-1", Reference: "AWB-1", Title: "Cold chain", Status: "open", Severity: "high", SLACurrency: "EUR"},
		Decisions: []apiclient.DecisionRecord{{ID: "decision-1", CaseID: "case-1", SignalType: "late_inbound", SignalSource: "aviationstack", PredictionType: "handoff_failure", PredictedProbability: &probability, Recommendation: "Use approved alternative", State: "proposed", ApprovalRequired: true}},
		Audit:     []apiclient.CaseAuditEvent{{ID: "audit-1", CaseID: "case-1", Type: "decision.proposed", ActorType: "agent", CreatedAt: now}},
	}
	var output bytes.Buffer
	if err := CaseDetailPage(detail, "", "csrf").Render(context.Background(), &output); err != nil {
		t.Fatalf("render case detail: %v", err)
	}
	for _, expected := range []string{"Prediction → decision → outcome", "82.0%", "Use approved alternative", "Approve", "decision.proposed", "Agent"} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("case detail HTML missing %q", expected)
		}
	}
}
