package controller

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	airport "github.com/FACorreiaa/Aviation-tracker/controller/html/airports"
	svg2 "github.com/FACorreiaa/Aviation-tracker/controller/svg"

	"github.com/a-h/templ"

	"github.com/FACorreiaa/Aviation-tracker/controller/models"
)

func (h *Handlers) getAirports(_ http.ResponseWriter, r *http.Request) (int, []models.Airport, error) {
	pageSize := 15
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	orderBy := r.URL.Query().Get("orderBy")
	sortBy := r.URL.Query().Get("sortBy")

	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	a, err := h.core.airports.GetAirports(context.Background(), page, pageSize, orderBy, sortBy)
	if err != nil {
		return 0, nil, err
	}

	return page, a, nil
}

func (h *Handlers) getAirportsLocationService() ([]models.Airport, error) {
	a, err := h.core.airports.GetAirportsLocation(context.Background())
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (h *Handlers) getAllAirportsService() (int, error) {
	total, err := h.core.airports.GetSum(context.Background())
	pageSize := 15
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (h *Handlers) getAirportDetails(_ http.ResponseWriter, r *http.Request) (models.Airport, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["airport_id"]
	if !ok {
		err := errors.New("airport_id not found in path")
		HandleError(err, "Error getting airport_id")
		return models.Airport{}, err
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		HandleError(err, "Error converting airport_id to integer")
		return models.Airport{}, err
	}

	ap, err := h.core.airports.GetAirportByID(context.Background(), id)
	if err != nil {
		HandleError(err, "Error fetching airport details")
		return models.Airport{}, err
	}

	return ap, nil
}

func (h *Handlers) getAirportByName(_ http.ResponseWriter, r *http.Request) (int, []models.Airport, error) {
	param := r.FormValue("search")
	pageSize := 15
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	orderBy := r.URL.Query().Get("orderBy")
	sortBy := r.URL.Query().Get("sortBy")
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}
	ap, err := h.core.airports.GetAirportByName(context.Background(), param, page, pageSize, orderBy, sortBy)
	if err != nil {
		return 0, nil, err
	}
	return page, ap, err
}

func (h *Handlers) renderAirportTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	var ap = make([]models.Airport, 0)
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
		{Title: "Airport Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Country Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Phone Number", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Timezone", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "GMT", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Latitude", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Longitude", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
	}

	fullPage, airportList, _ := h.getAirports(w, r)
	filteredPage, filteredAirport, _ := h.getAirportByName(w, r)

	if len(param) > 0 {
		ap = filteredAirport
		page = filteredPage
	} else {
		ap = airportList
		page = fullPage
	}

	if page-1 < 0 {
		return nil, nil
	}

	lastPage, err := h.getAllAirportsService()
	if err != nil {
		return nil, err
	}
	a := models.AirportTable{
		Column:      columnNames,
		Airports:    ap,
		PrevPage:    page - 1,
		NextPage:    page + 1,
		Page:        page,
		LastPage:    lastPage,
		SearchParam: param,
		OrderParam:  orderBy,
		SortParam:   sortAux,
	}
	airportTable := airport.AirportTableComponent(a)

	return airportTable, nil
}

func (h *Handlers) renderSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home", Icon: svg2.HomeIcon()},
		{Path: "/airports", Label: "Airports", Icon: svg2.BuildingOfficeIcon()},
		{Path: "/airports/locations", Label: "Airport locations", Icon: svg2.MapIcon()},
		{Path: "/settings", Label: "Settings", Icon: svg2.SettingsIcon()},
		{Path: "/log-out", Label: "Log out", Icon: svg2.LogoutIcon()},
	}
	return sidebar
}

func (h *Handlers) airportPage(w http.ResponseWriter, r *http.Request) error {
	at, err := h.renderAirportTable(w, r)
	al, err := h.getAirportsLocationService()

	sidebar := h.renderSidebar()
	if err != nil {
		HandleError(err, "Error fetching airport data table")
		return err
	}
	a := airport.AirportPage(at, sidebar, "Airports", "Check airport locations", al)
	return h.CreateLayout(w, r, "Airport Page", a).Render(context.Background(), w)
}

func (h *Handlers) airportLocationPage(w http.ResponseWriter, r *http.Request) error {
	al, err := h.getAirportsLocationService()
	sidebar := h.renderSidebar()
	if err != nil {
		HandleError(err, "Error fetching airport location table")
		return err
	}
	a := airport.AirportLocationsPage(sidebar, al, "Airports", "Check airport locations")
	return h.CreateLayout(w, r, "Airport Locations Page", a).Render(context.Background(), w)
}

func (h *Handlers) airportDetailsPage(w http.ResponseWriter, r *http.Request) error {
	ad, err := h.getAirportDetails(w, r)
	sidebar := h.renderSidebar()
	if err != nil {
		HandleError(err, "Error fetching airport details page")
		return err
	}
	a := airport.AirportDetailsPage(sidebar, ad, "Airports", "Check airport information")
	return h.CreateLayout(w, r, "Airport Details Page", a).Render(context.Background(), w)
}
