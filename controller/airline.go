package controller

import (
	"context"
	"math"
	"net/http"
	"strconv"

	"github.com/a-h/templ"

	airline "github.com/FACorreiaa/Aviation-tracker/controller/html/airlines"
	svg2 "github.com/FACorreiaa/Aviation-tracker/controller/svg"

	"github.com/FACorreiaa/Aviation-tracker/controller/models"
)

func (h *Handlers) renderAirlineSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home", Icon: svg2.HomeIcon()},
		{Path: "/airlines", Label: "Airline", Icon: svg2.ArrowRightIcon()},
		{Path: "/airlines/map", Label: "Airline map temporary", Icon: svg2.MapIcon()},
		{Path: "/airlines/tax", Label: "Airline Tax", Icon: svg2.MapIcon()},
		{Path: "/airlines/aircraft", Label: "Aircraft", Icon: svg2.MapIcon()},
		{Path: "/airlines/airplane", Label: "Airplane", Icon: svg2.MapIcon()},
		{Path: "/settings", Label: "Settings", Icon: svg2.SettingsIcon()},
		{Path: "/log-out", Label: "Log out", Icon: svg2.LogoutIcon()},
	}
	return sidebar
}

func (h *Handlers) getAirlinesLocations() ([]models.Airline, error) {
	a, err := h.core.airlines.GetAirlinesLocations(context.Background())
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (h *Handlers) getTotalAirline() (int, error) {
	total, err := h.core.airlines.GetAirlineSum(context.Background())
	pageSize := 10
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (h *Handlers) getAirline(_ http.ResponseWriter, r *http.Request) (int, []models.Airline, error) {
	pageSize := 10
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	a, err := h.core.airlines.GetAirlines(context.Background(), page, pageSize)
	if err != nil {
		return 0, nil, err
	}

	return page, a, nil
}

func (h *Handlers) renderAirlineTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	columnNames := []string{"Airline Name", "Date Founded", "Fleet Average Age", "Fleet Size",
		"Call Sign", "Hub Code", "Status", "Type", "Country name",
	}

	page, al, _ := h.getAirline(w, r)
	nextPage := page + 1
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	lastPage, err := h.getTotalAirline()
	if err != nil {
		return nil, err
	}
	a := models.AirlineTable{
		Column:   columnNames,
		Airline:  al,
		PrevPage: prevPage,
		NextPage: nextPage,
		Page:     page,
		LastPage: lastPage,
	}
	airlineTable := airline.AirlineTable(a)

	return airlineTable, nil
}

func (h *Handlers) airlineMainPage(w http.ResponseWriter, r *http.Request) error {
	taxTable, err := h.renderAirlineTable(w, r)
	sidebar := h.renderAirlineSidebar()
	if err != nil {
		return err
	}
	a := airline.AirlineLayoutPage("Airline", "Check models about aircrafts", taxTable, sidebar)
	return h.CreateLayout(w, r, "Airline Tax Page", a).Render(context.Background(), w)
}

func (h *Handlers) airlineLocationPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderAirlineSidebar()
	airlines, err := h.getAirlinesLocations()
	if err != nil {
		return err
	}
	a := airline.AirlineLocationsPage(sidebar, airlines)
	return h.CreateLayout(w, r, "Airline Locations Page", a).Render(context.Background(), w)
}
