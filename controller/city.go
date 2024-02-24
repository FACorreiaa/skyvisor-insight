package controller

import (
	"context"
	"errors"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/FACorreiaa/Aviation-tracker/controller/html/locations"
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
	svg2 "github.com/FACorreiaa/Aviation-tracker/controller/svg"
	"github.com/a-h/templ"
)

func (h *Handlers) renderLocationsBar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home", Icon: svg2.HomeIcon()},
		{
			Label: "Cities",
			Icon:  svg2.GlobeEuropeIcon(),
			SubItems: []models.SidebarItem{
				{Path: "/locations/city", Label: "City data", Icon: svg2.GlobeEuropeIcon()},
				{Path: "/locations/city/map", Label: "City locations", Icon: svg2.MapIcon()},
			},
		},
		{
			Label: "Countries",
			Icon:  svg2.GlobeIcon(),
			SubItems: []models.SidebarItem{
				{Path: "/locations/country", Label: "Country data", Icon: svg2.GlobeIcon()},
				{Path: "/locations/country/map", Label: "Country locations", Icon: svg2.MapIcon()},
			},
		},
		{Path: "/settings", Label: "Settings", Icon: svg2.SettingsIcon()},
		{Path: "/log-out", Label: "Log out", Icon: svg2.LogoutIcon()},
	}
	return sidebar
}

func (h *Handlers) getCityLocationsService() ([]models.City, error) {
	c, err := h.core.locations.GetCityLocation(context.Background())
	if err != nil {
		HandleError(err, "Error fetching cities")
		return nil, err
	}

	return c, nil
}

func (h *Handlers) getAllCitiesService() (int, error) {
	total, err := h.core.locations.GetCitySum(context.Background())
	pageSize := 15
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		HandleError(err, "Error fetching sum of cities")
		return 0, err
	}
	return lastPage, nil
}

func (h *Handlers) getCities(_ http.ResponseWriter, r *http.Request) (int, []models.City, error) {
	pageSize := 15
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	param := r.FormValue("search")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		HandleError(err, "page can't go lower than 1")
		page = 1
	}

	c, err := h.core.locations.GetCity(context.Background(), page, pageSize, orderBy, sortBy, param)
	if err != nil {
		return 0, nil, err
	}

	return page, c, nil
}

func (h *Handlers) getCityDetails(_ http.ResponseWriter, r *http.Request) (models.City, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["city_id"]
	log.Printf("Vars: %+v", vars)

	if !ok {
		err := errors.New("city_id not found in path")
		HandleError(err, "Error fetching city_id")
		return models.City{}, err
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		HandleError(err, "Error converting city_id to integer")
		return models.City{}, err
	}

	c, err := h.core.locations.GetCityByID(context.Background(), id)
	if err != nil {
		HandleError(err, "Error fetching city details")
		return models.City{}, err
	}

	return c, nil
}

func (h *Handlers) renderCityTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	var sortAux string

	param := r.FormValue("search")
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")

	if sortBy == ASC {
		sortAux = DESC
	} else {
		sortAux = ASC
	}

	//
	columnNames := []models.ColumnItems{
		{Title: "City Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Country Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Continent", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Currency Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Timezone", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "GMT", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Phone Prefix", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Latitude", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Longitude", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
	}

	page, c, _ := h.getCities(w, r)

	nextPage := page + 1
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	lastPage, err := h.getAllCitiesService()
	if err != nil {
		HandleError(err, "Error fetching cities")
		return nil, err
	}
	ct := models.CityTable{
		Column:      columnNames,
		City:        c,
		PrevPage:    prevPage,
		NextPage:    nextPage,
		Page:        page,
		LastPage:    lastPage,
		SearchParam: param,
		OrderParam:  orderBy,
		SortParam:   sortAux,
	}
	cityTable := locations.CityTable(ct)

	return cityTable, nil
}

func (h *Handlers) cityMainPage(w http.ResponseWriter, r *http.Request) error {
	taxTable, err := h.renderCityTable(w, r)
	sidebar := h.renderLocationsBar()
	if err != nil {
		HandleError(err, "Error rendering table")
		return err
	}
	cities, err := h.getCityLocationsService()
	if err != nil {
		HandleError(err, "Error fetching locations")
		return err
	}
	c := locations.CityLayoutPage("Cities", "Check cities data around the world", taxTable, sidebar, cities)
	return h.CreateLayout(w, r, "City Page", c).Render(context.Background(), w)
}

func (h *Handlers) cityLocationsPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderLocationsBar()
	c, err := h.getCityLocationsService()
	if err != nil {
		HandleError(err, "Error fetching locations")
		return err
	}
	cl := locations.CityLocations(sidebar, c, "City Locations", "Check location of cities around the world")
	return h.CreateLayout(w, r, "City locations page", cl).Render(context.Background(), w)
}

func (h *Handlers) cityDetailsPage(w http.ResponseWriter, r *http.Request) error {
	c, err := h.getCityDetails(w, r)
	sidebar := h.renderLocationsBar()
	if err != nil {
		HandleError(err, "Error fetching city details page")
		return err
	}
	a := locations.CityDetailsPage(sidebar, c, "City", "Check information about cities")
	return h.CreateLayout(w, r, "City Details Page", a).Render(context.Background(), w)
}
