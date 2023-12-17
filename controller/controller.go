package controller

import (
	"embed"
	"errors"
	"github.com/FACorreiaa/go-ollama/core/account"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
)

//go:embed static
var staticFS embed.FS

//go:embed html
var htmlFS embed.FS

type core struct {
	accounts *account.Accounts
}

type Handlers struct {
	pgpool      *pgxpool.Pool
	formDecoder *form.Decoder
	validator   *validator.Validate
	translator  ut.Translator
	sessions    *sessions.CookieStore
	core        *core
	redisClient *redis.Client
}

func Router(pool *pgxpool.Pool, sessionSecret []byte, redisClient *redis.Client) http.Handler {
	validate := validator.New()
	translator, _ := ut.New(en.New(), en.New()).GetTranslator("en")
	if err := en_translations.RegisterDefaultTranslations(validate, translator); err != nil {
		slog.Error("Error registering translations", "error", err)
	}

	formDecoder := form.NewDecoder()

	r := mux.NewRouter()
	h := Handlers{
		pgpool:      pool,
		formDecoder: formDecoder,
		validator:   validate,
		translator:  translator,
		sessions:    sessions.NewCookieStore(sessionSecret),
		redisClient: redisClient,
		core: &core{
			accounts: account.NewAccounts(pool, redisClient, validate),
		},
	}

	// Static files
	r.PathPrefix("/static/").Handler(http.FileServer(http.FS(staticFS)))
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, _ *http.Request) {
		file, _ := staticFS.ReadFile("static/favicon.ico")
		w.Header().Set("Content-Type", "image/x-icon")
		w.Write(file)
	})

	// Public routes, authentication is optional
	optAuth := r.NewRoute().Subrouter()
	optAuth.Use(h.authMiddleware)
	optAuth.HandleFunc("/", handler(h.homePage)).Methods(http.MethodGet)

	// Routes that shouldn't be available to authenticated users
	noAuth := r.NewRoute().Subrouter()
	noAuth.Use(h.authMiddleware)
	noAuth.Use(h.redirectIfAuth)

	noAuth.HandleFunc("/login", handler(h.loginPage)).Methods(http.MethodGet)
	noAuth.HandleFunc("/login", handler(h.loginPost)).Methods(http.MethodPost)
	noAuth.HandleFunc("/register", handler(h.registerPage)).Methods(http.MethodGet)
	noAuth.HandleFunc("/register", handler(h.registerPost)).Methods(http.MethodPost)

	// Authenticated routes
	auth := r.NewRoute().Subrouter()
	auth.Use(h.authMiddleware)
	auth.Use(h.requireAuth)

	auth.HandleFunc("/logout", handler(h.logout)).Methods(http.MethodPost)
	auth.HandleFunc("/settings", handler(h.settingsPage)).Methods(http.MethodGet)

	return r
}

func handler(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			slog.Error("Error handling request", "error", err)
		}
	}
}

func (h *Handlers) formErrors(err error) []string {
	var decodeErrors form.DecodeErrors
	isDecodeError := errors.As(err, &decodeErrors)
	if isDecodeError {
		var errs []string
		for _, decodeError := range decodeErrors {
			errs = append(errs, decodeError.Error())
		}

		return errs
	}

	validateErrors, isValidateError := err.(validator.ValidationErrors)
	if isValidateError {
		var errs []string
		for _, validateError := range validateErrors {
			errs = append(errs, validateError.Translate(h.translator))
		}
		return errs
	}

	return []string{err.Error()}
}
