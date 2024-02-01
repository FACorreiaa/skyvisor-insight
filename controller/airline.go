package controller

import (
	"context"
	"net/http"
	"strconv"

	airline "github.com/FACorreiaa/Aviation-tracker/controller/html/airlines"
	svg2 "github.com/FACorreiaa/Aviation-tracker/controller/svg"

	"github.com/FACorreiaa/Aviation-tracker/controller/models"
	"github.com/a-h/templ"
)

// tax

func (h *Handlers) getAirlineTax(w http.ResponseWriter, r *http.Request) (int, []models.Tax, error) {
	pageSize := 10
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	tax, err := h.core.airlines.GetTax(context.Background(), page, pageSize)
	if err != nil {
		return 0, nil, err
	}

	return page, tax, nil
}

func (h *Handlers) getAirlineAircraft(w http.ResponseWriter, r *http.Request) (int, []models.Aircraft, error) {
	pageSize := 10
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	aircraft, err := h.core.airlines.GetAircraft(context.Background(), page, pageSize)
	if err != nil {
		return 0, nil, err
	}

	return page, aircraft, nil
}

func (h *Handlers) getTotalTax() (int, error) {
	total, err := h.core.airlines.GetSum(context.Background())
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (h *Handlers) getTotalAircraft() (int, error) {
	total, err := h.core.airlines.GetAircraftSum(context.Background())
	if err != nil {
		return 0, err
	}
	return total, nil
}

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

func (h *Handlers) renderAirlineTaxTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	columnNames := []string{"Tax Name", "Airline Name", "Country Name", "City Name"}

	page, tax, _ := h.getAirlineTax(w, r)
	nextPage := page + 1
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	lastPage, err := h.getTotalTax()
	if err != nil {
		return nil, err
	}
	taxData := models.TaxTable{
		Column:   columnNames,
		Tax:      tax,
		PrevPage: prevPage,
		NextPage: nextPage,
		Page:     page,
		LastPage: lastPage,
	}
	taxTable := airline.AirlineTaxTable(taxData)

	return taxTable, nil
}

func (h *Handlers) renderAirlineAircraftTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	columnNames := []string{"Aircraft Name", "Model Name", "Construction Number", "Number of Engines",
		"Type of Engine", "Date of first flight", "Line Number", "Model Code", "Plane Age", "Plane Class", "Plane Owner",
		"Plane Series", "Plane Status",
	}

	page, aircraft, _ := h.getAirlineAircraft(w, r)
	nextPage := page + 1
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	lastPage, err := h.getTotalAircraft()
	if err != nil {
		return nil, err
	}
	aircraftData := models.AircraftTable{
		Column:   columnNames,
		Aircraft: aircraft,
		PrevPage: prevPage,
		NextPage: nextPage,
		Page:     page,
		LastPage: lastPage,
	}
	taxTable := airline.AirlineAircraftTable(aircraftData)

	return taxTable, nil
}

func (h *Handlers) airlineTaxPage(w http.ResponseWriter, r *http.Request) error {
	taxTable, err := h.renderAirlineTaxTable(w, r)
	sidebar := h.renderAirlineSidebar()
	if err != nil {
		return err
	}
	airport := airline.AirlineLayoutPage("Airline Tax", "Check data about tax", taxTable, sidebar)
	return h.CreateLayout(w, r, "Airline Tax Page", airport).Render(context.Background(), w)
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
