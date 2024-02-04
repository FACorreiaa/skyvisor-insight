package controller

import (
	"context"
	"math"
	"net/http"
	"strconv"

	"github.com/FACorreiaa/Aviation-tracker/controller/html/locations"
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
	"github.com/a-h/templ"
)

func (h *Handlers) getCountryLocations() ([]models.Country, error) {
	c, err := h.core.locations.GetCountryLocation(context.Background())
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (h *Handlers) getTotalCountries() (int, error) {
	total, err := h.core.locations.GetCountrySum(context.Background())
	pageSize := 10
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (h *Handlers) getCountries(_ http.ResponseWriter, r *http.Request) (int, []models.Country, error) {
	pageSize := 10
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	c, err := h.core.locations.GetCountry(context.Background(), page, pageSize)
	if err != nil {
		return 0, nil, err
	}

	return page, c, nil
}

func (h *Handlers) renderCountryTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	columnNames := []string{"Country Name", "Capital", "Continent", "Currency Name",
		"Population", "Currency Code", "Currency Name", "Latitude", "Longitude",
	}

	page, c, _ := h.getCountries(w, r)
	nextPage := page + 1
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	lastPage, err := h.getTotalCountries()
	if err != nil {
		return nil, err
	}

	ct := models.CountryTable{
		Column:   columnNames,
		Country:  c,
		PrevPage: prevPage,
		NextPage: nextPage,
		Page:     page,
		LastPage: lastPage,
	}
	cityTable := locations.CountryTable(ct)

	return cityTable, nil
}

func (h *Handlers) countryMainPage(w http.ResponseWriter, r *http.Request) error {
	taxTable, err := h.renderCountryTable(w, r)
	sidebar := h.renderLocationsBar()
	if err != nil {
		return err
	}
	c := locations.LocationsLayoutPage("Countries", "Check countries of the world", taxTable, sidebar)
	return h.CreateLayout(w, r, "Country Page", c).Render(context.Background(), w)
}

func (h *Handlers) countryLocationPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderLocationsBar()
	c, err := h.getCountryLocations()
	if err != nil {
		return err
	}

	cl := locations.CountryLocations(sidebar, c, "Country capital Locations",
		"Check data of country capitals around the world")
	return h.CreateLayout(w, r, "Country locations page", cl).Render(context.Background(), w)
}
