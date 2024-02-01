package controller

import (
	"context"
	"github.com/FACorreiaa/Aviation-tracker/controller/html/layout"
	"net/http"
)

func (h *Handlers) homePage(w http.ResponseWriter, r *http.Request) error {
	home := layout.HomePage()
	return h.CreateLayout(w, r, "Home Page", home).Render(context.Background(), w)
}
