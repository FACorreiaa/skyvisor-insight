package session

import (
	"context"
	"net/http"
)

type ctxKey int

const (
	CtxKeyAuthUser ctxKey = iota
)

// AuthMiddleware to set the current logged in user in the context.
// AuthMiddleware See `Handlers.requireAuth` or `Handlers.redirectIfAuth` middleware.
func (h *AccountRepository) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, _ := h.sessions.Get(r, "auth")

		token := s.Values["token"]
		if token != nil {
			if token, ok := token.(string); ok {
				user, err := h.UserFromSessionToken(r.Context(), Token(token))

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

func (h *AccountRepository) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(CtxKeyAuthUser)
		if user == nil {
			http.Redirect(w, r, "/login?return_to="+r.URL.Path, http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *AccountRepository) RedirectIfAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(CtxKeyAuthUser)
		if user != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
