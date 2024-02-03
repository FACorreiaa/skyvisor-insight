package controller

import (
	"context"
	"math"
	"net/http"
	"strconv"

	airport "github.com/FACorreiaa/Aviation-tracker/controller/html/airports"
	svg2 "github.com/FACorreiaa/Aviation-tracker/controller/svg"

	"github.com/a-h/templ"

	"github.com/FACorreiaa/Aviation-tracker/controller/models"
)

func (h *Handlers) getAirports(w http.ResponseWriter, r *http.Request) (int, []models.Airport, error) {
	pageSize := 10
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	a, err := h.core.airports.GetAirports(context.Background(), page, pageSize)
	if err != nil {
		return 0, nil, err
	}

	return page, a, nil
}

func (h *Handlers) getAirportsLocation() ([]models.Airport, error) {

	a, err := h.core.airports.GetAirportsLocation(context.Background())
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (h *Handlers) getTotalAirports() (int, error) {
	total, err := h.core.airports.GetSum(context.Background())
	pageSize := 10
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (h *Handlers) renderAirportTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	columnNames := []string{"Airport Name", "Country Name", "Phone Number",
		"Timezone", "GMT", "Latitude", "Longitude",
	}

	page, airports, _ := h.getAirports(w, r)

	nextPage := page + 1
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	lastPage, err := h.getTotalAirports()
	if err != nil {
		return nil, err
	}
	a := models.AirportTable{
		Column:   columnNames,
		Airports: airports,
		PrevPage: prevPage,
		NextPage: nextPage,
		Page:     page,
		LastPage: lastPage,
	}
	airportTable := airport.AirportTableComponent(a)

	return airportTable, nil
}

func (h *Handlers) renderSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home", Icon: svg2.HomeIcon()},
		{Path: "/airports", Label: "Airports", Icon: svg2.ArrowRightIcon()},
		{Path: "/airports/locations", Label: "Show Airports", Icon: svg2.MapIcon()},
		{Path: "/settings", Label: "Settings", Icon: svg2.SettingsIcon()},
		{Path: "/log-out", Label: "Log out", Icon: svg2.LogoutIcon()},
	}
	return sidebar
}

func (h *Handlers) airportPage(w http.ResponseWriter, r *http.Request) error {
	airportTable, err := h.renderAirportTable(w, r)
	sidebar := h.renderSidebar()
	if err != nil {
		return err
	}
	a := airport.AirportPage(airportTable, sidebar)
	return h.CreateLayout(w, r, "Airport Page", a).Render(context.Background(), w)
}

func (h *Handlers) airportLocationPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderSidebar()
	airportsLocations, err := h.getAirportsLocation()
	if err != nil {
		return err
	}
	a := airport.AirportLocationsPage(sidebar, airportsLocations)
	return h.CreateLayout(w, r, "Airport Locations Page", a).Render(context.Background(), w)
}
