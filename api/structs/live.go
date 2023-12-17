package structs

import (
	"github.com/google/uuid"
)

type FlightStatus string

const (
	Scheduled FlightStatus = "scheduled"
	Active    FlightStatus = "active"
	Landed    FlightStatus = "landed"
	Cancelled FlightStatus = "cancelled"
	Incident  FlightStatus = "incident"
	Diverted  FlightStatus = "diverted"
)

type LiveFlights struct {
	ID           uuid.UUID    `db:"id"`
	FlightDate   string       `json:"flight_date,omitempty"`
	FlightStatus FlightStatus `json:"flight_status,omitempty"`
	Departure    struct {
		Airport         string      `json:"airport"`
		Timezone        string      `json:"timezone"`
		Iata            string      `json:"iata"`
		Icao            string      `json:"icao"`
		Terminal        string      `json:"terminal"`
		Gate            interface{} `json:"gate"`
		Delay           *int        `json:"delay"`
		Scheduled       string      `json:"scheduled"`
		Estimated       string      `json:"estimated"`
		Actual          interface{} `json:"actual"`
		EstimatedRunway interface{} `json:"estimated_runway"`
		ActualRunway    interface{} `json:"actual_runway"`
	} `json:"departure,omitempty"`
	Arrival struct {
		Airport         string      `json:"airport"`
		Timezone        string      `json:"timezone"`
		Iata            string      `json:"iata"`
		Icao            string      `json:"icao"`
		Terminal        interface{} `json:"terminal"`
		Gate            interface{} `json:"gate"`
		Baggage         interface{} `json:"baggage"`
		Delay           *int        `json:"delay"`
		Scheduled       string      `json:"scheduled"`
		Estimated       string      `json:"estimated"`
		Actual          interface{} `json:"actual"`
		EstimatedRunway interface{} `json:"estimated_runway"`
		ActualRunway    interface{} `json:"actual_runway"`
	} `json:"arrival,omitempty"`
	Airline struct {
		Name string `json:"name"`
		Iata string `json:"iata"`
		Icao string `json:"icao"`
	} `json:"airline,omitempty"`
	Flight struct {
		Number     string `json:"number"`
		Iata       string `json:"iata"`
		Icao       string `json:"icao"`
		Codeshared struct {
			AirlineName  string `json:"airline_name"`
			AirlineIata  string `json:"airline_iata"`
			AirlineIcao  string `json:"airline_icao"`
			FlightNumber string `json:"flight_number"`
			FlightIata   string `json:"flight_iata"`
			FlightIcao   string `json:"flight_icao"`
		} `json:"codeshared,omitempty"`
	} `json:"flight"`
	Aircraft struct {
		AircraftRegistration string `json:"registration"`
		AircraftIata         string `json:"iata"`
		AircraftIcao         string `json:"icao"`
		AircraftIcao24       string `json:"icao24"`
	} `json:"aircraft,omitempty"`
	Live struct {
		LiveUpdated         string  `json:"updated"`
		LiveLatitude        float32 `json:"latitude,omitempty"`
		LiveLongitude       float32 `json:"longitude,omitempty"`
		LiveAltitude        int     `json:"altitude"`
		LiveDirection       float32 `json:"direction"`
		LiveSpeedHorizontal int     `json:"speed_horizontal"`
		LiveSpeedVertical   int     `json:"speed_vertical"`
		LiveIsGround        bool    `json:"is_ground"`
	} `json:"live,omitempty"`
	CreatedAt CustomTime `json:"created_at"`
}

type FlightApiData struct {
	Pagination Pagination    `json:"pagination"`
	Data       []LiveFlights `json:"data"`
}
