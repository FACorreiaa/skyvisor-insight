package structs

// create an intermediate type & then convert to a concrete one
type Airport struct {
	ID           string      `json:"id"`
	GMT          string      `json:"gmt"`
	AirportId    int         `json:"airport_id,string,omitempty"`
	IataCode     string      `json:"iata_code"`
	CityIataCode string      `json:"city_iata_code"`
	IcaoCode     string      `json:"icao_code"`
	CountryISO2  string      ` json:"country_iso2"`
	GeonameID    string      `json:"geoname_id,omitempty"`
	Latitude     float64     `json:"latitude,string,omitempty"`
	Longitude    float64     `json:"longitude,string,omitempty"`
	AirportName  string      `json:"airport_name"`
	CountryName  string      ` json:"country_name"`
	PhoneNumber  interface{} ` json:"phone_number"`
	Timezone     string      ` json:"timezone"`
	CreatedAt    CustomTime  `db:"created_at" json:"created_at"`
}

type AirportApiData struct {
	Pagination Pagination `json:"pagination"`
	Data       []Airport  `json:"data"`
}

//create an intermediate type & then convert to a concrete one
// type Deez struct {
//   Data map[string]string `json:".data"`
// }

// func (d Deez) MyBeautifulWellFormattedStruct() (Airport, error) {...}
