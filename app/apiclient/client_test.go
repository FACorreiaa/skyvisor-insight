package apiclient

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientCallsAPI(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":{"code":"unauthorized","message":"nope"}}`))
			return
		}
		switch r.Method + " " + r.URL.Path {
		case "GET /v1/trips":
			_, _ = w.Write([]byte(`{"data":[{"id":"a","name":"Lisbon","segments":[{"flight_number":"TP1363","departure_iata":"LIS"}]}]}`))
		case "POST /v1/trips":
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":"b","name":"Porto","segments":[]}`))
		case "POST /v1/trips/import":
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":"c","name":"Imported","segments":[{"flight_number":"BA492"}]}`))
		case "DELETE /v1/trips/missing":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":{"code":"trip_not_found","message":"Trip not found"}}`))
		case "DELETE /v1/trips/a":
			w.WriteHeader(http.StatusNoContent)
		case "GET /v1/watches":
			_, _ = w.Write([]byte(`{"data":[{"id":"w1","flight_number":"TP1363","status":"active"}]}`))
		case "POST /v1/watches":
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":"w2","flight_number":"BA492","status":"active"}`))
		case "DELETE /v1/watches/w1":
			w.WriteHeader(http.StatusNoContent)
		case "POST /v1/billing/checkout":
			_, _ = w.Write([]byte(`{"id":"cs_1","url":"https://checkout.stripe.test/session"}`))
		default:
			w.WriteHeader(http.StatusTeapot)
		}
	}))
	t.Cleanup(server.Close)

	client, err := New(server.URL)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	ctx := context.Background()

	trips, err := client.ListTrips(ctx, "test-token")
	if err != nil {
		t.Fatalf("ListTrips() error = %v", err)
	}
	if len(trips) != 1 || trips[0].Name != "Lisbon" ||
		len(trips[0].Segments) != 1 || trips[0].Segments[0].FlightNumber != "TP1363" {
		t.Fatalf("ListTrips() = %#v", trips)
	}

	created, err := client.CreateTrip(ctx, "test-token", CreateTrip{Name: "Porto"})
	if err != nil {
		t.Fatalf("CreateTrip() error = %v", err)
	}
	if created.ID != "b" {
		t.Fatalf("CreateTrip() = %#v", created)
	}

	imported, err := client.ImportTrip(ctx, "test-token", "Booking text")
	if err != nil {
		t.Fatalf("ImportTrip() error = %v", err)
	}
	if imported.ID != "c" || len(imported.Segments) != 1 {
		t.Fatalf("ImportTrip() = %#v", imported)
	}

	if err := client.DeleteTrip(ctx, "test-token", "a"); err != nil {
		t.Fatalf("DeleteTrip() error = %v", err)
	}
	if err := client.DeleteTrip(ctx, "test-token", "missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("DeleteTrip(missing) error = %v, want ErrNotFound", err)
	}
	if _, err := client.ListTrips(ctx, "wrong-token"); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("ListTrips(bad token) error = %v, want ErrUnauthorized", err)
	}

	watches, err := client.ListWatches(ctx, "test-token")
	if err != nil {
		t.Fatalf("ListWatches() error = %v", err)
	}
	if len(watches) != 1 || watches[0].FlightNumber != "TP1363" {
		t.Fatalf("ListWatches() = %#v", watches)
	}
	createdWatch, err := client.CreateWatch(ctx, "test-token", CreateWatch{FlightNumber: "BA492"})
	if err != nil {
		t.Fatalf("CreateWatch() error = %v", err)
	}
	if createdWatch.ID != "w2" {
		t.Fatalf("CreateWatch() = %#v", createdWatch)
	}
	if err := client.DeleteWatch(ctx, "test-token", "w1"); err != nil {
		t.Fatalf("DeleteWatch() error = %v", err)
	}
	checkout, err := client.CreateCheckout(ctx, "test-token")
	if err != nil || checkout.URL == "" {
		t.Fatalf("CreateCheckout() = %#v, err = %v", checkout, err)
	}
}

func TestStreamEvents(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/events" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("event: ready\ndata: {\"status\":\"connected\"}\n\n"))
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		_, _ = w.Write([]byte("event: flight.delayed\ndata: {\"message\":\"TP1363 is delayed by 30 minutes\"}\n\n"))
	}))
	t.Cleanup(server.Close)

	client, err := New(server.URL)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	body, err := client.StreamEvents(context.Background(), "test-token")
	if err != nil {
		t.Fatalf("StreamEvents() error = %v", err)
	}
	t.Cleanup(func() { _ = body.Close() })
	data, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("read stream: %v", err)
	}
	if !strings.Contains(string(data), "event: ready") || !strings.Contains(string(data), "flight.delayed") {
		t.Fatalf("stream body = %q", string(data))
	}

	if _, err := client.StreamEvents(context.Background(), "wrong"); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("StreamEvents(bad token) error = %v, want ErrUnauthorized", err)
	}
}

func TestNewRejectsInvalidBaseURL(t *testing.T) {
	t.Parallel()
	for _, invalid := range []string{"", "not-a-url", "ftp://example.com"} {
		if _, err := New(invalid); err == nil {
			t.Fatalf("New(%q) expected an error", invalid)
		}
	}
}
