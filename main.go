package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/api"
	"github.com/FACorreiaa/Aviation-tracker/app"
	"github.com/FACorreiaa/Aviation-tracker/config"
	"github.com/FACorreiaa/Aviation-tracker/db"
	"github.com/redis/go-redis/v9"
)

func run(ctx context.Context) error {
	//go:generate npx tailwindcss build -c tailwind.config.js -o ./controller/static/css/style.css -
	//go:generate ./tailwindcss -i controller/static/css/main.css -o controller/static/css/output.css --minify
	cfg, err := config.NewConfig()

	if err != nil {
		return err
	}

	c, err := config.InitConfig()

	var logHandler slog.Handler

	logHandlerOptions := slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.Log.Level,
	}
	if cfg.Log.Format == "json" {
		logHandler = slog.NewJSONHandler(os.Stdout, &logHandlerOptions)
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &logHandlerOptions)
	}
	slog.SetDefault(slog.New(logHandler))

	pool, err := db.Init(cfg.Database.ConnectionURL)
	if err != nil {
		log.Println(err)
	}
	defer pool.Close()

	db.WaitForDB(pool)

	redisClient, err := db.InitRedis(cfg.Redis.Host, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		return fmt.Errorf("failed to set config: %w", err)

	}
	defer func(redisClient *redis.Client) {
		err = redisClient.Close()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}(redisClient)

	if err = db.Migrate(pool); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)

	}

	startTime := time.Now()

	tableDataMigration := api.NewRepository(pool)
	if err = tableDataMigration.MigrateAirlineAPIData(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	if err = tableDataMigration.MigrateAircraftAPIData(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	if err = tableDataMigration.MigrateTaxAPIData(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	if err = tableDataMigration.MigrateAirplaneAPIData(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	if err = tableDataMigration.MigrateAirportAPIData(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	if err = tableDataMigration.MigrateCountryAPIData(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	if err = tableDataMigration.MigrateCityAPIData(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	if err = tableDataMigration.MigrateFlightAPIData(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	slog.Info("This operation took: ", time.Since(startTime))

	srv := &http.Server{
		Addr:         cfg.Server.Addr,
		WriteTimeout: cfg.Server.WriteTimeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		Handler:      app.Router(pool, []byte(cfg.Server.SessionKey), redisClient),
	}

	jobRepo := api.NewRepositoryJob(pool)
	jobService := api.NewServiceJob(jobRepo)
	jobService.StartAPICheckCronJob()

	go func() {
		slog.Info("Starting server " + cfg.Server.Addr)
		if err = srv.ListenAndServe(); err != nil {
			slog.Error("ListenAndServe", "error", err)
		}
	}()

	err = config.InitPprof(c.Pprof.Addr, strconv.Itoa(c.Pprof.Port))
	if err != nil {
		fmt.Printf("Error initializing pprof config: %s", err)
		panic(err)
	}

	<-ctx.Done() // Wait for cancellation signal

	// Shutdown server
	ctxShutdown, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulTimeout)
	defer cancel()

	if err = srv.Shutdown(ctxShutdown); err != nil {
		slog.Error("Error shutting down server", err)
	}

	slog.Info("Shutting down")
	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		cancel()
		os.Exit(1)
	}
}
