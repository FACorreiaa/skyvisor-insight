package handlers

import (
	"net/http"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
	"github.com/FACorreiaa/Aviation-tracker/app/view/onboarding"
	"github.com/gorilla/csrf"
)

func (h *Handler) WelcomePage(w http.ResponseWriter, r *http.Request) error {
	name := ""
	if currentUser, ok := r.Context().Value(models.CtxKeyAuthUser).(*models.UserSession); ok && currentUser != nil {
		name = currentUser.Username
	}
	plan := r.URL.Query().Get("plan")
	if plan != "pro" && plan != "business" {
		plan = ""
	}
	page := onboarding.WelcomePage(name, plan, csrf.Token(r))
	return h.CreateLayout(w, r, "Welcome", page).Render(r.Context(), w)
}
