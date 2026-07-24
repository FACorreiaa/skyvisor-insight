package flightui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
)

type MapAttrs struct {
	Latitude      float64
	Longitude     float64
	Zoom          float64
	DepLat        float64
	DepLon        float64
	ArrLat        float64
	ArrLon        float64
	LiveLat       float64
	LiveLon       float64
	FlightNumber  string
	HasRoute      bool
	HasLive       bool
	ProgressPct   int
}

func enrichFlightCoords(flight *models.LiveFlights) {
	if flight == nil {
		return
	}
	if flight.DepartureLatitude == "" {
		if coord, ok := airportCoord(flight.Departure.Iata); ok {
			flight.DepartureLatitude = coord[0]
			flight.DepartureLongitude = coord[1]
		}
	}
	if flight.ArrivalLatitude == "" {
		if coord, ok := airportCoord(flight.Arrival.Iata); ok {
			flight.ArrivalLatitude = coord[0]
			flight.ArrivalLongitude = coord[1]
		}
	}
}

func airportCoord(iata string) ([2]string, bool) {
	table := map[string][2]string{
		"LIS": {"38.7813", "-9.1359"},
		"LHR": {"51.4700", "-0.4543"},
		"JFK": {"40.6413", "-73.7781"},
		"LAX": {"33.9425", "-118.4085"},
		"CDG": {"49.0097", "2.5479"},
		"FRA": {"50.0379", "8.5622"},
		"DXB": {"25.2532", "55.3657"},
		"SIN": {"1.3644", "103.9915"},
		"MAD": {"40.4983", "-3.5676"},
		"FCO": {"41.8003", "12.2389"},
	}
	coord, ok := table[strings.ToUpper(strings.TrimSpace(iata))]
	return coord, ok
}

// EnrichFlightCoords fills missing airport coordinates from a small IATA table.
func EnrichFlightCoords(flight *models.LiveFlights) {
	enrichFlightCoords(flight)
}

func MapAttrsFromFlight(flightNumber string, flight *models.LiveFlights) MapAttrs {
	if flight != nil {
		clone := *flight
		enrichFlightCoords(&clone)
		flight = &clone
	}
	attrs := MapAttrs{
		Latitude:     28,
		Longitude:    -12,
		Zoom:         1.25,
		FlightNumber: strings.ToUpper(strings.TrimSpace(flightNumber)),
	}
	if flight == nil {
		return attrs
	}
	if attrs.FlightNumber == "" {
		attrs.FlightNumber = strings.ToUpper(strings.TrimSpace(flight.Flight.Number))
	}

	depLat, depLon, depOK := parseCoord(flight.DepartureLatitude, flight.DepartureLongitude)
	arrLat, arrLon, arrOK := parseCoord(flight.ArrivalLatitude, flight.ArrivalLongitude)
	liveLat := float64(flight.Live.LiveLatitude)
	liveLon := float64(flight.Live.LiveLongitude)
	liveOK := liveLat != 0 || liveLon != 0

	if depOK && arrOK {
		attrs.DepLat, attrs.DepLon = depLat, depLon
		attrs.ArrLat, attrs.ArrLon = arrLat, arrLon
		attrs.HasRoute = true
		attrs.Latitude = (depLat + arrLat) / 2
		attrs.Longitude = (depLon + arrLon) / 2
		attrs.Zoom = 3.2
	}
	if liveOK {
		attrs.LiveLat, attrs.LiveLon = liveLat, liveLon
		attrs.HasLive = true
		attrs.Latitude = liveLat
		attrs.Longitude = liveLon
		attrs.Zoom = 4.5
	}
	attrs.ProgressPct = progressPercent(flight)
	return attrs
}

func progressPercent(flight *models.LiveFlights) int {
	if flight == nil {
		return 0
	}
	dep := firstTime(flight.Departure.Actual, flight.Departure.Estimated, flight.Departure.Scheduled)
	arr := firstTime(flight.Arrival.Actual, flight.Arrival.Estimated, flight.Arrival.Scheduled)
	if dep.IsZero() || arr.IsZero() || !arr.After(dep) {
		return 0
	}
	now := time.Now()
	if now.Before(dep) {
		return 0
	}
	if now.After(arr) {
		return 100
	}
	pct := float64(now.Sub(dep)) / float64(arr.Sub(dep)) * 100
	return int(math.Round(math.Min(100, math.Max(0, pct))))
}

func firstTime(values ...time.Time) time.Time {
	for _, value := range values {
		if !value.IsZero() {
			return value
		}
	}
	return time.Time{}
}

func parseCoord(latRaw, lonRaw string) (float64, float64, bool) {
	lat, latErr := strconv.ParseFloat(strings.TrimSpace(latRaw), 64)
	lon, lonErr := strconv.ParseFloat(strings.TrimSpace(lonRaw), 64)
	if latErr != nil || lonErr != nil {
		return 0, 0, false
	}
	return lat, lon, true
}

func FormatCoord(value float64) string {
	if value == 0 {
		return ""
	}
	return fmt.Sprintf("%.6f", value)
}

func FormatZoom(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

func StatusLabel(flight *models.LiveFlights) string {
	if flight == nil {
		return ""
	}
	return strings.ReplaceAll(string(flight.FlightStatus), "_", " ")
}

func DelayMinutes(flight *models.LiveFlights) int {
	if flight == nil || flight.Departure.Delay == nil {
		return 0
	}
	return *flight.Departure.Delay
}

func DisplayFlightNumber(query string, flight *models.LiveFlights) string {
	if flight != nil && flight.Flight.Number != "" {
		return flight.Flight.Number
	}
	return strings.ToUpper(strings.TrimSpace(query))
}

func AirlineName(flight *models.LiveFlights) string {
	if flight == nil {
		return ""
	}
	return flight.Airline.Name
}

func DepartureIATA(flight *models.LiveFlights) string {
	if flight == nil {
		return ""
	}
	return fallback(flight.Departure.Iata, flight.Departure.CityCode)
}

func ArrivalIATA(flight *models.LiveFlights) string {
	if flight == nil {
		return ""
	}
	return fallback(flight.Arrival.Iata, flight.Arrival.CityCode)
}

func DepartureTime(flight *models.LiveFlights) string {
	if flight == nil {
		return "—"
	}
	return displayTime(flight.Departure.Estimated, flight.Departure.Scheduled)
}

func ArrivalTime(flight *models.LiveFlights) string {
	if flight == nil {
		return "—"
	}
	return displayTime(flight.Arrival.Estimated, flight.Arrival.Scheduled)
}

func displayTime(preferred, fallbackTime time.Time) string {
	if !preferred.IsZero() {
		return preferred.Format("15:04")
	}
	if !fallbackTime.IsZero() {
		return fallbackTime.Format("15:04")
	}
	return "—"
}

func fallback(value, fallbackValue string) string {
	if strings.TrimSpace(value) == "" {
		return fallbackValue
	}
	return value
}

func GateText(flight *models.LiveFlights) string {
	if flight == nil {
		return ""
	}
	dep := fmt.Sprint(flight.Departure.Gate)
	arr := fmt.Sprint(flight.Arrival.Gate)
	if dep == "" && arr == "" {
		return ""
	}
	if dep == "" {
		dep = "—"
	}
	if arr == "" {
		arr = "—"
	}
	return fmt.Sprintf("Gate %s → %s", dep, arr)
}

func AircraftText(flight *models.LiveFlights) string {
	if flight == nil || flight.Aircraft.AircraftRegistration == "" {
		return ""
	}
	return flight.Aircraft.AircraftRegistration
}

func TerminalText(flight *models.LiveFlights, leg string) string {
	if flight == nil {
		return ""
	}
	if leg == "dep" {
		return strings.TrimSpace(flight.Departure.Terminal)
	}
	return strings.TrimSpace(fmt.Sprint(flight.Arrival.Terminal))
}

func BaggageText(flight *models.LiveFlights) string {
	if flight == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(flight.Arrival.Baggage))
}

func DelayStory(flight *models.LiveFlights) string {
	if flight == nil {
		return ""
	}
	dep := 0
	if flight.Departure.Delay != nil {
		dep = *flight.Departure.Delay
	}
	arr := 0
	if flight.Arrival.Delay != nil {
		arr = *flight.Arrival.Delay
	}
	if dep == 0 && arr == 0 {
		return ""
	}
	parts := make([]string, 0, 2)
	if dep > 0 {
		parts = append(parts, fmt.Sprintf("+%dm dep", dep))
	}
	if arr > 0 {
		parts = append(parts, fmt.Sprintf("+%dm arr", arr))
	}
	return strings.Join(parts, " · ")
}

func StaleChipText(flight *models.LiveFlights) string {
	if flight == nil || flight.Freshness == nil {
		return ""
	}
	switch flight.Freshness.Status {
	case "stale":
		if flight.Freshness.AgeSeconds > 0 {
			mins := flight.Freshness.AgeSeconds / 60
			if mins < 1 {
				mins = 1
			}
			return fmt.Sprintf("Stale %dm", mins)
		}
		return "Stale data"
	case "unknown":
		return "Freshness unknown"
	default:
		return ""
	}
}
