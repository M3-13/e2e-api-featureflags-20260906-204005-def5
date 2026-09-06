package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
)

// maxBodyBytes is the request-body limit enforced by DecodeJSON.
const maxBodyBytes = 1 << 20 // 1 MiB

// WriteJSON writes v as a JSON response with the given status and a
// Content-Type of application/json.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteError writes {"error": msg} with the given status.
func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"error": msg})
}

// DecodeJSON decodes the request body into v, enforcing a 1 MiB limit. On a
// body larger than the limit it writes 413 and returns the error; on broken
// JSON it writes 400 and returns the error. It returns nil on success.
func DecodeJSON(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
		} else {
			WriteError(w, http.StatusBadRequest, "invalid JSON")
		}
		return err
	}
	return nil
}
