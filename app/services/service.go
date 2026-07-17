package services

import (
	"log"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
	"github.com/FACorreiaa/Aviation-tracker/app/auth"
	"github.com/FACorreiaa/Aviation-tracker/app/repository"
)

type Service struct {
	airlineRepo  *repository.AirlineRepository
	airportRepo  *repository.AirportRepository
	locationRepo *repository.LocationsRepository
	flightRepo   *repository.FlightsRepository
	accountRepo  *repository.AccountRepository
	oidc         *auth.Client
	api          *apiclient.Client
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
	accountRepo *repository.AccountRepository,
	oidc *auth.Client,
	api *apiclient.Client) *Service {

	return &Service{
		airlineRepo:  airlineRepo,
		airportRepo:  airportRepo,
		locationRepo: locationRepo,
		flightRepo:   flightRepo,
		accountRepo:  accountRepo,
		oidc:         oidc,
		api:          api,
	}
}

// API exposes the skyvisor-api client. Nil when SKYVISOR_API_URL is not set.
func (h *Service) API() *apiclient.Client {
	return h.api
}
