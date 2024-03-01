package airport

import (
	"context"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
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
	page, pageSize int, orderBy string, sortBy string) ([]models.Airport, error) {
	query := `SELECT   id,
			         gmt,
			         airport_id,
			         iata_code,
			         city_iata_code,
			         icao_code,
			         country_iso2,
			         geoname_id,
			         latitude,
			         longitude,
			         airport_name AS "Airport Name",
			         country_name,
			         phone_number,
			         timezone,
			         created_at
			FROM     airport
			WHERE    airport_name IS NOT NULL
			AND      Trim(Upper(airport_name)) != ''
			ORDER BY
			         CASE
			                  WHEN $1 = 'Airport Name'
			                  AND      $2 = 'ASC' THEN airport_name::text
			         END ASC,
			         CASE
			                  WHEN $1 = 'Airport Name'
			                  AND      $2 = 'DESC' THEN airport_name::text
			         END DESC,
			         CASE
			                  WHEN $1 = 'Country Name'
			                  AND      $2 = 'ASC' THEN country_name::text
			         END ASC,
			         CASE
			                  WHEN $1 = 'Country Name'
			                  AND      $2 = 'DESC' THEN country_name::text
			         END DESC,
			         CASE
			                  WHEN $1 = 'Phone Number'
			                  AND      $2 = 'ASC' THEN phone_number::text
			         END ASC,
			         CASE
			                  WHEN $1 = 'Phone Number'
			                  AND      $2 = 'DESC' THEN phone_number::text
			         END DESC,
			         CASE
			                  WHEN $1 = 'Timezone'
			                  AND      $2 = 'ASC' THEN timezone::text
			         END ASC,
			         CASE
			                  WHEN $1 = 'Timezone'
			                  AND      $2 = 'DESC' THEN timezone::text
			         END DESC,
			         CASE
			                  WHEN $1 = 'GMT'
			                  AND      $2 = 'ASC' THEN timezone::text
			         END ASC,
			         CASE
			                  WHEN $1 = 'GMT'
			                  AND      $2 = 'DESC' THEN timezone::text
			         END DESC,
			         CASE
			                  WHEN $1 = 'Latitude'
			                  AND      $2 = 'ASC' THEN timezone::text
			         END ASC,
			         CASE
			                  WHEN $1 = 'Latitude'
			                  AND      $2 = 'DESC' THEN timezone::text
			         END DESC,
			         CASE
			                  WHEN $1 = 'Longitude'
			                  AND      $2 = 'ASC' THEN timezone::text
			         END ASC,
			         CASE
			                  WHEN $1 = 'Longitude'
			                  AND      $2 = 'DESC' THEN timezone::text
			         END DESC
			         offset $3 limit $4`
	offset := (page - 1) * pageSize

	return r.getAirportData(ctx, query, orderBy, sortBy, offset, pageSize)
}

func (r *RepositoryAirport) GetSum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT Count(id)
										FROM   airport
										WHERE  airport_name IS NOT NULL
										       AND Trim(Upper(airport_name)) != ''`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RepositoryAirport) GetAirportsLocation(ctx context.Context) ([]models.Airport, error) {
	var airport []models.Airport

	rows, err := r.pgpool.Query(ctx, `SELECT a.id,
											       a.gmt,
											       a.latitude,
											       a.longitude,
											       a.airport_name,
											       c.city_name,
											       a.country_name,
											       a.phone_number,
											       a.timezone
											FROM   airport a
											JOIN city c
											ON a.city_iata_code = c.iata_code
											WHERE  a.airport_name IS NOT NULL
											       AND Trim(Upper(a.airport_name)) != ''
											ORDER  BY id `)
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
	page, pageSize int, orderBy string, sortBy string) ([]models.Airport, error) {
	query := `SELECT   id,
				         gmt,
				         airport_id,
				         iata_code,
				         city_iata_code,
				         icao_code,
				         country_iso2,
				         geoname_id,
				         latitude,
				         longitude,
				         airport_name AS "Airport Name",
				         country_name,
				         phone_number,
				         timezone,
				         created_at
				FROM     airport
				WHERE    Trim(Upper(airport_name)) ilike trim(upper('%'
				                  || $1
				                  || '%'))
				ORDER BY
			         CASE
			                  WHEN $2 = 'Airport Name'
			                  AND      $3 = 'ASC' THEN airport_name::text
			         END ASC,
			         CASE
			                  WHEN $2 = 'Airport Name'
			                  AND      $3 = 'DESC' THEN airport_name::text
			         END DESC,
			         CASE
			                  WHEN $2 = 'Country Name'
			                  AND      $3 = 'ASC' THEN country_name::text
			         END ASC,
			         CASE
			                  WHEN $2 = 'Country Name'
			                  AND      $3 = 'DESC' THEN country_name::text
			         END DESC,
			         CASE
			                  WHEN $2 = 'Phone Number'
			                  AND      $3 = 'ASC' THEN phone_number::text
			         END ASC,
			         CASE
			                  WHEN $2 = 'Phone Number'
			                  AND      $3 = 'DESC' THEN phone_number::text
			         END DESC,
			         CASE
			                  WHEN $2 = 'Timezone'
			                  AND      $3 = 'ASC' THEN timezone::text
			         END ASC,
			         CASE
			                  WHEN $2 = 'Timezone'
			                  AND      $3 = 'DESC' THEN timezone::text
			         END DESC,
			         CASE
			                  WHEN $2 = 'GMT'
			                  AND      $3 = 'ASC' THEN timezone::text
			         END ASC,
			         CASE
			                  WHEN $2 = 'GMT'
			                  AND      $3 = 'DESC' THEN timezone::text
			         END DESC,
			         CASE
			                  WHEN $2 = 'Latitude'
			                  AND      $3 = 'ASC' THEN timezone::text
			         END ASC,
			         CASE
			                  WHEN $2 = 'Latitude'
			                  AND      $3 = 'DESC' THEN timezone::text
			         END DESC,
			         CASE
			                  WHEN $2 = 'Longitude'
			                  AND      $3 = 'ASC' THEN timezone::text
			         END ASC,
			         CASE
			                  WHEN $2 = 'Longitude'
			                  AND      $3 = 'DESC' THEN timezone::text
			         END DESC
				offset $4 limit $5`
	offset := (page - 1) * pageSize

	return r.getAirportData(ctx, query, name, orderBy, sortBy, offset, pageSize)
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
		&ap.AirportID,
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
