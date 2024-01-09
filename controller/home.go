package controller

import (
	"context"
	controller "github.com/FACorreiaa/go-ollama/controller/html"
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

func (h *Handlers) homePageTempl(w http.ResponseWriter, r *http.Request) error {
	return controller.HomeTemplExample("Home").Render(context.Background(), w)
}
