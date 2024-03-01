package services

import (
	"context"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
	"github.com/FACorreiaa/Aviation-tracker/app/session"
)

func (h *Service) RegisterNewAccount(ctx context.Context, form models.RegisterForm) (*session.Token, error) {
	return h.accountRepo.RegisterNewAccount(ctx, form)
}

func (h *Service) Login(ctx context.Context, form session.LoginForm) (*session.Token, error) {
	return h.accountRepo.Login(ctx, form)
}

func (h *Service) Logout(ctx context.Context, token session.Token) error {
	return h.accountRepo.Logout(ctx, token)
}
