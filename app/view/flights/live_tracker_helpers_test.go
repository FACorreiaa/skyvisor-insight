package flights

import (
	"testing"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
)

func TestLiveFlightPresentation(t *testing.T) {
	t.Parallel()
	delay := 27
	tests := []struct {
		name          string
		flight        apiclient.Flight
		wantDelay     string
		wantFreshness string
		wantMovement  string
	}{
		{
			name: "provider freshness and movement",
			flight: apiclient.Flight{
				DepartureDelayMinutes: &delay,
				Freshness:             &apiclient.DataFreshness{Status: "fresh"},
				Live:                  &apiclient.FlightLive{Altitude: 31000, Speed: 442},
			},
			wantDelay: "+27m", wantFreshness: "fresh", wantMovement: "31000 ft · 442 kt",
		},
		{
			name:      "unknown record",
			flight:    apiclient.Flight{},
			wantDelay: "—", wantFreshness: "unknown", wantMovement: "No position",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := liveDelay(test.flight); got != test.wantDelay {
				t.Fatalf("delay = %q, want %q", got, test.wantDelay)
			}
			if got := liveFreshness(test.flight); got != test.wantFreshness {
				t.Fatalf("freshness = %q, want %q", got, test.wantFreshness)
			}
			if got := movementLabel(test.flight); got != test.wantMovement {
				t.Fatalf("movement = %q, want %q", got, test.wantMovement)
			}
		})
	}
}

func TestLiveFreshnessFallsBackToTimestamp(t *testing.T) {
	t.Parallel()
	flight := apiclient.Flight{UpdatedAt: time.Now().Add(-3 * time.Minute)}
	if got := liveFreshness(flight); got != "stale" {
		t.Fatalf("freshness = %q, want stale", got)
	}
}
