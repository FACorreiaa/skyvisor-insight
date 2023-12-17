package controller

import (
	"github.com/FACorreiaa/go-ollama/core/account"
	"net/http"
)

type NavItem struct {
	Path  string
	Icon  string
	Label string
}

type Layout[T any] struct {
	Title     string
	Nav       []NavItem
	ActiveNav string
	User      *account.User
	Page      T
}

func CreateLayout[T any](r *http.Request, title string, data T) Layout[T] {
	var user *account.User
	userCtx := r.Context().Value(ctxKeyAuthUser)
	if userCtx != nil {
		user = userCtx.(*account.User)
	}

	var nav []NavItem

	if user == nil {
		nav = []NavItem{
			{Path: "/", Label: "Home"},
			{Path: "/login", Label: "Sign in"},
			{Path: "/register", Label: "Sign up"},
		}
	} else {
		nav = []NavItem{
			{Path: "/", Label: "Home"},
			{Path: "/editor", Label: "New Article", Icon: "ion-compose"},
			{Path: "/settings", Label: "Settings", Icon: "ion-gear-a"},
		}
	}

	return Layout[T]{
		Title:     title,
		Nav:       nav,
		ActiveNav: r.URL.Path,
		Page:      data,
		User:      user,
	}
}
