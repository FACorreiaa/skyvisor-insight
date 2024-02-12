package controller

import (
	"context"
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
	svg2 "github.com/FACorreiaa/Aviation-tracker/controller/svg"
	"net/http"

	"github.com/FACorreiaa/Aviation-tracker/controller/html/flights"
)

// https://openlayers.org/en/latest/examples/feature-move-animation.html future feature
// https://openlayers.org/en/latest/examples/flight-animation.html
// future feature on this branch for flights with destination

func (h *Handlers) renderLiveLocationsSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home", Icon: svg2.HomeIcon()},
		{
			Label: "Cities",
			Icon:  svg2.GlobeEuropeIcon(),
			SubItems: []models.SidebarItem{
				{Path: "/locations/city", Label: "City data", Icon: svg2.GlobeEuropeIcon()},
				{Path: "/locations/city/map", Label: "City locations", Icon: svg2.MapIcon()},
			},
		},
		{
			Label: "Countries",
			Icon:  svg2.GlobeIcon(),
			SubItems: []models.SidebarItem{
				{Path: "/locations/country", Label: "Country data", Icon: svg2.GlobeIcon()},
				{Path: "/locations/country/map", Label: "Country locations", Icon: svg2.MapIcon()},
			},
		},
		{Path: "/settings", Label: "Settings", Icon: svg2.SettingsIcon()},
		{Path: "/log-out", Label: "Log out", Icon: svg2.LogoutIcon()},
	}
	return sidebar
}

func (h *Handlers) liveFlightsPage(w http.ResponseWriter, r *http.Request) error {
	s := h.renderLiveLocationsSidebar()

	f := flights.LiveFlightsPage(s, "Live Flights", "check live flights data")
	return h.CreateLayout(w, r, "Live Flights", f).Render(context.Background(), w)
}
