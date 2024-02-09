package controller

import (
	"context"
	"math"
	"net/http"
	"strconv"

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

func (h *Handlers) getCityLocations() ([]models.City, error) {
	c, err := h.core.locations.GetCityLocation(context.Background())
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (h *Handlers) getTotalCities() (int, error) {
	total, err := h.core.locations.GetCitySum(context.Background())
	pageSize := 10
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (h *Handlers) getCities(_ http.ResponseWriter, r *http.Request) (int, []models.City, error) {
	pageSize := 10
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	c, err := h.core.locations.GetCity(context.Background(), page, pageSize, orderBy, sortBy)
	if err != nil {
		return 0, nil, err
	}

	return page, c, nil
}

func (h *Handlers) getCityByName(_ http.ResponseWriter, r *http.Request) (int, []models.City, error) {
	param := r.FormValue("search")
	pageSize := 10
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")

	if err != nil {
		// Handle error or set a default page number
		page = 1
	}
	c, err := h.core.locations.GetCityByName(context.Background(), page, pageSize, param, orderBy, sortBy)
	if err != nil {
		return 0, nil, err
	}
	return page, c, err
}

func (h *Handlers) renderCityTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	c := make([]models.City, 0)
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
		{Title: "City Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Timezone", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "GMT", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Continent", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Country Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Currency Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Phone Prefix", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Latitude", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Longitude", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
	}

	fullPage, cityList, _ := h.getCities(w, r)
	filteredPage, filteredCity, _ := h.getCityByName(w, r)

	if len(param) > 0 {
		c = filteredCity
		page = filteredPage
	} else {
		c = cityList
		page = fullPage
	}

	nextPage := page + 1
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	lastPage, err := h.getTotalCities()
	if err != nil {
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
		return err
	}
	c := locations.LocationsLayoutPage("Cities", "Check cities data around the world", taxTable, sidebar)
	return h.CreateLayout(w, r, "City Page", c).Render(context.Background(), w)
}

func (h *Handlers) cityLocationsPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderLocationsBar()
	c, err := h.getCityLocations()
	if err != nil {
		return err
	}
	cl := locations.CityLocations(sidebar, c, "City Locations", "Check location of cities around the world")
	return h.CreateLayout(w, r, "City locations page", cl).Render(context.Background(), w)
}
