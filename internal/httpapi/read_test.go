package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"featureflags/internal/store"
)

func TestListHandlerEmpty(t *testing.T) {
	s := store.New()
	h := ListHandler(s)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/flags", nil)
	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", ct)
	}
	if got := rec.Body.String(); got != "[]\n" {
		t.Fatalf("body = %q, want an empty JSON array", got)
	}
}

func TestListHandlerMultipleFlags(t *testing.T) {
	s := store.New()
	s.Create(store.Flag{Key: "a", Description: "first", Enabled: true, RolloutPercent: 50})
	s.Create(store.Flag{Key: "b", Description: "second", Enabled: false, RolloutPercent: 0})

	h := ListHandler(s)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/flags", nil)
	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var flags []store.Flag
	if err := json.Unmarshal(rec.Body.Bytes(), &flags); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if len(flags) != 2 {
		t.Fatalf("len = %d, want 2", len(flags))
	}

	byKey := map[string]store.Flag{}
	for _, f := range flags {
		byKey[f.Key] = f
	}
	if f, ok := byKey["a"]; !ok || !f.Enabled || f.Description != "first" {
		t.Fatalf("flag a missing or wrong: %+v", f)
	}
	if f, ok := byKey["b"]; !ok || f.Enabled || f.Description != "second" {
		t.Fatalf("flag b missing or wrong: %+v", f)
	}
}

func TestGetHandlerKnownKey(t *testing.T) {
	s := store.New()
	s.Create(store.Flag{Key: "feature-x", Description: "desc", Enabled: true, RolloutPercent: 100})

	h := GetHandler(s)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/flags/feature-x", nil)
	req.SetPathValue("key", "feature-x")
	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var flag store.Flag
	if err := json.Unmarshal(rec.Body.Bytes(), &flag); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if flag.Key != "feature-x" || !flag.Enabled || flag.Description != "desc" || flag.RolloutPercent != 100 {
		t.Fatalf("flag = %+v", flag)
	}
}

func TestGetHandlerUnknownKey(t *testing.T) {
	s := store.New()
	h := GetHandler(s)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/flags/missing", nil)
	req.SetPathValue("key", "missing")
	h(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if body["error"] == "" {
		t.Fatalf("expected an error object, got %v", body)
	}
}
