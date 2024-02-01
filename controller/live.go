package controller

import (
	"context"
	"github.com/FACorreiaa/Aviation-tracker/controller/html/flights"
	"net/http"
)

// https://openlayers.org/en/latest/examples/feature-move-animation.html future feature
// https://openlayers.org/en/latest/examples/flight-animation.html future feature on this branch for flights with destination

func (h *Handlers) liveFlightsPage(w http.ResponseWriter, r *http.Request) error {
	f := flights.LiveFlightsPage()
	return h.CreateLayout(w, r, "Live Flights", f).Render(context.Background(), w)
}
