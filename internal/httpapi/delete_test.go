package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"featureflags/internal/store"
)

func TestDeleteHandlerExistingKey(t *testing.T) {
	s := store.New()
	s.Create(store.Flag{Key: "feature-x", Enabled: true})

	req := httptest.NewRequest(http.MethodDelete, "/flags/feature-x", nil)
	req.SetPathValue("key", "feature-x")
	rr := httptest.NewRecorder()

	DeleteHandler(s)(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
	if _, ok := s.Get("feature-x"); ok {
		t.Fatal("expected flag to be removed")
	}
}

func TestDeleteHandlerSecondDeleteReturns404(t *testing.T) {
	s := store.New()
	s.Create(store.Flag{Key: "feature-x", Enabled: true})

	h := DeleteHandler(s)

	first := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodDelete, "/flags/feature-x", nil)
	req1.SetPathValue("key", "feature-x")
	h(first, req1)
	if first.Code != http.StatusNoContent {
		t.Fatalf("expected first delete 204, got %d", first.Code)
	}

	second := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodDelete, "/flags/feature-x", nil)
	req2.SetPathValue("key", "feature-x")
	h(second, req2)

	if second.Code != http.StatusNotFound {
		t.Fatalf("expected second delete 404, got %d", second.Code)
	}
	if ct := second.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %q", ct)
	}
	var body map[string]string
	if err := json.Unmarshal(second.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON error body: %v", err)
	}
	if body["error"] == "" {
		t.Fatalf("expected non-empty error message, got %v", body)
	}
}
