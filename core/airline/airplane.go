package airline

import (
	"context" //nolint:goimports //delete later
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
)

// Airplane

func (r *RepositoryAirline) getAirplaneData(ctx context.Context, query string,
	args ...interface{}) ([]models.Airplane, error) {
	var ap []models.Airplane

	rows, err := r.pgpool.Query(ctx, query, args...)
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
		ap = append(ap, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ap, nil
}

func (r *RepositoryAirline) GetAirplanes(ctx context.Context, page, pageSize int,
	orderBy string, sortBy string, name string) ([]models.Airplane, error) {

	offset := (page - 1) * pageSize
	query := `SELECT ap.id, ap.model_name, al.airline_name, ap.plane_series, ap.plane_owner,
       												ap.plane_class,
													ap.plane_age, ap.plane_status, ap.line_number, ap.first_flight_date,
													ap.engines_type, ap.engines_count, ap.construction_number,
													ap.production_line, ap.test_registration_number,
													ap.registration_date, ap.registration_number
				FROM airplane ap
				JOIN airline al on al.airline_id = ap.airplane_id
				WHERE ap.model_name IS NOT NULL AND TRIM(UPPER(ap.model_name)) != ''
       	AND    Trim(Upper(model_name))
		           ILIKE trim(upper('%' || $1 || '%'))
		ORDER BY
	    CASE WHEN $2 = 'Model Name'
	                  AND $3 = 'ASC' THEN ap.model_name::text END ASC,
	    CASE WHEN $2 = 'Model Name'
	                  AND $3 = 'DESC' THEN ap.model_name::text END DESC,
	    CASE WHEN $2 = 'Airline Name'
	                  AND $3 = 'ASC' THEN al.airline_name::text END ASC,
	    CASE WHEN $2 = 'Airline Name'
	                  AND $3 = 'DESC' THEN al.airline_name::text END DESC,
	    CASE WHEN $2 = 'Plane Series'
	                  AND $3 = 'ASC' THEN ap.plane_series::text END ASC,
	    CASE WHEN $2 = 'Plane Series'
	                  AND $3 = 'DESC' THEN ap.plane_series::text END DESC,
		CASE WHEN $2 = 'Plane Owner'
	                  AND $3 = 'ASC' THEN ap.plane_owner::text END ASC,
	    CASE WHEN $2 = 'Plane Owner'
	                  AND $3 = 'DESC' THEN ap.plane_owner::text END DESC,
		CASE WHEN $2 = 'Plane Class'
	                  AND $3 = 'ASC' THEN ap.plane_class::text END ASC,
	    CASE WHEN $2 = 'Plane Class'
	                  AND $3 = 'DESC' THEN ap.plane_class::text END DESC,
		CASE WHEN $2 = 'Plane Age'
	                  AND $3 = 'ASC' THEN ap.plane_age::text END ASC,
	    CASE WHEN $2 = 'Plane Age'
	                  AND $3 = 'DESC' THEN ap.plane_age::text END DESC,
		CASE WHEN $2 = 'Plane Status'
	                  AND $3 = 'ASC' THEN ap.plane_status::text END ASC,
	    CASE WHEN $2 = 'Plane Status'
	                  AND $3 = 'DESC' THEN ap.plane_status::text END DESC,
		CASE WHEN $2 = 'Line Number'
	                  AND $3 = 'ASC' THEN ap.line_number::text END ASC,
	    CASE WHEN $2 = 'Line Number'
	                  AND $3 = 'DESC' THEN ap.line_number::text END DESC,
		CASE WHEN $2 = 'First Flight Date'
	                  AND $3 = 'ASC' THEN ap.first_flight_date::text END ASC,
	    CASE WHEN $2 = 'First Flight Date'
	                  AND $3 = 'DESC' THEN ap.first_flight_date::text END DESC,
		CASE WHEN $2 = 'Engine Type'
	                  AND $3 = 'ASC' THEN ap.engines_type::text END ASC,
	    CASE WHEN $2 = 'Engine Type'
	                  AND $3 = 'DESC' THEN ap.engines_type::text END DESC,
		CASE WHEN $2 = 'Engine Count'
	                  AND $3 = 'ASC' THEN ap.engines_count::text END ASC,
	    CASE WHEN $2 = 'Engine Count'
	                  AND $3 = 'DESC' THEN ap.engines_count::text END DESC,
		CASE WHEN $2 = 'Construction Number'
	                  AND $3 = 'ASC' THEN ap.construction_number::text END ASC,
	    CASE WHEN $2 = 'Construction Number'
	                  AND $3 = 'DESC' THEN ap.construction_number::text END DESC,
		CASE WHEN $2 = 'Production Line'
	                  AND $3 = 'ASC' THEN ap.production_line::text END ASC,
	    CASE WHEN $2 = 'Production Line'
	                  AND $3 = 'DESC' THEN ap.production_line::text END DESC,
		CASE WHEN $2 = 'Test Registration Date'
	                  AND $3 = 'ASC' THEN ap.test_registration_number::text END ASC,
	    CASE WHEN $2 = 'Test Registration Date'
	                  AND $3 = 'DESC' THEN ap.test_registration_number::text END DESC,
		CASE WHEN $2 = 'Registration Date'
	                  AND $3 = 'ASC' THEN ap.registration_date::text END ASC,
	    CASE WHEN $2 = 'Registration Date'
	                  AND $3 = 'DESC' THEN ap.registration_date::text END DESC,
		CASE WHEN $2 = 'Registration Number'
	                  AND $3 = 'ASC' THEN ap.registration_number::text END ASC,
	    CASE WHEN $2 = 'Registration Number'
	                  AND $3 = 'DESC' THEN ap.registration_number::text END DESC
       	OFFSET $4 LIMIT $5`

	return r.getAirplaneData(ctx, query, name, orderBy, sortBy, offset, pageSize)

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
