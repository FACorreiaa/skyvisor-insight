package services

import (
	"log"

	"github.com/FACorreiaa/Aviation-tracker/app/repository"
)

type AirlineService struct {
	repo *repository.AirlineRepository
}

func HandleError(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
}

func NewAirlineService(repo *repository.AirlineRepository) *AirlineService {
	return &AirlineService{repo: repo}
}
