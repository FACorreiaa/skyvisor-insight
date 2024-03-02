package services

import (
	"log"

	"github.com/FACorreiaa/Aviation-tracker/app/repository"
)

type Service struct {
	airlineRepo  *repository.AirlineRepository
	airportRepo  *repository.AirportRepository
	locationRepo *repository.LocationsRepository
	flightRepo   *repository.FlightsRepository
	accountRepo  *repository.AccountRepository
}

func HandleError(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
}

func NewService(
	airlineRepo *repository.AirlineRepository,
	airportRepo *repository.AirportRepository,
	locationRepo *repository.LocationsRepository,
	flightRepo *repository.FlightsRepository,
	accountRepo *repository.AccountRepository) *Service {

	return &Service{
		airlineRepo:  airlineRepo,
		airportRepo:  airportRepo,
		locationRepo: locationRepo,
		flightRepo:   flightRepo,
		accountRepo:  accountRepo,
	}
}
