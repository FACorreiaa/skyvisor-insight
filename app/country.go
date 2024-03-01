package app

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
	svg2 "github.com/FACorreiaa/Aviation-tracker/app/svg"
	"github.com/FACorreiaa/Aviation-tracker/app/view/locations"
	"github.com/a-h/templ"
)

func (h *Handlers) getCountryLocationsService() ([]models.Country, error) {
	c, err := h.core.locations.GetCountryLocation(context.Background())
	if err != nil {
		HandleError(err, "Error fetching locations")
		return nil, err
	}

	return c, nil
}

func (h *Handlers) getAllCountriesService() (int, error) {
	total, err := h.core.locations.GetCountrySum(context.Background())
	pageSize := 30
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (h *Handlers) getCountries(_ http.ResponseWriter, r *http.Request) (int, []models.Country, error) {
	pageSize := 30
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	param := r.FormValue("search")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	c, err := h.core.locations.GetCountry(context.Background(), page, pageSize, orderBy, sortBy, param)
	if err != nil {
		return 0, nil, err
	}

	return page, c, nil
}

func (h *Handlers) getCountryDetails(_ http.ResponseWriter, r *http.Request) (models.Country, error) {
	vars := mux.Vars(r)
	country, ok := vars["country_name"]
	if !ok {
		err := errors.New("country_name not found in path")
		HandleError(err, "Error fetching country_name")
		return models.Country{}, err
	}

	c, err := h.core.locations.GetCountryByName(context.Background(), country)
	if err != nil {
		HandleError(err, "Error fetching country details")
		return models.Country{}, err
	}

	return c, nil
}

func (h *Handlers) renderCountryTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
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
		{Title: "Country Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Capital", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Continent", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Currency Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
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

	lastPage, err := h.getAllCountriesService()
	if err != nil {
		return nil, err
	}

	ct := models.CountryTable{
		Column:      columnNames,
		Country:     c,
		PrevPage:    prevPage,
		NextPage:    nextPage,
		Page:        page,
		LastPage:    lastPage,
		SearchParam: param,
		OrderParam:  orderBy,
		SortParam:   sortAux,
	}
	cityTable := locations.CountryTable(ct)

	return cityTable, nil
}

func (h *Handlers) countryMainPage(w http.ResponseWriter, r *http.Request) error {
	taxTable, err := h.renderCountryTable(w, r)
	country, err := h.getCountryLocationsService()

	sidebar := h.renderLocationsBar()
	if err != nil {
		HandleError(err, "Error rendering table")
		return err
	}
	c := locations.CountryLayoutPage("Countries", "Check countries of the world", taxTable, sidebar, country)
	return h.CreateLayout(w, r, "Country Page", c).Render(context.Background(), w)
}

func (h *Handlers) countryLocationPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderLocationsBar()
	c, err := h.getCountryLocationsService()
	if err != nil {
		HandleError(err, "Error fetching countries")
		return err
	}

	cl := locations.CountryLocations(sidebar, c, "Country capital Locations",
		"Check data of country capitals around the world")
	return h.CreateLayout(w, r, "Country locations page", cl).Render(context.Background(), w)
}

func (h *Handlers) countryDetailsPage(w http.ResponseWriter, r *http.Request) error {
	c, err := h.getCountryDetails(w, r)
	sidebar := h.renderLocationsBar()
	if err != nil {
		HandleError(err, "Error fetching country details page")
		return err
	}
	a := locations.CountryDetailsPage(sidebar, c, "Country", "Check information about countries")
	return h.CreateLayout(w, r, "Country Details Page", a).Render(context.Background(), w)
}
