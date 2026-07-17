package handlers

import (
	"bufio"
	"errors"
	"log/slog"
	"net/http"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
)

// EventsStream proxies the account's skyvisor-api SSE stream to the browser.
// The browser authenticates with its session cookie; this handler attaches the
// OIDC access token upstream, which EventSource cannot do on its own.
func (h *Handler) EventsStream(w http.ResponseWriter, r *http.Request) error {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming is unavailable.", http.StatusInternalServerError)
		return nil
	}
	if h.service.API() == nil {
		http.Error(w, "Live updates are not available on this server.", http.StatusServiceUnavailable)
		return nil
	}
	accessToken, err := h.apiAccessToken(r)
	if err != nil {
		http.Error(w, "Sign in to receive live updates.", http.StatusUnauthorized)
		return nil
	}

	body, err := h.service.API().StreamEvents(r.Context(), accessToken)
	if err != nil {
		if errors.Is(err, apiclient.ErrUnauthorized) {
			http.Error(w, "Sign in to receive live updates.", http.StatusUnauthorized)
			return nil
		}
		slog.ErrorContext(r.Context(), "open events stream", "error", err)
		http.Error(w, "Live updates are temporarily unavailable.", http.StatusBadGateway)
		return nil
	}
	defer func() { _ = body.Close() }()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	flusher.Flush()

	// Copy line-by-line so each SSE frame is flushed to the browser promptly.
	reader := bufio.NewReader(body)
	for {
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			if _, writeErr := w.Write(line); writeErr != nil {
				return nil
			}
			if len(line) == 1 && line[0] == '\n' {
				flusher.Flush()
			}
		}
		if err != nil {
			return nil
		}
	}
}
