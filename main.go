package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/api"
	"github.com/FACorreiaa/Aviation-tracker/config"
	"github.com/FACorreiaa/Aviation-tracker/controller"
	"github.com/FACorreiaa/Aviation-tracker/db"
	"github.com/redis/go-redis/v9"
)

func main() {
	//go:generate npx tailwindcss build -c tailwind.config.js -o ./controller/static/css/style.css -

	cfg, err := config.NewConfig()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

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
		os.Exit(1)
	}
	defer pool.Close()

	db.WaitForDB(pool)

	redisClient, err := db.InitRedis(cfg.Redis.Host, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer func(redisClient *redis.Client) {
		err = redisClient.Close()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}(redisClient)
	// db.WaitForRedis(redisClient)

	if err = db.Migrate(pool); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	startTime := time.Now()

	tableDataMigration := api.NewRepository(pool)
	if err = tableDataMigration.MigrateAirlineAPIData(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err = tableDataMigration.MigrateAircraftAPIData(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err = tableDataMigration.MigrateTaxAPIData(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err = tableDataMigration.MigrateAirplaneAPIData(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err = tableDataMigration.MigrateAirportAPIData(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err = tableDataMigration.MigrateCountryAPIData(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err = tableDataMigration.MigrateCityAPIData(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err = tableDataMigration.MigrateFlightAPIData(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("This operation took: ", time.Since(startTime))

	// if err = db.MigrateRedis(redisClient); err != nil {
	//	log.Println(err)
	//	os.Exit(1)
	//}

	srv := &http.Server{
		Addr:         cfg.Server.Addr,
		WriteTimeout: cfg.Server.WriteTimeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		Handler:      controller.Router(pool, []byte(cfg.Server.SessionKey), redisClient),
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulTimeout)
	defer cancel()
	srv.Shutdown(ctx)
	slog.Info("shutting down")
	os.Exit(0)
}
