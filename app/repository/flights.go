package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FlightsRepository struct {
	pgpool *pgxpool.Pool
}

func NewFlightsRepository(db *pgxpool.Pool) *FlightsRepository {
	return &FlightsRepository{pgpool: db}
}

const millis = 1000
const minutes = 60

func (r *FlightsRepository) getFlightsData(ctx context.Context, query string,
	args ...interface{}) ([]models.LiveFlights, error) {
	var lf []models.LiveFlights

	rows, err := r.pgpool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var f models.LiveFlights

		err = rows.Scan(
			&f.Flight.Number,
			&f.FlightDate,
			&f.FlightStatus,
			&f.Live.LiveUpdated,
			&f.ID,
			&f.Arrival.Actual,
			&f.Arrival.ActualRunway,
			&f.Arrival.Airport,
			&f.Arrival.Baggage,
			&f.Arrival.Delay,
			&f.Arrival.Estimated,
			&f.Arrival.Terminal,
			&f.Arrival.Gate,
			&f.Arrival.Timezone,
			&f.Departure.Scheduled,
			&f.Departure.EstimatedRunway,
			&f.Departure.Timezone,
			&f.Departure.Terminal,
			&f.Departure.Gate,
			&f.Departure.Actual,
			&f.Departure.ActualRunway,
			&f.Departure.Airport,
			&f.Departure.Estimated,
			&f.Departure.Delay,
			&f.Airline.Name,
		)
		if err != nil {
			return nil, err
		}

		// if f.Arrival.Delay != nil {
		//	minutes := *f.Arrival.Delay / (millis * minutes)
		//	f.Arrival.Delay = &minutes
		// }

		lf = append(lf, f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lf, nil
}

func (r *FlightsRepository) getFlightsLocationsData(ctx context.Context, query string,
	args ...interface{}) ([]models.LiveFlights, error) {
	var lf []models.LiveFlights

	rows, err := r.pgpool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var f models.LiveFlights

		err = rows.Scan(
			&f.Flight.Number,
			&f.Airline.Name,
			&f.DepartureAirportName,
			&f.DepartureLatitude,
			&f.DepartureLongitude,
			&f.ArrivalAirportName,
			&f.ArrivalLatitude,
			&f.ArrivalLongitude,
			&f.FlightDate,
			&f.FlightStatus,
			&f.Live.LiveUpdated,
			&f.ID,
			&f.Arrival.Actual,
			&f.Arrival.ActualRunway,
			&f.Arrival.Airport,
			&f.Arrival.Baggage,
			&f.Arrival.Delay,
			&f.Arrival.Estimated,
			&f.Arrival.Terminal,
			&f.Arrival.Gate,
			&f.Arrival.Timezone,
			&f.Departure.Scheduled,
			&f.Departure.EstimatedRunway,
			&f.Departure.Timezone,
			&f.Departure.Terminal,
			&f.Departure.Gate,
			&f.Departure.Actual,
			&f.Departure.ActualRunway,
			&f.Departure.Airport,
			&f.Departure.Estimated,
			&f.Departure.Delay,
			&f.Departure.CityCode,
			&f.Departure.CountryCode,
			&f.Arrival.CityCode,
			&f.Arrival.CityCode)
		if err != nil {
			return nil, err
		}

		if f.Arrival.Delay != nil {
			minutes := *f.Arrival.Delay / (millis * minutes)
			f.Arrival.Delay = &minutes
		}

		lf = append(lf, f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lf, nil
}

func (r *FlightsRepository) GetAllFlights(ctx context.Context,
	page, pageSize int, orderBy, sortBy string) ([]models.LiveFlights, error) {
	query := `SELECT DISTINCT ON (f.flight_number)
							       f.flight_number,
							       f.flight_date,
							       f.flight_status,
							       f.live_updated,
							       f.id,
							       f.arrival_actual,
							       f.arrival_actual_runway,
							       f.arrival_airport,
							    	COALESCE(f.arrival_baggage, 'N/A') as arrival_baggage,
							       COALESCE(FLOOR(f.arrival_delay / (1000 * 60)), 0) as arrival_delay,
									f.arrival_estimated,
    								COALESCE(f.arrival_terminal, 'N/A') as arrival_terminal,
							        COALESCE(f.arrival_gate, 'N/A') as arrival_gate,
							        f.arrival_timezone,
							        f.arrival_estimated,
							        f.departure_estimated_runway,
							       f.departure_timezone,
							       f.departure_terminal,
							       COALESCE(f.departure_gate, 'N/A') as departure_gate,
							       f.departure_actual,
							       f.departure_actual_runway,
							       f.departure_airport,
							       f.departure_estimated,
							       COALESCE(FLOOR(f.departure_delay / (1000 * 60)), 0) as departure_delay,
							       COALESCE(f.codeshared_airline_name, 'N/A') as airline_name

							from flights f
							WHERE Trim(Upper(f.flight_number)) != '' AND
							Trim(Upper(f.codeshared_airline_name)) != '' AND
							Trim(Upper(f.codeshared_airline_name)) != 'N/A'
-- 							WHERE	Trim(Upper(f.flight_number)) ILIKE trim(upper('%' || $5 || '%'))
-- 							AND Trim(Upper(f.airline_name)) ILIKE trim(upper('%' || $6 || '%'))
-- 							AND Trim(Upper(f.flight_status)) ILIKE trim(upper('%' || $7 || '%'))

							ORDER BY
							    f.flight_number,
							CASE
									WHEN $3 = 'Flight Number'
									AND      $4 = 'ASC' THEN flight_number::text
							END ASC,
							CASE
									WHEN $3 = 'Flight Number'
									AND      $4 = 'DESC' THEN flight_number::text
							END DESC,
							CASE
									WHEN $3 = 'Flight Status'
									AND      $4 = 'ASC' THEN flight_status::text
							END ASC,
							CASE
									WHEN $3 = 'Flight Status'
									AND      $4 = 'DESC' THEN flight_status::text
							END DESC,
							CASE
									WHEN $3 = 'Flight Date'
									AND      $4 = 'ASC' THEN flight_date::text
							END ASC,
							CASE
									WHEN $3 = 'Flight Date'
									AND      $4 = 'DESC' THEN flight_date::text
							END DESC,
							departure_estimated
							offset $1 limit $2`

	offset := (page - 1) * pageSize

	return r.getFlightsData(ctx, query, offset, pageSize, orderBy, sortBy)
}

func (r *FlightsRepository) GetAllFlightsSum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT Count(DISTINCT flight_number)
										FROM   flights`)

	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *FlightsRepository) GetFlightByID(ctx context.Context, flightNumber string) (models.LiveFlights, error) {
	var f models.LiveFlights

	query := `
				SELECT
			    f.flight_number,
			    COALESCE(al.airline_name, 'N/A') as airline_name,
			    COALESCE(ad.airport_name, 'N/A') AS departure_airport,
			    COALESCE(ad.latitude, 0.0) AS departure_latitude,
			    COALESCE(ad.longitude, 0.0) AS departure_longitude,
			    COALESCE(aa.airport_name, 'N/A') AS arrival_airport,
			    COALESCE(aa.latitude, 0.0) AS arrival_latitude,
			    COALESCE(aa.longitude, 0.0) AS arrival_longitude,
			    f.flight_date,
			    f.flight_status,
			    f.live_updated,
			    f.id,
			    f.arrival_actual,
				arrival_actual_runway,
			    f.arrival_airport,
			    COALESCE(f.arrival_baggage, 'N/A'),
			    COALESCE(FLOOR(f.arrival_delay / (1000 * 60)), 0) as arrival_delay,
				f.arrival_estimated,
			    COALESCE(f.arrival_terminal, 'N/A'),
			    COALESCE(f.arrival_gate, 'N/A'),
			    f.arrival_timezone,
			    f.departure_scheduled,
			    f.departure_estimated_runway,
			    f.departure_timezone,
			    f.departure_terminal,
			    COALESCE(f.departure_gate, 'N/A'),
			    f.departure_actual,
			    f.departure_actual_runway,
			    f.departure_airport,
			    f.departure_estimated,
			    COALESCE(FLOOR(f.departure_delay / (1000 * 60)), 0) as departure_delay,
				f.departure_icao,
				f.arrival_icao,
				COALESCE(ad.city_iata_code, 'N/A') AS departure_city_code,
			    COALESCE(ad.country_iso2 , 'N/A') AS departure_country_code,
			    COALESCE(aa.city_iata_code, 'N/A') AS arrival_city_code,
			    COALESCE(aa.country_iso2, 'N/A') AS arrival_country_code
			FROM
			    flights f
			        LEFT JOIN
			    airport ad ON f.departure_iata = ad.iata_code
			        LEFT JOIN
			    airport aa ON f.arrival_iata = aa.iata_code
			        LEFT JOIN
			    airline al on f.codeshared_airline_iata = al.iata_code
			WHERE
			    ad.latitude IS NOT NULL
			  AND ad.longitude IS NOT NULL AND
				Trim(Upper(f.flight_number)) != '' AND
							Trim(Upper(f.codeshared_airline_name)) != '' AND
							Trim(Upper(f.codeshared_airline_name)) != 'N/A'
			  AND flight_number = $1;`

	err := r.pgpool.QueryRow(ctx, query, flightNumber).Scan(
		&f.Flight.Number,
		&f.Airline.Name,
		&f.DepartureAirportName,
		&f.DepartureLatitude,
		&f.DepartureLongitude,
		&f.ArrivalAirportName,
		&f.ArrivalLatitude,
		&f.ArrivalLongitude,
		&f.FlightDate,
		&f.FlightStatus,
		&f.Live.LiveUpdated,
		&f.ID,
		&f.Arrival.Actual,
		&f.Arrival.ActualRunway,
		&f.Arrival.Airport,
		&f.Arrival.Baggage,
		&f.Arrival.Delay,
		&f.Arrival.Estimated,
		&f.Arrival.Terminal,
		&f.Arrival.Gate,
		&f.Arrival.Timezone,
		&f.Departure.Scheduled,
		&f.Departure.EstimatedRunway,
		&f.Departure.Timezone,
		&f.Departure.Terminal,
		&f.Departure.Gate,
		&f.Departure.Actual,
		&f.Departure.ActualRunway,
		&f.Departure.Airport,
		&f.Departure.Estimated,
		&f.Departure.Delay,
		&f.Departure.Icao,
		&f.Arrival.Icao,
		&f.Departure.CityCode,
		&f.Departure.CountryCode,
		&f.Arrival.CityCode,
		&f.Arrival.CountryCode,
	)
	if err != nil {
		return models.LiveFlights{}, err
	}

	return f, nil
}

func (r *FlightsRepository) GetAllFlightsLocation(ctx context.Context) ([]models.LiveFlights, error) {
	query := `
				SELECT
				DISTINCT ON (f.flight_number)
			    f.flight_number,
			    COALESCE(al.airline_name, 'N/A') as airline_name,
			    COALESCE(ad.airport_name, 'N/A') AS departure_airport,
			    COALESCE(ad.latitude, 0.0) AS departure_latitude,
			    COALESCE(ad.longitude, 0.0) AS departure_longitude,
			    COALESCE(aa.airport_name, 'N/A') AS arrival_airport,
			    COALESCE(aa.latitude, 0.0) AS arrival_latitude,
			    COALESCE(aa.longitude, 0.0) AS arrival_longitude,
			    f.flight_date,
			    f.flight_status,
			    f.live_updated,
			    f.id,
			    f.arrival_actual,
			    f.arrival_actual_runway,
			    f.arrival_airport,
			    COALESCE(f.arrival_baggage, 'N/A'),
			    COALESCE(FLOOR(f.arrival_delay / (1000 * 60)), 0) as arrival_delay,
				f.arrival_estimated,
			    COALESCE(f.arrival_terminal, 'N/A'),
			    COALESCE(f.arrival_gate, 'N/A'),
			    f.arrival_timezone,
			    f.departure_scheduled,
			    departure_estimated_runway,
			    f.departure_timezone,
			    f.departure_terminal,
			    COALESCE(f.departure_gate, 'N/A'),
			    f.departure_actual,
			    f.departure_actual_runway,
				f.departure_airport,
				f.departure_estimated,
			    COALESCE(FLOOR(f.departure_delay / (1000 * 60)), 0) as departure_delay,
				COALESCE(ad.city_iata_code, 'N/A') AS departure_city_code,
			    COALESCE(ad.country_iso2 , 'N/A') AS departure_country_code,
			    COALESCE(aa.city_iata_code, 'N/A') AS arrival_city_code,
			    COALESCE(aa.country_iso2, 'N/A') AS arrival_country_code
			FROM
			    flights f
			        LEFT JOIN
			    airport ad ON f.departure_iata = ad.iata_code
			        LEFT JOIN
			    airport aa ON f.arrival_iata = aa.iata_code
			        LEFT JOIN
			    airline al on f.airline_iata = al.iata_code
			WHERE Trim(Upper(f.flight_number)) != '' AND
							Trim(Upper(f.codeshared_airline_name)) != '' AND
							Trim(Upper(f.codeshared_airline_name)) != 'N/A' AND
			    ad.latitude IS NOT NULL
			  AND ad.longitude IS NOT NULL
			  AND ((date_trunc('day', TO_DATE(f.flight_date, 'YYYY-MM-DD')) = CURRENT_DATE))
			  ORDER BY f.flight_number`

	return r.getFlightsLocationsData(ctx, query)
}

func (r *FlightsRepository) GetAllFlightsByStatus(ctx context.Context,
	page, pageSize int, orderBy, sortBy, flightNumber, flightStatus string) ([]models.LiveFlights, error) {
	query := `select
							       f.flight_number,
							       f.flight_date,
							       f.flight_status,
							       f.live_updated,
							       f.id,
							       f.arrival_actual,
							       f.arrival_actual_runway,
							       f.arrival_airport,
							    	COALESCE(f.arrival_baggage, 'N/A') as arrival_baggage,
							       COALESCE(FLOOR(f.arrival_delay / (1000 * 60)), 0) as arrival_delay,
							       f.arrival_estimated,
							        COALESCE(f.arrival_terminal, 'N/A') as arrival_terminal,
							        COALESCE(f.arrival_gate, 'N/A') as arrival_gate,
							        f.arrival_timezone,
							        f.departure_scheduled,
							        f.departure_estimated_runway,
							       f.departure_timezone,
							       f.departure_terminal,
							       COALESCE(f.departure_gate, 'N/A') as departure_gate,
							       f.departure_actual,
							       f.departure_actual_runway,
							       f.departure_airport,
							       f.departure_estimated,
							       COALESCE(FLOOR(f.departure_delay / (1000 * 60)), 0) as departure_delay,
							       COALESCE(f.airline_name, 'N/A') as airline_name
							from flights f
							WHERE	Trim(Upper(f.flight_number))
							          ILIKE trim(upper('%' || $5 || '%'))
							          AND flight_status = $6
									  AND ((date_trunc('day', TO_DATE(f.flight_date, 'YYYY-MM-DD')) = CURRENT_DATE))

-- 							AND (
-- 						        $6 = '' -- No status provided, include all flights
-- 						        OR ($6 = 'active' AND f.flight_status = 'active')
-- 						        OR ($6 = 'canceled' AND f.flight_status = 'canceled')
-- 						        -- Add more conditions for other statuses as needed
--     						)
							ORDER BY
							CASE
									WHEN $3 = 'Flight Number'
									AND      $4 = 'ASC' THEN flight_number::text
							END ASC,
							CASE
									WHEN $3 = 'Flight Number'
									AND      $4 = 'DESC' THEN flight_number::text
							END DESC,
							CASE
									WHEN $3 = 'Flight Status'
									AND      $4 = 'ASC' THEN flight_status::text
							END ASC,
							CASE
									WHEN $3 = 'Flight Status'
									AND      $4 = 'DESC' THEN flight_status::text
							END DESC,
							CASE
									WHEN $3 = 'Flight Date'
									AND      $4 = 'ASC' THEN flight_date::text
							END ASC,
							CASE
									WHEN $3 = 'Flight Date'
									AND      $4 = 'DESC' THEN flight_date::text
							END DESC,
							departure_estimated
							offset $1 limit $2`

	offset := (page - 1) * pageSize

	return r.getFlightsData(ctx, query, offset, pageSize, orderBy, sortBy, flightNumber, flightStatus)
}

func (r *FlightsRepository) GetAllFlightsLocationsByStatus(ctx context.Context, flightStatus string) ([]models.LiveFlights, error) {
	query := `
				SELECT
				DISTINCT ON (f.flight_number)
			    f.flight_number,
			    COALESCE(al.airline_name, 'N/A') as airline_name,
			    COALESCE(ad.airport_name, 'N/A') AS departure_airport,
			    COALESCE(ad.latitude, 0.0) AS departure_latitude,
			    COALESCE(ad.longitude, 0.0) AS departure_longitude,
			    COALESCE(aa.airport_name, 'N/A') AS arrival_airport,
			    COALESCE(aa.latitude, 0.0) AS arrival_latitude,
			    COALESCE(aa.longitude, 0.0) AS arrival_longitude,
			    f.flight_date,
			    f.flight_status,
			    f.live_updated,
			    f.id,
			    f.arrival_actual,
			    f.arrival_actual_runway,
			    f.arrival_airport,
			    COALESCE(f.arrival_baggage, 'N/A'),
			    COALESCE(FLOOR(f.arrival_delay / (1000 * 60)), 0) as arrival_delay,
				f.arrival_estimated,
			    COALESCE(f.arrival_terminal, 'N/A'),
			    COALESCE(f.arrival_gate, 'N/A'),
			    f.arrival_timezone,
			    f.departure_scheduled,
			    departure_estimated_runway,
			    f.departure_timezone,
			    f.departure_terminal,
			    COALESCE(f.departure_gate, 'N/A'),
			    f.departure_actual,
			    f.departure_actual_runway,
				f.departure_airport,
				f.departure_estimated,
			    COALESCE(FLOOR(f.departure_delay / (1000 * 60)), 0) as departure_delay,
				COALESCE(ad.city_iata_code, 'N/A') AS departure_city_code,
			    COALESCE(ad.country_iso2 , 'N/A') AS departure_country_code,
			    COALESCE(aa.city_iata_code, 'N/A') AS arrival_city_code,
			    COALESCE(aa.country_iso2, 'N/A') AS arrival_country_code
			FROM
			    flights f
			        LEFT JOIN
			    airport ad ON f.departure_iata = ad.iata_code
			        LEFT JOIN
			    airport aa ON f.arrival_iata = aa.iata_code
			        LEFT JOIN
			    airline al on f.airline_iata = al.iata_code
			WHERE Trim(Upper(f.flight_number)) != '' AND
							Trim(Upper(f.codeshared_airline_name)) != '' AND
							Trim(Upper(f.codeshared_airline_name)) != 'N/A' AND
			    ad.latitude IS NOT NULL
			  AND ad.longitude IS NOT NULL
			  AND flight_status = $1
			  AND ((date_trunc('day', f.flight_date) = TO_CHAR(CURRENT_DATE, 'YYYY-MM-DD')))
			  ORDER BY f.flight_number`

	return r.getFlightsLocationsData(ctx, query, flightStatus)
}

func (r *FlightsRepository) getLiveFlightsData(ctx context.Context, query string,
	args ...interface{}) ([]models.LiveFlights, error) {
	var lf []models.LiveFlights

	rows, err := r.pgpool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var f models.LiveFlights

		err = rows.Scan(
			&f.Flight.Number,
			&f.FlightStatus,
			&f.Live.LiveUpdated,
			&f.ID,
			&f.Arrival.Actual,
			&f.Arrival.ActualRunway,
			&f.Arrival.Airport,
			&f.Arrival.Baggage,
			&f.Arrival.Delay,
			&f.Arrival.Estimated,
			&f.Arrival.Terminal,
			&f.Arrival.Gate,
			&f.Arrival.Timezone,
			&f.Departure.Scheduled,
			&f.Departure.EstimatedRunway,
			&f.Departure.Timezone,
			&f.Departure.Terminal,
			&f.Departure.Gate,
			&f.Departure.Actual,
			&f.Departure.ActualRunway,
			&f.Departure.Airport,
			&f.Departure.Estimated,
			&f.Departure.Delay,
			&f.Airline.Name,
			&f.Live.LiveLatitude,
			&f.Live.LiveLongitude,
		)
		if err != nil {
			return nil, err
		}

		// if f.Arrival.Delay != nil {
		//	minutes := *f.Arrival.Delay / (millis * minutes)
		//	f.Arrival.Delay = &minutes
		// }

		lf = append(lf, f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lf, nil
}

func (r *FlightsRepository) GetLiveFlights(ctx context.Context,
	page, pageSize int, orderBy, sortBy string) ([]models.LiveFlights, error) {
	query := `select
							       f.flight_number,
							       f.flight_status,
							       f.live_updated,
							       f.id,
							       f.arrival_actual,
							       f.arrival_actual_runway,
							       f.arrival_airport,
							    	COALESCE(f.arrival_baggage, 'N/A') as arrival_baggage,
							       COALESCE(FLOOR(f.arrival_delay / (1000 * 60)), 0) as arrival_delay,
									f.arrival_estimated,
    								COALESCE(f.arrival_terminal, 'N/A') as arrival_terminal,
							        COALESCE(f.arrival_gate, 'N/A') as arrival_gate,
							        f.arrival_timezone,
							        f.arrival_estimated,
							        f.departure_estimated_runway,
							       f.departure_timezone,
							       f.departure_terminal,
							       COALESCE(f.departure_gate, 'N/A') as departure_gate,
							       f.departure_actual,
							       f.departure_actual_runway,
							       f.departure_airport,
							       f.departure_estimated,
							       COALESCE(FLOOR(f.departure_delay / (1000 * 60)), 0) as departure_delay,
							       COALESCE(f.codeshared_airline_name, 'N/A') as airline_namel,
							       f.live_latitude,
									f.live_longitude

							from flights f
							WHERE f.live_latitude != 0 AND f.live_longitude != 0
							AND ((date_trunc('day', TO_DATE(f.flight_date, 'YYYY-MM-DD')) = CURRENT_DATE))
							AND Trim(Upper(f.flight_number)) != '' AND
							Trim(Upper(f.codeshared_airline_name)) != '' AND
							Trim(Upper(f.codeshared_airline_name)) != 'N/A'
							ORDER BY
							CASE
									WHEN $3 = 'Flight Number'
									AND      $4 = 'ASC' THEN flight_number::text
							END ASC,
							CASE
									WHEN $3 = 'Flight Number'
									AND      $4 = 'DESC' THEN flight_number::text
							END DESC,
							CASE
									WHEN $3 = 'Flight Status'
									AND      $4 = 'ASC' THEN flight_status::text
							END ASC,
							CASE
									WHEN $3 = 'Flight Status'
									AND      $4 = 'DESC' THEN flight_status::text
							END DESC,
							CASE
									WHEN $3 = 'Flight Date'
									AND      $4 = 'ASC' THEN flight_date::text
							END ASC,
							CASE
									WHEN $3 = 'Flight Date'
									AND      $4 = 'DESC' THEN flight_date::text
							END DESC,
							departure_estimated
							offset $1 limit $2`

	offset := (page - 1) * pageSize

	return r.getLiveFlightsData(ctx, query, offset, pageSize, orderBy, sortBy)
}

func (r *FlightsRepository) GetLiveFlightsLocations(ctx context.Context) ([]models.LiveFlights, error) {
	query := `select
							       f.flight_number,
							       f.flight_status,
							       f.live_updated,
							       f.id,
							       f.arrival_actual,
							       f.arrival_actual_runway,
							       f.arrival_airport,
							    	COALESCE(f.arrival_baggage, 'N/A') as arrival_baggage,
							       COALESCE(FLOOR(f.arrival_delay / (1000 * 60)), 0) as arrival_delay,
									f.arrival_estimated,
    								COALESCE(f.arrival_terminal, 'N/A') as arrival_terminal,
							        COALESCE(f.arrival_gate, 'N/A') as arrival_gate,
							        f.arrival_timezone,
							        f.arrival_estimated,
							        f.departure_estimated_runway,
							       f.departure_timezone,
							       f.departure_terminal,
							       COALESCE(f.departure_gate, 'N/A') as departure_gate,
							       f.departure_actual,
							       f.departure_actual_runway,
							       f.departure_airport,
							       f.departure_estimated,
							       COALESCE(FLOOR(f.departure_delay / (1000 * 60)), 0) as departure_delay,
							       COALESCE(f.codeshared_airline_name, 'N/A') as airline_namel,
							       f.live_latitude,
									f.live_longitude

							from flights f
							WHERE f.live_latitude != 0 AND f.live_longitude != 0
							AND ((date_trunc('day', TO_DATE(f.flight_date, 'YYYY-MM-DD')) = CURRENT_DATE))
							AND Trim(Upper(f.flight_number)) != '' AND
							Trim(Upper(f.codeshared_airline_name)) != '' AND
							Trim(Upper(f.codeshared_airline_name)) != 'N/A'
`

	return r.getLiveFlightsData(ctx, query)
}

func (r *FlightsRepository) GetFlightResumeByStatus(ctx context.Context, flightStatus string) (models.LiveFlightsResume, error) {
	query := `
        SELECT
            DISTINCT f.flight_status,
            COALESCE(f.airline_name, 'N/A') AS airline_name,
            COUNT(*) AS num_flights
        FROM
            flights f
        WHERE
            flight_status = $1
		AND
            COALESCE(f.airline_name, 'N/A') != 'N/A'
        	AND ((date_trunc('day', TO_DATE(f.flight_date, 'YYYY-MM-DD')) = CURRENT_DATE))
        	AND Trim(Upper(f.flight_number)) != '' AND
							Trim(Upper(f.codeshared_airline_name)) != '' AND
							Trim(Upper(f.codeshared_airline_name)) != 'N/A'
        GROUP BY
            f.flight_status,
            f.airline_name
        ORDER BY
            num_flights DESC
        LIMIT 1;
    `

	var lfr models.LiveFlightsResume
	err := r.pgpool.QueryRow(ctx, query, flightStatus).Scan(
		&lfr.Flight,
		&lfr.AirlineName,
		&lfr.NumFlights,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.LiveFlightsResume{}, nil
		}
		return models.LiveFlightsResume{}, err
	}
	return lfr, nil
}

func (r *FlightsRepository) GetFlightsResume(ctx context.Context) ([]models.LiveFlightsResume, error) {
	query := `WITH RankedFlights AS (
		    SELECT
		        f.flight_status,
		        COALESCE(f.airline_name, 'N/A') AS airline_name,
		        COUNT(*) AS num_flights,
		        ROW_NUMBER() OVER (PARTITION BY f.flight_status ORDER BY COUNT(*) DESC) AS rank_max
		    FROM
		        flights f
		    WHERE
		        COALESCE(f.airline_name, 'N/A') != 'N/A'
			AND
				((date_trunc('day', TO_DATE(f.flight_date, 'YYYY-MM-DD')) != CURRENT_DATE))
			AND Trim(Upper(f.flight_number)) != '' AND
							Trim(Upper(f.codeshared_airline_name)) != '' AND
							Trim(Upper(f.codeshared_airline_name)) != 'N/A'
		    GROUP BY
		        f.flight_status,
		        f.airline_name
		)
		SELECT
		    flight_status,
		    airline_name,
		    num_flights
		FROM
		    RankedFlights
		WHERE
		    rank_max = 1;
    `

	var lfr []models.LiveFlightsResume
	rows, err := r.pgpool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var l models.LiveFlightsResume
		err = rows.Scan(
			&l.Flight,
			&l.AirlineName,
			&l.NumFlights)
		if err != nil {
			return nil, err
		}

		lfr = append(lfr, l)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lfr, nil
}
