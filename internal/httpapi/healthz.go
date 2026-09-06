package httpapi

import "net/http"

// HealthzHandler handles GET /healthz and reports the service as alive.
func HealthzHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}
