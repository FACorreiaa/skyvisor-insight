package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	if err := db.Migrate(pool); err != nil {
		return fmt.Errorf("run migrations: %w", err)
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
