package models

import (
	"github.com/FACorreiaa/go-ollama/core/account"
	"github.com/a-h/templ"
)

type NavItem struct {
	Path  string
	Icon  string
	Label string
}

type LayoutTempl struct {
	Title     string
	Nav       []NavItem
	ActiveNav string
	User      *account.User
	Content   templ.Component
}

type SettingsPage struct {
	Updated bool
	Errors  []string
	User    account.User
}

type LoginPage struct {
	Errors []string
}

type RegisterPage struct {
	Errors []string
	Values map[string]string
}
