package services

import (
	"context"
	"math"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
)

func (h *Service) GetAllFlightsSum() (int, error) {
	total, err := h.flightRepo.GetAllFlightsSum(context.Background())
	pageSize := 15
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (h *Service) GetAllFlightsLocation() ([]models.LiveFlights, error) {
	lf, err := h.flightRepo.GetAllFlightsLocation(context.Background())
	if err != nil {
		HandleError(err, "Error flights details")
		return nil, err
	}

	return lf, nil
}

func (h *Service) GetFlightByID(ctx context.Context, flightNumber string) (models.LiveFlights, error) {
	return h.flightRepo.GetFlightByID(ctx, flightNumber)
}

func (h *Service) GetAllFlightsByStatus(ctx context.Context,
	page, pageSize int, orderBy, sortBy, flightNumber, flightStatus string) ([]models.LiveFlights, error) {

	return h.flightRepo.GetAllFlightsByStatus(ctx, page, pageSize, orderBy, sortBy, flightNumber, flightStatus)
}

func (h *Service) GetAllFlights(ctx context.Context,
	page, pageSize int, orderBy, sortBy string) ([]models.LiveFlights, error) {

	return h.flightRepo.GetAllFlights(ctx, page, pageSize, orderBy, sortBy)
}

func (h *Service) GetAllFlightsLocationsByStatus(ctx context.Context,
	flightStatus string) ([]models.LiveFlights, error) {
	return h.flightRepo.GetAllFlightsLocationsByStatus(ctx, flightStatus)
}

func (h *Service) GetLiveFlights(ctx context.Context,
	page, pageSize int, orderBy, sortBy string) ([]models.LiveFlights, error) {

	return h.flightRepo.GetLiveFlights(ctx, page, pageSize, orderBy, sortBy)
}

func (h *Service) GetLiveFlightsLocations(ctx context.Context) ([]models.LiveFlights, error) {
	return h.flightRepo.GetLiveFlightsLocations(ctx)
}

func (h *Service) GetFlightResumeByStatus(ctx context.Context, flightStatus string) (models.LiveFlightsResume, error) {
	return h.flightRepo.GetFlightResumeByStatus(ctx, flightStatus)
}

func (h *Service) GetFlightsResume(ctx context.Context) ([]models.LiveFlightsResume, error) {
	return h.flightRepo.GetFlightsResume(ctx)
}
