package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/FACorreiaa/Aviation-tracker/controller/html/components"
	"github.com/FACorreiaa/Aviation-tracker/controller/html/pages"
	"github.com/FACorreiaa/Aviation-tracker/controller/html/svg"
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
	"github.com/a-h/templ"
)

func (h *Handlers) getAirlineTax(w http.ResponseWriter, r *http.Request) (int, []models.Tax, error) {
	pageSize := 10
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	tax, err := h.core.airlines.GetTax(context.Background(), page, pageSize)
	if err != nil {
		return 0, nil, err
	}

	return page, tax, nil
}

func (h *Handlers) getTotalTax() (int, error) {
	total, err := h.core.airlines.GetSum(context.Background())
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (h *Handlers) renderAirlineSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home", Icon: svg.HomeIcon()},
		{Path: "/airlines", Label: "Airlines", Icon: svg.ArrowRightIcon()},
		{Path: "/airlines/tax", Label: "Show Airline Tax", Icon: svg.MapIcon()},
		{Path: "/settings", Label: "Settings", Icon: svg.SettingsIcon()},
		{Path: "/log-out", Label: "Log out", Icon: svg.LogoutIcon()},
	}
	return sidebar
}

func (h *Handlers) airlinePage(w http.ResponseWriter, r *http.Request) error {
	//airportTable, err := h.renderAirportTable(w, r)
	sidebar := h.renderAirlineSidebar()

	airport := pages.AirlinePage(sidebar)
	return h.CreateLayout(w, r, "Airport Page", airport).Render(context.Background(), w)
}

func (h *Handlers) renderAirlineTaxTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	columnNames := []string{"Tax Name", "Airline Name", "Country Name", "City Name"}

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
		Column:   columnNames,
		Tax:      tax,
		PrevPage: prevPage,
		NextPage: nextPage,
		Page:     page,
		LastPage: lastPage,
	}
	taxTable := components.TaxTableComponent(taxData)

	return taxTable, nil
}

func (h *Handlers) airlineTaxPage(w http.ResponseWriter, r *http.Request) error {
	taxTable, err := h.renderAirlineTaxTable(w, r)
	sidebar := h.renderSidebar()
	if err != nil {
		return err
	}
	airport := pages.AirlineTaxPage(taxTable, sidebar)
	return h.CreateLayout(w, r, "Airline Tax Page", airport).Render(context.Background(), w)
}
