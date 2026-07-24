package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
	airportview "github.com/FACorreiaa/Aviation-tracker/app/view/airports"
	"github.com/go-chi/chi/v5"
)

// AirportBoardPage shows live arrivals/departures via skyvisor-api (Phase 5).
func (h *Handler) AirportBoardPage(w http.ResponseWriter, r *http.Request) error {
	accessToken, err := h.apiAccessToken(r)
	if err != nil || accessToken == "" || h.service.API() == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil
	}

	iata := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "iata")))
	if iata == "" {
		iata = strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("iata")))
	}
	q := apiclient.AirportBoardQuery{
		Direction: r.URL.Query().Get("direction"),
		Status:    r.URL.Query().Get("status"),
	}
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if n, convErr := strconv.Atoi(raw); convErr == nil {
			q.Limit = n
		}
	}
	if q.Direction == "" {
		q.Direction = "both"
	}

	message := ""
	var board apiclient.AirportBoard
	if len(iata) != 3 {
		message = "Enter a 3-letter IATA code such as LIS."
		board = apiclient.AirportBoard{IATA: iata, Arrivals: nil, Departures: nil}
	} else {
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()
		board, err = h.service.API().AirportBoard(ctx, accessToken, iata, q)
		if err != nil {
			message = "Airport board is temporarily unavailable."
			board = apiclient.AirportBoard{IATA: iata, Arrivals: nil, Departures: nil}
		}
	}

	if strings.EqualFold(r.Header.Get("HX-Request"), "true") && r.URL.Query().Get("partial") == "board" {
		return airportview.AirportBoardResults(board, message).Render(r.Context(), w)
	}
	page := airportview.AirportBoardPage(board, q, message)
	return h.CreateLayout(w, r, "Airport board", page).Render(r.Context(), w)
}
