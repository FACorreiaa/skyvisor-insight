package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"context"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
	svg2 "github.com/FACorreiaa/Aviation-tracker/app/svg"
	"github.com/FACorreiaa/Aviation-tracker/app/view/components"
	"github.com/FACorreiaa/Aviation-tracker/app/view/flights"
	"github.com/a-h/templ"
	"github.com/gorilla/mux"
)

// https://openlayers.org/en/latest/examples/feature-move-animation.html future feature
// https://openlayers.org/en/latest/examples/flight-animation.html
// future feature on this branch for flights with destination

// need to change sql query later

func (h *Handler) renderLiveLocationsSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home", Icon: svg2.HomeIcon()},
		{
			Label: "All Flights",
			Icon:  svg2.AcademicCapIcon(),
			SubItems: []models.SidebarItem{
				{Path: "/flights/flight", Label: "Flights", Icon: svg2.HomeIcon()},
				{Path: "/flights/flight/location", Label: "Preview Flights", Icon: svg2.HomeIcon()},
			},
		},
		{
			Label: "Live Flights",
			Icon:  svg2.GlobeEuropeIcon(),
			SubItems: []models.SidebarItem{
				{Path: "/flights/flight/live", Label: "Live Flights", Icon: svg2.GlobeEuropeIcon()},
				{Path: "/flights/flight/location/air/live", Label: "Live Flights Locations", Icon: svg2.MapIcon()},
			},
		},

		{
			Label: "Active Flights",
			Icon:  svg2.GlobeEuropeIcon(),
			SubItems: []models.SidebarItem{
				{Path: "/flights/flight/status/active", Label: "Active Flights", Icon: svg2.GlobeEuropeIcon()},
				{Path: "/flights/flight/status/location/active", Label: "Active Flights Locations", Icon: svg2.MapIcon()},
			},
		},
		{
			Label: "Landed Flights",
			Icon:  svg2.GlobeEuropeIcon(),
			SubItems: []models.SidebarItem{
				{Path: "/flights/flight/status/landed", Label: "Landed Flights", Icon: svg2.GlobeEuropeIcon()},
				{Path: "/flights/flight/status/location/landed", Label: "Landed Flights Location", Icon: svg2.MapIcon()},
			},
		},
		{Path: "/flights/flight/status/scheduled", Label: "Scheduled Flights", Icon: svg2.HomeIcon()},
		{Path: "/flights/flight/status/cancelled", Label: "Cancelled Flights", Icon: svg2.HomeIcon()},
		{Path: "/flights/flight/status/incident", Label: "Incident Flights", Icon: svg2.HomeIcon()},
		{Path: "/flights/flight/status/diverted", Label: "Diverted Flights", Icon: svg2.HomeIcon()},
		{Path: "/settings", Label: "Settings", Icon: svg2.SettingsIcon()},
		{Path: "/log-out", Label: "Log out", Icon: svg2.LogoutIcon()},
	}
	return sidebar
}

func (h *Handler) getFlights(_ http.ResponseWriter, r *http.Request) (int, []models.LiveFlights, error) {
	pageSize := 30
	vars := mux.Vars(r)
	flightStatus := vars["flight_status"]

	airlineName := r.FormValue("airline_name")
	flightNumber := r.FormValue("flight_number")
	flightStats := r.FormValue("flight_status")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	orderBy := r.URL.Query().Get("orderBy")
	sortBy := r.URL.Query().Get("sortBy")

	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	var lf []models.LiveFlights

	if flightStatus == "" {
		lf, err = h.service.GetAllFlights(context.Background(), page, pageSize, orderBy,
			sortBy, flightNumber, airlineName, flightStats)
	} else {
		lf, err = h.service.GetAllFlightsByStatus(context.Background(),
			page, pageSize, orderBy, sortBy, flightNumber, flightStatus)
	}

	if err != nil {
		return 0, nil, err
	}

	return page, lf, nil
}

func (h *Handler) getLiveFlights(_ http.ResponseWriter, r *http.Request) (int, []models.LiveFlights, error) {
	pageSize := 30
	//vars := mux.Vars(r)
	//flightStatus := vars["flight_status"]

	airlineName := r.FormValue("airline_name")
	flightNumber := r.FormValue("flight_number")
	flightStats := r.FormValue("flight_status")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	orderBy := r.URL.Query().Get("orderBy")
	sortBy := r.URL.Query().Get("sortBy")

	if err != nil {
		// Handle error or set a default page number
		page = 1
	}

	var lf []models.LiveFlights

	lf, err = h.service.GetLiveFlights(context.Background(), page, pageSize, orderBy, sortBy,
		flightNumber, airlineName, flightStats)

	if err != nil {
		return 0, nil, err
	}

	return page, lf, nil
}

func (h *Handler) renderFlightsTable(w http.ResponseWriter,
	r *http.Request) (templ.Component, error) {

	// vars := mux.Vars(r)
	// flightStatusRoute := vars["flight_status"]

	var sortAux string

	airlineName := r.FormValue("airline_name")
	flightNumber := r.FormValue("flight_number")
	flightStatus := r.FormValue("flight_status")

	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")

	if sortBy == ASC {
		sortAux = DESC
	} else {
		sortAux = ASC
	}

	columnNames := []models.ColumnItems{
		{Title: "Flight Number", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Airline", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Flight Status", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Flight Date", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Airport Departure", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Estimated Departure", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Airport Arrival", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Estimated Arrival", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Arrival Delay", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Departure Delay", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Departure Terminal", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Departure Gate", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
	}

	page, lf, err := h.getFlights(w, r)
	if err != nil {
		HandleError(err, "Error fetching total flights")
		return nil, err
	}

	if len(lf) == 0 {
		// If empty, return a message component
		message := components.EmptyPageComponent()
		return message, nil
	}

	nextPage := page + 1
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	lastPage, err := h.service.GetAllFlightsSum()
	if err != nil {
		HandleError(err, "Error fetching total flights")
		return nil, err
	}
	f := models.FlightsTable{
		Column:             columnNames,
		Flights:            lf,
		PrevPage:           prevPage,
		NextPage:           nextPage,
		Page:               page,
		LastPage:           lastPage,
		FilterFlightNumber: flightNumber,
		FilterFlightStatus: flightStatus,
		FilterAirlineName:  airlineName,
		OrderParam:         orderBy,
		SortParam:          sortAux,
	}
	flightsTable := flights.AllFlightsTableComponent(f)

	return flightsTable, nil
}

func (h *Handler) renderLiveFlightsTable(w http.ResponseWriter,
	r *http.Request) (templ.Component, error) {

	var sortAux string

	airlineName := r.FormValue("airline_name")
	flightNumber := r.FormValue("flight_number")
	flightStatus := r.FormValue("flight_status")

	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")

	if sortBy == ASC {
		sortAux = DESC
	} else {
		sortAux = ASC
	}

	columnNames := []models.ColumnItems{
		{Title: "Flight Number", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Airline", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Airport Departure", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Estimated Departure", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Airport Arrival", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Estimated Arrival", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Arrival Delay", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Departure Delay", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Live Latitude", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Live Longitude", Icon: svg2.ArrowOrderIcon(), SortParam: sortAux},
	}

	page, lf, err := h.getLiveFlights(w, r)
	if err != nil {
		HandleError(err, "Error fetching total flights")
		return nil, err
	}

	if len(lf) == 0 {
		// If empty, return a message component
		message := components.EmptyPageComponent()
		return message, nil
	}

	nextPage := page + 1
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	lastPage, err := h.service.GetAllFlightsSum()
	if err != nil {
		HandleError(err, "Error fetching total flights")
		return nil, err
	}
	f := models.FlightsTable{
		Column:             columnNames,
		Flights:            lf,
		PrevPage:           prevPage,
		NextPage:           nextPage,
		Page:               page,
		LastPage:           lastPage,
		FilterFlightNumber: flightNumber,
		FilterFlightStatus: flightStatus,
		FilterAirlineName:  airlineName,
		OrderParam:         orderBy,
		SortParam:          sortAux,
	}
	flightsTable := flights.LiveFlightsTableComponent(f)

	return flightsTable, nil
}

func (h *Handler) getFlightsDetails(_ http.ResponseWriter, r *http.Request) (models.LiveFlights, error) {
	vars := mux.Vars(r)
	flightNumber, ok := vars["flight_number"]
	if !ok {
		err := errors.New("flight_number not found in path")
		HandleError(err, "Error fetching flight_number")
		return models.LiveFlights{}, err
	}

	lf, err := h.service.GetFlightByID(context.Background(), flightNumber)
	if err != nil {
		HandleError(err, "Error flights details")
		return models.LiveFlights{}, err
	}

	return lf, nil
}

func (h *Handler) AllFlightsPage(w http.ResponseWriter, r *http.Request) error {
	table, err := h.renderFlightsTable(w, r)
	if err != nil {
		HandleError(err, "Error fetching flights table")
		return err
	}
	s := h.renderLiveLocationsSidebar()

	f := flights.AllFlightsPage(table, s, "Live Flights", "Check all flights going on")
	return h.CreateLayout(w, r, "Live Flights", f).Render(context.Background(), w)
}

func (h *Handler) DetailedFlightsPage(w http.ResponseWriter, r *http.Request) error {
	s := h.renderLiveLocationsSidebar()
	fd, err := h.getFlightsDetails(w, r)

	if err != nil {
		HandleError(err, "Error fetching flights details page")
		return err
	}

	f := flights.DetailedFlightsPage(s, fd, "Live Flights", "Detailed flight data")
	return h.CreateLayout(w, r, "Live Flights", f).Render(context.Background(), w)
}

func (h *Handler) FlightsLocation(w http.ResponseWriter, r *http.Request) error {
	s := h.renderLiveLocationsSidebar()
	fd, err := h.service.GetAllFlightsLocation()

	if err != nil {
		HandleError(err, "Error fetching flights details page")
		return err
	}

	f := flights.FLightsPreviewPage(s, fd, "Live Flights", "Detailed flight data")
	return h.CreateLayout(w, r, "Live Flights", f).Render(context.Background(), w)
}

func (h *Handler) FilteredFlightsPage(w http.ResponseWriter, r *http.Request) error {
	table, err := h.renderFlightsTable(w, r)
	if err != nil {
		HandleError(err, "Error fetching flights table")
		return err
	}
	s := h.renderLiveLocationsSidebar()

	f := flights.AllFlightsPage(table, s, "Live Flights", "Check all flights going on")
	return h.CreateLayout(w, r, "Live Flights", f).Render(context.Background(), w)
}

func (h *Handler) FlightsLocationsByStatus(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	flightStatus := vars["flight_status"]
	s := h.renderLiveLocationsSidebar()
	fd, err := h.service.GetAllFlightsLocationsByStatus(context.Background(), flightStatus)

	if err != nil {
		HandleError(err, "Error fetching flights details page")
		return err
	}

	f := flights.FLightsPreviewPage(s, fd, "Live Flights", "Detailed flight data")
	return h.CreateLayout(w, r, "Live Flights", f).Render(context.Background(), w)
}

func (h *Handler) LiveFlightsPage(w http.ResponseWriter, r *http.Request) error {
	table, err := h.renderLiveFlightsTable(w, r)
	if err != nil {
		HandleError(err, "Error fetching flights table")
		return err
	}
	s := h.renderLiveLocationsSidebar()

	f := flights.AllFlightsPage(table, s, "Live Flights", "Check on air flights")
	return h.CreateLayout(w, r, "Live Flights", f).Render(context.Background(), w)
}

func (h *Handler) LiveFlightsLocationsPage(w http.ResponseWriter, r *http.Request) error {
	s := h.renderLiveLocationsSidebar()
	fd, err := h.service.GetLiveFlightsLocations(context.Background())

	if err != nil {
		HandleError(err, "Error fetching flights details page")
		return err
	}

	f := flights.LiveFlightsLocationPage(s, fd, "Live Flights", "Detailed flight data")
	return h.CreateLayout(w, r, "Live Flights", f).Render(context.Background(), w)
}
