package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
	operationsview "github.com/FACorreiaa/Aviation-tracker/app/view/operations"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func (h *Handler) OperationalCasesPage(w http.ResponseWriter, r *http.Request) error {
	return h.renderOperationalCases(w, r, "")
}

func (h *Handler) OperationalCasesCreate(w http.ResponseWriter, r *http.Request) error {
	token, ok := h.operationsToken(w, r, "/operations/cases")
	if !ok {
		return nil
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	input := apiclient.CreateOperationalCase{
		Kind:             strings.TrimSpace(r.PostFormValue("kind")),
		Reference:        strings.TrimSpace(r.PostFormValue("reference")),
		Title:            strings.TrimSpace(r.PostFormValue("title")),
		Description:      strings.TrimSpace(r.PostFormValue("description")),
		FlightNumbers:    splitFlightNumbers(r.PostFormValue("flight_numbers")),
		Severity:         strings.TrimSpace(r.PostFormValue("severity")),
		Owner:            strings.TrimSpace(r.PostFormValue("owner")),
		EscalationPolicy: strings.TrimSpace(r.PostFormValue("escalation_policy")),
		SLACurrency:      strings.TrimSpace(r.PostFormValue("sla_currency")),
		SLAValueMinor:    parseInt64(r.PostFormValue("sla_value_minor")),
		FailureCostMinor: parseInt64(r.PostFormValue("failure_cost_minor")),
	}
	ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
	defer cancel()
	item, err := h.service.API().CreateOperationalCase(ctx, token, input)
	if err != nil {
		return h.renderOperationalCases(w, r, friendlyAPIError(err, "Unable to create operational case."))
	}
	http.Redirect(w, r, "/operations/cases/"+item.ID, http.StatusSeeOther)
	return nil
}

func (h *Handler) OperationalCasePage(w http.ResponseWriter, r *http.Request) error {
	return h.renderOperationalCase(w, r, chi.URLParam(r, "id"), "")
}

func (h *Handler) OperationalDecisionCreate(w http.ResponseWriter, r *http.Request) error {
	caseID := chi.URLParam(r, "id")
	token, ok := h.operationsToken(w, r, "/operations/cases/"+caseID)
	if !ok {
		return nil
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	probability := parseOptionalFloat(r.PostFormValue("predicted_probability"))
	confidence := parseOptionalFloat(r.PostFormValue("confidence"))
	input := apiclient.CreateDecisionRecord{
		SignalType: r.PostFormValue("signal_type"), SignalSource: r.PostFormValue("signal_source"),
		PredictionType: r.PostFormValue("prediction_type"), PredictedProbability: probability, Confidence: confidence,
		HorizonMinutes: int(parseInt64(r.PostFormValue("horizon_minutes"))), ModelVersion: r.PostFormValue("model_version"),
		ScopeAirport: r.PostFormValue("scope_airport"), ScopeAirline: r.PostFormValue("scope_airline"), ScopeRoute: r.PostFormValue("scope_route"),
		RecommendationType: r.PostFormValue("recommendation_type"), Recommendation: r.PostFormValue("recommendation"),
		Rationale: r.PostFormValue("rationale"), ProposedAction: r.PostFormValue("proposed_action"), ApprovalRequired: r.PostFormValue("approval_required") == "on",
	}
	ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
	defer cancel()
	_, err := h.service.API().CreateDecisionRecord(ctx, token, caseID, input)
	if err != nil {
		return h.renderOperationalCase(w, r, caseID, friendlyAPIError(err, "Unable to record decision."))
	}
	http.Redirect(w, r, "/operations/cases/"+caseID, http.StatusSeeOther)
	return nil
}

func (h *Handler) OperationalDecisionAction(w http.ResponseWriter, r *http.Request) error {
	caseID, decisionID := chi.URLParam(r, "id"), chi.URLParam(r, "decisionID")
	token, ok := h.operationsToken(w, r, "/operations/cases/"+caseID)
	if !ok {
		return nil
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
	defer cancel()
	_, err := h.service.API().RecordDecisionAction(ctx, token, caseID, decisionID, apiclient.RecordDecisionAction{
		Decision: r.PostFormValue("decision"), ActionTaken: r.PostFormValue("action_taken"),
		Notes: r.PostFormValue("notes"), ExternalRef: r.PostFormValue("external_ref"),
	})
	if err != nil {
		return h.renderOperationalCase(w, r, caseID, friendlyAPIError(err, "Unable to record action."))
	}
	http.Redirect(w, r, "/operations/cases/"+caseID, http.StatusSeeOther)
	return nil
}

func (h *Handler) OperationalDecisionOutcome(w http.ResponseWriter, r *http.Request) error {
	caseID, decisionID := chi.URLParam(r, "id"), chi.URLParam(r, "decisionID")
	token, ok := h.operationsToken(w, r, "/operations/cases/"+caseID)
	if !ok {
		return nil
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
	defer cancel()
	_, err := h.service.API().RecordDecisionOutcome(ctx, token, caseID, decisionID, apiclient.RecordDecisionOutcome{
		PredictionResult: r.PostFormValue("prediction_result"), ActionResult: r.PostFormValue("action_result"),
		AvoidedCostMinor: parseInt64(r.PostFormValue("avoided_cost_minor")), Notes: r.PostFormValue("notes"),
	})
	if err != nil {
		return h.renderOperationalCase(w, r, caseID, friendlyAPIError(err, "Unable to record outcome."))
	}
	http.Redirect(w, r, "/operations/cases/"+caseID, http.StatusSeeOther)
	return nil
}

func (h *Handler) renderOperationalCases(w http.ResponseWriter, r *http.Request, message string) error {
	token, ok := h.operationsToken(w, r, "/operations/cases")
	if !ok {
		return nil
	}
	ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
	defer cancel()
	items, err := h.service.API().ListOperationalCases(ctx, token, r.URL.Query().Get("status"))
	var trust apiclient.DecisionTrustMetrics
	if errors.Is(err, apiclient.ErrPaymentRequired) {
		message = "Operational cases require a Business plan."
	} else if err != nil && message == "" {
		message = "Operational cases are temporarily unavailable."
	}
	if err == nil {
		trust, _ = h.service.API().DecisionTrustMetrics(ctx, token)
	}
	page := operationsview.CasesPage(items, trust, r.URL.Query().Get("status"), message, csrf.Token(r))
	return h.CreateLayout(w, r, "Operational cases", page).Render(r.Context(), w)
}

func (h *Handler) renderOperationalCase(w http.ResponseWriter, r *http.Request, caseID, message string) error {
	token, ok := h.operationsToken(w, r, "/operations/cases/"+caseID)
	if !ok {
		return nil
	}
	ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
	defer cancel()
	detail, err := h.service.API().GetOperationalCase(ctx, token, caseID)
	if err != nil && message == "" {
		message = friendlyAPIError(err, "Unable to load operational case.")
	}
	page := operationsview.CaseDetailPage(detail, message, csrf.Token(r))
	return h.CreateLayout(w, r, "Operational case", page).Render(r.Context(), w)
}

func (h *Handler) operationsToken(w http.ResponseWriter, r *http.Request, returnTo string) (string, bool) {
	if h.service.API() == nil {
		http.Error(w, "Skyvisor API is not configured", http.StatusServiceUnavailable)
		return "", false
	}
	token, err := h.apiAccessToken(r)
	if err != nil || token == "" {
		http.Redirect(w, r, "/login?return_to="+returnTo, http.StatusSeeOther)
		return "", false
	}
	return token, true
}

func splitFlightNumbers(value string) []string {
	return strings.FieldsFunc(value, func(r rune) bool { return r == ',' || r == ' ' || r == '\n' })
}

func parseInt64(value string) int64 {
	parsed, _ := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	return parsed
}

func parseOptionalFloat(value string) *float64 {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil
	}
	return &parsed
}
