package httpapi

import (
	"net/http"

	"featureflags/internal/store"
)

// UpdateHandler handles PUT /flags/{key}. Scaffold stub: 501 until the
// update ticket fills it in.
func UpdateHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, http.StatusNotImplemented, "not implemented")
	}
}
