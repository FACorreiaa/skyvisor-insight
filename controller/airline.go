package controller

import (
	"context"
	"net/http"

	airline "github.com/FACorreiaa/Aviation-tracker/controller/html/airlines"
	svg2 "github.com/FACorreiaa/Aviation-tracker/controller/svg"

	"github.com/FACorreiaa/Aviation-tracker/controller/models"
)

func (h *Handlers) renderAirlineSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home", Icon: svg2.HomeIcon()},
		{Path: "/airlines", Label: "Airlines", Icon: svg2.ArrowRightIcon()},
		{Path: "/airlines/tax", Label: "Show Airline Tax", Icon: svg2.MapIcon()},
		{Path: "/airlines/aircraft", Label: "Aircraft", Icon: svg2.MapIcon()},
		{Path: "/settings", Label: "Settings", Icon: svg2.SettingsIcon()},
		{Path: "/log-out", Label: "Log out", Icon: svg2.LogoutIcon()},
	}
	return sidebar
}

func (h *Handlers) airlinePage(w http.ResponseWriter, r *http.Request) error {
	//airportTable, err := h.renderAirportTable(w, r)
	sidebar := h.renderAirlineSidebar()

	//mock

	a, _ := h.renderAirlineTaxTable(w, r)
	airport := airline.AirlineLayoutPage("Airline Page", "Check data about airlines", a, sidebar)
	return h.CreateLayout(w, r, "Airport Page", airport).Render(context.Background(), w)
}

func (h *Handlers) airlineAircraftPage(w http.ResponseWriter, r *http.Request) error {
	taxTable, err := h.renderAirlineAircraftTable(w, r)
	sidebar := h.renderAirlineSidebar()
	if err != nil {
		return err
	}
	airport := airline.AirlineLayoutPage("Aircrafts", "Check models about aircrafts", taxTable, sidebar)
	return h.CreateLayout(w, r, "Airline Tax Page", airport).Render(context.Background(), w)
}
