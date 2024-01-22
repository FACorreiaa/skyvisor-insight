package controller

import (
	"context"
	"net/http"

	"github.com/FACorreiaa/go-ollama/controller/html/pages"
)

func (h *Handlers) airportPage(w http.ResponseWriter, r *http.Request) error {
	// change templ layout and add data to the templates
	airports, err := h.core.airports.GetAirports(context.Background())
	if err != nil {
		return err
	}
	for _, a := range airports {
		// pass data to the table and to the map.
		println(a)
	}

	airport := pages.AirportPage()
	return h.CreateLayout(w, r, "Airport Page", airport).Render(context.Background(), w)
}
