package controller

import (
	"context"
	"net/http"

	"github.com/FACorreiaa/Aviation-tracker/controller/models"
	svg2 "github.com/FACorreiaa/Aviation-tracker/controller/svg"
	"github.com/a-h/templ"

	"github.com/FACorreiaa/Aviation-tracker/controller/html/flights"
)

// https://openlayers.org/en/latest/examples/feature-move-animation.html future feature
// https://openlayers.org/en/latest/examples/flight-animation.html
// future feature on this branch for flights with destination

// need to change sql query later

func (h *Handlers) renderLiveLocationsSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home", Icon: svg2.HomeIcon()},
		{Path: "/flights/live", Label: "All Flights", Icon: svg2.HomeIcon()},
		{
			Label: "Live Flights",
			Icon:  svg2.GlobeEuropeIcon(),
			SubItems: []models.SidebarItem{
				{Path: "/flights/active/data", Label: "Live Flights", Icon: svg2.GlobeEuropeIcon()},
				{Path: "/flights/active/map", Label: "Live Flights Locations", Icon: svg2.MapIcon()},
			},
		},
		{
			Label: "Landed Flights",
			Icon:  svg2.GlobeEuropeIcon(),
			SubItems: []models.SidebarItem{
				{Path: "/flights/landed/data", Label: "Landed Flights", Icon: svg2.GlobeEuropeIcon()},
				{Path: "/flights/landed/map", Label: "Landed Flights Location", Icon: svg2.MapIcon()},
			},
		},
		{Path: "/flights/scheduled", Label: "Scheduled Flights", Icon: svg2.HomeIcon()},
		{Path: "/flights/cancelled", Label: "Cancelled Flights", Icon: svg2.HomeIcon()},
		{Path: "/flights/incident", Label: "Incident Flights", Icon: svg2.HomeIcon()},
		{Path: "/flights/diverted", Label: "Diverted Flights", Icon: svg2.HomeIcon()},
		{Path: "/settings", Label: "Settings", Icon: svg2.SettingsIcon()},
		{Path: "/log-out", Label: "Log out", Icon: svg2.LogoutIcon()},
	}
	return sidebar
}

func (h *Handlers) renderFlightsTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	lf := make([]models.LiveFlights, 0)
	var page int
	var sortAux string

	param := r.FormValue("search")
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")

	if sortBy == ASC {
		sortAux = DESC
	} else {
		sortAux = ASC
	}

	columnNames := []models.ColumnItems{
		{Title: "Flight Number", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
	}

	//fullPage, airportList, _ := h.getAirports(w, r)
	//filteredPage, filteredAirport, _ := h.getAirportByName(w, r)

	//if len(param) > 0 {
	//	lf = filteredAirport
	//	page = filteredPage
	//} else {
	//	lf = airportList
	//	page = fullPage
	//}

	nextPage := page + 1
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	lastPage, err := h.getTotalAirports()
	if err != nil {
		return nil, err
	}
	f := models.FlightsTable{
		Column:      columnNames,
		Flights:     lf,
		PrevPage:    prevPage,
		NextPage:    nextPage,
		Page:        page,
		LastPage:    lastPage,
		SearchParam: param,
		OrderParam:  orderBy,
		SortParam:   sortAux,
	}
	flightsTable := flights.FlightsTableComponent(f)

	return flightsTable, nil
}

func (h *Handlers) liveFlightsPage(w http.ResponseWriter, r *http.Request) error {
	table, err := h.renderFlightsTable(w, r)
	if err != nil {
		HandleError(err, "Error fetching flights table")
		return err
	}
	s := h.renderLiveLocationsSidebar()

	f := flights.LiveFlightsPage(table, s, "Live Flights", "check live flights data")
	return h.CreateLayout(w, r, "Live Flights", f).Render(context.Background(), w)
}
