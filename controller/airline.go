package controller

import (
	"context"
	"github.com/a-h/templ"
	"math"
	"net/http"
	"strconv"

	airline "github.com/FACorreiaa/Aviation-tracker/controller/html/airlines"
	svg2 "github.com/FACorreiaa/Aviation-tracker/controller/svg"

	"github.com/FACorreiaa/Aviation-tracker/controller/models"
)

func (h *Handlers) renderAirlineSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home", Icon: svg2.HomeIcon()},
		{Path: "/airlines", Label: "Airline", Icon: svg2.ArrowRightIcon()},
		{Path: "/airlines/tax", Label: "Airline Tax", Icon: svg2.MapIcon()},
		{Path: "/airlines/aircraft", Label: "Aircraft", Icon: svg2.MapIcon()},
		{Path: "/airlines/airplane", Label: "Airplane", Icon: svg2.MapIcon()},
		{Path: "/settings", Label: "Settings", Icon: svg2.SettingsIcon()},
		{Path: "/log-out", Label: "Log out", Icon: svg2.LogoutIcon()},
	}
	return sidebar
}

func (h *Handlers) airlinePage(w http.ResponseWriter, r *http.Request) error {
	//airportTable, err := h.renderAirportTable(w, r)
	sidebar := h.renderAirlineSidebar()

	a, _ := h.renderAirlineTaxTable(w, r)
	airport := airline.AirlineLayoutPage("Airline Page", "Check data about airlines", a, sidebar)
	return h.CreateLayout(w, r, "Airline Page", airport).Render(context.Background(), w)
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

func (h *Handlers) getAirline(w http.ResponseWriter, r *http.Request) (int, []models.Airline, error) {
	pageSize := 10
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	aircraft, err := h.core.airlines.GetAirlines(context.Background(), page, pageSize)
	if err != nil {
		return 0, nil, err
	}

	return page, aircraft, nil
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
	airport := airline.AirlineLayoutPage("Airline", "Check models about aircrafts", taxTable, sidebar)
	return h.CreateLayout(w, r, "Airline Tax Page", airport).Render(context.Background(), w)
}
