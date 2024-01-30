package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/a-h/templ"

	"github.com/FACorreiaa/Aviation-tracker/controller/html/components"
	"github.com/FACorreiaa/Aviation-tracker/controller/html/pages"
	"github.com/FACorreiaa/Aviation-tracker/controller/html/svg"
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
)

func (h *Handlers) getAirports(w http.ResponseWriter, r *http.Request) (int, []models.Airport, error) {
	pageSize := 10
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	airports, err := h.core.airports.GetAirports(context.Background(), page, pageSize)
	if err != nil {
		return 0, nil, err
	}

	return page, airports, nil
}

func (h *Handlers) getTotalAirports(w http.ResponseWriter, r *http.Request) (int, error) {
	total, err := h.core.airports.GetSum(context.Background())
	if err != nil {
		return 0, err
	}
	return total, nil
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

	lastPage, err := h.getTotalAirports(w, r)
	if err != nil {
		return nil, err
	}
	airport := models.AirportTable{
		Column:   columnNames,
		Airports: airports,
		PrevPage: prevPage,
		NextPage: nextPage,
		Page:     page,
		LastPage: lastPage,
	}
	airportTable := components.TableDaisyComponent(airport)

	return airportTable, nil
}

func (h *Handlers) renderSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home", Icon: svg.HomeIcon()},
		{Path: "/airports", Label: "Airports", Icon: svg.ArrowRightIcon()},
		{Path: "/airports/locations", Label: "Show Airports", Icon: svg.MapIcon()},
		{Path: "/settings", Label: "Settings", Icon: svg.SettingsIcon()},
		{Path: "/log-out", Label: "Log out", Icon: svg.LogoutIcon()},
	}
	return sidebar
}

func (h *Handlers) airportPage(w http.ResponseWriter, r *http.Request) error {
	airportTable, err := h.renderAirportTable(w, r)
	sidebar := h.renderSidebar()
	if err != nil {
		return err
	}
	airport := pages.AirportPage(airportTable, sidebar)
	return h.CreateLayout(w, r, "Airport Page", airport).Render(context.Background(), w)
}

func (h *Handlers) airportLocationPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderSidebar()

	airport := pages.AirportLocationsPage(sidebar)
	return h.CreateLayout(w, r, "Airport Locations Page", airport).Render(context.Background(), w)
}
