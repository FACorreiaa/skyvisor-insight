package controller

import (
	"context"
	"math"
	"net/http"
	"strconv"

	airline "github.com/FACorreiaa/Aviation-tracker/controller/html/airlines"
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
	"github.com/a-h/templ"
)

func (h *Handlers) getAirlineAircraft(_ http.ResponseWriter, r *http.Request) (int, []models.Aircraft, error) {
	pageSize := 10
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	a, err := h.core.airlines.GetAircraft(context.Background(), page, pageSize)
	if err != nil {
		return 0, nil, err
	}

	return page, a, nil
}

func (h *Handlers) getTotalAircraft() (int, error) {
	total, err := h.core.airlines.GetAircraftSum(context.Background())
	pageSize := 10
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (h *Handlers) renderAirlineAircraftTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	columnNames := []string{"Aircraft Name", "Model Name", "Construction Number", "Number of Engines",
		"Type of Engine", "Date of first flight", "Line Number", "Model Code", "Plane Age", "Plane Class", "Plane Owner",
		"Plane Series", "Plane Status",
	}

	page, a, _ := h.getAirlineAircraft(w, r)
	nextPage := page + 1
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	lastPage, err := h.getTotalAircraft()
	if err != nil {
		return nil, err
	}
	data := models.AircraftTable{
		Column:   columnNames,
		Aircraft: a,
		PrevPage: prevPage,
		NextPage: nextPage,
		Page:     page,
		LastPage: lastPage,
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
