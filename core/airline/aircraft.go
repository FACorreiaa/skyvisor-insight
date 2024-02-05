package airline

import (
	"context"
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
)

// aircraft

func (r *RepositoryAirline) GetAircraft(ctx context.Context, page, pageSize int) ([]models.Aircraft, error) {
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

func (r *RepositoryAirline) GetAircraftSum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT Count(ac.id)
										FROM aircraft ac
										LEFT JOIN airplane ap ON ac.plane_type_id = ap.airplane_id
										WHERE ac.plane_type_id != 0 AND TRIM(UPPER(ac.aircraft_name)) != ''
`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Airlines
