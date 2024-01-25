package models

import (
	"database/sql"
	"time"

	"github.com/FACorreiaa/go-ollama/core/account"
	"github.com/a-h/templ"
)

type NavItem struct {
	Path  string
	Icon  string
	Label string
}

type LayoutTempl struct {
	Title     string
	Nav       []NavItem
	ActiveNav string
	User      *account.User
	Content   templ.Component
}

type SettingsPage struct {
	Updated bool
	Errors  []string
	User    account.User
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
	GMT          string         `json:"gmt"`
	AirportId    int            `json:"airport_id,string,omitempty"`
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
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`
}

func (a *Airport) GetPhoneNumber() string {
	if !a.PhoneNumber.Valid {
		return "Phone not available"

	}
	return a.PhoneNumber.String
}

type AirportTable struct {
	Column   []string
	Airports []Airport
	PrevPage int
	NextPage int
	Page     int
}
