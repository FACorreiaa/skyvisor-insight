package handlers

import (
	"errors"
	"net/http"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
	"github.com/FACorreiaa/Aviation-tracker/app/view/user"
	"github.com/gorilla/csrf"
)

func (h *Handler) SettingsPage(w http.ResponseWriter, r *http.Request) error {
	currentUser, ok := r.Context().Value(models.CtxKeyAuthUser).(*models.UserSession)
	if !ok || currentUser == nil {
		return errors.New("authenticated user missing from request")
	}

	settings := user.SettingsPage(models.SettingsPage{User: *currentUser}, csrf.Token(r))
	return h.CreateLayout(w, r, "Account", settings).Render(r.Context(), w)
}
