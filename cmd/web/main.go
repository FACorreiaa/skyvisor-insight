package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app"
	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
	"github.com/FACorreiaa/Aviation-tracker/app/auth"
	"github.com/FACorreiaa/Aviation-tracker/config"
	"github.com/FACorreiaa/Aviation-tracker/db"
)

func run(ctx context.Context) error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	logOptions := slog.HandlerOptions{Level: cfg.Log.Level}
	var logHandler slog.Handler
	if cfg.Log.Format == "json" {
		logHandler = slog.NewJSONHandler(os.Stdout, &logOptions)
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &logOptions)
	}
	slog.SetDefault(slog.New(logHandler))

	pool, err := db.Init(cfg.Database.ConnectionURL)
	if err != nil {
		return fmt.Errorf("initialize database pool: %w", err)
	}
	defer pool.Close()

	startupCtx, startupCancel := context.WithTimeout(ctx, 15*time.Second)
	defer startupCancel()
	if err := db.WaitForDB(startupCtx, pool); err != nil {
		return err
	}

	redisClient, err := db.InitRedis(cfg.Redis.Host, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		return fmt.Errorf("initialize redis client: %w", err)
	}
	defer func() {
		if closeErr := redisClient.Close(); closeErr != nil {
			slog.Error("close redis client", "error", closeErr)
		}
	}()
	if err := redisClient.Ping(startupCtx).Err(); err != nil {
		return fmt.Errorf("connect to redis: %w", err)
	}

	var oidcClient *auth.Client
	if cfg.OIDC != nil {
		oidcClient, err = auth.New(startupCtx, auth.Config{
			IssuerURL:    cfg.OIDC.IssuerURL,
			ClientID:     cfg.OIDC.ClientID,
			ClientSecret: cfg.OIDC.ClientSecret,
			RedirectURL:  cfg.OIDC.RedirectURL,
			Audience:     cfg.OIDC.Audience,
		})
		if err != nil {
			return fmt.Errorf("initialize OIDC client: %w", err)
		}
	} else {
		slog.Warn("OIDC is not configured; sign-in is disabled")
	}

	var apiClient *apiclient.Client
	if cfg.API != nil {
		apiClient, err = apiclient.New(cfg.API.BaseURL)
		if err != nil {
			return fmt.Errorf("initialize skyvisor-api client: %w", err)
		}
	} else {
		slog.Warn("SKYVISOR_API_URL is not configured; API-backed features are disabled")
	}

	server := &http.Server{
		Addr:           cfg.Server.Addr,
		WriteTimeout:   cfg.Server.WriteTimeout,
		ReadTimeout:    cfg.Server.ReadTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
		MaxHeaderBytes: 1 << 20,
		Handler:        app.Router(pool, []byte(cfg.Server.SessionKey), cfg.Server.CookieSecure, redisClient, oidcClient, apiClient),
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("web server listening", "address", cfg.Server.Addr)
		if listenErr := server.ListenAndServe(); listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			serverErrors <- listenErr
		}
	}()

	if err := config.InitPprof(os.Getenv("PPROF_ADDR"), os.Getenv("PPROF_PORT")); err != nil {
		return fmt.Errorf("initialize pprof: %w", err)
	}

	select {
	case <-ctx.Done():
	case err = <-serverErrors:
		return fmt.Errorf("serve HTTP: %w", err)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.GracefulTimeout)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}
	slog.Info("web server stopped")
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
