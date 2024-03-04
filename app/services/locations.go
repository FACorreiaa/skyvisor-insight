package services

import (
	"math"

	"context"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
)

// city

func (h *Service) GetCityLocations() ([]models.City, error) {
	c, err := h.locationRepo.GetCityLocation(context.Background())
	if err != nil {
		HandleError(err, "Error fetching cities")
		return nil, err
	}

	return c, nil
}

func (h *Service) GetAllCities() (int, error) {
	total, err := h.locationRepo.GetCitySum(context.Background())
	pageSize := 30
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		HandleError(err, "Error fetching sum of cities")
		return 0, err
	}
	return lastPage, nil
}

func (h *Service) GetCity(ctx context.Context, page, pageSize int,
	orderBy string, sortBy string, name string) ([]models.City, error) {

	return h.locationRepo.GetCity(ctx, page, pageSize, orderBy, sortBy, name)
}

func (h *Service) GetCityByID(ctx context.Context, cityID int) (models.City, error) {
	return h.locationRepo.GetCityByID(ctx, cityID)
}

// Country

func (h *Service) GetCountryLocations() ([]models.Country, error) {
	c, err := h.locationRepo.GetCountryLocation(context.Background())
	if err != nil {
		HandleError(err, "Error fetching locations")
		return nil, err
	}

	return c, nil
}

func (h *Service) GetAllCountries() (int, error) {
	total, err := h.locationRepo.GetCountrySum(context.Background())
	pageSize := 30
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if err != nil {
		return 0, err
	}
	return lastPage, nil
}

func (h *Service) GetCountry(ctx context.Context, page, pageSize int,
	orderBy, sortBy, name string) ([]models.Country, error) {

	return h.locationRepo.GetCountry(ctx, page, pageSize, orderBy, sortBy, name)
}

func (h *Service) GetCountryByName(ctx context.Context, name string) (models.Country, error) {
	return h.locationRepo.GetCountryByName(ctx, name)
}
