package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"slices"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/api/structs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

type RepositoryJob struct {
	Conn *pgxpool.Pool
}

func NewRepositoryJob(db *pgxpool.Pool) *RepositoryJob {
	return &RepositoryJob{Conn: db}
}

func NewServiceJob(repo *RepositoryJob) *ServiceJob {
	return &ServiceJob{repo: repo}
}

type ServiceJob struct {
	repo *RepositoryJob
}

type Model struct {
	City structs.City
}

func (model Model) GetID() int {
	return model.City.CityID
}

// getExistingID retrieves existing table_id from the database.
func (s *ServiceJob) getExistingID(query string, id int, tableData []int) ([]int, error) { //nolint
	rows, err := s.repo.Conn.Query(context.Background(), query)
	if err != nil {
		handleError(err, "Error querying DB")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			handleError(err, "Error scanning IDs")
			return nil, err
		}
		// existingIDs = append(existingIDs, id)
	}

	return tableData, nil
}

// findNewCityData slcies version.
func (s *ServiceJob) findNewCityData(apiData []structs.City, tableData []int) []structs.City {
	var newData []structs.City

	for _, a := range apiData {
		if hasKey := slices.ContainsFunc(tableData, func(cityID int) bool {
			return cityID == a.CityID
		}); !hasKey {
			newData = append(newData, a)
		}
	}

	return newData
}

func (s *ServiceJob) findNewCountryData(apiData []structs.Country, tableData []int) []structs.Country {
	var newData []structs.Country

	for _, a := range apiData {
		if hasKey := slices.ContainsFunc(tableData, func(c int) bool {
			return c == a.CountryIsoNumeric
		}); !hasKey {
			newData = append(newData, a)
		}
	}

	return newData
}

func (s *ServiceJob) findNewAirportData(apiData []structs.Airport, tableData []int) []structs.Airport {
	var newData []structs.Airport

	for _, a := range apiData {
		if hasKey := slices.ContainsFunc(tableData, func(airportID int) bool {
			return airportID == a.AirportID
		}); !hasKey {
			newData = append(newData, a)
		}
	}

	return newData
}

func (s *ServiceJob) findNewAirplaneData(apiData []structs.Airplane, tableData []int) []structs.Airplane {
	var newData []structs.Airplane

	for _, a := range apiData {
		if hasKey := slices.ContainsFunc(tableData, func(airplaneID int) bool {
			return airplaneID == a.AirplaneID
		}); !hasKey {
			newData = append(newData, a)
		}
	}

	return newData
}

func (s *ServiceJob) findNewTaxData(apiData []structs.Tax, tableData []int) []structs.Tax {
	var newData []structs.Tax

	for _, a := range apiData {
		if hasKey := slices.ContainsFunc(tableData, func(taxID int) bool {
			return taxID == a.TaxID
		}); !hasKey {
			newData = append(newData, a)
		}
	}

	return newData
}

func (s *ServiceJob) findNewAirlineData(apiData []structs.Airline, tableData []int) []structs.Airline {
	var newData []structs.Airline

	for _, a := range apiData {
		if hasKey := slices.ContainsFunc(tableData, func(airlineID int) bool {
			return airlineID == a.AirlineID
		}); !hasKey {
			newData = append(newData, a)
		}
	}

	return newData
}

func (s *ServiceJob) findNewAircraftData(apiData []structs.Aircraft, tableData []int) []structs.Aircraft {
	var newData []structs.Aircraft

	for _, a := range apiData {
		if hasKey := slices.ContainsFunc(tableData, func(p int) bool {
			return p == a.PlaneTypeID
		}); !hasKey {
			newData = append(newData, a)
		}
	}

	return newData
}

func (s *ServiceJob) insertNewCities() error {
	apiData, err := fetchAviationStackData("cities")
	query := `select city_id from city`
	var tableData []int
	var cityID int

	// apiData, err := os.ReadFile("./api/data//cities.json")

	if err != nil {
		handleError(err, "error fetching data")
		return err
	}

	// Unmarshal the API response
	apiRes := new(structs.CityAPIData)
	if err = json.NewDecoder(bytes.NewReader(apiData)).Decode(&apiRes); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	// var models []Model
	// for _, city := range apiRes.Data {
	// 	models = append(models, Model{City: city})
	// }

	// Check for existing data in the database
	existingData, err := s.getExistingID(query, cityID, tableData)

	if err != nil {
		handleError(err, "error fetching existing data from the database")
		return err
	}

	// Identify new data that is not already in the database
	newDataMap := s.findNewCityData(apiRes.Data, existingData)
	// Insert only the new data into the database
	if len(newDataMap) > 0 {
		if _, err = s.repo.Conn.CopyFrom(
			context.Background(),
			pgx.Identifier{"city"},
			[]string{"gmt", "city_id", "iata_code", "country_iso2", "geoname_id",
				"latitude", "longitude", "city_name", "timezone", "created_at",
			},
			pgx.CopyFromSlice(len(newDataMap), func(i int) ([]interface{}, error) {
				city := newDataMap[i]
				return []interface{}{
					city.GMT,
					city.CityID,
					city.IataCode,
					city.CountryISO2,
					city.GeonameID,
					city.Latitude,
					city.Longitude,
					city.CityName,
					city.Timezone,
					formatTime(time.Now()),
				}, nil
			}),
		); err != nil {
			handleError(err, "error inserting new data into cities table")
			return err
		}

		slog.Info("New data inserted into the city table")
	} else {
		slog.Info("No new data to insert into the city table")
	}

	return nil
}

func (s *ServiceJob) insertNewCountries() error {
	apiData, err := fetchAviationStackData("countries")
	query := `select country_iso_numeric from country`
	tableData := make([]int, 0)
	var countryIsoNumeric int

	// apiData, err := os.ReadFile("./api/data//countries.json")

	if err != nil {
		handleError(err, "error fetching data")
		return err
	}

	// Unmarshal the API response
	apiRes := new(structs.CountryAPIData)
	if err = json.NewDecoder(bytes.NewReader(apiData)).Decode(&apiRes); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	// Check for existing data in the database
	existingData, err := s.getExistingID(query, countryIsoNumeric, tableData)

	if err != nil {
		handleError(err, "error fetching existing data from the database")
		return err
	}

	// Identify new data that is not already in the database
	newDataMap := s.findNewCountryData(apiRes.Data, existingData)

	// Insert only the new data into the database
	if len(newDataMap) > 0 {
		if _, err = s.repo.Conn.CopyFrom(
			context.Background(),
			pgx.Identifier{"country"},
			[]string{"country_name", "country_iso2", "country_iso3", "country_iso_numeric", "population",
				"capital", "continent", "currency_name", "currency_code", "fips_code",
				"phone_prefix", "created_at",
			},
			pgx.CopyFromSlice(len(newDataMap), func(i int) ([]interface{}, error) {
				country := newDataMap[i]
				return []interface{}{
					country.CountryName,
					country.CountryISO2,
					country.CountryIso3,
					country.CountryIsoNumeric,
					country.Population,
					country.Capital,
					country.Continent,
					country.CurrencyName,
					country.CurrencyCode,
					country.FipsCode,
					country.PhonePrefix,
					formatTime(time.Now()),
				}, nil
			}),
		); err != nil {
			handleError(err, "error inserting new data into country table")
			return err
		}

		slog.Info("New data inserted into the country table")
	} else {
		slog.Info("No new data to insert into the country table")
	}

	return nil
}

func (s *ServiceJob) insertNewAirports() error {
	apiData, err := fetchAviationStackData("countries")
	query := `select airport_id from airport`
	tableData := make([]int, 0)
	var airportID int

	// apiData, err := os.ReadFile("./api/data//airports.json")

	if err != nil {
		handleError(err, "error fetching data")
		return err
	}

	// Unmarshal the API response
	apiRes := new(structs.AirportAPIData)
	if err = json.NewDecoder(bytes.NewReader(apiData)).Decode(&apiRes); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	// Check for existing data in the database
	existingData, err := s.getExistingID(query, airportID, tableData)

	if err != nil {
		handleError(err, "error fetching existing data from the database")
		return err
	}

	// Identify new data that is not already in the database
	newDataMap := s.findNewAirportData(apiRes.Data, existingData)

	// Insert only the new data into the database
	if len(newDataMap) > 0 {
		if _, err = s.repo.Conn.CopyFrom(
			context.Background(),
			pgx.Identifier{"airport"},
			[]string{"gmt", "airport_id", "iata_code", "city_iata_code", "icao_code",
				"country_iso2", "geoname_id", "latitude", "longitude", "airport_name",
				"country_name", "phone_number", "timezone", "created_at",
			},
			pgx.CopyFromSlice(len(newDataMap), func(i int) ([]interface{}, error) {
				airport := newDataMap[i]
				return []interface{}{
					airport.GMT,
					airport.AirportID,
					airport.IataCode,
					airport.CityIataCode,
					airport.IcaoCode,
					airport.CountryISO2,
					airport.GeonameID,
					airport.Latitude,
					airport.Longitude,
					airport.AirportName,
					airport.CountryName,
					airport.PhoneNumber,
					airport.Timezone, formatTime(time.Now()),
				}, nil
			}),
		); err != nil {
			handleError(err, "error inserting new data into airport table")
			return err
		}

		slog.Info("New data inserted into the airport table")
	} else {
		slog.Info("No new data to insert into the airport table")
	}

	return nil
}

func (s *ServiceJob) insertNewAirplanes() error {
	apiData, err := fetchAviationStackData("countries")
	query := `select airplane_id from airplane`
	tableData := make([]int, 0)
	var airplaneID int

	// apiData, err := os.ReadFile("./api/data//airplane.json")

	if err != nil {
		handleError(err, "error fetching data")
		return err
	}

	// Unmarshal the API response
	apiRes := new(structs.AirplaneAPIData)
	if err = json.NewDecoder(bytes.NewReader(apiData)).Decode(&apiRes); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	// Check for existing data in the database
	existingData, err := s.getExistingID(query, airplaneID, tableData)

	if err != nil {
		handleError(err, "error fetching existing data from the database")
		return err
	}

	// Identify new data that is not already in the database
	newDataMap := s.findNewAirplaneData(apiRes.Data, existingData)

	// Insert only the new data into the database
	if len(newDataMap) > 0 {
		if _, err = s.repo.Conn.CopyFrom(
			context.Background(),
			pgx.Identifier{"airplane"},
			[]string{"iata_type", "airplane_id", "airline_iata_code", "iata_code_long", "iata_code_short",
				"airline_icao_code", "construction_number", "delivery_date", "engines_count", "engines_type",
				"first_flight_date", "icao_code_hex", "line_number", "model_code", "registration_number",
				"test_registration_number", "plane_age", "plane_class", "model_name", "plane_owner", "plane_series",
				"plane_status", "production_line", "registration_date", "rollout_date", "created_at",
			},
			pgx.CopyFromSlice(len(newDataMap), func(i int) ([]interface{}, error) {
				airplane := newDataMap[i]
				return []interface{}{
					airplane.IataType,
					airplane.AirplaneID,
					airplane.AirlineIataCode,
					airplane.IataCodeLong,
					airplane.IataCodeShort,
					airplane.AirlineIcaoCode,
					airplane.ConstructionNumber,
					airplane.DeliveryDate.Time,
					airplane.EnginesCount,
					airplane.EnginesType,
					airplane.FirstFlightDate.Time,
					airplane.IcaoCodeHex,
					airplane.LineNumber,
					airplane.ModelCode,
					airplane.RegistrationNumber,
					airplane.TestRegistrationNumber,
					airplane.PlaneAge,
					airplane.PlaneClass,
					airplane.ModelName,
					airplane.PlaneOwner,
					airplane.PlaneSeries,
					airplane.PlaneStatus,
					airplane.ProductionLine,
					airplane.RegistrationDate.Time,
					airplane.RolloutDate.Time,
					formatTime(time.Now()),
				}, nil
			}),
		); err != nil {
			handleError(err, "error inserting new data into airport table")
			return err
		}

		slog.Info("New data inserted into the airport table")
	} else {
		slog.Info("No new data to insert into the airport table")
	}

	return nil
}

func (s *ServiceJob) insertNewTax() error {
	apiData, err := fetchAviationStackData("countries")
	query := `select tax_id from tax`
	tableData := make([]int, 0)
	var taxID int

	// apiData, err := os.ReadFile("./api/data//tax.json")
	if err != nil {
		handleError(err, "error fetching data")
		return err
	}

	// Unmarshal the API response
	apiRes := new(structs.TaxAPIData)
	if err = json.NewDecoder(bytes.NewReader(apiData)).Decode(&apiRes); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	// Check for existing data in the database
	existingData, err := s.getExistingID(query, taxID, tableData)

	if err != nil {
		handleError(err, "error fetching existing data from the database")
		return err
	}

	// Identify new data that is not already in the database
	newDataMap := s.findNewTaxData(apiRes.Data, existingData)

	// Insert only the new data into the database
	if len(newDataMap) > 0 {
		if _, err = s.repo.Conn.CopyFrom(
			context.Background(),
			pgx.Identifier{"tax"},
			[]string{"tax_id", "tax_name", "iata_code", "created_at"},
			pgx.CopyFromSlice(len(newDataMap), func(i int) ([]interface{}, error) {
				tax := newDataMap[i]
				return []interface{}{
					tax.TaxID, tax.TaxName, tax.IataCode,
					formatTime(time.Now()),
				}, nil
			}),
		); err != nil {
			handleError(err, "error inserting new data into tax table")
			return err
		}

		slog.Info("New data inserted into the tax table")
	} else {
		slog.Info("No new data to insert into the tax table")
	}

	return nil
}

func (s *ServiceJob) insertNewAirline() error {
	apiData, err := fetchAviationStackData("countries")
	query := `select airline_id from airline`
	tableData := make([]int, 0)
	var airlineID int

	// apiData, err := os.ReadFile("./api/data//airline.json")

	if err != nil {
		handleError(err, "error fetching data")
		return err
	}

	// Unmarshal the API response
	apiRes := new(structs.AirlineAPIData)
	if err = json.NewDecoder(bytes.NewReader(apiData)).Decode(&apiRes); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	// Check for existing data in the database
	existingData, err := s.getExistingID(query, airlineID, tableData)

	if err != nil {
		handleError(err, "error fetching existing data from the database")
		return err
	}

	// Identify new data that is not already in the database
	newDataMap := s.findNewAirlineData(apiRes.Data, existingData)

	// Insert only the new data into the database
	if len(newDataMap) > 0 {
		if _, err := s.repo.Conn.CopyFrom(
			context.Background(),
			pgx.Identifier{"airline"},
			[]string{"fleet_average_age", "airline_id", "callsign", "hub_code", "iata_code", "icao_code", "country_iso2",
				"date_founded", "iata_prefix_accounting", "airline_name", "country_name", "fleet_size", "status", "type",
				"created_at",
			}, pgx.CopyFromSlice(len(newDataMap), func(i int) ([]interface{}, error) {
				airline := newDataMap[i]
				return []interface{}{
					airline.FleetAverageAge,
					airline.AirlineID,
					airline.CallSign,
					airline.HubCode,
					airline.IataCode,
					airline.IcaoCode,
					airline.CountryISO2,
					airline.DateFounded,
					airline.IataPrefixAccounting,
					airline.AirlineName,
					airline.CountryName,
					airline.FleetSize,
					airline.Status,
					airline.Type,
					formatTime(time.Now()),
				}, nil
			}),
		); err != nil {
			handleError(err, "error inserting new data into airline table")
			return err
		}

		slog.Info("New data inserted into the airline table")
	} else {
		slog.Info("No new data to insert into the airline table")
	}

	return nil
}

func (s *ServiceJob) insertNewAircraft() error {
	apiData, err := fetchAviationStackData("countries")
	query := `select plane_type_id from aircraft`
	tableData := make([]int, 0)
	var planeTypeID int

	// apiData, err := os.ReadFile("./api/data//aircraft.json")

	if err != nil {
		handleError(err, "error fetching data")
		return err
	}

	// Unmarshal the API response
	apiRes := new(structs.AircraftAPIData)
	if err := json.NewDecoder(bytes.NewReader(apiData)).Decode(&apiRes); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	// Check for existing data in the database
	existingData, err := s.getExistingID(query, planeTypeID, tableData)

	if err != nil {
		handleError(err, "error fetching existing data from the database")
		return err
	}

	// Identify new data that is not already in the database
	newDataMap := s.findNewAircraftData(apiRes.Data, existingData)

	// Insert only the new data into the database
	if len(newDataMap) > 0 {
		if _, err := s.repo.Conn.CopyFrom(
			context.Background(),
			pgx.Identifier{"aircraft"},
			[]string{"iata_code", "aircraft_name", "plane_type_id", "created_at"},
			pgx.CopyFromSlice(len(newDataMap), func(i int) ([]interface{}, error) {
				aircraft := newDataMap[i]
				return []interface{}{
					aircraft.IataCode,
					aircraft.AircraftName,
					aircraft.PlaneTypeID,
					formatTime(time.Now()),
				}, nil
			}),
		); err != nil {
			handleError(err, "error inserting new data into aircraft table")
			return err
		}

		slog.Info("New data inserted into the aircraft table")
	} else {
		slog.Info("No new data to insert into the aircraft table")
	}

	return nil
}

func (s *ServiceJob) insertNewFlight() error {
	data, err := fetchAviationStackData("flights", "limit=100")
	if err != nil {
		handleError(err, "error fetching data")
		return err
	}
	res := new(structs.FlightAPIData)
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&res); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	// Insert data from the JSON
	if _, err := s.repo.Conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"flights"},
		[]string{"flight_date", "flight_status", "departure_airport", "departure_timezone", "departure_iata",
			"departure_icao", "departure_terminal", "departure_gate", "departure_delay", "departure_scheduled",
			"departure_estimated", "departure_actual", "departure_estimated_runway", "departure_actual_runway",
			"arrival_airport", "arrival_timezone", "arrival_iata", "arrival_icao", "arrival_terminal",
			"arrival_gate", "arrival_baggage", "arrival_delay", "arrival_scheduled", "arrival_estimated",
			"arrival_actual", "arrival_estimated_runway", "arrival_actual_runway", "airline_name",
			"airline_iata", "airline_icao", "flight_number", "flight_iata",
			"flight_icao", "codeshared_airline_name", "codeshared_airline_iata", "codeshared_airline_icao",
			"codeshared_flight_number", "codeshared_flight_iata", "codeshared_flight_icao", "aircraft_registration",
			"aircraft_icao24", "aircraft_iata", "aircraft_icao", "live_updated", "live_latitude", "live_longitude",
			"live_altitude", "live_direction", "live_speed_horizontal", "live_speed_vertical", "live_is_ground",
			"created_at",
		},
		pgx.CopyFromSlice(len(res.Data), func(i int) ([]interface{}, error) {
			// id := uuid.New()
			departure := res.Data[i].Departure
			arrival := res.Data[i].Arrival
			airline := res.Data[i].Airline
			flight := res.Data[i].Flight
			aircraft := res.Data[i].Aircraft
			live := res.Data[i].Live
			return []interface{}{
				// id,
				res.Data[i].FlightDate,
				res.Data[i].FlightStatus,
				departure.Airport,
				departure.Timezone,
				departure.Iata,
				departure.Icao,
				departure.Terminal,
				departure.Gate,
				departure.Delay,
				departure.Scheduled,
				departure.Estimated,
				departure.Actual,
				departure.EstimatedRunway,
				departure.ActualRunway,
				arrival.Airport,
				arrival.Timezone,
				arrival.Iata,
				arrival.Icao,
				arrival.Terminal,
				arrival.Gate,
				arrival.Baggage,
				arrival.Delay,
				arrival.Scheduled,
				arrival.Estimated,
				arrival.Actual,
				arrival.EstimatedRunway,
				arrival.ActualRunway,
				airline.Name,
				airline.Iata,
				airline.Icao,
				flight.Number,
				flight.Iata,
				flight.Icao,
				flight.Codeshared.AirlineName,
				flight.Codeshared.AirlineIata,
				flight.Codeshared.AirlineIcao,
				flight.Codeshared.FlightNumber,
				flight.Codeshared.FlightIata,
				flight.Codeshared.FlightIcao,
				aircraft.AircraftRegistration,
				aircraft.AircraftIcao24,
				aircraft.AircraftIata,
				aircraft.AircraftIcao,
				live.LiveUpdated,
				live.LiveLatitude,
				live.LiveLongitude,
				live.LiveAltitude,
				live.LiveDirection,
				live.LiveSpeedHorizontal,
				live.LiveSpeedVertical,
				live.LiveIsGround,
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

func (s *ServiceJob) StartAPICheckCronJob() {
	c := cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger), // or use cron.DefaultLogger
	))
	slog.Info("Insert api check job")
	_, err := c.AddFunc("@weekly", func() {
		startTime := time.Now()
		err := s.insertNewCities()
		duration := time.Since(startTime)
		valueFromDuration := slog.DurationValue(duration)
		slog.Info("Cities job finished", slog.Attr{Key: "duration", Value: valueFromDuration})
		handleError(err, "Error checking for new cities")
	})
	handleError(err, "Error running cron job")

	_, err = c.AddFunc("@weekly", func() {
		startTime := time.Now()
		err := s.insertNewCountries()
		duration := time.Since(startTime)
		valueFromDuration := slog.DurationValue(duration)
		slog.Info("Country job finished", slog.Attr{Key: "duration", Value: valueFromDuration})
		handleError(err, "Error checking for new countries")
	})
	handleError(err, "Error running cron job")

	_, err = c.AddFunc("@weekly", func() {
		startTime := time.Now()
		err := s.insertNewAirports()
		duration := time.Since(startTime)
		valueFromDuration := slog.DurationValue(duration)
		slog.Info("Airport job finished", slog.Attr{Key: "duration", Value: valueFromDuration})
		handleError(err, "Error checking for new airports")
	})
	handleError(err, "Error running cron job")

	_, err = c.AddFunc("@weekly", func() {
		startTime := time.Now()
		err := s.insertNewAirplanes()
		duration := time.Since(startTime)
		valueFromDuration := slog.DurationValue(duration)
		slog.Info("Airplane job finished", slog.Attr{Key: "duration", Value: valueFromDuration})
		handleError(err, "Error checking for new airplanes")
	})
	handleError(err, "Error running cron job")

	_, err = c.AddFunc("@weekly", func() {
		startTime := time.Now()
		err := s.insertNewTax()
		duration := time.Since(startTime)
		valueFromDuration := slog.DurationValue(duration)
		slog.Info("Tax job finished", slog.Attr{Key: "duration", Value: valueFromDuration})
		handleError(err, "Error checking for new tax")
	})
	handleError(err, "Error running cron job")

	_, _ = c.AddFunc("@weekly", func() {
		startTime := time.Now()
		err := s.insertNewAirline()
		duration := time.Since(startTime)
		valueFromDuration := slog.DurationValue(duration)
		slog.Info("Airline job finished", slog.Attr{Key: "duration", Value: valueFromDuration})
		handleError(err, "Error checking for new airline")
	})
	_, _ = c.AddFunc("@weekly", func() {
		startTime := time.Now()
		err := s.insertNewAircraft()
		duration := time.Since(startTime)
		valueFromDuration := slog.DurationValue(duration)
		slog.Info("Aircraft job finished", slog.Attr{Key: "duration", Value: valueFromDuration})
		handleError(err, "Error checking for new aircraft")
	})

	_, err = c.AddFunc("@every 3h", func() {
		startTime := time.Now()
		err := s.insertNewFlight()
		duration := time.Since(startTime)
		valueFromDuration := slog.DurationValue(duration)
		slog.Info("Live flights job finished", slog.Attr{Key: "duration", Value: valueFromDuration})
		handleError(err, "Error checking live flights")
	})
	handleError(err, "Error running cron job")

	c.Start()
}
