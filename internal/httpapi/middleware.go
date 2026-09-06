package httpapi

import (
	"log"
	"net/http"
)

// statusRecorder wraps http.ResponseWriter and intercepts WriteHeader so the
// status code of the response can be captured for access logging.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs each request with its method, path (without query
// string) and response status code. It never logs the query string or any
// query parameter value (e.g. user).
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Printf("%s %s %d", r.Method, r.URL.Path, rec.status)
	})
}
