package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

type ctxKey int

const (
	CtxKeyAuthUser ctxKey = iota
)

type LoginForm struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required"`
}

type UserSession struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash []byte
	Bio          string
	Image        *string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

type UserSessionToken struct {
	Token     string
	CreatedAt *time.Time
	User      *UserSession
}

type RegisterForm struct {
	Username        string `form:"username" validate:"required"`
	Email           string `form:"email" validate:"required,email"`
	Password        string `form:"password" validate:"required,min=8,max=72"`
	PasswordConfirm string `form:"password_confirm" validate:"required,eqfield=Password"`
}

type FlightStatus string

const (
	Scheduled FlightStatus = "scheduled"
	Active    FlightStatus = "active"
	Landed    FlightStatus = "landed"
	Cancelled FlightStatus = "cancelled"
	Incident  FlightStatus = "incident"
	Diverted  FlightStatus = "diverted"
)

type NavItem struct {
	Path  string
	Icon  templ.Component
	Label string
}

type TabItem struct {
	Path  string
	Icon  string
	Label string
}

type SidebarItem struct {
	Path       string
	Icon       templ.Component
	Label      string
	ActivePath string
	SubItems   []SidebarItem
}

type LayoutTempl struct {
	Title     string
	Nav       []NavItem
	ActiveNav string
	User      *UserSession
	Content   templ.Component
}

type SettingsPage struct {
	Updated bool
	Errors  []string
	User    UserSession
}

type LoginPage struct {
	Errors []string
}

type RegisterPage struct {
	Errors []string
	Values map[string]string
}

type Columns struct {
	Title string
}

type Airport struct {
	ID           string         `json:"id"`
	CityName     string         `json:"city_name"`
	GMT          string         `json:"gmt"`
	AirportID    int            `json:"airport_id,string,omitempty"`
	IataCode     string         `json:"iata_code"`
	CityIataCode string         `json:"city_iata_code"`
	IcaoCode     string         `json:"icao_code"`
	CountryISO2  string         ` json:"country_iso2"`
	GeonameID    string         `json:"geoname_id,omitempty"`
	Latitude     float64        `json:"latitude,string,omitempty"`
	Longitude    float64        `json:"longitude,string,omitempty"`
	AirportName  string         `json:"airport_name"`
	CountryName  string         ` json:"country_name"`
	PhoneNumber  sql.NullString ` json:"phone_number,omitempty"`
	Timezone     string         ` json:"timezone"`
	CreatedAt    CustomTime     `db:"created_at" json:"created_at"`
}

type City struct {
	ID           string     `json:"id"`
	GMT          string     `json:"gmt,omitempty"`
	CityID       int        `json:"city_id,string,omitempty"`
	IataCode     string     `json:"iata_code"`
	CountryISO2  string     `json:"country_iso2"`
	GeonameID    string     `json:"geoname_id,omitempty"`
	Latitude     float64    `json:"latitude,string,omitempty"`
	Longitude    float64    `json:"longitude,string,omitempty"`
	CityName     string     `json:"city_name"`
	Timezone     string     `json:"timezone"`
	Continent    string     `json:"continent"`
	CountryName  string     `json:"country_name"`
	CurrencyName string     `json:"currency_name"`
	PhonePrefix  string     `json:"phone_prefix"`
	CreatedAt    CustomTime `db:"created_at" json:"created_at"`
}

type Country struct {
	ID                string     `json:"id"`
	CountryName       string     `json:"country_name"`
	CountryISO2       string     `json:"country_iso2"`
	CountryIso3       string     `json:"country_iso3"`
	CountryIsoNumeric int        `json:"country_iso_numeric,string"`
	Population        int        `json:"population,string"`
	Capital           string     `json:"capital"`
	Continent         string     `json:"continent"`
	CurrencyName      string     `json:"currency_name"`
	CurrencyCode      string     `json:"currency_code"`
	FipsCode          string     `json:"fips_code"`
	PhonePrefix       string     `json:"phone_prefix"`
	Latitude          float64    `json:"latitude,string,omitempty"`
	Longitude         float64    `json:"longitude,string,omitempty"`
	CreatedAt         CustomTime `db:"created_at" json:"created_at"`
}

type Tax struct {
	ID          string     `json:"id"`
	TaxID       int        `json:"tax_id,string,omitempty"`
	TaxName     string     `json:"tax_name"`
	IataCode    string     `json:"iata_code"`
	AirlineName string     `json:"airline_name"`
	CountryName string     `json:"country_name"`
	CityName    string     `json:"city_name"`
	CreatedAt   CustomTime `db:"created_at" json:"created_at"`
}

type Aircraft struct {
	ID                     string      `json:"id"`
	AircraftName           string      `json:"aircraft_name"`
	IataType               string      `json:"iata_type"`
	AirplaneID             int         `json:"airplane_id,string"`
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
	CreatedAt              CustomTime  `db:"created_at" json:"created_at"`
}

type Airline struct {
	ID                   string      `json:"id"`
	FleetAverageAge      float64     `json:"fleet_average_age,string"`
	AirlineID            int         `json:"airline_id"`
	CallSign             string      `json:"callsign"`
	HubCode              string      `json:"hub_code"`
	IataCode             string      `json:"iata_code"`
	IcaoCode             string      `json:"icao_code"`
	CountryISO2          string      `json:"country_iso2"`
	DateFounded          int         `json:"date_founded,string"`
	IataPrefixAccounting int         `json:"iata_prefix_accounting,string"`
	AirlineName          string      `json:"airline_name"`
	CountryName          string      `json:"country_name"`
	FleetSize            int         `json:"fleet_size,string"`
	Status               string      `json:"status"`
	Type                 string      `json:"type"`
	CityName             string      `json:"city_name"`
	AirportName          string      `json:"airport_name"`
	Timezone             string      `json:"timezone"`
	Latitude             float64     `json:"latitude,string,omitempty"`
	Longitude            float64     `json:"longitude,string,omitempty"`
	PlaneAge             int         `json:"plane_age,string"`
	PlaneClass           interface{} `json:"plane_class"`
	ModelName            string      `json:"model_name"`
	PlaneOwner           interface{} `json:"plane_owner"`
	RegistrationDate     CustomTime  `json:"registration_date"`
	Continent            string      `json:"continent"`
	CreatedAt            CustomTime  `db:"created_at" json:"created_at"`
}

type Airplane struct {
	ID                     string      `json:"id" `
	AirlineName            string      `json:"airline_name"`
	IataType               string      `json:"iata_type"`
	AirplaneID             int         `json:"airplane_id,string"`
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
	CreatedAt              CustomTime  `db:"created_at" json:"created_at"`
}

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
		Scheduled       time.Time   `json:"scheduled"`
		Estimated       time.Time   `json:"estimated"`
		Actual          time.Time   `json:"actual"`
		EstimatedRunway time.Time   `json:"estimated_runway"`
		ActualRunway    time.Time   `json:"actual_runway"`
		CityCode        string      `json:"departure_city_code"`
		CountryCode     string      `json:"departure_country_code"`
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
		Scheduled       time.Time   `json:"scheduled"`
		Estimated       time.Time   `json:"estimated"`
		Actual          time.Time   `json:"actual"`
		EstimatedRunway time.Time   `json:"estimated_runway"`
		ActualRunway    time.Time   `json:"actual_runway"`
		CityCode        string      `json:"arrival_city_code"`
		CountryCode     string      `json:"arrival_country_code"`
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
		LiveAltitude        float32 `json:"altitude"`
		LiveDirection       float32 `json:"direction"`
		LiveSpeedHorizontal float32 `json:"speed_horizontal"`
		LiveSpeedVertical   float32 `json:"speed_vertical"`
		LiveIsGround        bool    `json:"is_ground"`
	} `json:"live,omitempty"`
	CreatedAt            CustomTime `json:"created_at"`
	DepartureAirportName string     `json:"departure_airport"`
	DepartureLatitude    string     `json:"departure_latitude"`
	DepartureLongitude   string     `json:"departure_longitude"`
	ArrivalAirportName   string     `json:"arrival_airport"`
	ArrivalLatitude      string     `json:"arrival_latitude"`
	ArrivalLongitude     string     `json:"arrival_longitude"`
}

func (a *Airport) GetPhoneNumber() string {
	if !a.PhoneNumber.Valid {
		return "Phone not available"
	}
	return a.PhoneNumber.String
}

type PageLayout struct {
	Table     AirportTable
	Sidebar   []SidebarItem
	ActiveNav string
}

type Table[T any, V any] struct {
	Column      []T
	Data        []V
	PrevPage    int
	NextPage    int
	Page        int
	LastPage    int
	SearchParam string
	QueryParam  string
}

type AirportTable struct {
	Column      []ColumnItems
	Airports    []Airport
	PrevPage    int
	NextPage    int
	Page        int
	LastPage    int
	SearchParam string
	OrderParam  string
	SortParam   string
}

type TaxTable struct {
	Column        []ColumnItems
	Tax           []Tax
	PrevPage      int
	NextPage      int
	Page          int
	LastPage      int
	FilterTax     string
	FilterAirline string
	FilterCountry string
	OrderParam    string
	SortParam     string
}

type AircraftTable struct {
	Column             []ColumnItems
	Aircraft           []Aircraft
	PrevPage           int
	NextPage           int
	Page               int
	LastPage           int
	FilterAircraft     string
	FilterTypeOfEngine string
	FilterModelCode    string
	FilterPlaneOwner   string
	OrderParam         string
	SortParam          string
}

// ColumnItems Need to change this to improve all the sorting on the tables.
type ColumnItems struct {
	Title     string
	Icon      templ.Component
	SortParam string
}
type AirlineTable struct {
	Column            []ColumnItems
	Airline           []Airline
	PrevPage          int
	NextPage          int
	Page              int
	LastPage          int
	FilterName        string
	FilterCallSign    string
	FilterHubCode     string
	FilterCountryName string
	OrderParam        string
	SortParam         string
}

type AirplaneTable struct {
	Column                   []ColumnItems
	Airplane                 []Airplane
	PrevPage                 int
	NextPage                 int
	Page                     int
	LastPage                 int
	FilterAirlineName        string
	FilterModelName          string
	FilterProductionLine     string
	FilterRegistrationNumber string
	OrderParam               string
	SortParam                string
}

type CityTable struct {
	Column      []ColumnItems
	City        []City
	PrevPage    int
	NextPage    int
	Page        int
	LastPage    int
	SearchParam string
	OrderParam  string
	SortParam   string
}

type CountryTable struct {
	Column      []ColumnItems
	Country     []Country
	PrevPage    int
	NextPage    int
	Page        int
	LastPage    int
	SearchParam string
	OrderParam  string
	SortParam   string
}

type FlightsTable struct {
	Column      []ColumnItems
	Flights     []LiveFlights
	PrevPage    int
	NextPage    int
	Page        int
	LastPage    int
	SearchParam string
	OrderParam  string
	SortParam   string
}

type CustomTime struct {
	time.Time
}

type Timeable interface {
	GetTime() time.Time
}

func (ct *CustomTime) GetTime() time.Time {
	return ct.Time
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	var dateStr string
	err := json.Unmarshal(data, &dateStr)
	if err != nil {
		log.Println("Error parsing date: ", err)
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
		log.Println("Error parsing date:", err)
		return err
	}

	ct.Time = t
	return nil
}

//

// Value Implement driver.Valuer interface.
func (ct CustomTime) Value() (driver.Value, error) {
	// Return the underlying time value as a string in RFC3339 format
	return ct.Time.Format(time.RFC3339), nil
}

// Scan Implement sql.Scanner interface.
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
		return errors.New("unsupported Scan value for CustomTime")
	}
}
