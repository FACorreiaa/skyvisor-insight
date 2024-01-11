package controller

import (
	"context"
	"net/http"

	pages "github.com/FACorreiaa/go-ollama/controller/html/pages"
)

func (h *Handlers) liveFlightsPage(w http.ResponseWriter, r *http.Request) error {
	flights := pages.LiveFlightsPage()
	return h.CreateLayout(w, r, "Live Flights", flights).Render(context.Background(), w)
}
