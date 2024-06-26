package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"context"

	httperror "github.com/FACorreiaa/Aviation-tracker/app/errors"
	"github.com/FACorreiaa/Aviation-tracker/app/models"
	svg2 "github.com/FACorreiaa/Aviation-tracker/app/static/svg"
	airline "github.com/FACorreiaa/Aviation-tracker/app/view/airlines"
	"github.com/a-h/templ"
	"github.com/gorilla/mux"
)

// Airline
func (h *Handler) renderAirlineSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home", Icon: svg2.HomeIcon()},
		{Path: "/airlines/airline", Label: "Airlines", Icon: svg2.CreditCardIcon()},
		{Path: "/airlines/airline/location", Label: "Airline location", Icon: svg2.MapIcon()},
		{Path: "/airlines/tax", Label: "Airline Tax", Icon: svg2.CreditCardIcon()},
		{Path: "/airlines/aircraft", Label: "Aircraft", Icon: svg2.PaperAirplaneIcon()},
		{Path: "/airlines/airplane", Label: "Airplane", Icon: svg2.PaperAirplaneIcon()},
		{Path: "/settings", Label: "Settings", Icon: svg2.SettingsIcon()},

		{Path: "/log-out", Label: "Log out", Icon: svg2.LogoutIcon()},
	}
	return sidebar
}

func (h *Handler) getAirline(w http.ResponseWriter, r *http.Request) (int, []models.Airline, error) {
	pageSize := 20
	name := r.FormValue("airline_name")
	callSign := r.FormValue("call_sign")
	hubCode := r.FormValue("hub_code")
	countryName := r.FormValue("country_name")

	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	al, err := h.service.GetAirlines(context.Background(), page, pageSize, orderBy, sortBy,
		name, callSign, hubCode, countryName)
	if err != nil {
		HandleError(err, "Error fetching airlines")
		return 0, nil, err
	}

	return page, al, nil
}

func (h *Handler) getAirlineDetails(w http.ResponseWriter, r *http.Request) (models.Airline, error) {
	vars := mux.Vars(r)
	airlineName, ok := vars["airline_name"]
	if !ok {
		err := errors.New("airline_name not found in path")
		HandleError(err, "Error fetching airline_name")
		httperror.ErrNotFound.WriteError(w)
		return models.Airline{}, err
	}

	c, err := h.service.GetAirlineByName(context.Background(), airlineName)
	if err != nil {
		HandleError(err, "Error fetching airline_name details")
		httperror.ErrInvalidID.WriteError(w)
		return models.Airline{}, err
	}

	return c, nil
}

func (h *Handler) renderAirlineTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	name := r.FormValue("airline_name")
	callSign := r.FormValue("call_sign")
	hubCode := r.FormValue("hub_Code")
	countryName := r.FormValue("country_name")
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	var sortAux string

	if sortBy == ASC {
		sortAux = DESC
	} else {
		sortAux = ASC
	}
	columnNames := []models.ColumnItems{
		{Title: "Airline Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Date Founded", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Fleet Average Size", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Fleet Size", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Call Sign", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Hub Code", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Status", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Type", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Country Name", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
	}

	page, al, err := h.getAirline(w, r)

	if err != nil {
		HandleError(err, "Error fetching airlines")
		httperror.ErrNotFound.WriteError(w)
		return nil, err
	}

	nextPage := page + 1
	prevPage := page - 1
	if prevPage <= 1 {
		prevPage = 1
	}

	lastPage, err := h.service.GetAllAirline()

	if err != nil {
		HandleError(err, "error fetching total airline")
		httperror.ErrNotFound.WriteError(w)
		return nil, err
	}

	a := models.AirlineTable{
		Column:            columnNames,
		Airline:           al,
		PrevPage:          prevPage,
		NextPage:          nextPage,
		Page:              page,
		LastPage:          lastPage,
		FilterName:        name,
		FilterCallSign:    callSign,
		FilterHubCode:     hubCode,
		FilterCountryName: countryName,
		OrderParam:        orderBy,
		SortParam:         sortAux,
	}
	airlineTable := airline.AirlineTable(a)

	return airlineTable, nil
}

func (h *Handler) AirlineMainPage(w http.ResponseWriter, r *http.Request) error {
	var table, err = h.renderAirlineTable(w, r)
	if err != nil {
		HandleError(err, "Error rendering airline table")
		httperror.ErrInternalServer.WriteError(w)
		return err
	}

	al, err := h.service.GetAirlinesLocation()
	if err != nil {
		HandleError(err, "Error rendering airlines location")
		httperror.ErrNotFound.WriteError(w)
		return err
	}

	sidebar := h.renderAirlineSidebar()

	a := airline.AirlineMainPageLayout("Airline", "Check data about Airlines", table, sidebar, al)
	return h.CreateLayout(w, r, "Airline Tax Page", a).Render(context.Background(), w)
}

func (h *Handler) AirlineLocationPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderAirlineSidebar()
	al, err := h.service.GetAirlinesLocation()
	if err != nil {
		HandleError(err, "Error rendering locations")
		httperror.ErrInternalServer.WriteError(w)
		return err
	}
	a := airline.AirlineLocationsPage(sidebar, al, "Airline", "Check airline expanded locations")
	return h.CreateLayout(w, r, "Airline Details Page", a).Render(context.Background(), w)
}

func (h *Handler) AirlineDetailsPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderAirlineSidebar()
	al, err := h.getAirlineDetails(w, r)
	if err != nil {
		HandleError(err, "Error rendering details")
		httperror.ErrNotFound.WriteError(w)
		return err
	}
	a := airline.AirlineDetailsPage(sidebar, al, "Airline", "Check airport locations")
	return h.CreateLayout(w, r, "Airline Locations Page", a).Render(context.Background(), w)
}

// Aircraft
func (h *Handler) renderAirlineAircraftTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	var sortAux string
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	aircraftName := r.FormValue("aircraft_name")
	typeEngine := r.FormValue("type_engine")
	modelCode := r.FormValue("model_code")
	planeOwner := r.FormValue("plane_owner")

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
	if prevPage <= 1 {
		prevPage = 1
	}

	lastPage, err := h.service.GetAllAircraft()
	if err != nil {
		HandleError(err, "Error fetching last page")
		httperror.ErrNotFound.WriteError(w)
		return nil, err
	}
	data := models.AircraftTable{
		Column:             columnNames,
		Aircraft:           a,
		PrevPage:           prevPage,
		NextPage:           nextPage,
		Page:               page,
		LastPage:           lastPage,
		FilterAircraft:     aircraftName,
		FilterTypeOfEngine: typeEngine,
		FilterModelCode:    modelCode,
		FilterPlaneOwner:   planeOwner,
		OrderParam:         orderBy,
		SortParam:          sortAux,
	}
	taxTable := airline.AirlineAircraftTable(data)

	return taxTable, nil
}

func (h *Handler) getAircraft(w http.ResponseWriter, r *http.Request) (int, []models.Aircraft, error) {
	pageSize := 20
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	aircraftName := r.FormValue("aircraft_name")
	typeEngine := r.FormValue("type_engine")
	modelCode := r.FormValue("model_code")
	planeOwner := r.FormValue("plane_owner")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	a, err := h.service.GetAircraft(context.Background(), page, pageSize, aircraftName, orderBy, sortBy,
		typeEngine, modelCode, planeOwner)
	if err != nil {
		HandleError(err, "Error fetching aircraft")
		httperror.ErrNotFound.WriteError(w)
		return 0, nil, err
	}

	return page, a, nil
}

func (h *Handler) AirlineAircraftPage(w http.ResponseWriter, r *http.Request) error {
	taxTable, err := h.renderAirlineAircraftTable(w, r)
	sidebar := h.renderAirlineSidebar()
	if err != nil {
		httperror.ErrInternalServer.WriteError(w)
		return err
	}
	a := airline.AirlineLayoutPage("Aircraft", "Check models about aircraft", taxTable, sidebar)
	return h.CreateLayout(w, r, "Aircraft Tax Page", a).Render(context.Background(), w)
}

// Airplane

func (h *Handler) getAirplane(w http.ResponseWriter, r *http.Request) (int, []models.Airplane, error) {
	pageSize := 20
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	airlineName := r.FormValue("airline_name")
	modelName := r.FormValue("model_name")
	productionLine := r.FormValue("production_line")
	registrationNumber := r.FormValue("registration_number")
	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	a, err := h.service.GetAirplanes(context.Background(), page, pageSize, orderBy, sortBy,
		airlineName, modelName, productionLine, registrationNumber)
	if err != nil {
		httperror.ErrNotFound.WriteError(w)
		return 0, nil, err
	}

	return page, a, nil
}

func (h *Handler) renderAirlineAirplaneTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	var sortAux string

	airlineName := r.FormValue("airlineName")
	modelName := r.FormValue("model_name")
	productionLine := r.FormValue("production_line")
	registrationNumber := r.FormValue("registration_number")

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
	if prevPage <= 1 {
		prevPage = 1
	}

	lastPage, err := h.service.GetAllAirplanes()
	if err != nil {
		HandleError(err, "Error fetching last page")
		httperror.ErrNotFound.WriteError(w)
		return nil, err
	}
	a := models.AirplaneTable{
		Column:                   columnNames,
		Airplane:                 ap,
		PrevPage:                 prevPage,
		NextPage:                 nextPage,
		Page:                     page,
		LastPage:                 lastPage,
		FilterAirlineName:        airlineName,
		FilterModelName:          modelName,
		FilterProductionLine:     productionLine,
		FilterRegistrationNumber: registrationNumber,
		OrderParam:               orderBy,
		SortParam:                sortAux,
	}
	airlineTable := airline.AirplaneTable(a)

	return airlineTable, nil
}

func (h *Handler) AirlineAirplanePage(w http.ResponseWriter, r *http.Request) error {
	a, err := h.renderAirlineAirplaneTable(w, r)
	sidebar := h.renderAirlineSidebar()
	if err != nil {
		httperror.ErrInternalServer.WriteError(w)
		return err
	}
	al := airline.AirlineLayoutPage("Airplane", "Check models about airplanes", a, sidebar)
	return h.CreateLayout(w, r, "Airplane Page", al).Render(context.Background(), w)
}

// Tax

func (h *Handler) getAirlineTax(w http.ResponseWriter, r *http.Request) (int, []models.Tax, error) {
	pageSize := 20
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	taxName := r.FormValue("tax_name")
	countryName := r.FormValue("country_name")
	airlineName := r.FormValue("airline_name")

	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	t, err := h.service.GetTax(context.Background(), page, pageSize, orderBy, sortBy, taxName, countryName, airlineName)
	if err != nil {
		httperror.ErrNotFound.WriteError(w)
		return 0, nil, err
	}

	return page, t, nil
}

func (h *Handler) renderAirlineTaxTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	var page int
	var sortAux string

	taxName := r.FormValue("tax_name")
	airlineName := r.FormValue("airline_name")
	countryName := r.FormValue("country_name")
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
	if prevPage <= 1 {
		prevPage = 1
	}

	lastPage, err := h.service.GetSum()
	if err != nil {
		HandleError(err, "Error fetching tax")
		httperror.ErrNotFound.WriteError(w)
		return nil, err
	}
	taxData := models.TaxTable{
		Column:        columnNames,
		Tax:           tax,
		PrevPage:      prevPage,
		NextPage:      nextPage,
		Page:          page,
		LastPage:      lastPage,
		FilterTax:     taxName,
		FilterAirline: airlineName,
		FilterCountry: countryName,
		OrderParam:    orderBy,
		SortParam:     sortAux,
	}
	taxTable := airline.AirlineTaxTable(taxData)

	return taxTable, nil
}

func (h *Handler) AirlineTaxPage(w http.ResponseWriter, r *http.Request) error {
	taxTable, err := h.renderAirlineTaxTable(w, r)
	sidebar := h.renderAirlineSidebar()
	if err != nil {
		httperror.ErrNotFound.WriteError(w)
		HandleError(err, "Error rendering table")
		return err
	}
	t := airline.AirlineLayoutPage("Airline Tax", "Check data about tax", taxTable, sidebar)
	return h.CreateLayout(w, r, "Airline Tax Page", t).Render(context.Background(), w)
}
