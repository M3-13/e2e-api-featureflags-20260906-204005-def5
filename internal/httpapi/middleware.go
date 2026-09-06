package httpapi

import "net/http"

// LoggingMiddleware wraps next. In the scaffold it is a pure pass-through;
// the access-logging ticket fills in the actual logging.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
