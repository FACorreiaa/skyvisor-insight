package airline

import (
	"context"
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AirlineRepository struct {
	pgpool *pgxpool.Pool
}

func NewAirports(
	pgpool *pgxpool.Pool,

) *AirlineRepository {
	return &AirlineRepository{
		pgpool: pgpool,
	}
}

// tax

func (r *AirlineRepository) GetTax(ctx context.Context, page, pageSize int) ([]models.Tax, error) {
	var tax []models.Tax

	offset := (page - 1) * pageSize
	rows, err := r.pgpool.Query(ctx, `SELECT t.id, t.tax_name, a.airline_name, a.country_name,  ct.city_name
											FROM tax t
													 INNER JOIN airline a ON t.iata_code = a.iata_code
													 INNER JOIN country c ON a.country_name = c.country_name
													 INNER JOIN city ct ON ct.country_iso2 = c.country_iso2
											WHERE t.tax_name IS NOT NULL AND t.tax_name != ''
       										OFFSET $1 LIMIT $2`,
		offset, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Tax
		err := rows.Scan(
			&t.ID, &t.TaxName, &t.AirlineName,
			&t.CountryName, &t.CityName,
		)

		if err != nil {
			return nil, err
		}
		tax = append(tax, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tax, nil
}

func (r *AirlineRepository) GetSum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT Count(id)
										FROM tax
										WHERE tax_name != ''
`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// aircraft

func (r *AirlineRepository) GetAircraft(ctx context.Context, page, pageSize int) ([]models.Aircraft, error) {
	var aircraft []models.Aircraft

	offset := (page - 1) * pageSize
	rows, err := r.pgpool.Query(ctx, `SELECT ac.id, ac.aircraft_name, ap.model_name, ap.construction_number,
											ap.engines_count, ap.engines_type, ap.first_flight_date, ap.line_number,
											ap.model_code, ap.plane_age, ap.plane_class, ap.plane_owner, ap.plane_series, ap.plane_status
											FROM
																		aircraft ac
											LEFT JOIN airplane ap ON ac.plane_type_id = ap.airplane_id
											WHERE ac.plane_type_id != 0 AND TRIM(UPPER(ac.aircraft_name)) != ''
											ORDER BY ac.aircraft_name
       										OFFSET $1 LIMIT $2`,
		offset, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Aircraft
		err := rows.Scan(
			&a.ID, &a.AircraftName, &a.ModelName, &a.ConstructionNumber,
			&a.EnginesCount, &a.EnginesType, &a.FirstFlightDate, &a.LineNumber,
			&a.ModelCode, &a.PlaneAge, &a.PlaneClass, &a.PlaneOwner, &a.PlaneSeries,
			&a.PlaneStatus,
		)

		if err != nil {
			return nil, err
		}
		aircraft = append(aircraft, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return aircraft, nil
}

func (r *AirlineRepository) GetAircraftSum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT Count(id)
										FROM aircraft
										WHERE  TRIM(UPPER(aircraft_name)) != ''
`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Airlines

func (r *AirlineRepository) GetAirlines(ctx context.Context, page, pageSize int) ([]models.Airline, error) {
	var aircraft []models.Airline

	offset := (page - 1) * pageSize
	rows, err := r.pgpool.Query(ctx, `select al.id, al.airline_name, al.date_founded, al.fleet_average_age, al.fleet_size,
											al.callsign, al.hub_code, al.status, al.type, al.country_name
											from  airline al
											where al.airline_id != 0 AND TRIM(UPPER(al.airline_name)) != ''
											order by al.airline_name
       										OFFSET $1 LIMIT $2`,
		offset, pageSize)
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
		aircraft = append(aircraft, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return aircraft, nil
}

func (r *AirlineRepository) GetAirlineSum(ctx context.Context) (int, error) {
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

// Airplane

func (r *AirlineRepository) GetAirplanes(ctx context.Context, page, pageSize int) ([]models.Airplane, error) {
	var airplane []models.Airplane

	offset := (page - 1) * pageSize
	rows, err := r.pgpool.Query(ctx, `select ap.id, ap.model_name, al.airline_name, ap.plane_series, ap.plane_owner,
       												ap.plane_class,
													ap.plane_age, ap.plane_status, ap.line_number, ap.first_flight_date,
													ap.engines_type, ap.engines_count, ap.construction_number,
													ap.production_line, ap.test_registration_number,
													ap.registration_date, ap.registration_number
													from airplane ap
													left join airline al on al.airline_id = ap.airplane_id
													where ap.model_name IS NOT NULL AND TRIM(UPPER(ap.model_name)) != ''
       										OFFSET $1 LIMIT $2`,
		offset, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Airplane
		err := rows.Scan(
			&a.ID, &a.ModelName, &a.AirlineName, &a.PlaneSeries, &a.PlaneOwner, &a.PlaneClass,
			&a.PlaneAge, &a.PlaneStatus, &a.LineNumber, &a.FirstFlightDate,
			&a.EnginesType, &a.EnginesCount, &a.ConstructionNumber, &a.ProductionLine, &a.TestRegistrationNumber,
			&a.RegistrationDate, &a.RegistrationNumber,
		)

		if err != nil {
			return nil, err
		}
		airplane = append(airplane, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return airplane, nil
}

func (r *AirlineRepository) GetAirplaneSum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT Count(id)
										FROM airplane
										WHERE  TRIM(UPPER(model_name)) != ''
`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
