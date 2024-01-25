package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/a-h/templ"

	"github.com/FACorreiaa/go-ollama/controller/html/components"
	"github.com/FACorreiaa/go-ollama/controller/html/pages"
	"github.com/FACorreiaa/go-ollama/controller/models"
)

func (h *Handlers) renderAirportTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	columnNames := []string{"Airport Name", "Country Name", "Phone Number",
		"Timezone", "GMT", "Latitude", "Longitude",
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}
	pageSize := 15

	nextPage := page + 1
	prevPage := page - 1

	airports, err := h.core.airports.GetAirports(context.Background(), page, pageSize)
	if err != nil {
		return nil, err
	}

	airport := models.AirportTable{
		Column:   columnNames,
		Airports: airports,
		PrevPage: prevPage,
		NextPage: nextPage,
		Page:     page,
	}
	airportTable := components.TableDaisyComponent(airport)

	return airportTable, nil
}

func (h *Handlers) airportPage(w http.ResponseWriter, r *http.Request) error {
	airportTable, err := h.renderAirportTable(w, r)
	if err != nil {
		return err
	}
	airport := pages.AirportPage(airportTable)
	return h.CreateLayout(w, r, "Airport Page", airport).Render(context.Background(), w)
}
