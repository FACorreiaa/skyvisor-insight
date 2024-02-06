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
			&a.FleetSize, &a.Callsign, &a.HubCode, &a.Status,
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

func (r *RepositoryAirline) GetAirlines(ctx context.Context, page, pageSize int) ([]models.Airline, error) {
	offset := (page - 1) * pageSize
	query := `select al.id, al.airline_name, al.date_founded, al.fleet_average_age, al.fleet_size,
						al.callsign, al.hub_code, al.status, al.type, al.country_name
						from  airline al
						where al.airline_id != 0 AND TRIM(UPPER(al.airline_name)) != ''
						order by al.airline_name
						OFFSET $1 LIMIT $2`

	return r.getAirlineData(ctx, query, offset, pageSize)
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
											left join airport ap on ap.airport_id = airline_id
											left join city ct on ap.city_iata_code = ct.iata_code
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
			&a.FleetSize, &a.Callsign, &a.HubCode, &a.Status, &a.Type, &a.CountryName,
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

func (r *RepositoryAirline) GetAirlineByName(ctx context.Context, param string, page,
	pageSize int) ([]models.Airline, error) {
	offset := (page - 1) * pageSize
	query := `select al.id, al.airline_name, al.date_founded, al.fleet_average_age, al.fleet_size,
			       al.callsign, al.hub_code, al.status, al.type, al.country_name
			from  airline al
			where al.airline_id != 0
			  AND TRIM(UPPER(al.airline_name)) != ''
			  AND TRIM(UPPER(al.airline_name)) ILIKE TRIM(UPPER('%' || $1 || '%'))
			order by al.airline_name
			OFFSET $2 LIMIT $3`

	return r.getAirlineData(ctx, query, param, offset, pageSize)
}
