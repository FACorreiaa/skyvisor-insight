package controller

import (
	"embed"
	"errors"
	"log"
	"log/slog"
	"net/http"

	"github.com/FACorreiaa/Aviation-tracker/core/flights"

	"github.com/FACorreiaa/Aviation-tracker/core/location"

	"github.com/FACorreiaa/Aviation-tracker/core/airline"

	"github.com/FACorreiaa/Aviation-tracker/core/airport"

	"github.com/FACorreiaa/Aviation-tracker/core/account"
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
	accounts  *account.RepositoyAccount
	airports  *airport.RepositoryAirport
	airlines  *airline.RepositoryAirline
	locations *location.RepositoryLocation
	flights   *flights.RepositoryFlights
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

func HandlerErrorExtended(message string, values ...interface{}) {
	for _, v := range values {
		if v == nil {
			continue
		}

		log.Printf("%s: %v", message, v)
	}
}

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
			accounts:  account.NewAccounts(pool, redisClient, validate),
			airports:  airport.NewAirports(pool),
			airlines:  airline.NewAirlines(pool),
			locations: location.NewLocations(pool),
			flights:   flights.NewFlights(pool),
		},
	}

	// Static files
	r.PathPrefix("/static/").Handler(http.FileServer(http.FS(staticFS)))
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, _ *http.Request) {
		file, _ := staticFS.ReadFile("static/favicon.ico")

		w.Header().Set("Content-Type", "image/x-icon")
		_, err := w.Write(file)
		if err != nil {
			return
		}
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

	// Airlines Router
	alr := auth.PathPrefix("/airlines").Subrouter()
	alr.HandleFunc("/airline", handler(h.airlineMainPage)).Methods(http.MethodGet)
	alr.HandleFunc("/airline/location", handler(h.airlineLocationPage)).Methods(http.MethodGet)
	alr.HandleFunc("/airline/{airline_name}", handler(h.airlineDetailsPage)).Methods(http.MethodGet)

	// alr.HandleFunc("/airline/{name}", handler(h.airlineMainPage)).Methods(http.MethodGet)
	alr.HandleFunc("/aircraft", handler(h.airlineAircraftPage)).Methods(http.MethodGet)
	alr.HandleFunc("/airplane", handler(h.airlineAirplanePage)).Methods(http.MethodGet)
	alr.HandleFunc("/tax", handler(h.airlineTaxPage)).Methods(http.MethodGet)

	// locations
	lr := auth.PathPrefix("/locations").Subrouter()
	lr.HandleFunc("/city", handler(h.cityMainPage)).Methods(http.MethodGet)
	lr.HandleFunc("/city/map", handler(h.cityLocationsPage)).Methods(http.MethodGet)
	lr.HandleFunc("/city/details/{city_id}", handler(h.cityDetailsPage)).Methods(http.MethodGet)
	lr.HandleFunc("/country", handler(h.countryMainPage)).Methods(http.MethodGet)
	lr.HandleFunc("/country/map", handler(h.countryLocationPage)).Methods(http.MethodGet)
	lr.HandleFunc("/country/details/{country_name}", handler(h.countryDetailsPage)).Methods(http.MethodGet)

	// Airports router
	apr := auth.PathPrefix("/airports").Subrouter()
	apr.HandleFunc("", handler(h.airportPage)).Methods(http.MethodGet)
	apr.HandleFunc("/locations", handler(h.airportLocationPage)).Methods(http.MethodGet)
	apr.HandleFunc("/details/{airport_id}", handler(h.airportDetailsPage)).Methods(http.MethodGet)

	// Flights router
	fr := auth.PathPrefix("/flights").Subrouter()
	fr.HandleFunc("", handler(h.allFlightsPage)).Methods(http.MethodGet)
	fr.HandleFunc("/details/{flight_number}", handler(h.detailedFlightsPage)).Methods(http.MethodGet)
	fr.HandleFunc("/preview", handler(h.flightsPreview)).Methods(http.MethodGet)

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

	// var validateErrors validator.ValidationErrors
	// isValidateError := errors.As(err, &validateErrors)
	if isValidateError {
		var errs []string
		for _, validateError := range validateErrors {
			errs = append(errs, validateError.Translate(h.translator))
		}
		return errs
	}

	return []string{err.Error()}
}
