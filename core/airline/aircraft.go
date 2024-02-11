package airline

import (
	"context"
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
)

// aircraft

func (r *RepositoryAirline) getAircraftData(ctx context.Context, query string,
	args ...interface{}) ([]models.Aircraft, error) {
	var aircraft []models.Aircraft

	rows, err := r.pgpool.Query(ctx, query, args...)
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

func (r *RepositoryAirline) GetAircraft(ctx context.Context, page, pageSize int, name string,
	orderBy string, sortBy string) ([]models.Aircraft, error) {

	offset := (page - 1) * pageSize
	query := `SELECT ac.id, ac.aircraft_name, ap.model_name, ap.construction_number,
											ap.engines_count, ap.engines_type, ap.first_flight_date, ap.line_number,
											ap.model_code, ap.plane_age, ap.plane_class, ap.plane_owner, ap.plane_series, ap.plane_status
											FROM
																		aircraft ac
											JOIN airplane ap ON ac.plane_type_id = ap.airplane_id
											WHERE ac.plane_type_id != 0 AND TRIM(UPPER(ac.aircraft_name)) != ''
       										AND    Trim(Upper(aircraft_name))
											           ILIKE trim(upper('%' || $1 || '%'))
											ORDER BY
			    CASE WHEN $2 = 'Aircraft Name' AND $3 = 'ASC' THEN ac.aircraft_name::text END ASC,
			    CASE WHEN $2 = 'Aircraft Name' AND $3 = 'DESC' THEN ac.aircraft_name::text END DESC,
			    CASE WHEN $2 = 'Model Name' AND $3 = 'ASC' THEN ap.model_name::text END ASC,
			    CASE WHEN $2 = 'Model Name' AND $3 = 'DESC' THEN ap.model_name::text END DESC,
			    CASE WHEN $2 = 'Construction Number' AND $3 = 'ASC' THEN ap.construction_number::text END ASC,
			    CASE WHEN $2 = 'Construction Number' AND $3 = 'DESC' THEN ap.construction_number::text END DESC,
			    CASE WHEN $2 = 'Number of Engines' AND $3 = 'ASC' THEN ap.engines_count::text END ASC,
			    CASE WHEN $2 = 'Number of Engines' AND $3 = 'DESC' THEN ap.engines_count::text END DESC,
			    CASE WHEN $2 = 'Type of Engine' AND $3 = 'ASC' THEN ap.engines_type::text END ASC,
			    CASE WHEN $2 = 'Type of Engine' AND $3 = 'DESC' THEN ap.engines_type::text END DESC,
			    CASE WHEN $2 = 'Date of first flight' AND $3 = 'ASC' THEN ap.first_flight_date::text END ASC,
			    CASE WHEN $2 = 'Date of first flight' AND $3 = 'DESC' THEN ap.first_flight_date::text END DESC,
			    CASE WHEN $2 = 'Line Number' AND $3 = 'ASC' THEN ap.line_number::text END ASC,
			    CASE WHEN $2 = 'line Number' AND $3 = 'DESC' THEN ap.line_number::text END DESC,
			    CASE WHEN $2 = 'Model Code' AND $3 = 'ASC' THEN ap.model_name::text END ASC,
			    CASE WHEN $2 = 'Model Code' AND $3 = 'DESC' THEN ap.model_name::text END DESC,
			    CASE WHEN $2 = 'Plane Age' AND $3 = 'ASC' THEN ap.plane_age::text END ASC,
			    CASE WHEN $2 = 'Plane Age' AND $3 = 'DESC' THEN ap.plane_age::text END DESC,
			    CASE WHEN $2 = 'Plane Class' AND $3 = 'ASC' THEN ap.plane_class::text END ASC,
			    CASE WHEN $2 = 'Plane Class' AND $3 = 'DESC' THEN ap.plane_class::text END DESC,
			    CASE WHEN $2 = 'Plane Owner' AND $3 = 'ASC' THEN ap.plane_owner::text END ASC,
			    CASE WHEN $2 = 'Plane Owner' AND $3 = 'DESC' THEN ap.plane_owner::text END DESC,
			    CASE WHEN $2 = 'Plane Series' AND $3 = 'ASC' THEN ap.plane_series::text END ASC,
			    CASE WHEN $2 = 'Plane Series' AND $3 = 'DESC' THEN ap.plane_series::text END DESC
       										OFFSET $4 LIMIT $5`

	return r.getAircraftData(ctx, query, name, orderBy, sortBy, offset, pageSize)

}

func (r *RepositoryAirline) GetAircraftSum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT Count(ac.id)
										FROM aircraft ac
										JOIN airplane ap ON ac.plane_type_id = ap.airplane_id
										WHERE ac.plane_type_id != 0 AND TRIM(UPPER(ac.aircraft_name)) != ''
`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Airlines
