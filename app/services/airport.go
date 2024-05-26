package services

import (
	"math"

	"context"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
)

func (h *Service) GetAirports(ctx context.Context,
	page, pageSize int, orderBy string, sortBy string) ([]models.Airport, error) {

	return h.airportRepo.GetAirports(ctx, page, pageSize, orderBy, sortBy)
}
func (h *Service) GetAirportsLocation() ([]models.Airport, error) {
	a, err := h.airportRepo.GetAirportsLocation(context.Background())
	if err != nil {
		HandleError(err, "Error fetching airports locations")
		return nil, err
	}

	return a, nil
}

func (h *Service) GetAllAirports() (int, error) {
	total, err := h.airportRepo.GetSum(context.Background())
	pageSize := 20
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (h *Service) GetAirportByName(ctx context.Context, page, pageSize int,
	orderBy, sortBy, airportName, countryName, gmt string) ([]models.Airport, error) {

	return h.airportRepo.GetAirportByName(ctx, page, pageSize, orderBy, sortBy, airportName, countryName, gmt)
}

func (h *Service) GetAirportByID(ctx context.Context, id int) (models.Airport, error) {
	return h.airportRepo.GetAirportByID(ctx, id)
}
