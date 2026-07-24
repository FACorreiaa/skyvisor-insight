package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
	"github.com/FACorreiaa/Aviation-tracker/app/models"
	"github.com/FACorreiaa/Aviation-tracker/app/view/components"
	"github.com/FACorreiaa/Aviation-tracker/app/view/components/flightui"
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
	shareFlash := strings.TrimSpace(r.URL.Query().Get("share"))

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
			flightui.EnrichFlightCoords(result)
		case errors.Is(err, pgx.ErrNoRows):
			if apiFlight, apiErr := h.lookupFlightViaAPI(r, flightNumber); apiErr == nil {
				mapped := flightui.LiveFlightFromAPI(apiFlight)
				result = &mapped
			} else {
				message = "We could not find that flight in the current data. Check the number and try again."
			}
		default:
			slog.ErrorContext(r.Context(), "flight lookup failed", "error", err)
			message = "Flight data is temporarily unavailable. Please try again shortly."
		}
	}

	isAuthenticated := r.Context().Value(models.CtxKeyAuthUser) != nil
	canWatch := isAuthenticated && h.service.API() != nil
	watchID, shareURL := h.findWatchForFlight(r, flightNumber, result)

	lookup := components.FlightLookupResultAuth(flightNumber, result, message, canWatch, csrf.Token(r), isAuthenticated, watchID, shareURL, shareFlash)
	if r.URL.Query().Get("partial") == "flight" {
		target := strings.TrimSpace(r.Header.Get("HX-Target"))
		if target == "flight-result" || target == "#flight-result" {
			return components.FlightLookupResultContainer(lookup).Render(r.Context(), w)
		}
		return lookup.Render(r.Context(), w)
	}
	if strings.EqualFold(r.Header.Get("HX-Request"), "true") {
		return lookup.Render(r.Context(), w)
	}

	page := components.FlightLookupPage(flightNumber, result, lookup, isAuthenticated && result != nil)
	return h.CreateLayout(w, r, "Track a flight", page).Render(r.Context(), w)
}

func (h *Handler) findWatchForFlight(r *http.Request, query string, flight *models.LiveFlights) (watchID, shareURL string) {
	if h.service.API() == nil {
		return "", ""
	}
	accessToken, err := h.apiAccessToken(r)
	if err != nil || accessToken == "" {
		return "", ""
	}
	number := strings.ToUpper(strings.TrimSpace(query))
	if flight != nil && flight.Flight.Number != "" {
		number = flight.Flight.Number
	}
	if number == "" {
		return "", ""
	}
	ctx, cancel := context.WithTimeout(r.Context(), 6*time.Second)
	defer cancel()
	watches, err := h.service.API().ListWatches(ctx, accessToken)
	if err != nil {
		return "", ""
	}
	for _, watch := range watches {
		if strings.EqualFold(watch.FlightNumber, number) {
			if watch.Share != nil {
				return watch.ID, watch.Share.URLPath
			}
			return watch.ID, ""
		}
	}
	return "", ""
}

func (h *Handler) lookupFlightViaAPI(r *http.Request, number string) (apiclient.Flight, error) {
	if h.service.API() == nil {
		return apiclient.Flight{}, errors.New("api unavailable")
	}
	token, err := h.apiAccessToken(r)
	if err != nil || token == "" {
		return apiclient.Flight{}, errors.New("api auth required")
	}
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	return h.service.API().GetFlight(ctx, token, number)
}
