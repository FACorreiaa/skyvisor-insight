package controller

import (
	"context"
	"math"
	"net/http"
	"strconv"

	svg2 "github.com/FACorreiaa/Aviation-tracker/controller/svg"

	airline "github.com/FACorreiaa/Aviation-tracker/controller/html/airlines"
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
	"github.com/a-h/templ"
)

func (h *Handlers) getAircraft(_ http.ResponseWriter, r *http.Request) (int, []models.Aircraft, error) {
	pageSize := 25
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	param := r.FormValue("search")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	a, err := h.core.airlines.GetAircraft(context.Background(), page, pageSize, param, orderBy, sortBy)
	if err != nil {
		HandleError(err, "Error fetching aircrafts")
		return 0, nil, err
	}

	return page, a, nil
}

func (h *Handlers) getAllAircraftService() (int, error) {
	total, err := h.core.airlines.GetAircraftSum(context.Background())
	pageSize := 25
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (h *Handlers) renderAirlineAircraftTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	var sortAux string
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	param := r.FormValue("search")

	if sortBy == ASC {
		sortAux = DESC
	} else {
		sortAux = ASC
	}

	columnNames := []models.ColumnItems{
		{Title: "Aircraft Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Model Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Construction Number", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Number of Engines", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Type of Engine", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Date of first flight", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Line Number", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Model Code", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Plane Age", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Plane Class", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Plane Owner", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Plane Series", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Plane Plane Status", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
	}

	page, a, _ := h.getAircraft(w, r)
	nextPage := page + 1
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	lastPage, err := h.getAllAircraftService()
	if err != nil {
		HandleError(err, "Error fetching last page")
		return nil, err
	}
	data := models.AircraftTable{
		Column:      columnNames,
		Aircraft:    a,
		PrevPage:    prevPage,
		NextPage:    nextPage,
		Page:        page,
		LastPage:    lastPage,
		SearchParam: param,
		OrderParam:  orderBy,
		SortParam:   sortAux,
	}
	taxTable := airline.AirlineAircraftTable(data)

	return taxTable, nil
}

func (h *Handlers) airlineAircraftPage(w http.ResponseWriter, r *http.Request) error {
	taxTable, err := h.renderAirlineAircraftTable(w, r)
	sidebar := h.renderAirlineSidebar()
	if err != nil {
		return err
	}
	a := airline.AirlineLayoutPage("Aircrafts", "Check models about aircrafts", taxTable, sidebar)
	return h.CreateLayout(w, r, "Aircraft Tax Page", a).Render(context.Background(), w)
}
