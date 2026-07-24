package airport

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
)

func boardLimit(n int) string {
	if n <= 0 {
		return "40"
	}
	return strconv.Itoa(n)
}

func boardHXURL(iata string) string {
	return "/airports/board"
}

func flightCount(items []apiclient.Flight) string {
	return strconv.Itoa(len(items))
}

var _ = apiclient.AirportBoard{}

func boardFlights(board apiclient.AirportBoard) []apiclient.Flight {
	return append(append([]apiclient.Flight(nil), board.Arrivals...), board.Departures...)
}

func boardStatusCount(board apiclient.AirportBoard, status string) string {
	count := 0
	for _, flight := range boardFlights(board) {
		if strings.EqualFold(flight.Status, status) {
			count++
		}
	}
	return strconv.Itoa(count)
}

func boardDelayedCount(board apiclient.AirportBoard) string {
	count := 0
	for _, flight := range boardFlights(board) {
		if boardMaxDelay(flight) >= 15 {
			count++
		}
	}
	return strconv.Itoa(count)
}

func boardGateCount(board apiclient.AirportBoard) string {
	count := 0
	for _, flight := range boardFlights(board) {
		if flight.DepartureGate != "" || flight.ArrivalGate != "" {
			count++
		}
	}
	return strconv.Itoa(count)
}

func boardFreshCount(board apiclient.AirportBoard) string {
	count := 0
	for _, flight := range boardFlights(board) {
		if boardFreshness(flight) == "fresh" {
			count++
		}
	}
	return strconv.Itoa(count)
}

func boardTime(f apiclient.Flight, direction string) string {
	value := f.EstimatedArrival
	if direction == "departure" {
		value = f.EstimatedDeparture
		if value.IsZero() {
			value = f.ScheduledAt
		}
	}
	if value.IsZero() {
		return "—"
	}
	return value.Local().Format("15:04")
}

func boardRoute(f apiclient.Flight) string {
	if f.DepartureIATA == "" && f.ArrivalIATA == "" {
		return "Pending"
	}
	return boardValue(f.DepartureIATA) + " → " + boardValue(f.ArrivalIATA)
}

func boardValue(value string) string {
	if strings.TrimSpace(value) == "" {
		return "?"
	}
	return value
}

func boardMaxDelay(f apiclient.Flight) int {
	value := 0
	if f.DepartureDelayMinutes != nil && *f.DepartureDelayMinutes > value {
		value = *f.DepartureDelayMinutes
	}
	if f.ArrivalDelayMinutes != nil && *f.ArrivalDelayMinutes > value {
		value = *f.ArrivalDelayMinutes
	}
	return value
}

func boardDelay(f apiclient.Flight, direction string) string {
	var value *int
	if direction == "departure" {
		value = f.DepartureDelayMinutes
	} else {
		value = f.ArrivalDelayMinutes
	}
	if value == nil || *value <= 0 {
		return "—"
	}
	return fmt.Sprintf("+%dm", *value)
}

func boardPosition(f apiclient.Flight, direction string) string {
	if direction == "departure" {
		return "T" + boardValue(f.DepartureTerminal) + " · G" + boardValue(f.DepartureGate)
	}
	parts := "T" + boardValue(f.ArrivalTerminal) + " · G" + boardValue(f.ArrivalGate)
	if f.ArrivalBaggage != "" {
		parts += " · B" + f.ArrivalBaggage
	}
	return parts
}

func boardStatusClass(status string) string {
	switch strings.ToLower(status) {
	case "cancelled", "incident", "diverted":
		return "bg-destructive/10 text-destructive"
	case "active":
		return "bg-primary/10 text-primary"
	case "landed":
		return "bg-emerald-500/10 text-emerald-700 dark:text-emerald-400"
	default:
		return "bg-muted text-muted-foreground"
	}
}

func boardFreshness(f apiclient.Flight) string {
	if f.Freshness != nil && f.Freshness.Status != "" {
		return f.Freshness.Status
	}
	if f.UpdatedAt.IsZero() {
		return "unknown"
	}
	if time.Since(f.UpdatedAt) > 2*time.Minute {
		return "stale"
	}
	return "fresh"
}

func boardEvidence(f apiclient.Flight) string {
	provider := "provider"
	if f.Provenance != nil && f.Provenance.Provider != "" {
		provider = f.Provenance.Provider
	}
	return provider + " · " + boardFreshness(f)
}

func boardGenerated(value time.Time) string {
	d := time.Since(value)
	if d < 0 {
		d = 0
	}
	if d < time.Minute {
		return "just now"
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	}
	return value.Local().Format("15:04")
}
