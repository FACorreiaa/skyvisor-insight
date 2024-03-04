package services

import (
	"context"
	"math"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
)

func (h *Service) GetAirlinesLocation() ([]models.Airline, error) {
	al, err := h.airlineRepo.GetAirlinesLocations(context.Background())
	if err != nil {
		HandleError(err, "Error fetching locations")
		return nil, err
	}

	return al, nil
}

func (h *Service) GetAllAirline() (int, error) {
	total, err := h.airlineRepo.GetAirlineSum(context.Background())
	pageSize := 30
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (h *Service) GetAirlines(ctx context.Context, page,
	pageSize int, orderBy, sortBy, name, callSign, hubCode, countryName string) ([]models.Airline, error) {

	return h.airlineRepo.GetAirlines(ctx, page, pageSize, orderBy, sortBy, name, callSign, hubCode, countryName)
}

func (h *Service) GetAirlineByName(ctx context.Context, airlineName string) (models.Airline, error) {
	return h.airlineRepo.GetAirlineByName(ctx, airlineName)
}

// Aircraft

func (h *Service) GetAllAircraft() (int, error) {
	total, err := h.airlineRepo.GetAircraftSum(context.Background())
	pageSize := 25
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (h *Service) GetAircraft(ctx context.Context, page, pageSize int, name string,
	orderBy string, sortBy string) ([]models.Aircraft, error) {

	return h.airlineRepo.GetAircraft(ctx, page, pageSize, name, orderBy, sortBy)
}

// Airplane

func (h *Service) GetAllAirplanes() (int, error) {
	total, err := h.airlineRepo.GetAirplaneSum(context.Background())
	pageSize := 25
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (h *Service) GetAirplanes(ctx context.Context, page, pageSize int,
	orderBy string, sortBy string, name string) ([]models.Airplane, error) {
	return h.airlineRepo.GetAirplanes(ctx, page, pageSize, orderBy, sortBy, name)
}

// tax

func (h *Service) GetTax(ctx context.Context, page, pageSize int,
	orderBy, sortBy, taxName, countryName, airlineName string) ([]models.Tax, error) {

	return h.airlineRepo.GetTax(ctx, page, pageSize, orderBy, sortBy, taxName, countryName, airlineName)
}
func (h *Service) GetSum() (int, error) {
	total, err := h.airlineRepo.GetTaxSum(context.Background())
	pageSize := 30
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}
