package api

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
	"time"
)

func handleError(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05-07:00")
}

type MigrateInterface interface {
	MigrateAirlineAPIData() error
	MigrateAircraftAPIData() error
	MigrateTaxAPIData() error
	MigrateAirplaneAPIData() error
	MigrateAirportAPIData() error
	MigrateCountryAPIData() error
	MigrateCityAPIData() error
	MigrateFlightAPIData() error
}

type MigrateRepository struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) MigrateInterface {
	return &MigrateRepository{conn: conn}
}

/*Airline Migration function */

func (m *MigrateRepository) MigrateAirlineAPIData() error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := m.conn.QueryRow(ctx, "SELECT COUNT(*) FROM airline").Scan(&count); err != nil {
		handleError(err, "Error querying the table")
		return err
	}

	if count == 0 {
		// No data in the airline table, fetch from the external API
		if err := FetchAndInsertAirlineData(m.conn); err != nil {
			handleError(err, "Error inserting data")
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

/*Aircraft Migration function */

func (m *MigrateRepository) MigrateAircraftAPIData() error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := m.conn.QueryRow(ctx, "SELECT COUNT(*) FROM aircraft").Scan(&count); err != nil {
		handleError(err, "Error querying the table")
		return err
	}

	if count == 0 {
		// No data in the airline table, fetch from the external API
		if err := FetchAndInsertAircraftData(m.conn); err != nil {
			handleError(err, "Error inserting data")
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

/*Tax Migration function */

func (m *MigrateRepository) MigrateTaxAPIData() error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := m.conn.QueryRow(ctx, "SELECT COUNT(*) FROM tax").Scan(&count); err != nil {
		handleError(err, "Error querying the table")
		return err
	}

	if count == 0 {
		// No data in the airline table, fetch from the external API
		if err := FetchAndInsertTaxData(m.conn); err != nil {
			handleError(err, "Error inserting data")
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

/* Airplane */

func (m *MigrateRepository) MigrateAirplaneAPIData() error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := m.conn.QueryRow(ctx, "SELECT COUNT(*) FROM airplane").Scan(&count); err != nil {
		handleError(err, "Error querying the table")
		return err
	}

	if count == 0 {
		// No data in the airline table, fetch from the external API
		if err := FetchAndInsertAirplaneData(m.conn); err != nil {
			handleError(err, "Error inserting data")
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

/* Airports */

func (m *MigrateRepository) MigrateAirportAPIData() error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := m.conn.QueryRow(ctx, "SELECT COUNT(*) FROM airport").Scan(&count); err != nil {
		handleError(err, "Error querying the table")
		return err
	}

	if count == 0 {
		// No data in the airport table, fetch from the external API
		if err := FetchAndInsertAirportData(m.conn); err != nil {
			handleError(err, "Error inserting data")
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

/* Countries */

func (m *MigrateRepository) MigrateCountryAPIData() error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := m.conn.QueryRow(ctx, "SELECT COUNT(*) FROM country").Scan(&count); err != nil {
		handleError(err, "Error querying the table")
		return err
	}

	if count == 0 {
		// No data in the country table, fetch from the external API
		if err := FetchAndInsertCountryData(m.conn); err != nil {
			handleError(err, "Error inserting data")
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

/* Cities */

func (m *MigrateRepository) MigrateCityAPIData() error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := m.conn.QueryRow(ctx, "SELECT COUNT(*) FROM city").Scan(&count); err != nil {
		handleError(err, "Error querying the table")
		return err
	}

	if count == 0 {
		// No data in the airport table, fetch from the external API
		if err := FetchAndInsertCityData(m.conn); err != nil {
			handleError(err, "Error inserting data")
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

/* Flights */

func (m *MigrateRepository) MigrateFlightAPIData() error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := m.conn.QueryRow(ctx, "SELECT COUNT(*) FROM flights").Scan(&count); err != nil {
		handleError(err, "Error querying the table")
		return err
	}

	if count == 0 {
		// No data in the flights table, fetch from the external API
		if err := FetchAndInsertFlightData(m.conn); err != nil {
			handleError(err, "Error inserting data")
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}
