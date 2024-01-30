package airport

import (
	"context"

	"github.com/FACorreiaa/Aviation-tracker/controller/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AirportRepository struct {
	pgpool *pgxpool.Pool
}

func NewAirports(
	pgpool *pgxpool.Pool,

) *AirportRepository {
	return &AirportRepository{
		pgpool: pgpool,
	}
}

func (r *AirportRepository) GetAirports(ctx context.Context, page, pageSize int) ([]models.Airport, error) {
	var airport []models.Airport

	offset := (page - 1) * pageSize
	rows, err := r.pgpool.Query(ctx, `SELECT id, gmt, airport_id, iata_code,
       										city_iata_code, icao_code, country_iso2,
       										geoname_id, latitude, longitude, airport_name,
       										country_name, phone_number, timezone,
       										created_at
       								FROM airport ORDER BY id
       								OFFSET $1 LIMIT $2`,
		offset, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Airport
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

func (r *AirportRepository) GetSum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT Count(id) FROM airport`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
