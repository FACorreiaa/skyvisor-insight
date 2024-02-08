package airport

import (
	"context"
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryAirport struct {
	pgpool *pgxpool.Pool
}

func NewAirports(
	pgpool *pgxpool.Pool,

) *RepositoryAirport {
	return &RepositoryAirport{
		pgpool: pgpool,
	}
}

func (r *RepositoryAirport) getAirportData(ctx context.Context, query string,
	args ...interface{}) ([]models.Airport, error) {
	var al []models.Airport

	rows, err := r.pgpool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Airport
		err := rows.Scan(
			&a.ID, &a.GMT, &a.AirportID, &a.IataCode,
			&a.CityIataCode, &a.IcaoCode, &a.CountryISO2,
			&a.GeonameID, &a.Latitude, &a.Longitude,
			&a.AirportName, &a.CountryName, &a.PhoneNumber,
			&a.Timezone, &a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		al = append(al, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return al, nil
}

func (r *RepositoryAirport) GetAirports(ctx context.Context,
	page, pageSize int, orderBy string) ([]models.Airport, error) {
	query := `SELECT id, gmt, airport_id, iata_code,
       										city_iata_code, icao_code, country_iso2,
       										geoname_id, latitude, longitude, airport_name as "Airport Name",
       										country_name, phone_number, timezone,
       										created_at
       								FROM airport
       								WHERE  airport_name IS NOT NULL AND TRIM(UPPER(airport_name)) != ''
       								ORDER BY CASE
							         WHEN $1 = 'Airport Name' then airport_name::text
							         WHEN $1 = 'Country Name' then country_name::text
       								 WHEN $1 = 'Phone Number' then phone_number::text
       								 WHEN $1 = 'Timezone' then timezone::text   
       								 WHEN $1 = 'GMT' then gmt::text
       								 WHEN $1 = 'Latitude' then latitude::text
       								 WHEN $1 = 'Longitude' then longitude::text
							         END DESC 
       								OFFSET $2 LIMIT $3`
	offset := (page - 1) * pageSize

	return r.getAirportData(ctx, query, orderBy, offset, pageSize)
}

func (r *RepositoryAirport) GetSum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT Count(id)
										FROM airport
										WHERE  airport_name IS NOT NULL AND TRIM(UPPER(airport_name)) != ''
`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RepositoryAirport) GetAirportsLocation(ctx context.Context) ([]models.Airport, error) {
	var airport []models.Airport

	rows, err := r.pgpool.Query(ctx, `SELECT a.id, a.gmt, a.latitude, a.longitude,
       											 a.airport_name, c.city_name,
												 a.country_name, a.phone_number, a.timezone
												 FROM airport a
												 INNER JOIN
    											 	City c ON a.city_iata_code = c.iata_code
												 WHERE
    												a.airport_name IS NOT NULL AND TRIM(UPPER(a.airport_name)) != ''
												 ORDER BY id`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a models.Airport
		err := rows.Scan(
			&a.ID, &a.GMT, &a.Latitude, &a.Longitude,
			&a.AirportName, &a.CityName, &a.CountryName, &a.PhoneNumber,
			&a.Timezone,
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

func (r *RepositoryAirport) GetAirportByName(ctx context.Context, name string,
	page, pageSize int, orderBy string) ([]models.Airport, error) {
	query := `SELECT id, gmt, airport_id, iata_code,
       										city_iata_code, icao_code, country_iso2,
       										geoname_id, latitude, longitude, airport_name as "Airport Name",
       										country_name, phone_number, timezone,
       										created_at
       								FROM airport
       								WHERE  TRIM(UPPER(airport_name)) ILIKE TRIM(UPPER('%' || $1 || '%'))
       								ORDER BY $2 DESC
       								OFFSET $3 LIMIT $4`
	offset := (page - 1) * pageSize

	return r.getAirportData(ctx, query, name, orderBy, offset, pageSize)
}

func (r *RepositoryAirport) GetAirportByID(ctx context.Context, id int) (models.Airport, error) {
	var ap models.Airport
	query := `
		select distinct on(airport_id)
		    ap.airport_id,
		    ap.airport_name,
		    ap.phone_number,
		    ap.country_name,
		    ap.timezone,
		    ap.gmt,
		    ap.geoname_id,
		    ap.created_at,
		    ap.latitude,
		    ap.longitude,
		    ct.city_name,
		    ct.timezone
		from airport ap
		join city ct on ct.iata_code = ap.iata_code
		where ap.airport_id = $1;
	`
	err := r.pgpool.QueryRow(ctx, query, id).Scan(
		&ap.ID,
		&ap.AirportName,
		&ap.PhoneNumber,
		&ap.CountryName,
		&ap.Timezone,
		&ap.GMT,
		&ap.GeonameID,
		&ap.CreatedAt,
		&ap.Latitude,
		&ap.Longitude,
		&ap.CityName,
		&ap.Timezone,
	)

	if err != nil {
		return models.Airport{}, err
	}

	return ap, nil
}
