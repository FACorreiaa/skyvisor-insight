package repository

import (
	"net/http"

	"context"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type MiddlewareRepository struct {
	Pgpool      *pgxpool.Pool
	RedisClient *redis.Client
	Validator   *validator.Validate
	Sessions    *sessions.CookieStore
}

// middleware

// AuthMiddleware to set the current logged in user in the context.
// AuthMiddleware See `Handlers.requireAuth` or `Handlers.redirectIfAuth` middleware.
func (m *MiddlewareRepository) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, _ := m.Sessions.Get(r, "auth")

		token := s.Values["token"]
		if token != nil {
			if token, ok := token.(string); ok {
				user, err := m.UserFromSessionToken(r.Context(), token)

				if err == nil {
					ctx := context.WithValue(r.Context(), models.CtxKeyAuthUser, user)
					r = r.WithContext(ctx)
				}
			}
		} else {
			ctx := context.WithValue(r.Context(), models.CtxKeyAuthUser, nil)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func (m *MiddlewareRepository) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(models.CtxKeyAuthUser)
		if user == nil {
			http.Redirect(w, r, "/login?return_to="+r.URL.Path, http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *MiddlewareRepository) RedirectIfAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(models.CtxKeyAuthUser)
		if user != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
