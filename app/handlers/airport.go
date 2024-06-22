package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"context"

	httperror "github.com/FACorreiaa/Aviation-tracker/app/errors"
	"github.com/FACorreiaa/Aviation-tracker/app/models"
	svg2 "github.com/FACorreiaa/Aviation-tracker/app/static/svg"
	airport "github.com/FACorreiaa/Aviation-tracker/app/view/airports"
	"github.com/a-h/templ"
	"github.com/gorilla/mux"
)

func (h *Handler) getAirports(w http.ResponseWriter, r *http.Request) (int, []models.Airport, error) {
	pageSize := 20
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	orderBy := r.URL.Query().Get("orderBy")
	sortBy := r.URL.Query().Get("sortBy")

	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	a, err := h.service.GetAirports(context.Background(), page, pageSize, orderBy, sortBy)
	if err != nil {
		httperror.ErrNotFound.WriteError(w)
		return 0, nil, err
	}

	return page, a, nil
}

func (h *Handler) getAirportDetails(w http.ResponseWriter, r *http.Request) (models.Airport, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["airport_id"]
	if !ok {
		err := errors.New("airport_id not found in path")
		HandleError(err, "Error fetching airport_id")
		httperror.ErrInternalServer.WriteError(w)
		return models.Airport{}, err
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		HandleError(err, "Error converting airport_id to integer")
		return models.Airport{}, err
	}

	ap, err := h.service.GetAirportByID(context.Background(), id)
	if err != nil {
		HandleError(err, "Error fetching airport details")
		httperror.ErrNotFound.WriteError(w)
		return models.Airport{}, err
	}

	return ap, nil
}

func (h *Handler) getAirportByName(w http.ResponseWriter, r *http.Request) (int, []models.Airport, error) {
	airportName := r.FormValue("airport_name")
	countryName := r.FormValue("country_name")
	gmt := r.FormValue("gmt")
	pageSize := 20
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	orderBy := r.URL.Query().Get("orderBy")
	sortBy := r.URL.Query().Get("sortBy")
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}
	ap, err := h.service.GetAirportByName(context.Background(), page, pageSize, orderBy, sortBy, airportName, countryName, gmt)
	if err != nil {
		httperror.ErrNotFound.WriteError(w)
		return 0, nil, err
	}
	return page, ap, err
}

func (h *Handler) renderAirportTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	var ap []models.Airport
	var page int
	var sortAux string

	airportName := r.FormValue("airport_name")
	countryName := r.FormValue("country_name")
	gmt := r.FormValue("gmt")
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

	if len(airportName) > 0 {
		ap = filteredAirport
		page = filteredPage
	} else {
		ap = airportList
		page = fullPage
	}

	if page-1 < 0 {
		return nil, nil
	}

	lastPage, err := h.service.GetAllAirports()
	if err != nil {
		httperror.ErrInternalServer.WriteError(w)
		HandleError(err, "Error fetching airports")
		return nil, err
	}
	a := models.AirportTable{
		Column:            columnNames,
		Airports:          ap,
		PrevPage:          page - 1,
		NextPage:          page + 1,
		Page:              page,
		LastPage:          lastPage,
		FilterAirportName: airportName,
		FilterCountryName: countryName,
		FilterGMT:         gmt,
		OrderParam:        orderBy,
		SortParam:         sortAux,
	}
	airportTable := airport.AirportTableComponent(a)

	return airportTable, nil
}

func (h *Handler) renderSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home", Icon: svg2.HomeIcon()},
		{Path: "/airports", Label: "Airports", Icon: svg2.BuildingOfficeIcon()},
		{Path: "/airports/locations", Label: "Airport locations", Icon: svg2.MapIcon()},
		{Path: "/settings", Label: "Settings", Icon: svg2.SettingsIcon()},
		{Path: "/log-out", Label: "Log out", Icon: svg2.LogoutIcon()},
	}
	return sidebar
}

func (h *Handler) AirportPage(w http.ResponseWriter, r *http.Request) error {
	at, err := h.renderAirportTable(w, r)
	if err != nil {
		httperror.ErrInternalServer.WriteError(w)
		HandleError(err, "Error fetching airport data table")
		return err
	}
	al, err := h.service.GetAirportsLocation()
	if err != nil {
		httperror.ErrNotFound.WriteError(w)
		HandleError(err, "Error getting airport locations")
		return err
	}

	sidebar := h.renderSidebar()

	a := airport.AirportPage(at, sidebar, "Airports", "Check airport locations", al)
	return h.CreateLayout(w, r, "Airport Page", a).Render(context.Background(), w)
}

func (h *Handler) AirportLocationPage(w http.ResponseWriter, r *http.Request) error {
	al, err := h.service.GetAirportsLocation()
	sidebar := h.renderSidebar()
	if err != nil {
		httperror.ErrInternalServer.WriteError(w)
		HandleError(err, "Error fetching airport location table")
		return err
	}
	a := airport.AirportLocationsPage(sidebar, al, "Airports", "Check airport locations")
	return h.CreateLayout(w, r, "Airport Locations Page", a).Render(context.Background(), w)
}

func (h *Handler) AirportDetailsPage(w http.ResponseWriter, r *http.Request) error {
	ad, err := h.getAirportDetails(w, r)
	sidebar := h.renderSidebar()
	if err != nil {
		httperror.ErrInternalServer.WriteError(w)
		HandleError(err, "Error fetching airport details page")
		return err
	}
	a := airport.AirportDetailsPage(sidebar, ad, "Airports", "Check airport information")
	return h.CreateLayout(w, r, "Airport Details Page", a).Render(context.Background(), w)
}
