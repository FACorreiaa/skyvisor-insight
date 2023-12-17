package structs

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

//airline

type Airline struct {
	ID                   string     `json:"id"`
	FleetAverageAge      float64    `json:"fleet_average_age,string"`
	AirlineId            int        `json:"airline_id,string"`
	Callsign             string     `json:"callsign"`
	HubCode              string     `json:"hub_code"`
	IataCode             string     `json:"iata_code"`
	IcaoCode             string     `json:"icao_code"`
	CountryISO2          string     `json:"country_iso2"`
	DateFounded          int        `json:"date_founded,string"`
	IataPrefixAccounting int        `json:"iata_prefix_accounting,string"`
	AirlineName          string     `json:"airline_name"`
	CountryName          string     `json:"country_name"`
	FleetSize            int        `json:"fleet_size,string"`
	Status               string     `json:"status"`
	Type                 string     `json:"type"`
	CreatedAt            CustomTime `db:"created_at" json:"created_at"`
}

//aircrafts

type Aircraft struct {
	ID           string     `json:"id"`
	IataCode     string     `json:"iata_code"`
	AircraftName string     `json:"aircraft_name"`
	PlaneTypeId  int        `json:"plane_type_id,string"`
	CreatedAt    CustomTime `db:"created_at" json:"created_at"`
}

//airplanes

type Airplane struct {
	ID                     string      `json:"id" `
	IataType               string      `json:"iata_type"`
	AirplaneId             int         `json:"airplane_id,string"`
	AirlineIataCode        string      `json:"airline_iata_code"`
	IataCodeLong           string      `json:"iata_code_long"`
	IataCodeShort          string      `json:"iata_code_short"`
	AirlineIcaoCode        interface{} `json:"airline_icao_code"`
	ConstructionNumber     string      `json:"construction_number"`
	DeliveryDate           CustomTime  `json:"delivery_date"`
	EnginesCount           int         `json:"engines_count,string"`
	EnginesType            string      `json:"engines_type"`
	FirstFlightDate        CustomTime  `json:"first_flight_date"`
	IcaoCodeHex            string      `json:"icao_code_hex"`
	LineNumber             interface{} `json:"line_number"`
	ModelCode              string      `json:"model_code"`
	RegistrationNumber     string      `json:"registration_number"`
	TestRegistrationNumber interface{} `json:"test_registration_number"`
	PlaneAge               int         `json:"plane_age,string"`
	PlaneClass             interface{} `json:"plane_class"`
	ModelName              string      `json:"model_name"`
	PlaneOwner             interface{} `json:"plane_owner"`
	PlaneSeries            string      `json:"plane_series"`
	PlaneStatus            string      `json:"plane_status"`
	ProductionLine         string      `json:"production_line"`
	RegistrationDate       CustomTime  `json:"registration_date"`
	RolloutDate            CustomTime  `json:"rollout_date"`
	CreatedAt              CustomTime  `json:"created_at,string"`
}

type Tax struct {
	ID        string     `json:"id"`
	TaxId     int        `json:"tax_id,string,omitempty"`
	TaxName   string     `json:"tax_name"`
	IataCode  string     `json:"iata_code"`
	CreatedAt CustomTime `db:"created_at" json:"created_at"`
}

type TaxApiData struct {
	Data []Tax `json:"data"`
}

type AircraftApiData struct {
	Data []Aircraft `json:"data"`
}

type AirlineApiData struct {
	Data []Airline `json:"data"`
}

type AirplaneApiData struct {
	Data []Airplane `json:"data"`
}

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	var dateStr string
	err := json.Unmarshal(data, &dateStr)
	if err != nil {
		fmt.Println("Error parsing date: ", err)
		return err
	}

	// Check if the date is "0000-00-00" and set it to a default value
	if dateStr == "0000-00-00" {
		ct.Time = time.Time{} // Assign zero value of CustomTime
		return nil
	}

	// Check if the date string is empty and set it to a default value
	if dateStr == "" {
		ct.Time = time.Time{} // Assign zero value of CustomTime
		return nil
	}

	// Parse the date using the predefined time layout
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		fmt.Println("Error parsing date: ", err)
		return err
	}

	ct.Time = t
	return nil
}

//

// Implement driver.Valuer interface
func (ct CustomTime) Value() (driver.Value, error) {
	// Return the underlying time value as a string in RFC3339 format
	return ct.Time.Format(time.RFC3339), nil
}

// Implement sql.Scanner interface
func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		// Handle NULL values by setting the time to the zero value
		ct.Time = time.Time{}
		return nil
	}

	switch t := value.(type) {
	case time.Time:
		ct.Time = t
		return nil
	case []byte:
		parsedTime, err := time.Parse("2006-01-02", string(t))
		if err != nil {
			return err
		}
		ct.Time = parsedTime
		return nil
	case string:
		parsedTime, err := time.Parse("2006-01-02", t)
		if err != nil {
			return err
		}
		ct.Time = parsedTime
		return nil
	default:
		return fmt.Errorf("unsupported Scan value for CustomTime: %T", value)
	}
}
