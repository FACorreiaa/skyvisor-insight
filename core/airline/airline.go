package airline

import (
	"context"

	"github.com/FACorreiaa/Aviation-tracker/controller/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryAirline struct {
	pgpool *pgxpool.Pool
}

func NewAirlines(
	pgpool *pgxpool.Pool,

) *RepositoryAirline {
	return &RepositoryAirline{
		pgpool: pgpool,
	}
}

func (r *RepositoryAirline) getAirlineData(ctx context.Context, query string,
	args ...interface{}) ([]models.Airline, error) {
	var al []models.Airline

	rows, err := r.pgpool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Airline
		err := rows.Scan(
			&a.ID, &a.AirlineName, &a.DateFounded, &a.FleetAverageAge,
			&a.FleetSize, &a.CallSign, &a.HubCode, &a.Status,
			&a.Type, &a.CountryName,
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

func (r *RepositoryAirline) GetAirlines(ctx context.Context, page,
	pageSize int, orderBy, sortBy, name string) ([]models.Airline, error) {
	offset := (page - 1) * pageSize
	query := `select al.id, al.airline_name, al.date_founded, al.fleet_average_age, al.fleet_size,
						al.callsign, al.hub_code, al.status, al.type, al.country_name
						from  airline al
						where al.airline_id != 0 AND TRIM(UPPER(al.airline_name)) != ''
						AND TRIM(UPPER(al.airline_name)) ILIKE TRIM(UPPER('%' || $1 || '%'))
						order by
						    CASE WHEN $2 = 'Airline Name' AND $3 = 'ASC' THEN al.airline_name::text END ASC,
						    CASE WHEN $2 = 'Airline Name' AND $3 = 'DESC' THEN al.airline_name::text END DESC,
						    CASE WHEN $2 = 'Date Founded' AND $3 = 'ASC' THEN al.date_founded::text END ASC,
						    CASE WHEN $2 = 'Date Founded' AND $3 = 'DESC' THEN al.date_founded::text END DESC,
						    CASE WHEN $2 = 'Fleet Average Size' AND $3 = 'ASC' THEN al.fleet_average_age::text END ASC,
						    CASE WHEN $2 = 'Fleet Average Size' AND $3 = 'DESC' THEN al.fleet_average_age::text END DESC,
						    CASE WHEN $2 = 'Fleet Size' AND $3 = 'ASC' THEN al.fleet_size::text END ASC,
						    CASE WHEN $2 = 'Fleet Size' AND $3 = 'DESC' THEN al.fleet_size::text END DESC,
						    CASE WHEN $2 = 'Call Sign' AND $3 = 'ASC' THEN al.callsign::text END ASC,
						    CASE WHEN $2 = 'Call Sign' AND $3 = 'DESC' THEN al.callsign::text END DESC,
						    CASE WHEN $2 = 'Hub Code' AND $3 = 'ASC' THEN al.hub_code::text END ASC,
						    CASE WHEN $2 = 'Hub Code' AND $3 = 'DESC' THEN al.hub_code::text END DESC,
						    CASE WHEN $2 = 'Status' AND $3 = 'ASC' THEN al.status::text END ASC,
						    CASE WHEN $2 = 'Status' AND $3 = 'DESC' THEN al.status::text END DESC,
						    CASE WHEN $2 = 'Type' AND $3 = 'ASC' THEN al.type::text END ASC,
						    CASE WHEN $2 = 'Type' AND $3 = 'DESC' THEN al.type::text END DESC,
							CASE WHEN $2 = 'Country Name' AND $3 = 'ASC' THEN al.country_name::text END ASC,
						    CASE WHEN $2 = 'Country Name' AND $3 = 'DESC' THEN al.country_name::text END DESC
						OFFSET $4 LIMIT $5`

	return r.getAirlineData(ctx, query, name, orderBy, sortBy, offset, pageSize)
}

func (r *RepositoryAirline) GetAirlineSum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT Count(id)
										FROM airline
										WHERE  TRIM(UPPER(airline_name)) != ''
`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RepositoryAirline) GetAirlinesLocations(ctx context.Context) ([]models.Airline, error) {
	var airline []models.Airline

	rows, err := r.pgpool.Query(ctx, `select al.airline_id, al.airline_name, al.date_founded, al.fleet_average_age,
       										al.fleet_size, al.callsign, al.hub_code, al.status, al.type, al.country_name,
       										ct.city_name, ap.airport_name, ap.timezone,
       										ct.latitude, ct.longitude
											from  airline al
											join airport ap on ap.airport_id = airline_id
											join city ct on ap.city_iata_code = ct.iata_code
											where al.airline_id != 0
											  and TRIM(UPPER(al.airline_name)) != ''
											  and ct.longitude IS NOT NULL
											  and ct.longitude IS NOT NULL
											order by al.airline_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Airline
		err := rows.Scan(
			&a.AirlineID, &a.AirlineName, &a.DateFounded, &a.FleetAverageAge,
			&a.FleetSize, &a.CallSign, &a.HubCode, &a.Status, &a.Type, &a.CountryName,
			&a.CityName, &a.AirportName, &a.Timezone, &a.Latitude, &a.Longitude,
		)

		if err != nil {
			return nil, err
		}
		airline = append(airline, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return airline, nil
}

func (r *RepositoryAirline) GetAirlineByID(ctx context.Context, id int) (models.AirlineDetails, error) {
	var al models.AirlineDetails
	query := `
		select DISTINCT ON (al.airline_name) al.fleet_average_age, al.airline_id, al.callsign, al.hub_code,
                                     al.iata_code, al.icao_code, al.country_iso2, al.date_founded,
                                     al.iata_prefix_accounting, al.airline_name, al.country_name,
                                     al.fleet_size, al.status, al.type, al.created_at,
                                     ap.model_name, ap.plane_owner, ap.plane_age, ap.registration_date,
                                     c.continent
		from airline al
		         right join airplane ap on al.iata_code = ap.airline_iata_code
		         JOIN country c ON al.country_iso2 = c.country_iso2
		WHERE al.airline_name IS NOT NULL
		AND al.airline_name != ''
		AND al.airline_id = $1
;`
	err := r.pgpool.QueryRow(ctx, query, id).Scan(
		&al.FleetAverageAge,
		&al.AirlineID,
		&al.CallSign,
		&al.HubCode,
		&al.IataCode,
		&al.IcaoCode,
		&al.CountryISO2,
		&al.DateFounded,
		&al.IataPrefixAccounting,
		&al.AirlineName,
		&al.CountryName,
		&al.FleetSize,
		&al.Status,
		&al.Type,
		&al.CreatedAt,
		&al.ModelName,
		&al.PlaneOwner,
		&al.PlaneAge,
		&al.RegistrationDate,
		&al.Continent,
	)

	if err != nil {
		return models.AirlineDetails{}, err
	}

	return al, nil
}
