package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
	svg2 "github.com/FACorreiaa/Aviation-tracker/app/svg"
	"github.com/FACorreiaa/Aviation-tracker/app/view/locations"
	"github.com/a-h/templ"
	"github.com/gorilla/mux"
)

func (h *Handler) renderLocationsBar() []models.SidebarItem {
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

func (h *Handler) getCities(_ http.ResponseWriter, r *http.Request) (int, []models.City, error) {
	pageSize := 20
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	cityName := r.FormValue("city_name")
	currencyName := r.FormValue("currency_name")
	phonePrefix := r.FormValue("phone_prefix")
	gmt := r.FormValue("gmt")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		HandleError(err, "page can't go lower than 1")
		page = 1
	}

	c, err := h.service.GetCity(context.Background(), page, pageSize, orderBy, sortBy, cityName, currencyName, phonePrefix, gmt)
	if err != nil {
		return 0, nil, err
	}

	return page, c, nil
}

func (h *Handler) getCityDetails(_ http.ResponseWriter, r *http.Request) (models.City, error) {
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

	c, err := h.service.GetCityByID(context.Background(), id)
	if err != nil {
		HandleError(err, "Error fetching city details")
		return models.City{}, err
	}

	return c, nil
}

func (h *Handler) renderCityTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	var sortAux string

	cityName := r.FormValue("city_name")
	currencyName := r.FormValue("currency_name")
	phonePrefix := r.FormValue("phone_prefix")
	gmt := r.FormValue("gmt")
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

	lastPage, err := h.service.GetAllCities()
	if err != nil {
		HandleError(err, "Error fetching cities")
		return nil, err
	}
	ct := models.CityTable{
		Column:             columnNames,
		City:               c,
		PrevPage:           prevPage,
		NextPage:           nextPage,
		Page:               page,
		LastPage:           lastPage,
		FilterCityName:     cityName,
		FilterCurrencyName: currencyName,
		FilterGMT:          gmt,
		FilterPhonePrefix:  phonePrefix,
		OrderParam:         orderBy,
		SortParam:          sortAux,
	}
	cityTable := locations.CityTable(ct)

	return cityTable, nil
}

func (h *Handler) CityMainPage(w http.ResponseWriter, r *http.Request) error {
	table, err := h.renderCityTable(w, r)
	sidebar := h.renderLocationsBar()
	if err != nil {
		HandleError(err, "Error rendering table")
		return err
	}
	c := locations.CityLayoutPage("Cities", "Check cities data around the world", table, sidebar)
	return h.CreateLayout(w, r, "City Page", c).Render(context.Background(), w)
}

func (h *Handler) CityLocationsPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderLocationsBar()
	c, err := h.service.GetCityLocations()
	if err != nil {
		HandleError(err, "Error fetching locations")
		return err
	}
	cl := locations.CityLocations(sidebar, c, "City Locations", "Check location of cities around the world")
	return h.CreateLayout(w, r, "City locations page", cl).Render(context.Background(), w)
}

func (h *Handler) CityDetailsPage(w http.ResponseWriter, r *http.Request) error {
	c, err := h.getCityDetails(w, r)
	sidebar := h.renderLocationsBar()
	if err != nil {
		HandleError(err, "Error fetching city details page")
		return err
	}
	a := locations.CityDetailsPage(sidebar, c, "City", "Check information about cities")
	return h.CreateLayout(w, r, "City Details Page", a).Render(context.Background(), w)
}

// Country

func (h *Handler) getCountries(_ http.ResponseWriter, r *http.Request) (int, []models.Country, error) {
	pageSize := 20
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	countryName := r.FormValue("country_name")
	capital := r.FormValue("capital")
	continent := r.FormValue("continent")
	currencyCode := r.FormValue("currency_code")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	c, err := h.service.GetCountry(context.Background(), page, pageSize, orderBy, sortBy, countryName, capital, continent, currencyCode)
	if err != nil {
		return 0, nil, err
	}

	return page, c, nil
}

func (h *Handler) getCountryDetails(_ http.ResponseWriter, r *http.Request) (models.Country, error) {
	vars := mux.Vars(r)
	country, ok := vars["country_name"]
	if !ok {
		err := errors.New("country_name not found in path")
		HandleError(err, "Error fetching country_name")
		return models.Country{}, err
	}

	c, err := h.service.GetCountryByName(context.Background(), country)
	if err != nil {
		HandleError(err, "Error fetching country details")
		return models.Country{}, err
	}

	return c, nil
}

func (h *Handler) renderCountryTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	var page int
	var sortAux string

	countryName := r.FormValue("country_name")
	capital := r.FormValue("capital")
	continent := r.FormValue("continent")
	currencyCode := r.FormValue("currency_code")

	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")

	if sortBy == ASC {
		sortAux = DESC
	} else {
		sortAux = ASC
	}

	columnNames := []models.ColumnItems{
		{Title: "Country Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Capital", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Continent", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Population", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Currency Code", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Currency Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Phone Prefix", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Latitude", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Longitude", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
	}

	page, c, _ := h.getCountries(w, r)

	nextPage := page + 1
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	lastPage, err := h.service.GetAllCountries()
	if err != nil {
		return nil, err
	}

	ct := models.CountryTable{
		Column:             columnNames,
		Country:            c,
		PrevPage:           prevPage,
		NextPage:           nextPage,
		Page:               page,
		LastPage:           lastPage,
		FilterCountry:      countryName,
		FilterCapital:      capital,
		FilterContinent:    continent,
		FilterCurrencyCode: currencyCode,
		OrderParam:         orderBy,
		SortParam:          sortAux,
	}
	cityTable := locations.CountryTable(ct)

	return cityTable, nil
}

func (h *Handler) CountryMainPage(w http.ResponseWriter, r *http.Request) error {
	taxTable, err := h.renderCountryTable(w, r)
	if err != nil {
		HandleError(err, "Error rendering table")
		return err
	}
	country, err := h.service.GetCountryLocations()
	if err != nil {
		HandleError(err, "Error rendering country locations")
		return err
	}

	sidebar := h.renderLocationsBar()

	c := locations.CountryLayoutPage("Countries", "Check countries of the world", taxTable, sidebar, country)
	return h.CreateLayout(w, r, "Country Page", c).Render(context.Background(), w)
}

func (h *Handler) CountryLocationPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderLocationsBar()
	c, err := h.service.GetCountryLocations()
	if err != nil {
		HandleError(err, "Error fetching countries")
		return err
	}

	cl := locations.CountryLocations(sidebar, c, "Country capital Locations",
		"Check data of country capitals around the world")
	return h.CreateLayout(w, r, "Country locations page", cl).Render(context.Background(), w)
}

func (h *Handler) CountryDetailsPage(w http.ResponseWriter, r *http.Request) error {
	c, err := h.getCountryDetails(w, r)
	sidebar := h.renderLocationsBar()
	if err != nil {
		HandleError(err, "Error fetching country details page")
		return err
	}
	a := locations.CountryDetailsPage(sidebar, c, "Country", "Check information about countries")
	return h.CreateLayout(w, r, "Country Details Page", a).Render(context.Background(), w)
}
