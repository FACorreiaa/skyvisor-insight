package session

import (
	"context"
	"net/http"

	"github.com/FACorreiaa/Aviation-tracker/core/account"
)

type ctxKey int

const (
	CtxKeyAuthUser ctxKey = iota
)

// AuthMiddleware to set the current logged in user in the context.
// AuthMiddleware See `Handlers.requireAuth` or `Handlers.redirectIfAuth` middleware.
func (h *RepositoryAccount) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.sessions.Get(r, "auth")

		token := session.Values["token"]
		if token != nil {
			if token, ok := token.(string); ok {
				user, err := h.UserFromSessionToken(r.Context(), account.Token(token))

				if err == nil {
					ctx := context.WithValue(r.Context(), CtxKeyAuthUser, user)
					r = r.WithContext(ctx)
				}
			}
		} else {
			ctx := context.WithValue(r.Context(), CtxKeyAuthUser, nil)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func (h *RepositoryAccount) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(CtxKeyAuthUser)
		if user == nil {
			http.Redirect(w, r, "/login?return_to="+r.URL.Path, http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *RepositoryAccount) RedirectIfAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(CtxKeyAuthUser)
		if user != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
