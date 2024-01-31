package controller

import (
	"context"
	"net/http"

	"github.com/FACorreiaa/Aviation-tracker/controller/html/pages"
)

// https://openlayers.org/en/latest/examples/feature-move-animation.html future feature
// https://openlayers.org/en/latest/examples/flight-animation.html future feature on this branch for flights with destination
func (h *Handlers) liveFlightsPage(w http.ResponseWriter, r *http.Request) error {
	flights := pages.LiveFlightsPage()
	return h.CreateLayout(w, r, "Live Flights", flights).Render(context.Background(), w)
}
