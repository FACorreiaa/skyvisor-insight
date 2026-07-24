package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
	"github.com/FACorreiaa/Aviation-tracker/app/view/flights"
)

// LiveTrackerPage lists provider-backed flights via skyvisor-api (Phase 3).
func (h *Handler) LiveTrackerPage(w http.ResponseWriter, r *http.Request) error {
	accessToken, err := h.apiAccessToken(r)
	if err != nil || accessToken == "" || h.service.API() == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil
	}

	q := apiclient.LiveFlightQuery{
		Airline: r.URL.Query().Get("airline"),
		DepIATA: r.URL.Query().Get("dep_iata"),
		ArrIATA: r.URL.Query().Get("arr_iata"),
		Status:  r.URL.Query().Get("status"),
	}
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if n, convErr := strconv.Atoi(raw); convErr == nil {
			q.Limit = n
		}
	}
	if q.Status == "" && q.Airline == "" && q.DepIATA == "" && q.ArrIATA == "" {
		q.Status = "active"
	}

	ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
	defer cancel()
	items, listErr := h.service.API().ListLiveFlights(ctx, accessToken, q)
	message := ""
	if listErr != nil {
		message = "Live flight feed is temporarily unavailable."
		items = nil
	}

	if strings.EqualFold(r.Header.Get("HX-Request"), "true") {
		return flights.LiveTrackerResults(items, message).Render(r.Context(), w)
	}
	page := flights.LiveTrackerPage(items, q, message)
	return h.CreateLayout(w, r, "Live tracker", page).Render(r.Context(), w)
}
