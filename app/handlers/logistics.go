package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
	logisticsview "github.com/FACorreiaa/Aviation-tracker/app/view/logistics"
	"github.com/gorilla/csrf"
)

func (h *Handler) LogisticsPage(w http.ResponseWriter, r *http.Request) error {
	return h.renderLogistics(w, r, "")
}

func (h *Handler) LogisticsCreateTeam(w http.ResponseWriter, r *http.Request) error {
	if h.service.API() == nil {
		return h.renderLogistics(w, r, "API is not configured.")
	}
	token, err := h.apiAccessToken(r)
	if err != nil {
		http.Redirect(w, r, "/login?return_to=/logistics", http.StatusSeeOther)
		return nil
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	_, err = h.service.API().CreateTeam(ctx, token, strings.TrimSpace(r.PostFormValue("name")))
	if errors.Is(err, apiclient.ErrPaymentRequired) {
		return h.renderLogistics(w, r, "Business plan required for teams. Activate Business in billing (dev) or upgrade.")
	}
	if err != nil {
		return h.renderLogistics(w, r, friendlyAPIError(err, "Unable to create team."))
	}
	return h.renderLogistics(w, r, "Team created.")
}

func (h *Handler) LogisticsJoinTeam(w http.ResponseWriter, r *http.Request) error {
	if h.service.API() == nil {
		return h.renderLogistics(w, r, "API is not configured.")
	}
	token, err := h.apiAccessToken(r)
	if err != nil {
		http.Redirect(w, r, "/login?return_to=/logistics", http.StatusSeeOther)
		return nil
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	code := strings.TrimSpace(r.PostFormValue("invite_code"))
	if code == "" {
		return h.renderLogistics(w, r, "Invite code is required.")
	}
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	_, err = h.service.API().JoinTeam(ctx, token, code)
	if err != nil {
		return h.renderLogistics(w, r, friendlyAPIError(err, "Unable to join team."))
	}
	return h.renderLogistics(w, r, "Joined team.")
}

func (h *Handler) renderLogistics(w http.ResponseWriter, r *http.Request, message string) error {
	if h.service.API() == nil {
		page := logisticsview.LogisticsPage(apiclient.LogisticsOverview{}, nil, apiclient.LogisticsQuery{}, "Logistics requires skyvisor-api.", csrf.Token(r))
		return h.CreateLayout(w, r, "Logistics", page).Render(r.Context(), w)
	}
	token, err := h.apiAccessToken(r)
	if err != nil {
		http.Redirect(w, r, "/login?return_to=/logistics", http.StatusSeeOther)
		return nil
	}
	q := apiclient.LogisticsQuery{
		Airline: strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("airline"))),
		Airport: strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("airport"))),
	}
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var overview apiclient.LogisticsOverview
	overview, err = h.service.API().LogisticsOverview(ctx, token, q)
	if errors.Is(err, apiclient.ErrUnauthorized) {
		http.Redirect(w, r, "/login?return_to=/logistics", http.StatusSeeOther)
		return nil
	}
	if errors.Is(err, apiclient.ErrPaymentRequired) {
		if message == "" {
			message = "Logistics requires Business. Use billing (or POST /v1/billing/dev-activate-business in Bruno) to unlock."
		}
	} else if err != nil && message == "" {
		message = "Logistics overview is temporarily unavailable."
	}

	var teamPtr *apiclient.Team
	if team, teamErr := h.service.API().GetTeam(ctx, token); teamErr == nil {
		teamPtr = &team
	}

	page := logisticsview.LogisticsPage(overview, teamPtr, q, message, csrf.Token(r))
	return h.CreateLayout(w, r, "Logistics", page).Render(r.Context(), w)
}
