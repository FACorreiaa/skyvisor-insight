package flights

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
)

func limitValue(n int) string {
	if n <= 0 {
		return "50"
	}
	return strconv.Itoa(n)
}

func formatCoord(v float64) string {
	return strconv.FormatFloat(v, 'f', 3, 64)
}

func formatAlt(v float64) string {
	return strconv.FormatFloat(v, 'f', 0, 64) + " ft"
}

func formatSpeed(v float64) string {
	return strconv.FormatFloat(v, 'f', 0, 64) + " kt"
}

func trackerCount(n int) string { return strconv.Itoa(n) }

func statusCount(items []apiclient.Flight, status string) int {
	count := 0
	for _, item := range items {
		if strings.EqualFold(item.Status, status) {
			count++
		}
	}
	return count
}

func livePositionCount(items []apiclient.Flight) int {
	count := 0
	for _, item := range items {
		if item.Live != nil {
			count++
		}
	}
	return count
}

func delayedFlightCount(items []apiclient.Flight) int {
	count := 0
	for _, item := range items {
		if maxLiveDelay(item) >= 15 {
			count++
		}
	}
	return count
}

func freshnessCount(items []apiclient.Flight, status string) int {
	count := 0
	for _, item := range items {
		if liveFreshness(item) == status {
			count++
		}
	}
	return count
}

func routeLabelLive(f apiclient.Flight) string {
	if f.DepartureIATA == "" && f.ArrivalIATA == "" {
		return "Pending"
	}
	return valueOrLive(f.DepartureIATA, "?") + " → " + valueOrLive(f.ArrivalIATA, "?")
}

func valueOrLive(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func maxLiveDelay(f apiclient.Flight) int {
	value := 0
	if f.DepartureDelayMinutes != nil && *f.DepartureDelayMinutes > value {
		value = *f.DepartureDelayMinutes
	}
	if f.ArrivalDelayMinutes != nil && *f.ArrivalDelayMinutes > value {
		value = *f.ArrivalDelayMinutes
	}
	return value
}

func liveDelay(f apiclient.Flight) string {
	if maxLiveDelay(f) == 0 {
		return "—"
	}
	return fmt.Sprintf("+%dm", maxLiveDelay(f))
}

func movementLabel(f apiclient.Flight) string {
	if f.Live == nil {
		return "No position"
	}
	return formatAlt(f.Live.Altitude) + " · " + formatSpeed(f.Live.Speed)
}

func liveGates(f apiclient.Flight) string {
	if f.DepartureGate == "" && f.ArrivalGate == "" {
		return "—"
	}
	return valueOrLive(f.DepartureGate, "?") + " → " + valueOrLive(f.ArrivalGate, "?")
}

func liveFreshness(f apiclient.Flight) string {
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

func liveFreshnessClass(f apiclient.Flight) string {
	switch liveFreshness(f) {
	case "fresh":
		return "bg-emerald-500/10 text-emerald-700 dark:text-emerald-400"
	case "stale":
		return "bg-amber-500/10 text-amber-700 dark:text-amber-400"
	default:
		return "bg-muted text-muted-foreground"
	}
}

func liveStatusClass(status string) string {
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

func liveUpdated(f apiclient.Flight) string {
	if f.UpdatedAt.IsZero() {
		return "not observed"
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

// Ensure apiclient stays imported if helpers-only file is tree-shaken oddly.
var _ = apiclient.LiveFlightQuery{}
