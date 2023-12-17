package controller

import (
	"github.com/FACorreiaa/go-ollama/core/account"
	"html/template"
	"log/slog"
	"net/http"
)

var settingsPageTmpl = template.Must(template.ParseFS(
	htmlFS,
	"html/layout.html",
	"html/settings.html",
))

type SettingsPage struct {
	Updated bool
	Errors  []string
	User    *account.User
}

func (h *Handlers) settingsPage(w http.ResponseWriter, r *http.Request) error {
	data := CreateLayout[SettingsPage](r, "Settings", SettingsPage{})
	data.Page.User = data.User
	return settingsPageTmpl.Execute(w, data)
}

func (h *Handlers) logout(w http.ResponseWriter, r *http.Request) error {
	session, _ := h.sessions.Get(r, "auth")
	token := session.Values["token"]

	if token, ok := token.(string); ok {
		h.core.accounts.Logout(r.Context(), account.Token(token))
	}

	session.Values["token"] = ""
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		slog.Error("failed to clear auth session", "err", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
