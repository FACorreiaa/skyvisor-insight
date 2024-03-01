package handlers

import (
	"log"
	"net/http"

	"context"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
	"github.com/FACorreiaa/Aviation-tracker/app/services"
	svg2 "github.com/FACorreiaa/Aviation-tracker/app/svg"
	"github.com/FACorreiaa/Aviation-tracker/app/view/components"
	"github.com/FACorreiaa/Aviation-tracker/core/account"
	"github.com/a-h/templ"
)

const ASC = "ASC"
const DESC = "DESC"

type Handler struct {
	service *services.AirlineService
	ctx     context.Context
}

type ctxKey int

const (
	ctxKeyAuthUser ctxKey = iota
)

func NewAirlineHandler(s *services.AirlineService) *Handler {
	return &Handler{service: s, ctx: context.Background()}
}

func HandleError(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
}

func (h *Handler) CreateLayout(_ http.ResponseWriter, r *http.Request, title string,
	data templ.Component) templ.Component {
	var user *account.User
	userCtx := r.Context().Value(ctxKeyAuthUser)
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

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) error {
	home := components.HomePage()
	return h.CreateLayout(w, r, "Home Page", home).Render(context.Background(), w)
}
