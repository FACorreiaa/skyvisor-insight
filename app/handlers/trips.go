package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
	tripsview "github.com/FACorreiaa/Aviation-tracker/app/view/trips"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

// TripsPage lists the account's trips from skyvisor-api.
func (h *Handler) TripsPage(w http.ResponseWriter, r *http.Request) error {
	return h.renderTrips(w, r, "")
}

// TripsCreate creates a trip from the form and redirects back to the list.
func (h *Handler) TripsCreate(w http.ResponseWriter, r *http.Request) error {
	if h.service.API() == nil {
		return h.renderTrips(w, r, "Trips are not available on this server yet.")
	}
	accessToken, err := h.apiAccessToken(r)
	if err != nil {
		http.Redirect(w, r, "/login?return_to=/trips", http.StatusSeeOther)
		return nil
	}
	if err := r.ParseForm(); err != nil {
		return err
	}

	input := apiclient.CreateTrip{Name: strings.TrimSpace(r.PostFormValue("name"))}
	if starts := strings.TrimSpace(r.PostFormValue("starts_at")); starts != "" {
		startsAt, parseErr := time.Parse("2006-01-02", starts)
		if parseErr != nil {
			return h.renderTrips(w, r, "The start date must use the YYYY-MM-DD format.")
		}
		input.StartsAt = startsAt.UTC()
	}
	for _, line := range strings.FieldsFunc(r.PostFormValue("flights"), func(r rune) bool {
		return r == '\n' || r == '\r' || r == ','
	}) {
		if number := strings.TrimSpace(line); number != "" {
			input.Segments = append(input.Segments, apiclient.TripSegment{FlightNumber: number})
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	if _, err := h.service.API().CreateTrip(ctx, accessToken, input); err != nil {
		if errors.Is(err, apiclient.ErrUnauthorized) {
			http.Redirect(w, r, "/login?return_to=/trips", http.StatusSeeOther)
			return nil
		}
		slog.ErrorContext(r.Context(), "create trip", "error", err)
		return h.renderTrips(w, r, friendlyAPIError(err, "Unable to create the trip."))
	}
	http.Redirect(w, r, "/trips", http.StatusSeeOther)
	return nil
}

// TripsImport extracts an itinerary from pasted booking text via the API.
func (h *Handler) TripsImport(w http.ResponseWriter, r *http.Request) error {
	if h.service.API() == nil {
		return h.renderTrips(w, r, "Trips are not available on this server yet.")
	}
	accessToken, err := h.apiAccessToken(r)
	if err != nil {
		http.Redirect(w, r, "/login?return_to=/trips", http.StatusSeeOther)
		return nil
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	text := strings.TrimSpace(r.PostFormValue("text"))
	if text == "" {
		return h.renderTrips(w, r, "Paste your booking confirmation text first.")
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()
	if _, err := h.service.API().ImportTrip(ctx, accessToken, text); err != nil {
		if errors.Is(err, apiclient.ErrUnauthorized) {
			http.Redirect(w, r, "/login?return_to=/trips", http.StatusSeeOther)
			return nil
		}
		slog.ErrorContext(r.Context(), "import trip", "error", err)
		return h.renderTrips(w, r, friendlyAPIError(err, "We could not import that itinerary. Check the text or add the trip manually."))
	}
	http.Redirect(w, r, "/trips", http.StatusSeeOther)
	return nil
}

// TripsDelete removes one trip and redirects back to the list.
func (h *Handler) TripsDelete(w http.ResponseWriter, r *http.Request) error {
	if h.service.API() == nil {
		return h.renderTrips(w, r, "Trips are not available on this server yet.")
	}
	accessToken, err := h.apiAccessToken(r)
	if err != nil {
		http.Redirect(w, r, "/login?return_to=/trips", http.StatusSeeOther)
		return nil
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	err = h.service.API().DeleteTrip(ctx, accessToken, id)
	if err != nil && !errors.Is(err, apiclient.ErrNotFound) {
		if errors.Is(err, apiclient.ErrUnauthorized) {
			http.Redirect(w, r, "/login?return_to=/trips", http.StatusSeeOther)
			return nil
		}
		slog.ErrorContext(r.Context(), "delete trip", "error", err)
		return h.renderTrips(w, r, "Unable to delete the trip. Try again in a moment.")
	}
	http.Redirect(w, r, "/trips", http.StatusSeeOther)
	return nil
}

func (h *Handler) renderTrips(w http.ResponseWriter, r *http.Request, message string) error {
	var tripList []apiclient.Trip
	if h.service.API() == nil {
		if message == "" {
			message = "Trips are not available on this server yet."
		}
	} else if accessToken, err := h.apiAccessToken(r); err != nil {
		http.Redirect(w, r, "/login?return_to=/trips", http.StatusSeeOther)
		return nil
	} else {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		tripList, err = h.service.API().ListTrips(ctx, accessToken)
		if err != nil {
			if errors.Is(err, apiclient.ErrUnauthorized) {
				http.Redirect(w, r, "/login?return_to=/trips", http.StatusSeeOther)
				return nil
			}
			slog.ErrorContext(r.Context(), "list trips", "error", err)
			if message == "" {
				message = "Trips are temporarily unavailable. Try again shortly."
			}
		}
	}
	page := tripsview.TripsPage(tripList, message, csrf.Token(r))
	return h.CreateLayout(w, r, "Trips", page).Render(r.Context(), w)
}

func (h *Handler) apiAccessToken(r *http.Request) (string, error) {
	s, _ := h.sessions.Get(r, "auth")
	sessionToken, ok := s.Values[sessionKeyToken].(string)
	if !ok || sessionToken == "" {
		return "", errors.New("no active session")
	}
	return h.service.APIAccessToken(r.Context(), sessionToken)
}

// friendlyAPIError surfaces validation feedback (422s carry actionable
// messages) and hides everything else behind a generic sentence.
func friendlyAPIError(err error, fallback string) string {
	message := err.Error()
	if strings.Contains(message, "skyvisor-api 422") {
		if _, detail, found := strings.Cut(message, ": "); found && detail != "" {
			return strings.ToUpper(detail[:1]) + detail[1:] + "."
		}
	}
	return fallback
}
