package airline

import (
	"context" //nolint:goimports //delete later
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
)

// Airplane

func (r *RepositoryAirline) GetAirplanes(ctx context.Context, page, pageSize int) ([]models.Airplane, error) {
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

func (r *RepositoryAirline) GetAirplaneSum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT Count(ap.id)
										FROM airplane ap
										left join airline al on al.airline_id = ap.airplane_id
										where ap.model_name IS NOT NULL AND TRIM(UPPER(ap.model_name)) != ''
`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
