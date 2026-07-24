package flightui

import (
	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
	"github.com/FACorreiaa/Aviation-tracker/app/models"
)

// LiveFlightFromAPI maps a provider flight DTO into the web track model.
func LiveFlightFromAPI(f apiclient.Flight) models.LiveFlights {
	out := models.LiveFlights{
		FlightStatus: models.FlightStatus(f.Status),
	}
	out.Airline.Name = f.Airline
	out.Flight.Number = f.Number
	out.Departure.Iata = f.DepartureIATA
	out.Departure.Terminal = f.DepartureTerminal
	out.Departure.Gate = f.DepartureGate
	out.Departure.Delay = f.DepartureDelayMinutes
	out.Departure.Scheduled = f.ScheduledAt
	out.Departure.Estimated = f.EstimatedDeparture
	out.Departure.Actual = f.ActualDeparture
	out.Arrival.Iata = f.ArrivalIATA
	out.Arrival.Terminal = f.ArrivalTerminal
	out.Arrival.Gate = f.ArrivalGate
	out.Arrival.Baggage = f.ArrivalBaggage
	out.Arrival.Delay = f.ArrivalDelayMinutes
	out.Arrival.Scheduled = f.EstimatedArrival
	out.Arrival.Estimated = f.EstimatedArrival
	out.Arrival.Actual = f.ActualArrival
	out.Aircraft.AircraftRegistration = f.AircraftRegistration
	out.Aircraft.AircraftIata = f.AircraftIATA
	if f.Live != nil {
		out.Live.LiveLatitude = float32(f.Live.Latitude)
		out.Live.LiveLongitude = float32(f.Live.Longitude)
		out.Live.LiveAltitude = float32(f.Live.Altitude)
		out.Live.LiveSpeedHorizontal = float32(f.Live.Speed)
		out.Live.LiveIsGround = f.Live.IsGround
	}
	if f.Inbound != nil {
		out.Inbound = &models.TrackInbound{
			FlightNumber: f.Inbound.FlightNumber,
			LateRisk:     f.Inbound.LateRisk,
			Message:      f.Inbound.Message,
		}
	}
	if f.Freshness != nil {
		out.Freshness = &models.TrackFreshness{
			Status:     f.Freshness.Status,
			AgeSeconds: f.Freshness.AgeSeconds,
		}
	}
	EnrichFlightCoords(&out)
	return out
}

func MapAttrsFromAPIFlight(flightNumber string, f *apiclient.Flight) MapAttrs {
	if f == nil {
		return MapAttrs{FlightNumber: flightNumber}
	}
	live := LiveFlightFromAPI(*f)
	return MapAttrsFromFlight(flightNumber, &live)
}
