package repository

import (
	"context"

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
		)
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
	page, pageSize int, orderBy, sortBy, flightNumber string) ([]models.LiveFlights, error) {
	query := `select
							       f.flight_number,
							       f.flight_date,
							       f.flight_status,
							       f.live_updated,
							       f.id,
							       COALESCE(f.arrival_actual, 'Not defined') as arrival_actual,
							       COALESCE(f.arrival_actual_runway, 'Not defined') as arrival_actual,
							       f.arrival_airport,
							    	COALESCE(f.arrival_baggage, 'Not defined') as arrival_baggage,
							       COALESCE(FLOOR(f.arrival_delay / (1000 * 60)), 0) as arrival_delay,
							        f.arrival_estimated,
							        COALESCE(f.arrival_terminal, 'Not defined') as arrival_terminal,
							        COALESCE(f.arrival_gate, 'Not defined') as arrival_gate,
							        f.arrival_timezone,
							        f.departure_scheduled,
							        COALESCE(f.departure_estimated_runway, 'Not defined') as departure_estimated_runway,
							       f.departure_timezone,
							       f.departure_terminal,
							       COALESCE(f.departure_gate, 'Not defined') as departure_gate,
							       COALESCE(f.departure_actual, 'Not defined') as departure_actual,
							       COALESCE(f.departure_actual_runway, 'Not defined') as departure_actual_runway,
							       f.departure_airport,
							    	f.departure_estimated,
							       COALESCE(FLOOR(f.departure_delay / (1000 * 60)), 0) as departure_delay
							from flights f
							WHERE	Trim(Upper(f.flight_number)) ILIKE trim(upper('%' || $5 || '%'))
							
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
							END DESC
							offset $1 limit $2`

	offset := (page - 1) * pageSize

	return r.getFlightsData(ctx, query, offset, pageSize, orderBy, sortBy, flightNumber)
}

func (r *FlightsRepository) GetAllFlightsSum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT DISTINCT ON(flight_number) Count(flight_number)
										FROM   flights
										GROUP BY flight_number`)

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
			    COALESCE(al.airline_name, 'Not available') as airline_name,
			    COALESCE(ad.airport_name, 'Not available') AS departure_airport,
			    COALESCE(ad.latitude, 0.0) AS departure_latitude,
			    COALESCE(ad.longitude, 0.0) AS departure_longitude,
			    COALESCE(aa.airport_name, 'Not available') AS arrival_airport,
			    COALESCE(aa.latitude, 0.0) AS arrival_latitude,
			    COALESCE(aa.longitude, 0.0) AS arrival_longitude,
			    f.flight_date,
			    f.flight_status,
			    f.live_updated,
			    f.id,
			    COALESCE(f.arrival_actual, 'Not defined') as arrival_actual,
			    COALESCE(f.arrival_actual_runway, 'Not defined') as arrival_actual_runway,
			    f.arrival_airport,
			    COALESCE(f.arrival_baggage, 'Not defined'),
			    COALESCE(FLOOR(f.arrival_delay / (1000 * 60)), 0) as arrival_delay,
			    f.arrival_estimated,
			    COALESCE(f.arrival_terminal, 'Not defined'),
			    COALESCE(f.arrival_gate, 'Not defined'),
			    f.arrival_timezone,
			    f.departure_scheduled,
			    COALESCE(f.departure_estimated_runway, 'Not defined'),
			    f.departure_timezone,
			    f.departure_terminal,
			    COALESCE(f.departure_gate, 'Not defined'),
			    COALESCE(f.departure_actual, 'Not defined') as departure_actual,
			    COALESCE(f.departure_actual_runway, 'Not defined') as departure_actual_runway,
			    f.departure_airport,
			    f.departure_estimated,
			    COALESCE(FLOOR(f.departure_delay / (1000 * 60)), 0) as departure_delay,
				f.departure_icao,
				f.arrival_icao,
				COALESCE(ad.city_iata_code, 'Not available') AS departure_city_code,
			    COALESCE(ad.country_iso2 , 'Not available') AS departure_country_code,
			    COALESCE(aa.city_iata_code, 'Not available') AS arrival_city_code,
			    COALESCE(aa.country_iso2, 'Not available') AS arrival_country_code
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
			  AND ad.longitude IS NOT NULL
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

func (r *FlightsRepository) GetAllFlightsPreview(ctx context.Context) ([]models.LiveFlights, error) {
	query := `
				SELECT
				DISTINCT ON (f.flight_number)
			    f.flight_number,
			    COALESCE(al.airline_name, 'Not available') as airline_name,
			    COALESCE(ad.airport_name, 'Not available') AS departure_airport,
			    COALESCE(ad.latitude, 0.0) AS departure_latitude,
			    COALESCE(ad.longitude, 0.0) AS departure_longitude,
			    COALESCE(aa.airport_name, 'Not available') AS arrival_airport,
			    COALESCE(aa.latitude, 0.0) AS arrival_latitude,
			    COALESCE(aa.longitude, 0.0) AS arrival_longitude,
			    f.flight_date,
			    f.flight_status,
			    f.live_updated,
			    f.id,
			    COALESCE(f.arrival_actual, 'Not defined') as arrival_actual,
			    COALESCE(f.arrival_actual_runway, 'Not defined') as arrival_actual_runway,
			    f.arrival_airport,
			    COALESCE(f.arrival_baggage, 'Not defined'),
			    COALESCE(FLOOR(f.arrival_delay / (1000 * 60)), 0) as arrival_delay,
			    f.arrival_estimated,
			    COALESCE(f.arrival_terminal, 'Not defined'),
			    COALESCE(f.arrival_gate, 'Not defined'),
			    f.arrival_timezone,
			    f.departure_scheduled,
			    COALESCE(f.departure_estimated_runway, 'Not defined'),
			    f.departure_timezone,
			    f.departure_terminal,
			    COALESCE(f.departure_gate, 'Not defined'),
			    COALESCE(f.departure_actual, 'Not defined') as departure_actual,
			    COALESCE(f.departure_actual_runway, 'Not defined') as departure_actual_runway,
			    f.departure_airport,
			    f.departure_estimated,
			    COALESCE(FLOOR(f.departure_delay / (1000 * 60)), 0) as departure_delay,
				COALESCE(ad.city_iata_code, 'Not available') AS departure_city_code,
			    COALESCE(ad.country_iso2 , 'Not available') AS departure_country_code,
			    COALESCE(aa.city_iata_code, 'Not available') AS arrival_city_code,
			    COALESCE(aa.country_iso2, 'Not available') AS arrival_country_code
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
			  AND ad.longitude IS NOT NULL
			  ORDER BY f.flight_number;`

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
							       COALESCE(f.arrival_actual, 'Not defined') as arrival_actual,
							       COALESCE(f.arrival_actual_runway, 'Not defined') as arrival_actual,
							       f.arrival_airport,
							    	COALESCE(f.arrival_baggage, 'Not defined') as arrival_baggage,
							       COALESCE(FLOOR(f.arrival_delay / (1000 * 60)), 0) as arrival_delay,
							        f.arrival_estimated,
							        COALESCE(f.arrival_terminal, 'Not defined') as arrival_terminal,
							        COALESCE(f.arrival_gate, 'Not defined') as arrival_gate,
							        f.arrival_timezone,
							        f.departure_scheduled,
							        COALESCE(f.departure_estimated_runway, 'Not defined') as departure_estimated_runway,
							       f.departure_timezone,
							       f.departure_terminal,
							       COALESCE(f.departure_gate, 'Not defined') as departure_gate,
							       COALESCE(f.departure_actual, 'Not defined') as departure_actual,
							       COALESCE(f.departure_actual_runway, 'Not defined') as departure_actual_runway,
							       f.departure_airport,
							    	f.departure_estimated,
							       COALESCE(FLOOR(f.departure_delay / (1000 * 60)), 0) as departure_delay
							from flights f
							WHERE	Trim(Upper(f.flight_number))
							          ILIKE trim(upper('%' || $5 || '%'))
							          AND flight_status = $6
-- 							AND (
-- 						        $6 = '' -- No status provided, include all flights
-- 						        OR ($6 = 'active' AND f.flight_status = 'active')
-- 						        OR ($6 = 'cancelled' AND f.flight_status = 'cancelled')
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
							END DESC
							offset $1 limit $2`

	offset := (page - 1) * pageSize

	return r.getFlightsData(ctx, query, offset, pageSize, orderBy, sortBy, flightNumber, flightStatus)
}