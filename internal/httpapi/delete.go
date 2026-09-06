package httpapi

import (
	"net/http"

	"featureflags/internal/store"
)

// DeleteHandler handles DELETE /flags/{key}. Scaffold stub: 501 until the
// delete ticket fills it in.
func DeleteHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, http.StatusNotImplemented, "not implemented")
	}
}
