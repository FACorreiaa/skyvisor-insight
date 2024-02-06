package controller

import (
	"context"
	"math"
	"net/http"
	"strconv"

	svg2 "github.com/FACorreiaa/Aviation-tracker/controller/svg"

	"github.com/a-h/templ"

	airline "github.com/FACorreiaa/Aviation-tracker/controller/html/airlines"

	"github.com/FACorreiaa/Aviation-tracker/controller/models"
)

func (h *Handlers) renderAirlineSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home", Icon: svg2.HomeIcon()},
		{
			Label: "Airlines",
			Icon:  svg2.TicketIcon(),
			SubItems: []models.SidebarItem{
				{Path: "/airlines/airline", Label: "Airline", Icon: svg2.TicketIcon()},
				{Path: "/airlines/airline/location", Label: "Airline location", Icon: svg2.MapIcon()},
			},
		},
		{Path: "/airlines/tax", Label: "Airline Tax", Icon: svg2.CreditCardIcon()},
		{Path: "/airlines/aircraft", Label: "Aircraft", Icon: svg2.PaperAirplaneIcon()},
		{Path: "/airlines/airplane", Label: "Airplane", Icon: svg2.PaperAirplaneIcon()},
		{Path: "/settings", Label: "Settings", Icon: svg2.SettingsIcon()},

		{Path: "/log-out", Label: "Log out", Icon: svg2.LogoutIcon()},
	}
	return sidebar
}

func (h *Handlers) getAirlinesLocations() ([]models.Airline, error) {
	al, err := h.core.airlines.GetAirlinesLocations(context.Background())
	if err != nil {
		return nil, err
	}

	return al, nil
}

func (h *Handlers) getAirlineByName(_ http.ResponseWriter, r *http.Request) (int, []models.Airline, error) {
	param := r.FormValue("search")
	pageSize := 10
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}
	al, err := h.core.airlines.GetAirlineByName(context.Background(), param, page, pageSize)
	if err != nil {
		return 0, nil, err
	}
	return page, al, err
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

	al, err := h.core.airlines.GetAirlines(context.Background(), page, pageSize)
	if err != nil {
		return 0, nil, err
	}

	return page, al, nil
}

func (h *Handlers) renderAirlineTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	columnNames := []models.ColumnItems{
		{Title: "Airline Name", Icon: svg2.ArrowOrderIcon()},
		{Title: "Date Founder", Icon: svg2.ArrowOrderIcon()},
		{Title: "Fleet Average Size", Icon: svg2.ArrowOrderIcon()},
		{Title: "Fleet Size", Icon: svg2.ArrowOrderIcon()},
		{Title: "Call Sign", Icon: svg2.ArrowOrderIcon()},
		{Title: "Hub Code", Icon: svg2.ArrowOrderIcon()},
		{Title: "Status", Icon: svg2.ArrowOrderIcon()},
		{Title: "Type", Icon: svg2.ArrowOrderIcon()},
		{Title: "Country Name", Icon: svg2.ArrowOrderIcon()},
	}

	param := r.FormValue("search")
	var page int

	var al []models.Airline
	fullPage, airlineList, _ := h.getAirline(w, r)
	filteredPage, filteredAirline, _ := h.getAirlineByName(w, r)

	if len(param) > 0 {
		al = filteredAirline
		page = filteredPage
	} else {
		al = airlineList
		page = fullPage
	}

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
		Column:      columnNames,
		Airline:     al,
		PrevPage:    prevPage,
		NextPage:    nextPage,
		Page:        page,
		LastPage:    lastPage,
		SearchParam: param,
	}
	airlineTable := airline.AirlineTable(a)

	return airlineTable, nil
}

func (h *Handlers) airlineMainPage(w http.ResponseWriter, r *http.Request) error {
	table, err := h.renderAirlineTable(w, r)

	sidebar := h.renderAirlineSidebar()
	if err != nil {
		return err
	}
	a := airline.AirlineLayoutPage("Airline", "Check out data about Airlines", table, sidebar)
	return h.CreateLayout(w, r, "Airline Tax Page", a).Render(context.Background(), w)
}

func (h *Handlers) airlineLocationPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderAirlineSidebar()
	airlines, err := h.getAirlinesLocations()
	if err != nil {
		return err
	}
	a := airline.AirlineLocationsPage(sidebar, airlines, "Airline", "Check out airport locations")
	return h.CreateLayout(w, r, "Airline Locations Page", a).Render(context.Background(), w)
}
