package controller

import (
	"fmt"
	"github.com/FACorreiaa/go-ollama/core/account"
	"html/template"
	"net/http"
)

type LoginPage struct {
	Errors []string
}

var loginPageTmpl = template.Must(template.ParseFS(
	htmlFS,
	"html/layout.html",
	"html/login.html",
))

func (h *Handlers) loginPage(w http.ResponseWriter, r *http.Request) error {
	data := CreateLayout[LoginPage](r, "Sign in", LoginPage{})
	return loginPageTmpl.Execute(w, data)
}

func (h *Handlers) loginPost(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	var form account.LoginForm
	var token *account.Token

	err := h.formDecoder.Decode(&form, r.PostForm)
	if err == nil {
		token, err = h.core.accounts.Login(r.Context(), form)
	}

	if err != nil {
		return loginPageTmpl.Execute(
			w,
			CreateLayout[LoginPage](r, "Sign in", LoginPage{
				Errors: h.formErrors(err),
			}),
		)
	}

	session, _ := h.sessions.Get(r, "auth")
	session.Values["token"] = token
	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("error saving session: %w", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
