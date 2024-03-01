package app

import (
	"embed"
	"log"
	"log/slog"
	"net/http"

	"github.com/FACorreiaa/Aviation-tracker/app/handlers"
	"github.com/FACorreiaa/Aviation-tracker/app/repository"
	"github.com/FACorreiaa/Aviation-tracker/app/services"
	"github.com/FACorreiaa/Aviation-tracker/app/session"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

//go:embed static
var staticFS embed.FS

type core struct {
	accounts *session.AccountRepository
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

const ASC = "ASC"
const DESC = "DESC"

func HandleError(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
}

func Router(pool *pgxpool.Pool, sessionSecret []byte, redisClient *redis.Client) http.Handler {
	validate := validator.New()
	translator, _ := ut.New(en.New(), en.New()).GetTranslator("en")
	if err := enTranslations.RegisterDefaultTranslations(validate, translator); err != nil {
		slog.Error("Error registering translations", "error", err)
	}
	// var dir string
	//
	// flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	// flag.Parse()

	//formDecoder := form.NewDecoder()

	r := mux.NewRouter()
	//h := Handlers{
	//	pgpool:      pool,
	//	formDecoder: formDecoder,
	//	validator:   validate,
	//	translator:  translator,
	//	sessions:    sessions.NewCookieStore(sessionSecret),
	//	redisClient: redisClient,
	//	core: &core{
	//		accounts: session.NewAccounts(pool, redisClient, validate, sessions.NewCookieStore(sessionSecret)),
	//	},
	//}

	// Static files
	r.PathPrefix("/static/").Handler(http.FileServer(http.FS(staticFS)))
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, _ *http.Request) {
		file, _ := staticFS.ReadFile("static/favicon.ico")
		w.Header().Set("Content-Type", "image/x-icon")
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Type", "image/svg+xml")

		_, err := w.Write(file)
		if err != nil {
			return
		}
	})

	// session
	sessionsHandlers := session.NewAccounts(pool, redisClient, validate, sessions.NewCookieStore(sessionSecret))

	// business
	airlineRepo := repository.NewAirlineRepository(pool)
	airportRepo := repository.NewAirportRepository(pool)
	locationRepo := repository.NewLocationsRepository(pool)
	flightsRepo := repository.NewFlightsRepository(pool)
	authRepo := repository.NewAccounts(pool, redisClient, validate)

	service := services.NewService(airlineRepo, airportRepo, locationRepo, flightsRepo, authRepo)
	h := handlers.NewHandler(service)

	// r.HandleFunc("/icons/marker.png", func(w http.ResponseWriter, _ *http.Request) {
	//	file, _ := staticFS.ReadFile("icons/marker.png")
	//	w.Header().Set("Content-Type", "image/x-icon")
	//	w.Header().Set("Content-Type", "image/png")
	//	w.Header().Set("Content-Type", "image/jpeg")
	//	w.Header().Set("Content-Type", "image/svg+xml")
	//
	//	_, err := w.Write(file)
	//	if err != nil {
	//		return
	//	}
	//})
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))

	// Public routes, authentication is optional
	optAuth := r.NewRoute().Subrouter()
	optAuth.Use(sessionsHandlers.AuthMiddleware)
	optAuth.HandleFunc("/", handler(h.Homepage)).Methods(http.MethodGet)

	// Routes that shouldn't be available to authenticated users
	noAuth := r.NewRoute().Subrouter()
	noAuth.Use(sessionsHandlers.AuthMiddleware)
	noAuth.Use(sessionsHandlers.RedirectIfAuth)

	noAuth.HandleFunc("/login", handler(h.LoginPage)).Methods(http.MethodGet)
	noAuth.HandleFunc("/login", handler(h.LoginPost)).Methods(http.MethodPost)
	noAuth.HandleFunc("/register", handler(h.RegisterPage)).Methods(http.MethodGet)
	noAuth.HandleFunc("/register", handler(h.RegisterPost)).Methods(http.MethodPost)

	// Authenticated routes
	auth := r.NewRoute().Subrouter()
	auth.Use(sessionsHandlers.AuthMiddleware)
	auth.Use(sessionsHandlers.RequireAuth)

	auth.HandleFunc("/logout", handler(h.Logout)).Methods(http.MethodPost)
	auth.HandleFunc("/settings", handler(h.SettingsPage)).Methods(http.MethodGet)

	// Airlines Router
	alr := auth.PathPrefix("/airlines").Subrouter()

	alr.HandleFunc("/airline", handler(h.AirlineMainPage)).Methods(http.MethodGet)
	alr.HandleFunc("/airline/location", handler(h.AirlineLocationPage)).Methods(http.MethodGet)
	alr.HandleFunc("/airline/{airline_name}", handler(h.AirlineDetailsPage)).Methods(http.MethodGet)

	alr.HandleFunc("/aircraft", handler(h.AirlineAircraftPage)).Methods(http.MethodGet)
	alr.HandleFunc("/airplane", handler(h.AirlineAirplanePage)).Methods(http.MethodGet)
	alr.HandleFunc("/tax", handler(h.AirlineTaxPage)).Methods(http.MethodGet)

	// locations
	lr := auth.PathPrefix("/locations").Subrouter()
	lr.HandleFunc("/city", handler(h.CityMainPage)).Methods(http.MethodGet)
	lr.HandleFunc("/city/map", handler(h.CityLocationsPage)).Methods(http.MethodGet)
	lr.HandleFunc("/city/details/{city_id}", handler(h.CityDetailsPage)).Methods(http.MethodGet)
	lr.HandleFunc("/country", handler(h.CountryMainPage)).Methods(http.MethodGet)
	lr.HandleFunc("/country/map", handler(h.CountryLocationPage)).Methods(http.MethodGet)
	lr.HandleFunc("/country/details/{country_name}", handler(h.CountryDetailsPage)).Methods(http.MethodGet)

	// Airports router
	apr := auth.PathPrefix("/airports").Subrouter()
	apr.HandleFunc("", handler(h.AirportPage)).Methods(http.MethodGet)
	apr.HandleFunc("/locations", handler(h.AirportLocationPage)).Methods(http.MethodGet)
	apr.HandleFunc("/details/{airport_id}", handler(h.AirportDetailsPage)).Methods(http.MethodGet)

	// Flights router
	fr := auth.PathPrefix("/flights").Subrouter()
	fr.HandleFunc("", handler(h.AllFlightsPage)).Methods(http.MethodGet)
	fr.HandleFunc("/{flight_status}", handler(h.AllFlightsPage)).Methods(http.MethodGet)
	fr.HandleFunc("/details/{flight_number}", handler(h.DetailedFlightsPage)).Methods(http.MethodGet)
	fr.HandleFunc("/preview", handler(h.FlightsPreview)).Methods(http.MethodGet)

	return r
}

func handler(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			slog.Error("Error handling request", "error", err)
		}
	}
}
