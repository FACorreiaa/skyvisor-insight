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
	pageSize := 10
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	orderBy := r.URL.Query().Get("orderBy")
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	a, err := h.core.airports.GetAirports(context.Background(), page, pageSize, orderBy)
	if err != nil {
		return 0, nil, err
	}

	return page, a, nil
}

func (h *Handlers) getAirportsLocation() ([]models.Airport, error) {
	a, err := h.core.airports.GetAirportsLocation(context.Background())
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (h *Handlers) getTotalAirports() (int, error) {
	total, err := h.core.airports.GetSum(context.Background())
	pageSize := 10
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
	pageSize := 10
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	orderBy := r.URL.Query().Get("orderBy")

	if err != nil {
		// Handle error or set a default page number
		page = 1
	}
	ap, err := h.core.airports.GetAirportByName(context.Background(), param, page, pageSize, orderBy)
	if err != nil {
		return 0, nil, err
	}
	return page, ap, err
}

func (h *Handlers) renderAirportTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	columnNames := []models.ColumnItems{
		{Title: "Airport Name", Icon: svg2.ArrowOrderIcon()},
		{Title: "Country Name", Icon: svg2.ArrowOrderIcon()},
		{Title: "Phone Number", Icon: svg2.ArrowOrderIcon()},
		{Title: "Timezone", Icon: svg2.ArrowOrderIcon()},
		{Title: "GMT", Icon: svg2.ArrowOrderIcon()},
		{Title: "Latitude", Icon: svg2.ArrowOrderIcon()},
		{Title: "Longitude", Icon: svg2.ArrowOrderIcon()},
	}

	var ap []models.Airport
	var page int

	param := r.FormValue("search")
	orderBy := r.FormValue("orderBy")
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

	lastPage, err := h.getTotalAirports()
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
	}
	airportTable := airport.AirportTableComponent(a)

	return airportTable, nil
}

func (h *Handlers) renderSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home", Icon: svg2.HomeIcon()},
		{
			Label: "Airports",
			Icon:  svg2.BuildingOfficeIcon(),
			SubItems: []models.SidebarItem{
				{Path: "/airports", Label: "Airport data", Icon: svg2.BuildingOfficeIcon()},
				{Path: "/airports/locations", Label: "Airport locations", Icon: svg2.MapIcon()},
			},
		},
		{Path: "/settings", Label: "Settings", Icon: svg2.SettingsIcon()},
		{Path: "/log-out", Label: "Log out", Icon: svg2.LogoutIcon()},
	}
	return sidebar
}

func (h *Handlers) airportPage(w http.ResponseWriter, r *http.Request) error {
	at, err := h.renderAirportTable(w, r)
	sidebar := h.renderSidebar()
	if err != nil {
		HandleError(err, "Error fetching airport data table")
		return err
	}
	a := airport.AirportPage(at, sidebar, "Airports", "Check out airport locations")
	return h.CreateLayout(w, r, "Airport Page", a).Render(context.Background(), w)
}

func (h *Handlers) airportLocationPage(w http.ResponseWriter, r *http.Request) error {
	al, err := h.getAirportsLocation()
	sidebar := h.renderSidebar()
	if err != nil {
		HandleError(err, "Error fetching airport location table")
		return err
	}
	a := airport.AirportLocationsPage(sidebar, al, "Airports", "Check out airport locations")
	return h.CreateLayout(w, r, "Airport Locations Page", a).Render(context.Background(), w)
}

func (h *Handlers) airportDetailsPage(w http.ResponseWriter, r *http.Request) error {
	ad, err := h.getAirportDetails(w, r)
	sidebar := h.renderSidebar()
	if err != nil {
		HandleError(err, "Error fetching airport details page")
		return err
	}
	a := airport.AirportDetailsPage(sidebar, ad, "Airports", "Check out airport information")
	return h.CreateLayout(w, r, "Airport Details Page", a).Render(context.Background(), w)
}
