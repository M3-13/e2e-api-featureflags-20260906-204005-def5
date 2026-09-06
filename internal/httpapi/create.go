package httpapi

import (
	"net/http"

	"featureflags/internal/store"
)

// CreateHandler handles POST /flags. Scaffold stub: 501 until the
// validation ticket fills it in.
func CreateHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, http.StatusNotImplemented, "not implemented")
	}
}
