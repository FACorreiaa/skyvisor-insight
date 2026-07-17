package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ErrNotFound = errors.New("resource not found")
var ErrUnauthorized = errors.New("API session is not authorized")

// Client is a thin, typed HTTP client for skyvisor-api. Every call takes the
// caller's OIDC access token; the client itself holds no user state.
type Client struct {
	baseURL string
	http    *http.Client
}

func New(baseURL string) (*Client, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return nil, errors.New("skyvisor-api base URL must be an absolute HTTP(S) URL")
	}
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http:    &http.Client{Timeout: 15 * time.Second},
	}, nil
}

type Me struct {
	Account struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"account"`
	Identity struct {
		Issuer  string `json:"issuer"`
		Subject string `json:"subject"`
		Email   string `json:"email"`
		Name    string `json:"name"`
	} `json:"identity"`
}

type TripSegment struct {
	FlightNumber     string    `json:"flight_number"`
	DepartureIATA    string    `json:"departure_iata,omitempty"`
	ArrivalIATA      string    `json:"arrival_iata,omitempty"`
	DepartsAt        time.Time `json:"departs_at,omitzero"`
	ArrivesAt        time.Time `json:"arrives_at,omitzero"`
	BookingReference string    `json:"booking_reference,omitempty"`
	Notes            string    `json:"notes,omitempty"`
}

type Trip struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Segments  []TripSegment `json:"segments"`
	StartsAt  time.Time     `json:"starts_at,omitzero"`
	CreatedAt time.Time     `json:"created_at"`
}

type CreateTrip struct {
	Name     string        `json:"name"`
	Segments []TripSegment `json:"segments,omitempty"`
	StartsAt time.Time     `json:"starts_at,omitzero"`
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

func (c *Client) DeleteTrip(ctx context.Context, accessToken, id string) error {
	return c.do(ctx, http.MethodDelete, "/v1/trips/"+url.PathEscape(id), accessToken, nil, nil)
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
	case response.StatusCode >= 400:
		return apiError(response)
	}
	if result == nil {
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
		return fmt.Errorf("skyvisor-api %d %s: %s", response.StatusCode, payload.Error.Code, payload.Error.Message)
	}
	return fmt.Errorf("skyvisor-api returned status %d", response.StatusCode)
}
