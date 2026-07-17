package trips

import (
	"testing"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
)

func TestBuildTimeline(t *testing.T) {
	t.Parallel()
	aug10 := time.Date(2026, time.August, 10, 15, 0, 0, 0, time.UTC)
	aug11 := time.Date(2026, time.August, 11, 9, 0, 0, 0, time.UTC)
	trips := []apiclient.Trip{
		{ID: "t1", Name: "Trip One", Segments: []apiclient.TripSegment{
			{FlightNumber: "BA492", DepartsAt: aug11},
			{FlightNumber: "TP1363", DepartsAt: aug10},
			{FlightNumber: "XX999"}, // undated
		}},
	}

	days, undated := BuildTimeline(trips, time.UTC)
	if len(days) != 2 {
		t.Fatalf("day count = %d, want 2", len(days))
	}
	if days[0].Entries[0].Segment.FlightNumber != "TP1363" {
		t.Fatalf("first day flight = %q, want TP1363 (chronological)", days[0].Entries[0].Segment.FlightNumber)
	}
	if days[1].Entries[0].Segment.FlightNumber != "BA492" {
		t.Fatalf("second day flight = %q, want BA492", days[1].Entries[0].Segment.FlightNumber)
	}
	if len(undated) != 1 || undated[0].Segment.FlightNumber != "XX999" {
		t.Fatalf("undated = %#v", undated)
	}
}

func TestBuildTimelineLocationGrouping(t *testing.T) {
	t.Parallel()
	// 23:30 UTC on Aug 10 is 00:30 on Aug 11 in Lisbon (summer, UTC+1),
	// so the two segments must land on different local days.
	lisbon, err := time.LoadLocation("Europe/Lisbon")
	if err != nil {
		t.Skipf("tz data unavailable: %v", err)
	}
	trips := []apiclient.Trip{{ID: "t1", Name: "Trip", Segments: []apiclient.TripSegment{
		{FlightNumber: "AA1", DepartsAt: time.Date(2026, time.August, 10, 12, 0, 0, 0, time.UTC)},
		{FlightNumber: "AA2", DepartsAt: time.Date(2026, time.August, 10, 23, 30, 0, 0, time.UTC)},
	}}}
	days, _ := BuildTimeline(trips, lisbon)
	if len(days) != 2 {
		t.Fatalf("day count in Lisbon = %d, want 2", len(days))
	}
}
