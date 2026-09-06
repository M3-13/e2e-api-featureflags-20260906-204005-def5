package httpapi

import (
	"net/http"

	"featureflags/internal/store"
)

// DeleteHandler handles DELETE /flags/{key}. It removes the flag with the
// given key and answers 204 on success, or 404 with a JSON error object when
// no such flag exists.
func DeleteHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		if !s.Delete(key) {
			WriteError(w, http.StatusNotFound, "flag not found")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
