package trips

import (
	"fmt"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
	"github.com/a-h/templ"
)

func flightsPlaceholder() string {
	return "TP1363\nBA492"
}

func tripDeleteAction(id string) templ.SafeURL {
	return templ.SafeURL("/trips/" + id + "/delete")
}

func tripSubtitle(trip apiclient.Trip) string {
	segments := "No flights yet"
	if count := len(trip.Segments); count == 1 {
		segments = "1 flight"
	} else if count > 1 {
		segments = fmt.Sprintf("%d flights", count)
	}
	if trip.StartsAt.IsZero() {
		return segments
	}
	return segments + " · starts " + trip.StartsAt.Format("Mon, 2 Jan 2006")
}

func segmentRoute(segment apiclient.TripSegment) string {
	departure := segment.DepartureIATA
	if departure == "" {
		departure = "???"
	}
	arrival := segment.ArrivalIATA
	if arrival == "" {
		arrival = "???"
	}
	if departure == "???" && arrival == "???" {
		return "Route not set"
	}
	return departure + " → " + arrival
}

func segmentTimes(segment apiclient.TripSegment) string {
	if segment.DepartsAt.IsZero() {
		return "Times not set"
	}
	departs := segment.DepartsAt.Format("2 Jan 15:04")
	if segment.ArrivesAt.IsZero() {
		return departs
	}
	return departs + " – " + segment.ArrivesAt.Format("15:04") + " UTC"
}
