package airport

import (
	"testing"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
)

func TestAirportBoardOperationalCounts(t *testing.T) {
	t.Parallel()
	delay := 25
	now := time.Now().UTC()
	board := apiclient.AirportBoard{
		Arrivals:   []apiclient.Flight{{Status: "landed", ArrivalDelayMinutes: &delay, ArrivalGate: "18", UpdatedAt: now}},
		Departures: []apiclient.Flight{{Status: "cancelled", DepartureGate: "A4", Freshness: &apiclient.DataFreshness{Status: "fresh"}}},
	}
	tests := []struct{ name, got, want string }{
		{name: "delayed", got: boardDelayedCount(board), want: "1"},
		{name: "cancelled", got: boardStatusCount(board, "cancelled"), want: "1"},
		{name: "gates", got: boardGateCount(board), want: "2"},
		{name: "fresh", got: boardFreshCount(board), want: "2"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if test.got != test.want {
				t.Fatalf("got %s, want %s", test.got, test.want)
			}
		})
	}
}

func TestBoardEvidenceUsesProvenance(t *testing.T) {
	t.Parallel()
	flight := apiclient.Flight{Provenance: &apiclient.SourceProvenance{Provider: "aviationstack"}, Freshness: &apiclient.DataFreshness{Status: "stale"}}
	if got := boardEvidence(flight); got != "aviationstack · stale" {
		t.Fatalf("evidence = %q", got)
	}
}
