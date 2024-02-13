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
		{Path: "/flights/live", Label: "All Flights", Icon: svg2.HomeIcon()},
		{
			Label: "Live Flights",
			Icon:  svg2.GlobeEuropeIcon(),
			SubItems: []models.SidebarItem{
				{Path: "/flights/active/data", Label: "Live Flights", Icon: svg2.GlobeEuropeIcon()},
				{Path: "/flights/active/map", Label: "Live Flights Locations", Icon: svg2.MapIcon()},
			},
		},
		{
			Label: "Landed Flights",
			Icon:  svg2.GlobeEuropeIcon(),
			SubItems: []models.SidebarItem{
				{Path: "/flights/landed/data", Label: "Landed Flights", Icon: svg2.GlobeEuropeIcon()},
				{Path: "/flights/landed/map", Label: "Landed Flights Location", Icon: svg2.MapIcon()},
			},
		},
		{Path: "/flights/scheduled", Label: "Scheduled Flights", Icon: svg2.HomeIcon()},
		{Path: "/flights/cancelled", Label: "Cancelled Flights", Icon: svg2.HomeIcon()},
		{Path: "/flights/incident", Label: "Incident Flights", Icon: svg2.HomeIcon()},
		{Path: "/flights/diverted", Label: "Diverted Flights", Icon: svg2.HomeIcon()},
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
