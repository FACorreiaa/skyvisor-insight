package controller

import (
	"net/http"

	pages "github.com/FACorreiaa/go-ollama/controller/html/pages"
	"github.com/FACorreiaa/go-ollama/controller/models"
	"github.com/a-h/templ"

	"github.com/FACorreiaa/go-ollama/core/account"
)

func (h *Handlers) CreateLayout(w http.ResponseWriter, r *http.Request, title string, data templ.Component) templ.Component {
	var user *account.User
	userCtx := r.Context().Value(ctxKeyAuthUser)
	if userCtx != nil {
		user = userCtx.(*account.User)
	}

	var nav []models.NavItem

	if user == nil {
		nav = []models.NavItem{
			{Path: "/", Label: "Home"},
			{Path: "/login", Label: "Sign in"},
			{Path: "/register", Label: "Sign up"},
		}
	} else {
		nav = []models.NavItem{
			{Path: "/", Label: "Home"},
			{Path: "/airlines", Label: "Airlines", Icon: "ion-paper-airplane"},
			{Path: "/airports", Label: "Airports", Icon: "ion-paper-airplane"},
			{Path: "/flights", Label: "Flights", Icon: "ion-paper-airplane"},
			{Path: "/locations", Label: "Locations", Icon: "ion-flag"},
			{Path: "/settings", Label: "Settings", Icon: "ion-gear-a"},
		}
	}

	layout := models.LayoutTempl{
		Title:     title,
		Nav:       nav,
		User:      user,
		ActiveNav: r.URL.Path,
		Content:   data,
	}

	return pages.LayoutPage(layout)
}
