package logistics

import (
	"fmt"
	"strings"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
)

func flightCount(items []apiclient.Flight) string {
	return fmt.Sprintf("%d", len(items))
}

func disruptionKindCount(items []apiclient.LogisticsDisruption, kind string) string {
	count := 0
	for _, item := range items {
		if item.Kind == kind {
			count++
		}
	}
	return fmt.Sprintf("%d", count)
}

func flightDelay(f apiclient.Flight) string {
	delay := 0
	if f.DepartureDelayMinutes != nil && *f.DepartureDelayMinutes > delay {
		delay = *f.DepartureDelayMinutes
	}
	if f.ArrivalDelayMinutes != nil && *f.ArrivalDelayMinutes > delay {
		delay = *f.ArrivalDelayMinutes
	}
	if delay == 0 {
		return "—"
	}
	return fmt.Sprintf("+%dm", delay)
}

func flightGates(f apiclient.Flight) string {
	if f.DepartureGate == "" && f.ArrivalGate == "" {
		return "—"
	}
	from, to := strings.TrimSpace(f.DepartureGate), strings.TrimSpace(f.ArrivalGate)
	if from == "" {
		from = "?"
	}
	if to == "" {
		to = "?"
	}
	return from + " → " + to
}

func flightUpdated(f apiclient.Flight) string {
	if f.UpdatedAt.IsZero() {
		return "unknown"
	}
	d := time.Since(f.UpdatedAt)
	if d < 0 {
		d = 0
	}
	if d < time.Minute {
		return "just now"
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	}
	return f.UpdatedAt.Local().Format("15:04")
}

func disruptionCount(items []apiclient.LogisticsDisruption) string {
	return fmt.Sprintf("%d", len(items))
}

func seatLabel(team *apiclient.Team) string {
	if team == nil {
		return ""
	}
	return fmt.Sprintf("%d / %d", team.MemberCount, team.SeatLimit)
}
