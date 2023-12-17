package controller

import (
	"html/template"
	"net/http"
)

type HomePage struct{}

var homePageTmpl = template.Must(template.ParseFS(
	htmlFS,
	"html/layout.html",
	"html/home.html",
))

func (h *Handlers) homePage(w http.ResponseWriter, r *http.Request) error {
	data := CreateLayout[HomePage](r, "Home", HomePage{})
	return homePageTmpl.Execute(w, data)
}
