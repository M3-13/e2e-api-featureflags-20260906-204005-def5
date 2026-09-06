package httpapi

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func captureLog(t *testing.T, f func()) string {
	t.Helper()
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)
	f()
	return buf.String()
}

func TestLoggingMiddlewareLogsMethodPathStatus(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})
	handler := LoggingMiddleware(next)

	req := httptest.NewRequest(http.MethodPost, "/flags", nil)
	rec := httptest.NewRecorder()

	out := captureLog(t, func() {
		handler.ServeHTTP(rec, req)
	})

	if !strings.Contains(out, "POST") {
		t.Errorf("log %q does not contain method POST", out)
	}
	if !strings.Contains(out, "/flags") {
		t.Errorf("log %q does not contain path /flags", out)
	}
	if !strings.Contains(out, "201") {
		t.Errorf("log %q does not contain status 201", out)
	}
}

func TestLoggingMiddlewareDoesNotLogQueryString(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := LoggingMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/flags/abc/evaluate?user=secretvalue", nil)
	rec := httptest.NewRecorder()

	out := captureLog(t, func() {
		handler.ServeHTTP(rec, req)
	})

	if strings.Contains(out, "secretvalue") {
		t.Errorf("log %q must not contain the user parameter value", out)
	}
	if strings.Contains(out, "user=") {
		t.Errorf("log %q must not contain the query string", out)
	}
	if strings.Contains(out, "?") {
		t.Errorf("log %q must not contain the query string", out)
	}
}

func TestLoggingMiddlewareStatusRecorderDefaultsToOK(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	handler := LoggingMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	out := captureLog(t, func() {
		handler.ServeHTTP(rec, req)
	})

	if !strings.Contains(out, "200") {
		t.Errorf("log %q does not contain default status 200", out)
	}
}
