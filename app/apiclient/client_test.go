package apiclient

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
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
			_, _ = w.Write([]byte(`{"data":[{"id":"a","name":"Lisbon","flights":["TP1363"]}]}`))
		case "POST /v1/trips":
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":"b","name":"Porto","flights":[]}`))
		case "DELETE /v1/trips/missing":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":{"code":"trip_not_found","message":"Trip not found"}}`))
		case "DELETE /v1/trips/a":
			w.WriteHeader(http.StatusNoContent)
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
	if len(trips) != 1 || trips[0].Name != "Lisbon" {
		t.Fatalf("ListTrips() = %#v", trips)
	}

	created, err := client.CreateTrip(ctx, "test-token", CreateTrip{Name: "Porto"})
	if err != nil {
		t.Fatalf("CreateTrip() error = %v", err)
	}
	if created.ID != "b" {
		t.Fatalf("CreateTrip() = %#v", created)
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
}

func TestNewRejectsInvalidBaseURL(t *testing.T) {
	t.Parallel()
	for _, invalid := range []string{"", "not-a-url", "ftp://example.com"} {
		if _, err := New(invalid); err == nil {
			t.Fatalf("New(%q) expected an error", invalid)
		}
	}
}
