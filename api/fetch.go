package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/FACorreiaa/go-ollama/api/structs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func fetchAviationStackData(endpoint string, queryParams ...string) ([]byte, error) {
	accessKey := os.Getenv("AVIATION_STACK_API_KEY")
	if accessKey == "" {
		return nil, fmt.Errorf("missing API access key")
	}

	baseURL := "http://api.aviationstack.com/v1/"

	// Parse the base URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	// Set the endpoint path
	parsedURL.Path += endpoint

	// Create a new query parameters object
	query := parsedURL.Query()

	// Add the access key parameter
	query.Set("access_key", accessKey)

	// Add additional query parameters
	if len(queryParams) > 0 {
		for _, param := range queryParams {
			parts := strings.SplitN(param, "=", 2)
			if len(parts) == 2 {
				query.Set(parts[0], parts[1])
			}
		}
	}

	parsedURL.RawQuery = query.Encode()

	finalURL := parsedURL.String()

	response, err := http.Get(finalURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %v", err)
	}

	if response.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("something is not ok")
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	defer response.Body.Close()

	return body, nil
}

func FetchAndInsertCityData(conn *pgxpool.Pool) error {
	data, err := fetchAviationStackData("cities", "limit=1000000")
	if err != nil {
		handleError(err, "error fetching data")
		return err
	}

	res := new(structs.CityApiData)

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&res); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	//Insert data from the json
	if _, err := conn.CopyFrom(

		context.Background(),
		pgx.Identifier{"city"},
		[]string{"gmt", "city_id", "iata_code", "country_iso2", "geoname_id",
			"latitude", "longitude", "city_name", "timezone", "created_at",
		},
		pgx.CopyFromSlice(len(res.Data), func(i int) ([]interface{}, error) {
			return []interface{}{
				res.Data[i].GMT,
				res.Data[i].CityID,
				res.Data[i].IataCode,
				res.Data[i].CountryISO2,
				res.Data[i].GeonameID,
				res.Data[i].Latitude,
				res.Data[i].Longitude,
				res.Data[i].CityName,
				res.Data[i].Timezone,
				formatTime(time.Now()),
			}, nil
		}),
	); err != nil {
		handleError(err, "error inserting data into cities table")
		return err
	}

	slog.Info("Data inserted into the city table")

	return nil
}

func FetchAndInsertCountryData(conn *pgxpool.Pool) error {
	res := new(structs.CountryApiData)
	data, err := fetchAviationStackData("countries", "limit=1000000")
	if err != nil {
		handleError(err, "error fetching data")
		return err
	}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&res); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	//Insert data from the json
	if _, err := conn.CopyFrom(

		context.Background(),
		pgx.Identifier{"country"},
		[]string{"country_name", "country_iso2", "country_iso3", "country_iso_numeric", "population",
			"capital", "continent", "currency_name", "currency_code", "fips_code",
			"phone_prefix", "created_at",
		},
		pgx.CopyFromSlice(len(res.Data), func(i int) ([]interface{}, error) {
			return []interface{}{
				res.Data[i].CountryName,
				res.Data[i].CountryISO2,
				res.Data[i].CountryIso3,
				res.Data[i].CountryIsoNumeric,
				res.Data[i].Population,
				res.Data[i].Capital,
				res.Data[i].Continent,
				res.Data[i].CurrencyName,
				res.Data[i].CurrencyCode,
				res.Data[i].FipsCode,
				res.Data[i].PhonePrefix,
				formatTime(time.Now()),
			}, nil
		}),
	); err != nil {
		handleError(err, "error inserting data into country table")
		return err
	}

	slog.Info("Data inserted into the country table")

	return nil
}

func FetchAndInsertAirportData(conn *pgxpool.Pool) error {
	res := new(structs.AirportApiData)
	data, err := fetchAviationStackData("airports", "limit=1000000")
	if err != nil {
		handleError(err, "error fetching data")
		return err
	}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&res); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	//Insert data from the json
	if _, err := conn.CopyFrom(

		context.Background(),
		pgx.Identifier{"airport"},
		[]string{"gmt", "airport_id", "iata_code", "city_iata_code", "icao_code",
			"country_iso2", "geoname_id", "latitude", "longitude", "airport_name",
			"country_name", "phone_number", "timezone", "created_at",
		},
		pgx.CopyFromSlice(len(res.Data), func(i int) ([]interface{}, error) {
			return []interface{}{
				res.Data[i].GMT, res.Data[i].AirportId, res.Data[i].IataCode,
				res.Data[i].CityIataCode, res.Data[i].IcaoCode, res.Data[i].CountryISO2,
				res.Data[i].GeonameID, res.Data[i].Latitude, res.Data[i].Longitude,
				res.Data[i].AirportName, res.Data[i].CountryName, res.Data[i].PhoneNumber,
				res.Data[i].Timezone, formatTime(time.Now()),
			}, nil
		}),
	); err != nil {
		handleError(err, "error inserting data into airports table")
		return err
	}

	slog.Info("Data inserted into the airport table")
	return nil
}

func FetchAndInsertAirplaneData(conn *pgxpool.Pool) error {
	data, err := fetchAviationStackData("airplanes", "limit=1000000")
	if err != nil {
		handleError(err, "error fetching data")
		return err
	}
	res := new(structs.AirplaneApiData)
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&res); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	//Insert data from the json
	if _, err := conn.CopyFrom(

		context.Background(),
		pgx.Identifier{"airplane"},
		[]string{"iata_type", "airplane_id", "airline_iata_code", "iata_code_long", "iata_code_short",
			"airline_icao_code", "construction_number", "delivery_date", "engines_count", "engines_type",
			"first_flight_date", "icao_code_hex", "line_number", "model_code", "registration_number",
			"test_registration_number", "plane_age", "plane_class", "model_name", "plane_owner", "plane_series",
			"plane_status", "production_line", "registration_date", "rollout_date", "created_at",
		},
		pgx.CopyFromSlice(len(res.Data), func(i int) ([]interface{}, error) {
			return []interface{}{
				res.Data[i].IataType,
				res.Data[i].AirplaneId,
				res.Data[i].AirlineIataCode,
				res.Data[i].IataCodeLong,
				res.Data[i].IataCodeShort,
				res.Data[i].AirlineIcaoCode,
				res.Data[i].ConstructionNumber,
				res.Data[i].DeliveryDate.Time,
				res.Data[i].EnginesCount,
				res.Data[i].EnginesType,
				res.Data[i].FirstFlightDate.Time,
				res.Data[i].IcaoCodeHex,
				res.Data[i].LineNumber,
				res.Data[i].ModelCode,
				res.Data[i].RegistrationNumber,
				res.Data[i].TestRegistrationNumber,
				res.Data[i].PlaneAge,
				res.Data[i].PlaneClass,
				res.Data[i].ModelName,
				res.Data[i].PlaneOwner,
				res.Data[i].PlaneSeries,
				res.Data[i].PlaneStatus,
				res.Data[i].ProductionLine,
				res.Data[i].RegistrationDate.Time,
				res.Data[i].RolloutDate.Time,
				formatTime(time.Now()),
			}, nil
		}),
	); err != nil {
		handleError(err, "error inserting data into airplane table")
		return err
	}

	slog.Info("Data inserted into the airplane table")

	return nil
}

func FetchAndInsertTaxData(conn *pgxpool.Pool) error {
	data, err := fetchAviationStackData("taxes", "limit=1000000")
	if err != nil {
		handleError(err, "error fetching data")
		return err
	}
	res := new(structs.TaxApiData)
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&res); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	//Insert data from the json
	if _, err := conn.CopyFrom(

		context.Background(),
		pgx.Identifier{"tax"},
		[]string{"tax_id", "tax_name", "iata_code", "created_at"},
		pgx.CopyFromSlice(len(res.Data), func(i int) ([]interface{}, error) {
			return []interface{}{
				res.Data[i].TaxId, res.Data[i].TaxName, res.Data[i].IataCode,
				formatTime(time.Now()),
			}, nil
		}),
	); err != nil {
		handleError(err, "error inserting data into tax table")
		return err
	}

	slog.Info("Data inserted into the aircraft table")
	return nil
}

func FetchAndInsertAircraftData(conn *pgxpool.Pool) error {
	res := new(structs.AircraftApiData)
	data, err := fetchAviationStackData("aircraft_types", "limit=1000000")
	if err != nil {
		handleError(err, "error fetching data")
		return err
	}

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&res); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	//Insert data from the json
	if _, err := conn.CopyFrom(

		context.Background(),
		pgx.Identifier{"aircraft"},
		[]string{"iata_code", "aircraft_name", "plane_type_id", "created_at"},
		pgx.CopyFromSlice(len(res.Data), func(i int) ([]interface{}, error) {
			return []interface{}{
				res.Data[i].IataCode,
				res.Data[i].AircraftName,
				res.Data[i].PlaneTypeId,
				formatTime(time.Now()),
			}, nil
		}),
	); err != nil {
		handleError(err, "error inserting data into aircraft table")
		return err
	}

	slog.Info("Data inserted into the aircraft table")

	return nil
}

func FetchAndInsertAirlineData(conn *pgxpool.Pool) error {
	data, err := fetchAviationStackData("airlines", "limit=1000000")
	if err != nil {
		handleError(err, "error fetching data")
		return err
	}
	res := new(structs.AirlineApiData)
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&res); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	// Insert data from the JSON
	if _, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"airline"},
		[]string{"fleet_average_age", "airline_id", "callsign", "hub_code", "iata_code", "icao_code", "country_iso2",
			"date_founded", "iata_prefix_accounting", "airline_name", "country_name", "fleet_size", "status", "type",
			"created_at",
		},
		pgx.CopyFromSlice(len(res.Data), func(i int) ([]interface{}, error) {
			return []interface{}{
				res.Data[i].FleetAverageAge,
				res.Data[i].AirlineId,
				res.Data[i].Callsign,
				res.Data[i].HubCode,
				res.Data[i].IataCode,
				res.Data[i].IcaoCode,
				res.Data[i].CountryISO2,
				res.Data[i].DateFounded,
				res.Data[i].IataPrefixAccounting,
				res.Data[i].AirlineName,
				res.Data[i].CountryName,
				res.Data[i].FleetSize,
				res.Data[i].Status,
				res.Data[i].Type,
				formatTime(time.Now()),
			}, nil
		}),
	); err != nil {
		handleError(err, "error inserting data into airline table")
		return err
	}

	slog.Info("Data inserted into the airline table")
	return nil
}

func FetchAndInsertFlightData(conn *pgxpool.Pool) error {
	//data, err := os.ReadFile("./api/data/flights.json")

	data, err := fetchAviationStackData("flights", "limit=1000000")
	if err != nil {
		handleError(err, "error fetching data")
		return err
	}
	res := new(structs.FlightApiData)
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&res); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	// Insert data from the JSON
	if _, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"flights"},
		[]string{"id", "flight_date", "flight_status", "departure_airport", "departure_timezone", "departure_iata",
			"departure_icao", "departure_terminal", "departure_gate", "departure_delay", "departure_scheduled",
			"departure_estimated", "departure_actual", "departure_estimated_runway", "departure_actual_runway",
			"arrival_airport", "arrival_timezone", "arrival_iata", "arrival_icao", "arrival_terminal",
			"arrival_gate", "arrival_baggage", "arrival_delay", "arrival_scheduled", "arrival_estimated",
			"arrival_actual", "arrival_estimated_runway", "arrival_actual_runway", "flight_number", "flight_iata",
			"flight_icao", "codeshared_airline_name", "codeshared_airline_iata", "codeshared_airline_icao",
			"codeshared_flight_number", "codeshared_flight_iata", "codeshared_flight_icao",
			"aircraft_registration", "aircraft_iata", "aircraft_icao", "aircraft_icao25", "live_updated",
			"live_latitude", "live_longitude", "live_altitude", "live_direction", "live_speed_horizontal",
			"live_speed_vertical", "live_is_ground", "created_at",
		},
		pgx.CopyFromSlice(len(res.Data), func(i int) ([]interface{}, error) {
			id := uuid.New()
			return []interface{}{
				id, res.Data[i].FlightDate, res.Data[i].FlightStatus, res.Data[i].Departure.Airport,
				res.Data[i].Departure.Timezone, res.Data[i].Departure.Iata, res.Data[i].Departure.Icao,
				res.Data[i].Departure.Terminal, res.Data[i].Departure.Gate, res.Data[i].Departure.Delay,
				res.Data[i].Departure.Scheduled, res.Data[i].Departure.Estimated, res.Data[i].Departure.Actual,
				res.Data[i].Departure.EstimatedRunway, res.Data[i].Departure.ActualRunway,
				res.Data[i].Arrival.Airport, res.Data[i].Arrival.Timezone, res.Data[i].Arrival.Iata,
				res.Data[i].Arrival.Icao, res.Data[i].Arrival.Terminal, res.Data[i].Arrival.Gate,
				res.Data[i].Arrival.Baggage, res.Data[i].Arrival.Delay, res.Data[i].Arrival.Scheduled, res.Data[i].Arrival.Estimated,
				res.Data[i].Arrival.Actual, res.Data[i].Arrival.EstimatedRunway, res.Data[i].Arrival.ActualRunway,
				res.Data[i].Flight.Number, res.Data[i].Flight.Iata, res.Data[i].Flight.Icao,
				res.Data[i].Flight.Codeshared.AirlineName,
				res.Data[i].Flight.Codeshared.AirlineIata, res.Data[i].Flight.Codeshared.AirlineIcao,
				res.Data[i].Flight.Codeshared.FlightNumber, res.Data[i].Flight.Codeshared.FlightIata,
				res.Data[i].Flight.Codeshared.FlightIcao, res.Data[i].Aircraft.AircraftRegistration,
				res.Data[i].Aircraft.AircraftIata, res.Data[i].Aircraft.AircraftIcao, res.Data[i].Aircraft.AircraftIcao24,
				res.Data[i].Live.LiveUpdated, res.Data[i].Live.LiveLatitude, res.Data[i].Live.LiveLongitude,
				res.Data[i].Live.LiveAltitude, res.Data[i].Live.LiveDirection, res.Data[i].Live.LiveSpeedHorizontal,
				res.Data[i].Live.LiveSpeedVertical, res.Data[i].Live.LiveIsGround,

				formatTime(time.Now()),
			}, nil
		}),
	); err != nil {
		handleError(err, "error inserting data into flights table")
		return err
	}

	slog.Info("Data inserted into the flights table")
	return nil
}
