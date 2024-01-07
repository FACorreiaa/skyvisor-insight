package controller

import (
	"html/template"
	"net/http"
)

type LiveFlightsPage struct{}

var liveFightsPageTmpl = template.Must(template.ParseFS(
	htmlFS,
	"html/layout.html",
	"html/live.html",
))

func (h *Handlers) liveFlightsPage(w http.ResponseWriter, r *http.Request) error {
	data := CreateLayout[LiveFlightsPage](r, "Home", LiveFlightsPage{})
	return liveFightsPageTmpl.Execute(w, data)
}
