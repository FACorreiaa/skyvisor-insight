package structs

// cities

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Count  int `json:"count"`
	Total  int `json:"total"`
}

type City struct {
	ID          string     `json:"id"`
	GMT         string     `json:"gmt,omitempty"`
	CityID      int        `json:"city_id,string,omitempty"`
	IataCode    string     `json:"iata_code"`
	CountryISO2 string     `json:"country_iso2"`
	GeonameID   string     `json:"geoname_id,omitempty"`
	Latitude    float64    `json:"latitude,string,omitempty"`
	Longitude   float64    `json:"longitude,string,omitempty"`
	CityName    string     `json:"city_name"`
	Timezone    string     `json:"timezone"`
	CreatedAt   CustomTime `json:"created_at"`
}

type CityApiData struct {
	Pagination Pagination `json:"pagination"`
	Data       []City     `json:"data"`
}

//countries

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
	CreatedAt         CustomTime `db:"created_at" json:"created_at"`
}

type CountryApiData struct {
	Pagination Pagination `json:"pagination"`
	Data       []Country  `json:"data"`
}

func ExtractCityID(c City) int {
	return c.CityID
}
