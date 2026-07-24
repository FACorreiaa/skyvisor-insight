package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
	analyticsview "github.com/FACorreiaa/Aviation-tracker/app/view/analytics"
)

// AnalyticsPage renders performance analytics via skyvisor-api (Phase 6).
func (h *Handler) AnalyticsPage(w http.ResponseWriter, r *http.Request) error {
	accessToken, err := h.apiAccessToken(r)
	if err != nil || accessToken == "" || h.service.API() == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil
	}

	q := apiclient.AnalyticsQuery{
		Airline: strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("airline"))),
		Airport: strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("airport"))),
	}
	if raw := r.URL.Query().Get("window_days"); raw != "" {
		if n, convErr := strconv.Atoi(raw); convErr == nil {
			q.WindowDays = n
		}
	}

	message := ""
	var report apiclient.AnalyticsReport
	if q.Airline == "" && q.Airport == "" {
		message = "Choose an airline IATA and/or airport IATA to run analytics."
	} else {
		ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
		defer cancel()
		report, err = h.service.API().Analytics(ctx, accessToken, q)
		if err != nil {
			message = "Analytics are temporarily unavailable."
		}
	}

	if strings.EqualFold(r.Header.Get("HX-Request"), "true") && r.URL.Query().Get("partial") == "report" {
		return analyticsview.AnalyticsResults(report, message).Render(r.Context(), w)
	}
	page := analyticsview.AnalyticsPage(report, q, message)
	return h.CreateLayout(w, r, "Analytics", page).Render(r.Context(), w)
}

// AnalyticsExport proxies the Pro CSV download from skyvisor-api.
func (h *Handler) AnalyticsExport(w http.ResponseWriter, r *http.Request) error {
	accessToken, err := h.apiAccessToken(r)
	if err != nil || accessToken == "" || h.service.API() == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil
	}
	q := apiclient.AnalyticsQuery{
		Airline: strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("airline"))),
		Airport: strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("airport"))),
	}
	if raw := r.URL.Query().Get("window_days"); raw != "" {
		if n, convErr := strconv.Atoi(raw); convErr == nil {
			q.WindowDays = n
		}
	}
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	body, err := h.service.API().AnalyticsExportCSV(ctx, accessToken, q)
	if errors.Is(err, apiclient.ErrPaymentRequired) {
		http.Redirect(w, r, "/analytics?airline="+q.Airline+"&airport="+q.Airport+"&upgrade=1", http.StatusSeeOther)
		return nil
	}
	if err != nil {
		http.Error(w, "Export failed", http.StatusBadGateway)
		return nil
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="skyvisor-analytics.csv"`)
	_, _ = w.Write(body)
	return nil
}
