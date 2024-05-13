package handlers

import (
	"net/http"

	"context"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
	"github.com/FACorreiaa/Aviation-tracker/app/view/user"
)

func (h *Handler) SettingsPage(w http.ResponseWriter, r *http.Request) error {
	settings := user.SettingsPage(models.SettingsPage{})
	data := h.CreateLayout(w, r, "Settings", settings).Render(context.Background(), w)
	return data
}
