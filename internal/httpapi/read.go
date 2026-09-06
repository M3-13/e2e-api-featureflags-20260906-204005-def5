package httpapi

import (
	"net/http"

	"featureflags/internal/store"
)

// ListHandler handles GET /flags. Scaffold stub: 501 until the read ticket
// fills it in.
func ListHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, http.StatusNotImplemented, "not implemented")
	}
}

// GetHandler handles GET /flags/{key}. Scaffold stub: 501 until the read
// ticket fills it in.
func GetHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, http.StatusNotImplemented, "not implemented")
	}
}
