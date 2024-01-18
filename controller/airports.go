package controller

import (
	"context"
	"net/http"

	pages "github.com/FACorreiaa/go-ollama/controller/html/pages"
)

func (h *Handlers) airportPage(w http.ResponseWriter, r *http.Request) error {
	airport := pages.AirportPage()
	return h.CreateLayout(w, r, "Airport Page", airport).Render(context.Background(), w)
}
