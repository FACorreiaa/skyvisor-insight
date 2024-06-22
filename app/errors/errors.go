package httperror

import (
	"encoding/json"
	"net/http"
)

type HTTPError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func (e *HTTPError) WriteError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.StatusCode)
	json.NewEncoder(w).Encode(e)
}

var (
	ErrInvalidID      = &HTTPError{Message: "Invalid ID", StatusCode: http.StatusBadRequest}
	ErrNotFound       = &HTTPError{Message: "Resource not found", StatusCode: http.StatusNotFound}
	ErrInternalServer = &HTTPError{Message: "Internal server error", StatusCode: http.StatusInternalServerError}
)

func NewHTTPError(message string, statusCode int) *HTTPError {
	return &HTTPError{
		Message:    message,
		StatusCode: statusCode,
	}
}
