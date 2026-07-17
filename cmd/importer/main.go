package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/api"
	"github.com/FACorreiaa/Aviation-tracker/config"
	"github.com/FACorreiaa/Aviation-tracker/db"
)

func run(ctx context.Context) error {
	if err := config.LoadEnvironment(); err != nil {
		return err
	}
	databaseConfig, err := config.NewDatabaseConfig()
	if err != nil {
		return err
	}
	pool, err := db.Init(databaseConfig.ConnectionURL)
	if err != nil {
		return fmt.Errorf("initialize database pool: %w", err)
	}
	defer pool.Close()

	waitCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := db.WaitForDB(waitCtx, pool); err != nil {
		return err
	}

	importer := api.NewRepository(pool)
	steps := []struct {
		name string
		run  func() error
	}{
		{name: "airlines", run: importer.MigrateAirlineAPIData},
		{name: "aircraft", run: importer.MigrateAircraftAPIData},
		{name: "taxes", run: importer.MigrateTaxAPIData},
		{name: "airplanes", run: importer.MigrateAirplaneAPIData},
		{name: "airports", run: importer.MigrateAirportAPIData},
		{name: "countries", run: importer.MigrateCountryAPIData},
		{name: "cities", run: importer.MigrateCityAPIData},
		{name: "flights", run: importer.MigrateFlightAPIData},
	}
	for _, step := range steps {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := step.run(); err != nil {
			return fmt.Errorf("import %s: %w", step.name, err)
		}
	}
	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
