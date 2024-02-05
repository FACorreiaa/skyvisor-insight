package controller

import (
	"context"
	"net/http"

	"github.com/FACorreiaa/Aviation-tracker/controller/html/components"
)

func (h *Handlers) homePage(w http.ResponseWriter, r *http.Request) error {
	home := components.HomePage()
	return h.CreateLayout(w, r, "Home Page", home).Render(context.Background(), w)
}
