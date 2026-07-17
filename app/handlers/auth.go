package handlers

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

const (
	sessionKeyToken        = "token"
	sessionKeyOIDCState    = "oidc_state"
	sessionKeyOIDCNonce    = "oidc_nonce"
	sessionKeyOIDCVerifier = "oidc_verifier"
	sessionKeyOIDCReturnTo = "oidc_return_to"
)

// LoginStart redirects the browser to the identity provider using the
// Authorization Code + PKCE flow. There is no local password form.
func (h *Handler) LoginStart(w http.ResponseWriter, r *http.Request) error {
	return h.startAuthFlow(w, r)
}

// RegisterStart is the same flow with a sign-up hint for providers that
// support it (Auth0's screen_hint).
func (h *Handler) RegisterStart(w http.ResponseWriter, r *http.Request) error {
	return h.startAuthFlow(w, r, oauth2.SetAuthURLParam("screen_hint", "signup"))
}

func (h *Handler) startAuthFlow(w http.ResponseWriter, r *http.Request, extra ...oauth2.AuthCodeOption) error {
	if h.oidc == nil {
		http.Error(w, "Sign-in is not configured on this server.", http.StatusServiceUnavailable)
		return nil
	}
	state, err := randomHex(16)
	if err != nil {
		return err
	}
	nonce, err := randomHex(16)
	if err != nil {
		return err
	}
	verifier := oauth2.GenerateVerifier()

	s, _ := h.sessions.Get(r, "auth")
	s.Values[sessionKeyOIDCState] = state
	s.Values[sessionKeyOIDCNonce] = nonce
	s.Values[sessionKeyOIDCVerifier] = verifier
	s.Values[sessionKeyOIDCReturnTo] = safeReturnPath(r.URL.Query().Get("return_to"))
	if err := s.Save(r, w); err != nil {
		return errors.New("error saving session")
	}

	http.Redirect(w, r, h.oidc.AuthCodeURL(state, nonce, verifier, extra...), http.StatusSeeOther)
	return nil
}

// AuthCallback completes the flow: state check, code exchange, ID token
// verification, local user upsert, and session creation.
func (h *Handler) AuthCallback(w http.ResponseWriter, r *http.Request) error {
	if h.oidc == nil {
		http.Error(w, "Sign-in is not configured on this server.", http.StatusServiceUnavailable)
		return nil
	}
	s, _ := h.sessions.Get(r, "auth")
	state, _ := s.Values[sessionKeyOIDCState].(string)
	nonce, _ := s.Values[sessionKeyOIDCNonce].(string)
	verifier, _ := s.Values[sessionKeyOIDCVerifier].(string)
	returnTo, _ := s.Values[sessionKeyOIDCReturnTo].(string)
	delete(s.Values, sessionKeyOIDCState)
	delete(s.Values, sessionKeyOIDCNonce)
	delete(s.Values, sessionKeyOIDCVerifier)
	delete(s.Values, sessionKeyOIDCReturnTo)

	if providerError := r.URL.Query().Get("error"); providerError != "" {
		slog.Warn("identity provider returned an error",
			"error", providerError,
			"description", r.URL.Query().Get("error_description"))
		if err := s.Save(r, w); err != nil {
			return errors.New("error saving session")
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	}

	query := r.URL.Query()
	if state == "" || verifier == "" || query.Get("state") != state || query.Get("code") == "" {
		http.Error(w, "The sign-in attempt expired. Start again from the login page.", http.StatusBadRequest)
		return nil
	}

	token, principal, err := h.oidc.Exchange(r.Context(), query.Get("code"), verifier, nonce)
	if err != nil {
		slog.Error("complete OIDC exchange", "error", err)
		http.Error(w, "Sign-in failed. Start again from the login page.", http.StatusBadGateway)
		return nil
	}

	sessionToken, err := h.service.CompleteOIDCLogin(r.Context(), principal, token)
	if err != nil {
		slog.Error("complete OIDC login", "error", err)
		http.Error(w, "Sign-in failed. Try again in a moment.", http.StatusInternalServerError)
		return nil
	}

	s.Values[sessionKeyToken] = *sessionToken
	if err := s.Save(r, w); err != nil {
		return errors.New("error saving session")
	}

	if returnTo == "" {
		returnTo = "/"
	}
	http.Redirect(w, r, returnTo, http.StatusSeeOther)
	return nil
}

// Logout clears the local session and the stored identity tokens.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) error {
	s, _ := h.sessions.Get(r, "auth")
	token := s.Values[sessionKeyToken]

	if token, ok := token.(string); ok {
		_ = h.service.Logout(r.Context(), token)
	}

	delete(s.Values, sessionKeyToken)
	delete(s.Values, "user")
	s.Options.MaxAge = -1
	if err := s.Save(r, w); err != nil {
		slog.Error("failed to clear auth session", "err", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

// safeReturnPath only allows same-site absolute paths so the callback can
// never redirect off-site (open redirect).
func safeReturnPath(path string) string {
	if !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") || strings.Contains(path, "\\") {
		return ""
	}
	return path
}

func randomHex(bytes int) (string, error) {
	buf := make([]byte, bytes)
	if _, err := rand.Read(buf); err != nil {
		return "", errors.New("generate random value")
	}
	return fmt.Sprintf("%x", buf), nil
}
