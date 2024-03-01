package services

import (
	"context"
	"math"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
)

func (h *AirlineService) GetAirlinesLocationService() ([]models.Airline, error) {
	al, err := h.repo.GetAirlinesLocations(context.Background())
	if err != nil {
		HandleError(err, "Error fetching locations")
		return nil, err
	}

	return al, nil
}

func (h *AirlineService) GetAllAirlineService() (int, error) {
	total, err := h.repo.GetAirlineSum(context.Background())
	pageSize := 30
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (h *AirlineService) GetAirlines(ctx context.Context, page,
	pageSize int, orderBy, sortBy, name string) ([]models.Airline, error) {

	return h.repo.GetAirlines(ctx, page, pageSize, orderBy, sortBy, name)
}

func (h *AirlineService) GetAirlineByName(ctx context.Context, airlineName string) (models.Airline, error) {
	return h.repo.GetAirlineByName(ctx, airlineName)
}
