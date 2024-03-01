package app

import (
	"log/slog"
	"net/http"

	"github.com/FACorreiaa/Aviation-tracker/app/view/user"

	"context"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
	"github.com/FACorreiaa/Aviation-tracker/core/account"
)

func (h *Handlers) settingsPage(w http.ResponseWriter, r *http.Request) error {
	settings := user.SettingsPage(models.SettingsPage{})
	data := h.CreateLayout(w, r, "Settings", settings).Render(context.Background(), w)
	return data
}

func (h *Handlers) logout(w http.ResponseWriter, r *http.Request) error {
	session, _ := h.sessions.Get(r, "auth")
	token := session.Values["token"]

	if token, ok := token.(string); ok {
		_ = h.core.accounts.Logout(r.Context(), account.Token(token))
	}

	session.Values["token"] = ""
	delete(session.Values, "token")
	delete(session.Values, "user")
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		slog.Error("failed to clear auth session", "err", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}