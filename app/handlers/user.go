package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
	"github.com/FACorreiaa/Aviation-tracker/app/models"
	"github.com/FACorreiaa/Aviation-tracker/app/view/user"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func (h *Handler) SettingsPage(w http.ResponseWriter, r *http.Request) error {
	return h.renderSettings(w, r, "")
}

func (h *Handler) SettingsAlerts(w http.ResponseWriter, r *http.Request) error {
	if h.service.API() == nil {
		return h.renderSettings(w, r, "Alert preferences require the SkyVisor API.")
	}
	accessToken, err := h.apiAccessToken(r)
	if err != nil {
		http.Redirect(w, r, "/login?return_to=/settings", http.StatusSeeOther)
		return nil
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	enabled := r.PostFormValue("email_alerts") == "1"
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	if err := h.service.API().SetEmailAlerts(ctx, accessToken, enabled); err != nil {
		if errors.Is(err, apiclient.ErrUnauthorized) {
			http.Redirect(w, r, "/login?return_to=/settings", http.StatusSeeOther)
			return nil
		}
		slog.ErrorContext(r.Context(), "update email alerts", "error", err)
		return h.renderSettings(w, r, "Unable to save alert preferences.")
	}
	return h.renderSettings(w, r, "Alert preferences saved.")
}

func (h *Handler) renderSettings(w http.ResponseWriter, r *http.Request, message string) error {
	return h.renderSettingsWithToken(w, r, message, "")
}

func (h *Handler) renderSettingsWithToken(w http.ResponseWriter, r *http.Request, message, newToken string) error {
	currentUser, ok := r.Context().Value(models.CtxKeyAuthUser).(*models.UserSession)
	if !ok || currentUser == nil {
		return errors.New("authenticated user missing from request")
	}

	pageModel := models.SettingsPage{User: *currentUser, Message: message, EmailAlerts: true, NewToken: newToken}
	if h.service.API() != nil {
		if accessToken, err := h.apiAccessToken(r); err == nil {
			ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
			defer cancel()
			if me, meErr := h.service.API().Me(ctx, accessToken); meErr == nil {
				pageModel.Email = me.Account.Email
				if pageModel.Email == "" {
					pageModel.Email = me.Identity.Email
				}
				pageModel.EmailAlerts = me.Account.EmailAlerts
				pageModel.Plan = me.Account.Plan
			}
			if tokens, patErr := h.service.API().ListPATs(ctx, accessToken); patErr == nil {
				pageModel.Tokens = tokens
			}
		}
	}

	settings := user.SettingsPage(pageModel, csrf.Token(r))
	return h.CreateLayout(w, r, "Account", settings).Render(r.Context(), w)
}

func (h *Handler) SettingsTokensCreate(w http.ResponseWriter, r *http.Request) error {
	if h.service.API() == nil {
		return h.renderSettings(w, r, "Connector tokens require the SkyVisor API.")
	}
	accessToken, err := h.apiAccessToken(r)
	if err != nil {
		http.Redirect(w, r, "/login?return_to=/settings", http.StatusSeeOther)
		return nil
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	name := strings.TrimSpace(r.PostFormValue("token_name"))
	if name == "" || len(name) > 80 {
		return h.renderSettings(w, r, "Token name must contain 1 to 80 characters.")
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	created, err := h.service.API().CreatePAT(ctx, accessToken, name)
	if err != nil {
		if errors.Is(err, apiclient.ErrUnauthorized) {
			http.Redirect(w, r, "/login?return_to=/settings", http.StatusSeeOther)
			return nil
		}
		slog.ErrorContext(r.Context(), "create pat", "error", err)
		return h.renderSettings(w, r, "Unable to create the connector token.")
	}
	return h.renderSettingsWithToken(w, r, "Token created. Copy it now — it will not be shown again.", created.Plaintext)
}

func (h *Handler) SettingsTokensRevoke(w http.ResponseWriter, r *http.Request) error {
	if h.service.API() == nil {
		return h.renderSettings(w, r, "Connector tokens require the SkyVisor API.")
	}
	accessToken, err := h.apiAccessToken(r)
	if err != nil {
		http.Redirect(w, r, "/login?return_to=/settings", http.StatusSeeOther)
		return nil
	}
	id := chi.URLParam(r, "id")
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	if err := h.service.API().RevokePAT(ctx, accessToken, id); err != nil {
		if errors.Is(err, apiclient.ErrUnauthorized) {
			http.Redirect(w, r, "/login?return_to=/settings", http.StatusSeeOther)
			return nil
		}
		if errors.Is(err, apiclient.ErrNotFound) {
			return h.renderSettings(w, r, "That token no longer exists.")
		}
		slog.ErrorContext(r.Context(), "revoke pat", "error", err)
		return h.renderSettings(w, r, "Unable to revoke the connector token.")
	}
	return h.renderSettings(w, r, "Token revoked. MCP sessions using it will stop working immediately.")
}
