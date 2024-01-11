package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/FACorreiaa/go-ollama/controller/html/pages"
	"github.com/FACorreiaa/go-ollama/controller/models"
	"github.com/FACorreiaa/go-ollama/core/account"
)

func (h *Handlers) loginPage(w http.ResponseWriter, r *http.Request) error {
	login := pages.LoginPage(models.LoginPage{})
	return h.CreateLayout(w, r, "Login", login).Render(context.Background(), w)
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
		login := pages.LoginPage(models.LoginPage{Errors: h.formErrors(err)})

		return h.CreateLayout(w, r, "Sign In", login).Render(context.Background(), w)
	}

	session, _ := h.sessions.Get(r, "auth")
	session.Values["token"] = token
	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("error saving session: %w", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
