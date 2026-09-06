package httpapi

import (
	"net/http"

	"featureflags/internal/store"
)

// ListHandler handles GET /flags and writes every stored flag as a JSON
// array. An empty store yields [] (never null).
func ListHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, s.List())
	}
}

// GetHandler handles GET /flags/{key}. It writes the flag for the path key,
// or 404 with a JSON error object when the key is unknown.
func GetHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		flag, ok := s.Get(key)
		if !ok {
			WriteError(w, http.StatusNotFound, "flag not found")
			return
		}
		WriteJSON(w, http.StatusOK, flag)
	}
}
