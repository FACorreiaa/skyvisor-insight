package repository

import (
	"context"
	"log"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
)

func HandleError(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
}

func (r *AirlineRepository) getAirlineData(ctx context.Context, query string,
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

func (r *AirlineRepository) GetAirlines(ctx context.Context, page,
	pageSize int, orderBy, sortBy, name, callSign, hubCode,
	countryName string) ([]models.Airline, error) {
	offset := (page - 1) * pageSize
	query := `select al.id, al.airline_name, al.date_founded, al.fleet_average_age, al.fleet_size,
						al.callsign, al.hub_code, al.status, al.type, al.country_name
						from  airline al
						where al.airline_id != 0 AND TRIM(UPPER(al.airline_name)) != ''
						AND TRIM(UPPER(al.airline_name)) ILIKE TRIM(UPPER('%' || $1 || '%'))
						AND TRIM(UPPER(al.callsign)) ILIKE TRIM(UPPER('%' || $6 || '%'))
						AND TRIM(UPPER(al.hub_code)) ILIKE TRIM(UPPER('%' || $7 || '%'))
						AND TRIM(UPPER(al.country_name)) ILIKE TRIM(UPPER('%' || $8 || '%'))
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

	return r.getAirlineData(ctx, query, name, orderBy, sortBy, offset,
		pageSize, callSign, hubCode, countryName)
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

func (r *AirlineRepository) GetAirlinesLocations(ctx context.Context) ([]models.Airline, error) {
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

func (r *AirlineRepository) GetAirlineByName(ctx context.Context, airlineName string) (models.Airline, error) {
	var al models.Airline
	query := `
		select DISTINCT ON (al.airline_name) al.fleet_average_age, al.airline_id, al.callsign, al.hub_code,
                                     al.iata_code, al.icao_code, al.country_iso2, al.date_founded,
                                     al.iata_prefix_accounting, al.airline_name, al.country_name,
                                     al.fleet_size, al.status, al.type, al.created_at,
                                     ap.model_name, ap.plane_owner, ap.plane_age, ap.registration_date,
                                     c.continent, ct.latitude, ct.longitude
		FROM airline al
		RIGHT JOIN airplane
		    ap ON al.iata_code = ap.airline_iata_code
		JOIN country
		    c ON al.country_iso2 = c.country_iso2
		join airport apt on apt.airport_id = airline_id
		join city ct on apt.city_iata_code = ct.iata_code

		WHERE
		    al.airline_name IS NOT NULL
		AND
		    al.airline_name != ''
		AND
		    Trim(Upper(airline_name)) ilike trim(upper('%'
				                  || $1
				                  || '%'))`
	err := r.pgpool.QueryRow(ctx, query, airlineName).Scan(
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
		&al.Latitude,
		&al.Longitude,
	)

	if err != nil {
		HandleError(err, "Error scanning airlines")
		return models.Airline{}, err
	}
	return al, nil
}

// Aircraft

func (r *AirlineRepository) getAircraftData(ctx context.Context, query string,
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

func (r *AirlineRepository) GetAircraft(ctx context.Context, page, pageSize int, aircraftName,
	orderBy, sortBy, typeEngine, modelCode, planeOwner string) ([]models.Aircraft, error) {

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
											AND    Trim(Upper(ap.engines_type))
											           ILIKE trim(upper('%' || $6 || '%'))
											AND    Trim(Upper(ap.model_code))
											           ILIKE trim(upper('%' || $7 || '%'))
											AND    Trim(Upper(ap.plane_owner))
											           ILIKE trim(upper('%' || $8 || '%'))
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

	return r.getAircraftData(ctx, query, aircraftName, orderBy, sortBy, offset,
		pageSize, typeEngine, modelCode, planeOwner)

}

func (r *AirlineRepository) GetAircraftSum(ctx context.Context) (int, error) {
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

// Airplane

func (r *AirlineRepository) getAirplaneData(ctx context.Context, query string,
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

func (r *AirlineRepository) GetAirplanes(ctx context.Context, page, pageSize int,
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

func (r *AirlineRepository) GetAirplaneSum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT Count(ap.id)
	FROM airplane ap
	JOIN airline al on al.airline_id = ap.airplane_id
										where ap.model_name IS NOT NULL AND TRIM(UPPER(ap.model_name)) != ''
`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// tax

func (r *AirlineRepository) getTaxData(ctx context.Context, query string,
	args ...interface{}) ([]models.Tax, error) {
	var tax []models.Tax

	rows, err := r.pgpool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Tax
		err := rows.Scan(
			&t.ID, &t.TaxName, &t.AirlineName,
			&t.CountryName,
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

func (r *AirlineRepository) GetTax(ctx context.Context, page, pageSize int,
	orderBy, sortBy, taxName, countryName, airlineName string) ([]models.Tax, error) {
	query := `SELECT
    										t.id, t.tax_name, a.airline_name, a.country_name
											FROM tax t
											         JOIN airline a ON t.iata_code = a.iata_code
											         JOIN country c ON a.country_name = c.country_name
											WHERE t.tax_name IS NOT NULL
											AND t.tax_name != ''
											AND    Trim(Upper(t.tax_name))
											           ILIKE trim(upper('%' || $1 || '%'))
											AND    Trim(Upper(a.country_name))
											           ILIKE trim(upper('%' || $6 || '%'))
											AND    Trim(Upper(a.airline_name))
											           ILIKE trim(upper('%' || $7 || '%'))
											ORDER BY
			    CASE WHEN $2 = 'Tax Name' AND $3 = 'ASC' THEN t.tax_name::text END ASC,
			    CASE WHEN $2 = 'Tax Name' AND $3 = 'DESC' THEN t.tax_name::text END DESC,
			    CASE WHEN $2 = 'Airline Name' AND $3 = 'ASC' THEN a.airline_name::text END ASC,
			    CASE WHEN $2 = 'Airline Name' AND $3 = 'DESC' THEN a.airline_name::text END DESC,
			    CASE WHEN $2 = 'Country Name' AND $3 = 'ASC' THEN a.country_name::text END ASC,
			    CASE WHEN $2 = 'Country Name' AND $3 = 'DESC' THEN a.country_name::text END DESC
       										OFFSET $4 LIMIT $5`

	offset := (page - 1) * pageSize
	return r.getTaxData(ctx, query, taxName, orderBy, sortBy, offset, pageSize, countryName, airlineName)
}

func (r *AirlineRepository) GetTaxSum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT Count(DISTINCT t.id)
										FROM tax t
										JOIN airline a ON t.iata_code = a.iata_code
										JOIN country c ON a.country_name = c.country_name
										WHERE t.tax_name IS NOT NULL AND t.tax_name != ''
`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
