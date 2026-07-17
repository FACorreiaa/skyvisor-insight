package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var ErrInvalidIDToken = errors.New("invalid identity token")

// Principal is the verified identity extracted from an OIDC ID token.
type Principal struct {
	Issuer        string
	Subject       string
	Email         string
	EmailVerified bool
	Name          string
	Picture       string
}

type Config struct {
	IssuerURL    string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	// Audience is the API identifier requested for the access token (Auth0's
	// `audience` parameter). Empty means the provider default.
	Audience string
}

// Client drives the Authorization Code + PKCE flow against one OIDC provider.
type Client struct {
	oauth    oauth2.Config
	verifier *oidc.IDTokenVerifier
	audience string
}

func New(ctx context.Context, cfg Config) (*Client, error) {
	ctx = oidc.ClientContext(ctx, &http.Client{Timeout: 10 * time.Second})
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("discover OIDC provider: %w", err)
	}
	return &Client{
		oauth: oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			Endpoint:     provider.Endpoint(),
			RedirectURL:  cfg.RedirectURL,
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "offline_access"},
		},
		verifier: provider.Verifier(&oidc.Config{ClientID: cfg.ClientID}),
		audience: cfg.Audience,
	}, nil
}

// AuthCodeURL builds the provider redirect for one login attempt. The state,
// nonce, and PKCE verifier must be single-use values bound to the session.
func (c *Client) AuthCodeURL(state, nonce, pkceVerifier string, extra ...oauth2.AuthCodeOption) string {
	opts := []oauth2.AuthCodeOption{
		oidc.Nonce(nonce),
		oauth2.S256ChallengeOption(pkceVerifier),
	}
	if c.audience != "" {
		opts = append(opts, oauth2.SetAuthURLParam("audience", c.audience))
	}
	return c.oauth.AuthCodeURL(state, append(opts, extra...)...)
}

// Exchange redeems the authorization code, verifies the ID token against the
// expected nonce, and returns the token set plus the verified principal.
func (c *Client) Exchange(ctx context.Context, code, pkceVerifier, nonce string) (*oauth2.Token, Principal, error) {
	token, err := c.oauth.Exchange(ctx, code, oauth2.VerifierOption(pkceVerifier))
	if err != nil {
		return nil, Principal{}, fmt.Errorf("exchange authorization code: %w", err)
	}
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok || rawIDToken == "" {
		return nil, Principal{}, ErrInvalidIDToken
	}
	idToken, err := c.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, Principal{}, ErrInvalidIDToken
	}
	if idToken.Nonce != nonce {
		return nil, Principal{}, ErrInvalidIDToken
	}
	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		Nickname      string `json:"nickname"`
		Picture       string `json:"picture"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return nil, Principal{}, ErrInvalidIDToken
	}
	subject := strings.TrimSpace(idToken.Subject)
	if subject == "" || len(subject) > 255 {
		return nil, Principal{}, ErrInvalidIDToken
	}
	name := claims.Name
	if name == "" {
		name = claims.Nickname
	}
	return token, Principal{
		Issuer:        idToken.Issuer,
		Subject:       subject,
		Email:         claims.Email,
		EmailVerified: claims.EmailVerified,
		Name:          name,
		Picture:       claims.Picture,
	}, nil
}

// TokenSource wraps the stored token so API calls refresh it when expired.
func (c *Client) TokenSource(ctx context.Context, token *oauth2.Token) oauth2.TokenSource {
	return c.oauth.TokenSource(ctx, token)
}
