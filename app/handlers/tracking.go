package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
	"github.com/FACorreiaa/Aviation-tracker/app/view/components"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
)

var flightNumberPattern = regexp.MustCompile(`^[A-Z0-9]{2,3}[0-9]{1,4}$`)

func normalizeFlightNumber(value string) (string, bool) {
	normalized := strings.ToUpper(strings.Join(strings.Fields(value), ""))
	return normalized, flightNumberPattern.MatchString(normalized)
}

func (h *Handler) TrackFlight(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query().Get("flight")
	flightNumber, valid := normalizeFlightNumber(query)

	var result *models.LiveFlights
	message := ""
	if query == "" {
		message = "Enter a flight number to begin."
	} else if !valid {
		message = "Use an airline code and flight number, for example TP1363."
	} else {
		ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
		defer cancel()

		flight, err := h.service.GetFlightByID(ctx, flightNumber)
		switch {
		case err == nil:
			result = &flight
		case errors.Is(err, pgx.ErrNoRows):
			message = "We could not find that flight in the current data. Check the number and try again."
		default:
			slog.ErrorContext(r.Context(), "flight lookup failed", "error", err)
			message = "Flight data is temporarily unavailable. Please try again shortly."
		}
	}

	canWatch := r.Context().Value(models.CtxKeyAuthUser) != nil && h.service.API() != nil
	lookup := components.FlightLookupResultAuth(flightNumber, result, message, canWatch, csrf.Token(r))
	if strings.EqualFold(r.Header.Get("HX-Request"), "true") {
		return lookup.Render(r.Context(), w)
	}

	page := components.FlightLookupPage(flightNumber, lookup)
	return h.CreateLayout(w, r, "Track a flight", page).Render(r.Context(), w)
}
