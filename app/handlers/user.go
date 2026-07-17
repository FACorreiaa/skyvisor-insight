package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
	"github.com/FACorreiaa/Aviation-tracker/app/models"
	"github.com/FACorreiaa/Aviation-tracker/app/view/user"
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
	currentUser, ok := r.Context().Value(models.CtxKeyAuthUser).(*models.UserSession)
	if !ok || currentUser == nil {
		return errors.New("authenticated user missing from request")
	}

	pageModel := models.SettingsPage{User: *currentUser, Message: message, EmailAlerts: true}
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
		}
	}

	settings := user.SettingsPage(pageModel, csrf.Token(r))
	return h.CreateLayout(w, r, "Account", settings).Render(r.Context(), w)
}
