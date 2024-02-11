package controller

import (
	"context"
	"math"
	"net/http"
	"strconv"

	airline "github.com/FACorreiaa/Aviation-tracker/controller/html/airlines"
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
	svg2 "github.com/FACorreiaa/Aviation-tracker/controller/svg"
	"github.com/a-h/templ"
)

func (h *Handlers) getAirlineTax(_ http.ResponseWriter, r *http.Request) (int, []models.Tax, error) {
	pageSize := 10
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	param := r.FormValue("search")

	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	t, err := h.core.airlines.GetTax(context.Background(), page, pageSize, orderBy, sortBy, param)
	if err != nil {
		return 0, nil, err
	}

	return page, t, nil
}

func (h *Handlers) getTotalTax() (int, error) {
	total, err := h.core.airlines.GetSum(context.Background())
	pageSize := 10
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (h *Handlers) renderAirlineTaxTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
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
		{Title: "Tax Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Airline Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Country Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
	}

	page, tax, _ := h.getAirlineTax(w, r)

	nextPage := page + 1
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	lastPage, err := h.getTotalTax()
	if err != nil {
		return nil, err
	}
	taxData := models.TaxTable{
		Column:      columnNames,
		Tax:         tax,
		PrevPage:    prevPage,
		NextPage:    nextPage,
		Page:        page,
		LastPage:    lastPage,
		SearchParam: param,
		OrderParam:  orderBy,
		SortParam:   sortAux,
	}
	taxTable := airline.AirlineTaxTable(taxData)

	return taxTable, nil
}

func (h *Handlers) airlineTaxPage(w http.ResponseWriter, r *http.Request) error {
	taxTable, err := h.renderAirlineTaxTable(w, r)
	sidebar := h.renderAirlineSidebar()
	if err != nil {
		return err
	}
	t := airline.AirlineLayoutPage("Airline Tax", "Check data about tax", taxTable, sidebar)
	return h.CreateLayout(w, r, "Airline Tax Page", t).Render(context.Background(), w)
}
