package controller

import (
	"context"
	"net/http"

	"github.com/FACorreiaa/go-ollama/controller/html/pages"
)

func (h *Handlers) airportPage(w http.ResponseWriter, r *http.Request) error {
	//change templ layout and add data to the templates

	//just hardcode the table values for now and improve solution later
	columnNames := []string{"Airport Name", "Country Name", "Phone Number",
		"Timezone", "GMT", "Latitude", "Longitude",
	}

	airports, err := h.core.airports.GetAirports(context.Background())
	if err != nil {
		return err
	}

	//for _, a := range airports {
	//	// pass data to the table and to the map.
	//	println(a)
	//}

	// doesn't make sense to do it like this

	//if len(airports) > 0 {
	//	airportType := reflect.TypeOf(airports[0])
	//
	//
	//	for i := 0; i < airportType.NumField(); i++ {
	//		field := airportType.Field(i)
	//		columnNames = append(columnNames, field.Name)
	//	}
	//}

	airport := pages.AirportPage(columnNames, airports)
	return h.CreateLayout(w, r, "Airport Page", airport).Render(context.Background(), w)
}
