package controller

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strconv"

	svg2 "github.com/FACorreiaa/Aviation-tracker/controller/svg"
	"github.com/gorilla/mux"

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

func (h *Handlers) getAirlinesLocationService() ([]models.Airline, error) {
	al, err := h.core.airlines.GetAirlinesLocations(context.Background())
	if err != nil {
		HandleError(err, "Error fetching locations")
		return nil, err
	}

	return al, nil
}

func (h *Handlers) getAllAirlineService() (int, error) {
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
	param := r.FormValue("search")
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	al, err := h.core.airlines.GetAirlines(context.Background(), page, pageSize, orderBy, sortBy, param)
	if err != nil {
		HandleError(err, "Error fetching airlines")
		return 0, nil, err
	}

	return page, al, nil
}

func (h *Handlers) getAirlineDetails(_ http.ResponseWriter, r *http.Request) (models.Airline, error) {
	vars := mux.Vars(r)
	airlineName, ok := vars["airline_name"]
	if !ok {
		err := errors.New("airline_name not found in path")
		HandleError(err, "Error getting airline_name")
		return models.Airline{}, err
	}

	c, err := h.core.airlines.GetAirlineByName(context.Background(), airlineName)
	if err != nil {
		HandleError(err, "Error fetching airline_name details")
		return models.Airline{}, err
	}

	return c, nil
}

func (h *Handlers) renderAirlineTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	param := r.FormValue("search")
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	var sortAux string

	if sortBy == ASC {
		sortAux = DESC
	} else {
		sortAux = ASC
	}
	columnNames := []models.ColumnItems{
		{Title: "Airline Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Date Founded", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Fleet Average Size", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Fleet Size", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Call Sign", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Hub Code", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Status", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Type", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Country Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
	}

	var page int

	page, al, _ := h.getAirline(w, r)

	nextPage := page + 1
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	lastPage, err := h.getAllAirlineService()
	if err != nil {
		HandleError(err, "error getting total airline")
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
		OrderParam:  orderBy,
		SortParam:   sortAux,
	}
	airlineTable := airline.AirlineTable(a)

	return airlineTable, nil
}

func (h *Handlers) airlineMainPage(w http.ResponseWriter, r *http.Request) error {
	table, err := h.renderAirlineTable(w, r)

	sidebar := h.renderAirlineSidebar()
	if err != nil {
		HandleError(err, "Error rendering airline table")
		return err
	}
	a := airline.AirlineLayoutPage("Airline", "Check data about Airlines", table, sidebar)
	return h.CreateLayout(w, r, "Airline Tax Page", a).Render(context.Background(), w)
}

func (h *Handlers) airlineLocationPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderAirlineSidebar()
	al, err := h.getAirlinesLocationService()
	if err != nil {
		HandleError(err, "Error rendering locations")
		return err
	}
	a := airline.AirlineLocationsPage(sidebar, al, "Airline", "Check airline details")
	return h.CreateLayout(w, r, "Airline Details Page", a).Render(context.Background(), w)
}

func (h *Handlers) airlineDetailsPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderAirlineSidebar()
	al, err := h.getAirlineDetails(w, r)
	if err != nil {
		HandleError(err, "Error rendering details")
		return err
	}
	a := airline.AirlineDetailsPage(sidebar, al, "Airline", "Check airport locations")
	return h.CreateLayout(w, r, "Airline Locations Page", a).Render(context.Background(), w)
}
