package controller

import (
	"context"
	"net/http"

	"github.com/FACorreiaa/go-ollama/controller/html/pages"
)

func (h *Handlers) airportPage(w http.ResponseWriter, r *http.Request) error {
	//change templ layout and add data to the templates

	//just hardcode the table values for now and improve solution later
	columnNames := []string{"Airport Name", "Country Name", "Phone Number",
		"Timezone", "GMT", "Latitude", "Longitude",
	}

	page := 1
	pageSize := 10
	nextPage := page + 1
	prevPage := page - 1
	airports, err := h.core.airports.GetAirports(context.Background(), page, pageSize)
	if err != nil {
		return err
	}

	airport := pages.AirportPage(columnNames, airports, prevPage, nextPage, page)
	return h.CreateLayout(w, r, "Airport Page", airport).Render(context.Background(), w)
}
