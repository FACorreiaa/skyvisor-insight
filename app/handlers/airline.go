package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"context"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
	svg2 "github.com/FACorreiaa/Aviation-tracker/app/svg"
	airline "github.com/FACorreiaa/Aviation-tracker/app/view/airlines"
	"github.com/a-h/templ"
	"github.com/gorilla/mux"
)

// Airline
func (h *Handler) renderAirlineSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home", Icon: svg2.HomeIcon()},
		{Path: "/airlines/airline", Label: "Airlines", Icon: svg2.CreditCardIcon()},
		{Path: "/airlines/airline/location", Label: "Airline location", Icon: svg2.MapIcon()},
		{Path: "/airlines/tax", Label: "Airline Tax", Icon: svg2.CreditCardIcon()},
		{Path: "/airlines/aircraft", Label: "Aircraft", Icon: svg2.PaperAirplaneIcon()},
		{Path: "/airlines/airplane", Label: "Airplane", Icon: svg2.PaperAirplaneIcon()},
		{Path: "/settings", Label: "Settings", Icon: svg2.SettingsIcon()},

		{Path: "/log-out", Label: "Log out", Icon: svg2.LogoutIcon()},
	}
	return sidebar
}

func (h *Handler) getAirline(_ http.ResponseWriter, r *http.Request) (int, []models.Airline, error) {
	pageSize := 30
	param := r.FormValue("search")
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	al, err := h.service.GetAirlines(context.Background(), page, pageSize, orderBy, sortBy, param)
	if err != nil {
		HandleError(err, "Error fetching airlines")
		return 0, nil, err
	}

	return page, al, nil
}

func (h *Handler) getAirlineDetails(_ http.ResponseWriter, r *http.Request) (models.Airline, error) {
	vars := mux.Vars(r)
	airlineName, ok := vars["airline_name"]
	if !ok {
		err := errors.New("airline_name not found in path")
		HandleError(err, "Error fetching airline_name")
		return models.Airline{}, err
	}

	c, err := h.service.GetAirlineByName(context.Background(), airlineName)
	if err != nil {
		HandleError(err, "Error fetching airline_name details")
		return models.Airline{}, err
	}

	return c, nil
}

func (h *Handler) renderAirlineTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
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

	lastPage, err := h.service.GetAllAirlineService()
	if err != nil {
		HandleError(err, "error fetching total airline")
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

func (h *Handler) AirlineMainPage(w http.ResponseWriter, r *http.Request) error {
	var table, err = h.renderAirlineTable(w, r)
	al, err := h.service.GetAirlinesLocationService()

	sidebar := h.renderAirlineSidebar()
	if err != nil {
		HandleError(err, "Error rendering airline table")
		return err
	}
	a := airline.AirlineMainPageLayout("Airline", "Check data about Airlines", table, sidebar, al)
	return h.CreateLayout(w, r, "Airline Tax Page", a).Render(context.Background(), w)
}

func (h *Handler) AirlineLocationPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderAirlineSidebar()
	al, err := h.service.GetAirlinesLocationService()
	if err != nil {
		HandleError(err, "Error rendering locations")
		return err
	}
	a := airline.AirlineLocationsPage(sidebar, al, "Airline", "Check airline expanded locations")
	return h.CreateLayout(w, r, "Airline Details Page", a).Render(context.Background(), w)
}

func (h *Handler) AirlineDetailsPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderAirlineSidebar()
	al, err := h.getAirlineDetails(w, r)
	if err != nil {
		HandleError(err, "Error rendering details")
		return err
	}
	a := airline.AirlineDetailsPage(sidebar, al, "Airline", "Check airport locations")
	return h.CreateLayout(w, r, "Airline Locations Page", a).Render(context.Background(), w)
}

// Aircraft

// Airplane

// Tax
