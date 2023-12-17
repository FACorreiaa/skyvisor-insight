package controller

import (
	"fmt"
	"github.com/FACorreiaa/go-ollama/core/account"
	"html/template"
	"net/http"
)

var registerPageTmpl = template.Must(template.ParseFS(
	htmlFS,
	"html/layout.html",
	"html/register.html",
))

type RegisterPage struct {
	Errors []string
	Values map[string]string
}

func (h *Handlers) registerPage(w http.ResponseWriter, r *http.Request) error {
	data := CreateLayout[RegisterPage](r, "Sign up", RegisterPage{})
	return registerPageTmpl.Execute(w, data)
}

func (h *Handlers) registerPost(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	var form account.RegisterForm
	var err error

	var token *account.Token
	err = h.formDecoder.Decode(&form, r.PostForm)
	if err == nil {
		token, err = h.core.accounts.RegisterNewAccount(r.Context(), form)
	}

	if err != nil {
		return registerPageTmpl.Execute(
			w,
			CreateLayout[RegisterPage](r, "Sign up", RegisterPage{
				Errors: h.formErrors(err),
			}),
		)
	}

	session, _ := h.sessions.Get(r, "auth")
	session.Values["token"] = token
	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	http.Redirect(w, r, "/settings", http.StatusSeeOther)
	return nil
}
