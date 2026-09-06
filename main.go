package main

import (
	"log"
	"net/http"

	"featureflags/internal/httpapi"
	"featureflags/internal/store"
)

func newHandler() http.Handler {
	s := store.New()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /flags", httpapi.CreateHandler(s))
	mux.HandleFunc("GET /flags", httpapi.ListHandler(s))
	mux.HandleFunc("GET /flags/{key}", httpapi.GetHandler(s))
	mux.HandleFunc("PUT /flags/{key}", httpapi.UpdateHandler(s))
	mux.HandleFunc("DELETE /flags/{key}", httpapi.DeleteHandler(s))
	mux.HandleFunc("GET /flags/{key}/evaluate", httpapi.EvaluateHandler(s))
	mux.HandleFunc("GET /healthz", httpapi.HealthzHandler())
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpapi.WriteError(w, http.StatusNotFound, "not found")
	})

	return httpapi.LoggingMiddleware(mux)
}

func main() {
	log.Println("featureflag-api listening on :8080")
	if err := http.ListenAndServe(":8080", newHandler()); err != nil {
		log.Fatal(err)
	}
}
