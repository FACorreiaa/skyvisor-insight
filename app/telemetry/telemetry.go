// Package telemetry configures OpenTelemetry tracing and Sentry error
// reporting. Both are opt-in via environment: tracing activates when
// OTEL_EXPORTER_OTLP_ENDPOINT is set (standard SDK variable, e.g.
// http://alloy.observability.svc.cluster.local:4317), Sentry when
// SENTRY_DSN is set.
package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// SetupTracing installs a global OTLP trace provider when
// OTEL_EXPORTER_OTLP_ENDPOINT is set. The returned shutdown function flushes
// pending spans and is safe to call even when tracing is disabled.
func SetupTracing(ctx context.Context, serviceName string) (func(context.Context) error, error) {
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" {
		return func(context.Context) error { return nil }, nil
	}

	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("create OTLP trace exporter: %w", err)
	}

	res, err := resource.Merge(resource.Default(), resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
	))
	if err != nil {
		return nil, fmt.Errorf("build OTel resource: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	))
	slog.Info("otel tracing enabled", "service", serviceName)
	return provider.Shutdown, nil
}

// SetupSentry initialises Sentry when SENTRY_DSN is set. The returned flush
// function drains the event queue and is safe to call when disabled.
func SetupSentry(serviceName string) (func(), error) {
	dsn := os.Getenv("SENTRY_DSN")
	if dsn == "" {
		return func() {}, nil
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		Environment: os.Getenv("SENTRY_ENVIRONMENT"),
		Release:     os.Getenv("SENTRY_RELEASE"),
		ServerName:  serviceName,
	})
	if err != nil {
		return nil, fmt.Errorf("initialize sentry: %w", err)
	}
	slog.Info("sentry error reporting enabled", "service", serviceName)
	return func() { sentry.Flush(2 * time.Second) }, nil
}
