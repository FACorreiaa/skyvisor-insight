package handlers

import (
	"net/http"

	"github.com/FACorreiaa/Aviation-tracker/app/view/legal"
)

func (h *Handler) TermsPage(w http.ResponseWriter, r *http.Request) error {
	return h.CreateLayout(w, r, "Terms of Service", legal.TermsPage()).Render(r.Context(), w)
}

func (h *Handler) PrivacyPage(w http.ResponseWriter, r *http.Request) error {
	return h.CreateLayout(w, r, "Privacy Policy", legal.PrivacyPage()).Render(r.Context(), w)
}
