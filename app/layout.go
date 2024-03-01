package app

import (
	"log"
	"net/http"

	"github.com/FACorreiaa/Aviation-tracker/app/session"
	svg2 "github.com/FACorreiaa/Aviation-tracker/app/svg"

	"github.com/FACorreiaa/Aviation-tracker/app/view/components"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
	"github.com/a-h/templ"

	"context"

	"github.com/FACorreiaa/Aviation-tracker/core/account"
)

func (h *Handlers) CreateLayout(_ http.ResponseWriter, r *http.Request, title string,
	data templ.Component) templ.Component {
	var user *account.User
	userCtx := r.Context().Value(session.CtxKeyAuthUser)
	if userCtx != nil {
		switch u := userCtx.(type) {
		case *account.User:
			user = u
		default:
			log.Printf("Unexpected type in userCtx: %T", userCtx)
		}
	}

	var nav []models.NavItem

	if user == nil {
		nav = []models.NavItem{
			{Path: "/", Label: "Home", Icon: svg2.HomeIcon()},
			{Path: "/login", Label: "Sign in", Icon: svg2.HomeIcon()},
			{Path: "/register", Label: "Sign up", Icon: svg2.HomeIcon()},
		}
	} else {
		nav = []models.NavItem{
			{Path: "/", Label: "Home", Icon: svg2.HomeIcon()},
			{Path: "/airlines/airline", Label: "Airlines", Icon: svg2.TicketIcon()},
			{Path: "/airports", Label: "Airports", Icon: svg2.BuildingOfficeIcon()},
			{Path: "/flights", Label: "Flights", Icon: svg2.PaperAirplaneIcon()},
			{Path: "/locations/city", Label: "Locations", Icon: svg2.LocationsIcon()},
			{Path: "/settings", Label: "Settings", Icon: svg2.SettingsIcon()},
		}
	}

	l := models.LayoutTempl{
		Title:     title,
		Nav:       nav,
		User:      user,
		ActiveNav: r.URL.Path,
		Content:   data,
	}

	return components.LayoutPage(l)
}

func (h *Handlers) homePage(w http.ResponseWriter, r *http.Request) error {
	home := components.HomePage()
	return h.CreateLayout(w, r, "Home Page", home).Render(context.Background(), w)
}
