package airport

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

// Structs

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
	CreatedAt    time.Time   `db:"created_at" json:"created_at"`
}

type Airports struct {
	pgpool *pgxpool.Pool
}

func NewAirports(
	pgpool *pgxpool.Pool,

) *Airports {
	return &Airports{
		pgpool: pgpool,
	}
}

func (r *Airports) GetAirports(ctx context.Context) ([]Airport, error) {
	var airport []Airport

	rows, err := r.pgpool.Query(ctx, `SELECT id, gmt, airport_id, iata_code,
       										city_iata_code, icao_code, country_iso2,
       										geoname_id, latitude, longitude, airport_name,
       										country_name, phone_number, timezone,
       										created_at
       								FROM airport ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a Airport
		err := rows.Scan(
			&a.ID, &a.GMT, &a.AirportId, &a.IataCode,
			&a.CityIataCode, &a.IcaoCode, &a.CountryISO2,
			&a.GeonameID, &a.Latitude, &a.Longitude,
			&a.AirportName, &a.CountryName, &a.PhoneNumber,
			&a.Timezone, &a.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		airport = append(airport, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return airport, nil
}
