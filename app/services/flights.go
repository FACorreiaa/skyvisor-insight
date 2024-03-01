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

func (h *Service) GetAllFlightsPreview() ([]models.LiveFlights, error) {
	lf, err := h.flightRepo.GetAllFlightsPreview(context.Background())
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
	page, pageSize int, orderBy, sortBy, flightNumber string) ([]models.LiveFlights, error) {

	return h.flightRepo.GetAllFlights(ctx, page, pageSize, orderBy, sortBy, flightNumber)
}
