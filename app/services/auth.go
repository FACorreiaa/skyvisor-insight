package services

import (
	"context"
	"errors"

	"github.com/FACorreiaa/Aviation-tracker/app/auth"
	"github.com/FACorreiaa/Aviation-tracker/app/repository"
	"golang.org/x/oauth2"
)

// CompleteOIDCLogin turns a verified OIDC principal into a local session:
// upsert the user row, mint a session token, and stash the OAuth2 tokens for
// API calls.
func (h *Service) CompleteOIDCLogin(ctx context.Context, principal auth.Principal, token *oauth2.Token) (*repository.Token, error) {
	user, err := h.accountRepo.UpsertOIDCUser(ctx, principal)
	if err != nil {
		return nil, err
	}
	sessionToken, err := h.accountRepo.CreateSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if err := h.accountRepo.StoreOIDCToken(ctx, sessionToken, token); err != nil {
		return nil, errors.New("unable to store identity tokens")
	}
	return &sessionToken, nil
}

// OIDCToken returns the stored OAuth2 token set for an active session.
func (h *Service) OIDCToken(ctx context.Context, sessionToken repository.Token) (*oauth2.Token, error) {
	return h.accountRepo.LoadOIDCToken(ctx, sessionToken)
}

// APIAccessToken returns a valid access token for skyvisor-api calls,
// refreshing and re-persisting the stored token set when it has expired.
func (h *Service) APIAccessToken(ctx context.Context, sessionToken repository.Token) (string, error) {
	token, err := h.accountRepo.LoadOIDCToken(ctx, sessionToken)
	if err != nil {
		return "", err
	}
	if h.oidc == nil {
		if token.Valid() {
			return token.AccessToken, nil
		}
		return "", errors.New("auth session expired")
	}
	fresh, err := h.oidc.TokenSource(ctx, token).Token()
	if err != nil {
		return "", errors.New("auth session expired")
	}
	if fresh.AccessToken != token.AccessToken {
		if err := h.accountRepo.StoreOIDCToken(ctx, sessionToken, fresh); err != nil {
			return "", err
		}
	}
	return fresh.AccessToken, nil
}

func (h *Service) Logout(ctx context.Context, token repository.Token) error {
	return h.accountRepo.Logout(ctx, token)
}
