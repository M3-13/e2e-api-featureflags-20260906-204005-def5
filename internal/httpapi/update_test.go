package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"featureflags/internal/store"
)

func doUpdate(s *store.Store, key, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPut, "/flags/"+key, strings.NewReader(body))
	req.SetPathValue("key", key)
	rr := httptest.NewRecorder()
	UpdateHandler(s).ServeHTTP(rr, req)
	return rr
}

func TestUpdateHandlerUpdatesValues(t *testing.T) {
	s := store.New()
	s.Create(store.Flag{Key: "feature-a", Enabled: false, Description: "old", RolloutPercent: 10})

	rr := doUpdate(s, "feature-a", `{"enabled":true,"description":"new desc","rollout_percent":75}`)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var flag store.Flag
	if err := json.Unmarshal(rr.Body.Bytes(), &flag); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if flag.Key != "feature-a" {
		t.Fatalf("expected key feature-a, got %q", flag.Key)
	}
	if !flag.Enabled {
		t.Fatalf("expected enabled true")
	}
	if flag.Description != "new desc" {
		t.Fatalf("expected description %q, got %q", "new desc", flag.Description)
	}
	if flag.RolloutPercent != 75 {
		t.Fatalf("expected rollout_percent 75, got %d", flag.RolloutPercent)
	}

	got, ok := s.Get("feature-a")
	if !ok || !got.Enabled || got.Description != "new desc" || got.RolloutPercent != 75 {
		t.Fatalf("store not updated: %+v", got)
	}
}

func TestUpdateHandlerPreservesOmittedFields(t *testing.T) {
	s := store.New()
	s.Create(store.Flag{Key: "feature-a", Enabled: false, Description: "keep", RolloutPercent: 40})

	rr := doUpdate(s, "feature-a", `{"enabled":true}`)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var flag store.Flag
	if err := json.Unmarshal(rr.Body.Bytes(), &flag); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if flag.Description != "keep" {
		t.Fatalf("expected description preserved as %q, got %q", "keep", flag.Description)
	}
	if flag.RolloutPercent != 40 {
		t.Fatalf("expected rollout_percent preserved as 40, got %d", flag.RolloutPercent)
	}
}

func TestUpdateHandlerUnknownKey(t *testing.T) {
	s := store.New()

	rr := doUpdate(s, "missing", `{"enabled":true}`)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if body["error"] == "" {
		t.Fatalf("expected error message, got %v", body)
	}
}

func TestUpdateHandlerMissingEnabled(t *testing.T) {
	s := store.New()
	s.Create(store.Flag{Key: "feature-a", Enabled: false})

	rr := doUpdate(s, "feature-a", `{"rollout_percent":50}`)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestUpdateHandlerRolloutOutOfRange(t *testing.T) {
	s := store.New()
	s.Create(store.Flag{Key: "feature-a", Enabled: false})

	rr := doUpdate(s, "feature-a", `{"enabled":true,"rollout_percent":101}`)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestUpdateHandlerBrokenJSON(t *testing.T) {
	s := store.New()

	rr := doUpdate(s, "feature-a", `{broken`)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}
