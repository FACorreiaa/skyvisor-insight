package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/FACorreiaa/Aviation-tracker/app/view/user"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
	"github.com/FACorreiaa/Aviation-tracker/core/account"
)

func (h *Handlers) loginPage(w http.ResponseWriter, r *http.Request) error {
	login := user.LoginPage(models.LoginPage{})
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
		login := user.LoginPage(models.LoginPage{Errors: h.formErrors(err)})

		return h.CreateLayout(w, r, "Sign In", login).Render(context.Background(), w)
	}

	session, _ := h.sessions.Get(r, "auth")
	session.Values["token"] = token

	if err := session.Save(r, w); err != nil {
		return errors.New("error saving session")
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
