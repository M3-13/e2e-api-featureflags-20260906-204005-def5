package httpapi

import (
	"net/http"

	"featureflags/internal/store"
)

// EvaluateHandler handles GET /flags/{key}/evaluate. Scaffold stub: 501 until
// the evaluate ticket fills it in.
func EvaluateHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, http.StatusNotImplemented, "not implemented")
	}
}
