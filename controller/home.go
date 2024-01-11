package controller

import (
	"context"
	"net/http"

	pages "github.com/FACorreiaa/go-ollama/controller/html/pages"
)

func (h *Handlers) homePage(w http.ResponseWriter, r *http.Request) error {
	home := pages.HomePage()
	return h.CreateLayout(w, r, "Home Page", home).Render(context.Background(), w)
}
