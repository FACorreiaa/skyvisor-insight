package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
	watchesview "github.com/FACorreiaa/Aviation-tracker/app/view/watches"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

// WatchesPage lists monitored flights and plan entitlements.
func (h *Handler) WatchesPage(w http.ResponseWriter, r *http.Request) error {
	return h.renderWatches(w, r, flashFromQuery(r))
}

// WatchesCreate starts watching a flight number.
func (h *Handler) WatchesCreate(w http.ResponseWriter, r *http.Request) error {
	if h.service.API() == nil {
		return h.renderWatches(w, r, "Watches are not available on this server yet.")
	}
	accessToken, err := h.apiAccessToken(r)
	if err != nil {
		http.Redirect(w, r, "/login?return_to=/watches", http.StatusSeeOther)
		return nil
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	input := apiclient.CreateWatch{
		FlightNumber: strings.TrimSpace(r.PostFormValue("flight_number")),
		TripID:       strings.TrimSpace(r.PostFormValue("trip_id")),
	}
	returnTo := strings.TrimSpace(r.PostFormValue("return_to"))
	if returnTo == "" {
		returnTo = "/watches"
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	if _, err := h.service.API().CreateWatch(ctx, accessToken, input); err != nil {
		if errors.Is(err, apiclient.ErrUnauthorized) {
			http.Redirect(w, r, "/login?return_to=/watches", http.StatusSeeOther)
			return nil
		}
		if errors.Is(err, apiclient.ErrPaymentRequired) {
			return h.renderWatches(w, r, "Free accounts can watch one flight at a time. Upgrade to Pro for unlimited watches.")
		}
		if errors.Is(err, apiclient.ErrConflict) {
			http.Redirect(w, r, "/watches", http.StatusSeeOther)
			return nil
		}
		slog.ErrorContext(r.Context(), "create watch", "error", err)
		return h.renderWatches(w, r, friendlyAPIError(err, "Unable to watch that flight."))
	}
	http.Redirect(w, r, returnTo, http.StatusSeeOther)
	return nil
}

// WatchesShare creates a public pickup link for a watched flight.
func (h *Handler) WatchesShare(w http.ResponseWriter, r *http.Request) error {
	if h.service.API() == nil {
		return h.renderWatches(w, r, "Shares are not available on this server yet.")
	}
	accessToken, err := h.apiAccessToken(r)
	if err != nil {
		http.Redirect(w, r, "/login?return_to=/watches", http.StatusSeeOther)
		return nil
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	link, err := h.service.API().CreateShare(ctx, accessToken, id)
	if err != nil {
		if errors.Is(err, apiclient.ErrUnauthorized) {
			http.Redirect(w, r, "/login?return_to=/watches", http.StatusSeeOther)
			return nil
		}
		slog.ErrorContext(r.Context(), "create share", "error", err)
		return h.renderWatches(w, r, "Unable to create a share link.")
	}
	return h.renderWatches(w, r, "Pickup link ready: "+link.URLPath)
}

// SharePage is the public, unauthenticated pickup view.
func (h *Handler) SharePage(w http.ResponseWriter, r *http.Request) error {
	token := strings.TrimSpace(chi.URLParam(r, "token"))
	if token == "" || h.service.API() == nil {
		http.NotFound(w, r)
		return nil
	}
	ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
	defer cancel()
	view, err := h.service.API().PublicShare(ctx, token)
	if err != nil {
		if errors.Is(err, apiclient.ErrNotFound) {
			http.NotFound(w, r)
			return nil
		}
		slog.ErrorContext(r.Context(), "public share", "error", err)
		page := watchesview.SharePage(apiclient.PublicShare{}, "This share link is unavailable or expired.")
		return h.CreateLayout(w, r, "Shared flight", page).Render(r.Context(), w)
	}
	page := watchesview.SharePage(view, "")
	return h.CreateLayout(w, r, "Shared flight "+view.FlightNumber, page).Render(r.Context(), w)
}

// WatchesDelete stops watching a flight.
func (h *Handler) WatchesDelete(w http.ResponseWriter, r *http.Request) error {
	if h.service.API() == nil {
		return h.renderWatches(w, r, "Watches are not available on this server yet.")
	}
	accessToken, err := h.apiAccessToken(r)
	if err != nil {
		http.Redirect(w, r, "/login?return_to=/watches", http.StatusSeeOther)
		return nil
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	err = h.service.API().DeleteWatch(ctx, accessToken, id)
	if err != nil && !errors.Is(err, apiclient.ErrNotFound) {
		if errors.Is(err, apiclient.ErrUnauthorized) {
			http.Redirect(w, r, "/login?return_to=/watches", http.StatusSeeOther)
			return nil
		}
		slog.ErrorContext(r.Context(), "delete watch", "error", err)
		return h.renderWatches(w, r, "Unable to stop watching that flight.")
	}
	http.Redirect(w, r, "/watches", http.StatusSeeOther)
	return nil
}

// BillingCheckout redirects to Stripe Checkout when configured.
func (h *Handler) BillingCheckout(w http.ResponseWriter, r *http.Request) error {
	if h.service.API() == nil {
		return h.renderWatches(w, r, "Billing is not available on this server yet.")
	}
	accessToken, err := h.apiAccessToken(r)
	if err != nil {
		http.Redirect(w, r, "/login?return_to=/watches", http.StatusSeeOther)
		return nil
	}
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	session, err := h.service.API().CreateCheckout(ctx, accessToken)
	if err != nil {
		if errors.Is(err, apiclient.ErrUnauthorized) {
			http.Redirect(w, r, "/login?return_to=/watches", http.StatusSeeOther)
			return nil
		}
		slog.ErrorContext(r.Context(), "billing checkout", "error", err)
		// Fall back to local dev upgrade when Stripe is not configured.
		if devErr := h.service.API().DevActivatePro(ctx, accessToken); devErr == nil {
			http.Redirect(w, r, "/watches?upgraded=1", http.StatusSeeOther)
			return nil
		}
		return h.renderWatches(w, r, "Upgrade is temporarily unavailable. Configure Stripe or enable billing dev mode on the API.")
	}
	http.Redirect(w, r, session.URL, http.StatusSeeOther)
	return nil
}

func (h *Handler) renderWatches(w http.ResponseWriter, r *http.Request, message string) error {
	var watchList []apiclient.Watch
	var entitlements apiclient.Entitlements
	if h.service.API() == nil {
		if message == "" {
			message = "Watches are not available on this server yet."
		}
	} else if accessToken, err := h.apiAccessToken(r); err != nil {
		http.Redirect(w, r, "/login?return_to=/watches", http.StatusSeeOther)
		return nil
	} else {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		me, meErr := h.service.API().Me(ctx, accessToken)
		if meErr != nil {
			if errors.Is(meErr, apiclient.ErrUnauthorized) {
				http.Redirect(w, r, "/login?return_to=/watches", http.StatusSeeOther)
				return nil
			}
			slog.ErrorContext(r.Context(), "load entitlements", "error", meErr)
		} else {
			entitlements = me.Entitlements
		}
		watchList, err = h.service.API().ListWatches(ctx, accessToken)
		if err != nil {
			if errors.Is(err, apiclient.ErrUnauthorized) {
				http.Redirect(w, r, "/login?return_to=/watches", http.StatusSeeOther)
				return nil
			}
			slog.ErrorContext(r.Context(), "list watches", "error", err)
			if message == "" {
				message = "Watches are temporarily unavailable. Try again shortly."
			}
		}
	}
	page := watchesview.WatchesPage(watchList, entitlements, message, csrf.Token(r))
	return h.CreateLayout(w, r, "Watches", page).Render(r.Context(), w)
}

func flashFromQuery(r *http.Request) string {
	switch {
	case r.URL.Query().Get("upgraded") == "1":
		return "You are on SkyVisor Pro. Watch as many flights as you need."
	case r.URL.Query().Get("checkout") == "cancelled":
		return "Checkout was cancelled. You can upgrade any time."
	default:
		return ""
	}
}
