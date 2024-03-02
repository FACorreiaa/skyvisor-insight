package services

import (
	"context"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
	"github.com/FACorreiaa/Aviation-tracker/app/repository"
)

func (h *Service) RegisterNewAccount(ctx context.Context, form models.RegisterForm) (*repository.Token, error) {
	return h.accountRepo.RegisterNewAccount(ctx, form)
}

func (h *Service) Login(ctx context.Context, form models.LoginForm) (*repository.Token, error) {
	return h.accountRepo.Login(ctx, form)
}

func (h *Service) Logout(ctx context.Context, token repository.Token) error {
	return h.accountRepo.Logout(ctx, token)
}
