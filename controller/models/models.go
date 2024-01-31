package models

import (
	"database/sql"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/core/account"
	"github.com/a-h/templ"
)

type NavItem struct {
	Path  string
	Icon  string
	Label string
}

type SidebarItem struct {
	Path       string
	Icon       templ.Component
	Label      string
	ActivePath string
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
	CityName     string         `json:"city_name"`
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

type City struct {
	ID          string  `json:"id"`
	GMT         string  `json:"gmt,omitempty"`
	CityID      int     `json:"city_id,string,omitempty"`
	IataCode    string  `json:"iata_code"`
	CountryISO2 string  `json:"country_iso2"`
	GeonameID   string  `json:"geoname_id,omitempty"`
	Latitude    float64 `json:"latitude,string,omitempty"`
	Longitude   float64 `json:"longitude,string,omitempty"`
	CityName    string  `json:"city_name"`
	Timezone    string  `json:"timezone"`
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

type AirportTable struct {
	Column   []string
	Airports []Airport
	PrevPage int
	NextPage int
	Page     int
	LastPage int
}
