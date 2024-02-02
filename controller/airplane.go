package controller

import (
	"context"
	airline "github.com/FACorreiaa/Aviation-tracker/controller/html/airlines"
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
	"github.com/a-h/templ"
	"math"
	"net/http"
	"strconv"
)

func (h *Handlers) getAirplane(w http.ResponseWriter, r *http.Request) (int, []models.Airplane, error) {
	pageSize := 10
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	a, err := h.core.airlines.GetAirplanes(context.Background(), page, pageSize)
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
	columnNames := []string{"Model Name", "Airline Name", "Plane Series", "Plane owner", "Plane class",
		"Plane age", "Plane status", "Line number", "First Flight Date", "Engine type",
		"Engine count", "Construction number", "Production line", "Test registration date",
		"Registration date", "Registration number",
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
		Column:   columnNames,
		Airplane: ap,
		PrevPage: prevPage,
		NextPage: nextPage,
		Page:     page,
		LastPage: lastPage,
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
