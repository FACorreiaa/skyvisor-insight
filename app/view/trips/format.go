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

func segmentDelay(segment apiclient.TripSegment) string {
	if segment.Live == nil || segment.Live.DepartureDelayMinutes == nil {
		return ""
	}
	return fmt.Sprintf("+%d min", *segment.Live.DepartureDelayMinutes)
}

func riskLabel(risk string) string {
	switch risk {
	case "ok":
		return "Connection OK"
	case "tight":
		return "Tight"
	case "miss_likely":
		return "Miss likely"
	case "unknown":
		return "Unknown"
	default:
		return risk
	}
}

func riskBadgeClass(risk string) string {
	switch risk {
	case "miss_likely":
		return "border-destructive/30 bg-destructive/10 text-destructive"
	case "tight":
		return "border-amber-500/30 bg-amber-500/10 text-amber-700 dark:text-amber-300"
	case "ok":
		return "border-emerald-500/30 bg-emerald-500/10 text-emerald-700 dark:text-emerald-300"
	default:
		return ""
	}
}

func connectionSummary(connection apiclient.ConnectionRisk) string {
	if connection.LayoverMinutes > 0 {
		return fmt.Sprintf("%d min layover · %s", connection.LayoverMinutes, connection.Reason)
	}
	return connection.Reason
}

func AutoWatchMessage(result *apiclient.AutoWatchResult) string {
	if result == nil {
		return ""
	}
	if len(result.Started) == 0 && len(result.Skipped) == 0 {
		return ""
	}
	msg := ""
	if len(result.Started) > 0 {
		msg = fmt.Sprintf("Watching %s.", joinFlights(result.Started))
	}
	for _, skipped := range result.Skipped {
		if skipped.Reason == "watch_limit_reached" {
			if msg != "" {
				msg += " "
			}
			msg += "Free plan watches one flight — upgrade to Pro for the rest."
			break
		}
	}
	return msg
}

func joinFlights(flights []string) string {
	if len(flights) == 0 {
		return ""
	}
	if len(flights) == 1 {
		return flights[0]
	}
	return fmt.Sprintf("%s and %d more", flights[0], len(flights)-1)
}
