package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ErrNotFound = errors.New("resource not found")
var ErrUnauthorized = errors.New("API session is not authorized")
var ErrPaymentRequired = errors.New("plan upgrade required")
var ErrConflict = errors.New("resource conflict")

// Client is a thin, typed HTTP client for skyvisor-api. Every call takes the
// caller's OIDC access token; the client itself holds no user state.
type Client struct {
	baseURL string
	http    *http.Client
	stream  *http.Client
}

func New(baseURL string) (*Client, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return nil, errors.New("skyvisor-api base URL must be an absolute HTTP(S) URL")
	}
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http:    &http.Client{Timeout: 15 * time.Second},
		// Streaming (SSE) connections are long-lived, so they use a client
		// without a request timeout; the caller's context bounds them.
		stream: &http.Client{},
	}, nil
}

// StreamEvents opens the authenticated SSE stream and returns its body. The
// caller owns closing the reader and cancels by cancelling ctx. It is meant to
// be proxied to a browser, which cannot set the Authorization header itself.
func (c *Client) StreamEvents(ctx context.Context, accessToken string) (io.ReadCloser, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/events", nil)
	if err != nil {
		return nil, fmt.Errorf("build events request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Accept", "text/event-stream")

	response, err := c.stream.Do(request)
	if err != nil {
		return nil, fmt.Errorf("open events stream: %w", err)
	}
	if response.StatusCode == http.StatusUnauthorized {
		_ = response.Body.Close()
		return nil, ErrUnauthorized
	}
	if response.StatusCode >= 400 {
		defer response.Body.Close()
		return nil, apiError(response)
	}
	return response.Body, nil
}

type Me struct {
	Account struct {
		ID                   string    `json:"id"`
		CreatedAt            time.Time `json:"created_at"`
		Plan                 string    `json:"plan"`
		WatchLimit           int       `json:"watch_limit"`
		ProTrialUsed         bool      `json:"pro_trial_used"`
		Email                string    `json:"email,omitempty"`
		EmailAlerts          bool      `json:"email_alerts"`
		StripeCustomerID     string    `json:"stripe_customer_id,omitempty"`
		StripeSubscriptionID string    `json:"stripe_subscription_id,omitempty"`
	} `json:"account"`
	Identity struct {
		Issuer  string `json:"issuer"`
		Subject string `json:"subject"`
		Email   string `json:"email"`
		Name    string `json:"name"`
	} `json:"identity"`
	Entitlements Entitlements `json:"entitlements"`
}

type Entitlements struct {
	Plan             string `json:"plan"`
	WatchLimit       int    `json:"watch_limit"`
	ActiveWatches    int    `json:"active_watches"`
	CanCreateWatch   bool   `json:"can_create_watch"`
	ProTrialEligible bool   `json:"pro_trial_eligible"`
	EmailAlerts      bool   `json:"email_alerts"`
}

type FlightLive struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Altitude  float64   `json:"altitude,omitempty"`
	Speed     float64   `json:"speed,omitempty"`
	IsGround  bool      `json:"is_ground"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
}

type Flight struct {
	Number                string      `json:"number"`
	Status                string      `json:"status"`
	Airline               string      `json:"airline"`
	DepartureIATA         string      `json:"departure_iata"`
	ArrivalIATA           string      `json:"arrival_iata"`
	DepartureTerminal     string      `json:"departure_terminal,omitempty"`
	DepartureGate         string      `json:"departure_gate,omitempty"`
	DepartureDelayMinutes *int        `json:"departure_delay_minutes,omitempty"`
	ArrivalTerminal       string      `json:"arrival_terminal,omitempty"`
	ArrivalGate           string      `json:"arrival_gate,omitempty"`
	ArrivalBaggage        string      `json:"arrival_baggage,omitempty"`
	ArrivalDelayMinutes   *int        `json:"arrival_delay_minutes,omitempty"`
	ScheduledAt           time.Time   `json:"scheduled_at,omitzero"`
	EstimatedDeparture    time.Time   `json:"estimated_departure,omitzero"`
	EstimatedArrival      time.Time   `json:"estimated_arrival,omitzero"`
	ActualDeparture       time.Time   `json:"actual_departure,omitzero"`
	ActualArrival         time.Time   `json:"actual_arrival,omitzero"`
	AircraftRegistration  string           `json:"aircraft_registration,omitempty"`
	AircraftIATA          string           `json:"aircraft_iata,omitempty"`
	Live                  *FlightLive      `json:"live,omitempty"`
	Inbound               *InboundAircraft `json:"inbound,omitempty"`
	UpdatedAt             time.Time        `json:"updated_at"`
}

type TripSegment struct {
	FlightNumber     string    `json:"flight_number"`
	DepartureIATA    string    `json:"departure_iata,omitempty"`
	ArrivalIATA      string    `json:"arrival_iata,omitempty"`
	DepartsAt        time.Time `json:"departs_at,omitzero"`
	ArrivesAt        time.Time `json:"arrives_at,omitzero"`
	BookingReference string    `json:"booking_reference,omitempty"`
	Notes            string    `json:"notes,omitempty"`
	Live             *Flight   `json:"live,omitempty"`
	ConnectionRisk   string    `json:"connection_risk,omitempty"`
}

type ConnectionRisk struct {
	FromFlight     string `json:"from_flight"`
	ToFlight       string `json:"to_flight"`
	FromIATA       string `json:"from_iata,omitempty"`
	ToIATA         string `json:"to_iata,omitempty"`
	LayoverMinutes int    `json:"layover_minutes,omitempty"`
	Risk           string `json:"risk"`
	Reason         string `json:"reason"`
}

type AutoWatchSkipped struct {
	FlightNumber string `json:"flight_number"`
	Reason       string `json:"reason"`
}

type AutoWatchResult struct {
	Started []string           `json:"started"`
	Skipped []AutoWatchSkipped `json:"skipped,omitempty"`
}

type Trip struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Segments    []TripSegment    `json:"segments"`
	StartsAt    time.Time        `json:"starts_at,omitzero"`
	CreatedAt   time.Time        `json:"created_at"`
	Connections []ConnectionRisk `json:"connections,omitempty"`
	AutoWatches *AutoWatchResult `json:"auto_watches,omitempty"`
}

type CreateTrip struct {
	Name     string        `json:"name"`
	Segments []TripSegment `json:"segments,omitempty"`
	StartsAt time.Time     `json:"starts_at,omitzero"`
}

type InboundAircraft struct {
	Registration        string    `json:"registration,omitempty"`
	FlightNumber        string    `json:"flight_number,omitempty"`
	Status              string    `json:"status,omitempty"`
	OriginIATA          string    `json:"origin_iata,omitempty"`
	DestinationIATA     string    `json:"destination_iata,omitempty"`
	EstimatedArrival    time.Time `json:"estimated_arrival,omitzero"`
	ArrivalDelayMinutes *int      `json:"arrival_delay_minutes,omitempty"`
	MinutesUntilArrival *int      `json:"minutes_until_arrival,omitempty"`
	LateRisk            bool      `json:"late_risk"`
	Message             string    `json:"message,omitempty"`
}

type ShareLink struct {
	Token     string    `json:"token"`
	WatchID   string    `json:"watch_id"`
	URLPath   string    `json:"url_path"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type PublicShare struct {
	FlightNumber string    `json:"flight_number"`
	Flight       *Flight   `json:"flight,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
	Label        string    `json:"label,omitempty"`
}

type Watch struct {
	ID            string     `json:"id"`
	FlightNumber  string     `json:"flight_number"`
	TripID        string     `json:"trip_id,omitempty"`
	Status        string     `json:"status"`
	Flight        *Flight    `json:"flight,omitempty"`
	LastEventType string     `json:"last_event_type,omitempty"`
	LastEventAt   time.Time  `json:"last_event_at,omitzero"`
	Share         *ShareLink `json:"share,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CreateWatch struct {
	FlightNumber string `json:"flight_number"`
	TripID       string `json:"trip_id,omitempty"`
}

type CheckoutSession struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// APIError is a structured skyvisor-api error response.
type APIError struct {
	Status  int
	Code    string
	Message string
}

func (e *APIError) Error() string {
	if e == nil {
		return "skyvisor-api error"
	}
	return fmt.Sprintf("skyvisor-api %d %s: %s", e.Status, e.Code, e.Message)
}

func (c *Client) Me(ctx context.Context, accessToken string) (Me, error) {
	var me Me
	err := c.do(ctx, http.MethodGet, "/v1/me", accessToken, nil, &me)
	return me, err
}

func (c *Client) ListTrips(ctx context.Context, accessToken string) ([]Trip, error) {
	var payload struct {
		Data []Trip `json:"data"`
	}
	if err := c.do(ctx, http.MethodGet, "/v1/trips", accessToken, nil, &payload); err != nil {
		return nil, err
	}
	return payload.Data, nil
}

func (c *Client) CreateTrip(ctx context.Context, accessToken string, input CreateTrip) (Trip, error) {
	var trip Trip
	err := c.do(ctx, http.MethodPost, "/v1/trips", accessToken, input, &trip)
	return trip, err
}

// ImportTrip sends pasted booking text; the API extracts an itinerary with
// AI and creates the trip.
func (c *Client) ImportTrip(ctx context.Context, accessToken, text string) (Trip, error) {
	var trip Trip
	err := c.do(ctx, http.MethodPost, "/v1/trips/import", accessToken, map[string]string{"text": text}, &trip)
	return trip, err
}

// ImportTripPDF uploads a PDF e-ticket; the API extracts the itinerary and
// creates the trip. filename is used only for the multipart part name.
func (c *Client) ImportTripPDF(ctx context.Context, accessToken, filename string, pdf io.Reader) (Trip, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return Trip{}, fmt.Errorf("build upload: %w", err)
	}
	if _, err := io.Copy(part, pdf); err != nil {
		return Trip{}, fmt.Errorf("copy upload: %w", err)
	}
	if err := writer.Close(); err != nil {
		return Trip{}, fmt.Errorf("finalize upload: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/trips/import/pdf", body)
	if err != nil {
		return Trip{}, fmt.Errorf("build request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Accept", "application/json")

	response, err := c.http.Do(request)
	if err != nil {
		return Trip{}, fmt.Errorf("call skyvisor-api: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 1<<16))
		_ = response.Body.Close()
	}()
	switch {
	case response.StatusCode == http.StatusUnauthorized:
		return Trip{}, ErrUnauthorized
	case response.StatusCode == http.StatusRequestEntityTooLarge:
		return Trip{}, &APIError{Status: response.StatusCode, Code: "file_too_large", Message: "The PDF must be at most 10MB"}
	case response.StatusCode >= 400:
		return Trip{}, apiError(response)
	}
	var trip Trip
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&trip); err != nil {
		return Trip{}, fmt.Errorf("decode skyvisor-api response: %w", err)
	}
	return trip, nil
}

func (c *Client) DeleteTrip(ctx context.Context, accessToken, id string) error {
	return c.do(ctx, http.MethodDelete, "/v1/trips/"+url.PathEscape(id), accessToken, nil, nil)
}

func (c *Client) GetFlight(ctx context.Context, accessToken, number string) (Flight, error) {
	var flight Flight
	err := c.do(ctx, http.MethodGet, "/v1/flights/"+url.PathEscape(number), accessToken, nil, &flight)
	return flight, err
}

func (c *Client) ListWatches(ctx context.Context, accessToken string) ([]Watch, error) {
	var payload struct {
		Data []Watch `json:"data"`
	}
	if err := c.do(ctx, http.MethodGet, "/v1/watches", accessToken, nil, &payload); err != nil {
		return nil, err
	}
	return payload.Data, nil
}

func (c *Client) CreateWatch(ctx context.Context, accessToken string, input CreateWatch) (Watch, error) {
	var watch Watch
	err := c.do(ctx, http.MethodPost, "/v1/watches", accessToken, input, &watch)
	return watch, err
}

func (c *Client) DeleteWatch(ctx context.Context, accessToken, id string) error {
	return c.do(ctx, http.MethodDelete, "/v1/watches/"+url.PathEscape(id), accessToken, nil, nil)
}

func (c *Client) CreateShare(ctx context.Context, accessToken, watchID string) (ShareLink, error) {
	var link ShareLink
	err := c.do(ctx, http.MethodPost, "/v1/watches/"+url.PathEscape(watchID)+"/share", accessToken, map[string]any{}, &link)
	return link, err
}

// PublicShare fetches an unauthenticated pickup page. Token is the path only;
// Authorization is intentionally omitted.
func (c *Client) PublicShare(ctx context.Context, token string) (PublicShare, error) {
	var view PublicShare
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/public/shares/"+url.PathEscape(token), nil)
	if err != nil {
		return PublicShare{}, err
	}
	request.Header.Set("Accept", "application/json")
	response, err := c.http.Do(request)
	if err != nil {
		return PublicShare{}, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNotFound {
		return PublicShare{}, ErrNotFound
	}
	if response.StatusCode == http.StatusGone {
		return PublicShare{}, fmt.Errorf("share link expired")
	}
	if response.StatusCode >= 400 {
		return PublicShare{}, apiError(response)
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&view); err != nil {
		return PublicShare{}, err
	}
	return view, nil
}

func (c *Client) CreateCheckout(ctx context.Context, accessToken string) (CheckoutSession, error) {
	var session CheckoutSession
	err := c.do(ctx, http.MethodPost, "/v1/billing/checkout", accessToken, map[string]any{}, &session)
	return session, err
}

func (c *Client) CreateBillingPortal(ctx context.Context, accessToken string) (CheckoutSession, error) {
	var session CheckoutSession
	err := c.do(ctx, http.MethodPost, "/v1/billing/portal", accessToken, map[string]any{}, &session)
	return session, err
}

func (c *Client) DevActivatePro(ctx context.Context, accessToken string) error {
	return c.do(ctx, http.MethodPost, "/v1/billing/dev-activate-pro", accessToken, nil, &map[string]any{})
}

func (c *Client) SetEmailAlerts(ctx context.Context, accessToken string, enabled bool) error {
	return c.do(ctx, http.MethodPatch, "/v1/me/preferences", accessToken, map[string]any{"email_alerts": enabled}, &map[string]any{})
}

type AssistantResponse struct {
	Answer   string `json:"answer"`
	Model    string `json:"model"`
	TripID   string `json:"trip_id,omitempty"`
	Grounded bool   `json:"grounded"`
}

func (c *Client) AskAssistant(ctx context.Context, accessToken, question, tripID string) (AssistantResponse, error) {
	var response AssistantResponse
	body := map[string]string{"question": question}
	if tripID != "" {
		body["trip_id"] = tripID
	}
	err := c.do(ctx, http.MethodPost, "/v1/assistant", accessToken, body, &response)
	return response, err
}

func (c *Client) do(ctx context.Context, method, path, accessToken string, body, result any) error {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encode request: %w", err)
		}
		reader = bytes.NewReader(payload)
	}
	request, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Accept", "application/json")
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := c.http.Do(request)
	if err != nil {
		return fmt.Errorf("call skyvisor-api: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 1<<16))
		_ = response.Body.Close()
	}()

	switch {
	case response.StatusCode == http.StatusUnauthorized:
		return ErrUnauthorized
	case response.StatusCode == http.StatusNotFound:
		return ErrNotFound
	case response.StatusCode == http.StatusPaymentRequired:
		return fmt.Errorf("%w: %v", ErrPaymentRequired, apiError(response))
	case response.StatusCode == http.StatusConflict:
		return fmt.Errorf("%w: %v", ErrConflict, apiError(response))
	case response.StatusCode >= 400:
		return apiError(response)
	}
	if result == nil || response.StatusCode == http.StatusNoContent {
		return nil
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(result); err != nil {
		return fmt.Errorf("decode skyvisor-api response: %w", err)
	}
	return nil
}

func apiError(response *http.Response) error {
	var payload struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<16)).Decode(&payload); err == nil && payload.Error.Message != "" {
		return &APIError{Status: response.StatusCode, Code: payload.Error.Code, Message: payload.Error.Message}
	}
	return fmt.Errorf("skyvisor-api returned status %d", response.StatusCode)
}
