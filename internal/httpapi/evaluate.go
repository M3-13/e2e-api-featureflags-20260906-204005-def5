package httpapi

import (
	"net/http"

	"featureflags/internal/rollout"
	"featureflags/internal/store"
)

// EvaluateHandler handles GET /flags/{key}/evaluate?user={id}. It resolves the
// flag by key (404 when unknown), requires a non-empty "user" query parameter
// (400 when missing or empty) and answers 200 with {"result": bool}, where
// result = flag.Enabled && Evaluate(key, user, flag.RolloutPercent).
func EvaluateHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		flag, ok := s.Get(key)
		if !ok {
			WriteError(w, http.StatusNotFound, "flag not found")
			return
		}

		user := r.URL.Query().Get("user")
		if user == "" {
			WriteError(w, http.StatusBadRequest, "missing user")
			return
		}

		result := flag.Enabled && rollout.Evaluate(key, user, flag.RolloutPercent)
		WriteJSON(w, http.StatusOK, map[string]bool{"result": result})
	}
}
