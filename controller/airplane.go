package controller

import (
	"context"
	svg2 "github.com/FACorreiaa/Aviation-tracker/controller/svg"
	"math"
	"net/http"
	"strconv"

	airline "github.com/FACorreiaa/Aviation-tracker/controller/html/airlines"
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
	"github.com/a-h/templ"
)

func (h *Handlers) getAirplane(_ http.ResponseWriter, r *http.Request) (int, []models.Airplane, error) {
	pageSize := 10
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	param := r.FormValue("search")
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	a, err := h.core.airlines.GetAirplanes(context.Background(), page, pageSize, orderBy, sortBy, param)
	if err != nil {
		return 0, nil, err
	}

	return page, a, nil
}

func (h *Handlers) getTotalAirplanes() (int, error) {
	total, err := h.core.airlines.GetAirplaneSum(context.Background())
	pageSize := 10
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (h *Handlers) renderAirlineAirplaneTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
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
		{Title: "Model Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Airline Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Plane Series", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Plane Owner", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Plane Class", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Plane Age", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Plane Status", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Line Number", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "First Flight Date", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Engine Type", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Engine Count", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Construction Number", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Production Line", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Test Registration Date", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Registration Date", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Registration Number", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
	}

	page, ap, _ := h.getAirplane(w, r)
	nextPage := page + 1
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	lastPage, err := h.getTotalAirplanes()
	if err != nil {
		return nil, err
	}
	a := models.AirplaneTable{
		Column:      columnNames,
		Airplane:    ap,
		PrevPage:    prevPage,
		NextPage:    nextPage,
		Page:        page,
		LastPage:    lastPage,
		SearchParam: param,
		OrderParam:  orderBy,
		SortParam:   sortAux,
	}
	airlineTable := airline.AirplaneTable(a)

	return airlineTable, nil
}

func (h *Handlers) airlineAirplanePage(w http.ResponseWriter, r *http.Request) error {
	a, err := h.renderAirlineAirplaneTable(w, r)
	sidebar := h.renderAirlineSidebar()
	if err != nil {
		return err
	}
	al := airline.AirlineLayoutPage("Airplane", "Check models about airplanes", a, sidebar)
	return h.CreateLayout(w, r, "Airplane Page", al).Render(context.Background(), w)
}
