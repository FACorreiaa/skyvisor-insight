package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
	tripsview "github.com/FACorreiaa/Aviation-tracker/app/view/trips"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

// TripsPage lists the account's trips from skyvisor-api.
func (h *Handler) TripsPage(w http.ResponseWriter, r *http.Request) error {
	return h.renderTrips(w, r, "", "", "", "")
}

// TripsCreate creates a trip from the form and redirects back to the list.
func (h *Handler) TripsCreate(w http.ResponseWriter, r *http.Request) error {
	if h.service.API() == nil {
		return h.renderTrips(w, r, "Trips are not available on this server yet.", "", "", "")
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
			return h.renderTrips(w, r, "The start date must use the YYYY-MM-DD format.", "", "", "")
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

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	created, err := h.service.API().CreateTrip(ctx, accessToken, input)
	if err != nil {
		if errors.Is(err, apiclient.ErrUnauthorized) {
			http.Redirect(w, r, "/login?return_to=/trips", http.StatusSeeOther)
			return nil
		}
		slog.ErrorContext(r.Context(), "create trip", "error", err)
		return h.renderTrips(w, r, friendlyAPIError(err, "Unable to create the trip."), "", "", "")
	}
	return h.renderTrips(w, r, tripsview.AutoWatchMessage(created.AutoWatches), "", "", "")
}

// TripsImport extracts an itinerary from pasted booking text via the API.
func (h *Handler) TripsImport(w http.ResponseWriter, r *http.Request) error {
	if h.service.API() == nil {
		return h.renderTrips(w, r, "Trips are not available on this server yet.", "", "", "")
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
		return h.renderTrips(w, r, "Paste your booking confirmation text first.", "", "", "")
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()
	created, err := h.service.API().ImportTrip(ctx, accessToken, text)
	if err != nil {
		if errors.Is(err, apiclient.ErrUnauthorized) {
			http.Redirect(w, r, "/login?return_to=/trips", http.StatusSeeOther)
			return nil
		}
		slog.ErrorContext(r.Context(), "import trip", "error", err)
		return h.renderTrips(w, r, friendlyAPIError(err, "We could not import that itinerary. Check the text or add the trip manually."), "", "", "")
	}
	return h.renderTrips(w, r, tripsview.AutoWatchMessage(created.AutoWatches), "", "", "")
}

// clientMaxPDFBytes mirrors the API's 10MB cap so oversized uploads fail fast.
const clientMaxPDFBytes = 10 << 20

// TripsImportPDF uploads a PDF e-ticket for AI extraction into a new trip.
func (h *Handler) TripsImportPDF(w http.ResponseWriter, r *http.Request) error {
	if h.service.API() == nil {
		return h.renderTrips(w, r, "Trips are not available on this server yet.", "", "", "")
	}
	accessToken, err := h.apiAccessToken(r)
	if err != nil {
		http.Redirect(w, r, "/login?return_to=/trips", http.StatusSeeOther)
		return nil
	}
	r.Body = http.MaxBytesReader(w, r.Body, clientMaxPDFBytes+1<<16)
	if err := r.ParseMultipartForm(clientMaxPDFBytes); err != nil {
		return h.renderTrips(w, r, "The PDF is too large. The limit is 10MB.", "", "", "")
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		return h.renderTrips(w, r, "Choose a PDF e-ticket to upload.", "", "", "")
	}
	defer func() { _ = file.Close() }()

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()
	created, err := h.service.API().ImportTripPDF(ctx, accessToken, header.Filename, file)
	if err != nil {
		if errors.Is(err, apiclient.ErrUnauthorized) {
			http.Redirect(w, r, "/login?return_to=/trips", http.StatusSeeOther)
			return nil
		}
		slog.ErrorContext(r.Context(), "import trip pdf", "error", err)
		return h.renderTrips(w, r, friendlyAPIError(err, "We could not read that PDF. If it is a scanned image, paste the text instead."), "", "", "")
	}
	return h.renderTrips(w, r, tripsview.AutoWatchMessage(created.AutoWatches), "", "", "")
}

// TripsTimeline renders every segment across trips in chronological order.
// An optional ?tz=Area/City query renders times in that location.
func (h *Handler) TripsTimeline(w http.ResponseWriter, r *http.Request) error {
	loc := time.UTC
	tzName := "UTC"
	if tz := strings.TrimSpace(r.URL.Query().Get("tz")); tz != "" {
		if parsed, err := time.LoadLocation(tz); err == nil {
			loc = parsed
			tzName = tz
		}
	}

	var tripList []apiclient.Trip
	if h.service.API() == nil {
		page := tripsview.TimelinePage(tripList, tzName, loc)
		return h.CreateLayout(w, r, "Timeline", page).Render(r.Context(), w)
	}
	accessToken, err := h.apiAccessToken(r)
	if err != nil {
		http.Redirect(w, r, "/login?return_to=/trips/timeline", http.StatusSeeOther)
		return nil
	}
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	tripList, err = h.service.API().ListTrips(ctx, accessToken)
	if err != nil {
		if errors.Is(err, apiclient.ErrUnauthorized) {
			http.Redirect(w, r, "/login?return_to=/trips/timeline", http.StatusSeeOther)
			return nil
		}
		slog.ErrorContext(r.Context(), "list trips for timeline", "error", err)
	}
	page := tripsview.TimelinePage(tripList, tzName, loc)
	return h.CreateLayout(w, r, "Timeline", page).Render(r.Context(), w)
}

// TripsDelete removes one trip and redirects back to the list.
func (h *Handler) TripsDelete(w http.ResponseWriter, r *http.Request) error {
	if h.service.API() == nil {
		return h.renderTrips(w, r, "Trips are not available on this server yet.", "", "", "")
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
		return h.renderTrips(w, r, "Unable to delete the trip. Try again in a moment.", "", "", "")
	}
	http.Redirect(w, r, "/trips", http.StatusSeeOther)
	return nil
}

// TripsAssistant asks the trip-grounded disruption assistant.
func (h *Handler) TripsAssistant(w http.ResponseWriter, r *http.Request) error {
	if h.service.API() == nil {
		return h.renderTrips(w, r, "Assistant is not available on this server yet.", "", "", "")
	}
	accessToken, err := h.apiAccessToken(r)
	if err != nil {
		http.Redirect(w, r, "/login?return_to=/trips", http.StatusSeeOther)
		return nil
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	question := strings.TrimSpace(r.PostFormValue("question"))
	tripID := strings.TrimSpace(r.PostFormValue("trip_id"))
	if question == "" || tripID == "" {
		return h.renderTrips(w, r, "Choose a trip and enter a question.", "", "", "")
	}
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	response, err := h.service.API().AskAssistant(ctx, accessToken, question, tripID)
	if err != nil {
		if errors.Is(err, apiclient.ErrUnauthorized) {
			http.Redirect(w, r, "/login?return_to=/trips", http.StatusSeeOther)
			return nil
		}
		slog.ErrorContext(r.Context(), "trip assistant", "error", err)
		return h.renderTrips(w, r, friendlyAPIError(err, "Assistant is temporarily unavailable."), "", "", "")
	}
	return h.renderTrips(w, r, "", response.Answer, "", "")
}

// TripsWhatIf runs a delay scenario against a trip via skyvisor-api.
func (h *Handler) TripsWhatIf(w http.ResponseWriter, r *http.Request) error {
	if h.service.API() == nil {
		return h.renderTrips(w, r, "What-if is not available on this server yet.", "", "", "")
	}
	accessToken, err := h.apiAccessToken(r)
	if err != nil {
		http.Redirect(w, r, "/login?return_to=/trips", http.StatusSeeOther)
		return nil
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	tripID := strings.TrimSpace(r.PostFormValue("trip_id"))
	if tripID == "" {
		return h.renderTrips(w, r, "Choose a trip for the what-if scenario.", "", "", "")
	}
	segIdx, _ := strconv.Atoi(strings.TrimSpace(r.PostFormValue("segment_index")))
	depDelay, _ := strconv.Atoi(strings.TrimSpace(r.PostFormValue("delay_departure_minutes")))
	arrDelay, _ := strconv.Atoi(strings.TrimSpace(r.PostFormValue("delay_arrival_minutes")))
	req := apiclient.WhatIfRequest{
		SegmentIndex:          segIdx,
		DelayDepartureMinutes: depDelay,
		DelayArrivalMinutes:   arrDelay,
		Question:              strings.TrimSpace(r.PostFormValue("question")),
	}
	ctx, cancel := context.WithTimeout(r.Context(), 35*time.Second)
	defer cancel()
	result, err := h.service.API().TripWhatIf(ctx, accessToken, tripID, req)
	if err != nil {
		if errors.Is(err, apiclient.ErrUnauthorized) {
			http.Redirect(w, r, "/login?return_to=/trips", http.StatusSeeOther)
			return nil
		}
		slog.ErrorContext(r.Context(), "trip what-if", "error", err)
		return h.renderTrips(w, r, friendlyAPIError(err, "What-if evaluation failed."), "", "", "")
	}
	return h.renderTrips(w, r, "", "", result.Summary, result.Assistant)
}

func (h *Handler) renderTrips(w http.ResponseWriter, r *http.Request, message, assistantAnswer, whatIfSummary, whatIfAnswer string) error {
	var tripList []apiclient.Trip
	if h.service.API() == nil {
		if message == "" {
			message = "Trips are not available on this server yet."
		}
	} else if accessToken, err := h.apiAccessToken(r); err != nil {
		http.Redirect(w, r, "/login?return_to=/trips", http.StatusSeeOther)
		return nil
	} else {
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
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
	page := tripsview.TripsPage(tripList, message, assistantAnswer, whatIfSummary, whatIfAnswer, csrf.Token(r))
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
