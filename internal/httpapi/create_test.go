package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"featureflags/internal/store"
)

func postFlags(body string) (*httptest.ResponseRecorder, *store.Store) {
	s := store.New()
	req := httptest.NewRequest(http.MethodPost, "/flags", strings.NewReader(body))
	rr := httptest.NewRecorder()
	CreateHandler(s)(rr, req)
	return rr, s
}

func TestCreateValid(t *testing.T) {
	rr, s := postFlags(`{"key":"k","enabled":true,"description":"d","rollout_percent":42}`)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
	var got store.Flag
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON body: %v", err)
	}
	want := store.Flag{Key: "k", Enabled: true, Description: "d", RolloutPercent: 42}
	if got != want {
		t.Fatalf("expected %+v, got %+v", want, got)
	}
	stored, ok := s.Get("k")
	if !ok || stored != want {
		t.Fatalf("flag not stored correctly: ok=%v stored=%+v", ok, stored)
	}
}

func TestCreateDefaultsRolloutPercentTo100(t *testing.T) {
	rr, _ := postFlags(`{"key":"k","enabled":false}`)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
	var got store.Flag
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON body: %v", err)
	}
	if got.RolloutPercent != 100 {
		t.Fatalf("expected rollout_percent default 100, got %d", got.RolloutPercent)
	}
	if got.Enabled {
		t.Fatalf("expected enabled false to be preserved")
	}
}

func TestCreateExplicitZeroRollout(t *testing.T) {
	rr, _ := postFlags(`{"key":"k","enabled":true,"rollout_percent":0}`)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
	var got store.Flag
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON body: %v", err)
	}
	if got.RolloutPercent != 0 {
		t.Fatalf("expected rollout_percent 0, got %d", got.RolloutPercent)
	}
}

func TestCreateBrokenJSON(t *testing.T) {
	rr, _ := postFlags(`{broken`)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestCreateMissingKey(t *testing.T) {
	rr, _ := postFlags(`{"enabled":true}`)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestCreateEmptyKey(t *testing.T) {
	rr, _ := postFlags(`{"key":"","enabled":true}`)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestCreateMissingEnabled(t *testing.T) {
	rr, _ := postFlags(`{"key":"k"}`)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestCreateInvalidEnabled(t *testing.T) {
	rr, _ := postFlags(`{"key":"k","enabled":"yes"}`)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestCreateRolloutBelowRange(t *testing.T) {
	rr, _ := postFlags(`{"key":"k","enabled":true,"rollout_percent":-1}`)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestCreateRolloutAboveRange(t *testing.T) {
	rr, _ := postFlags(`{"key":"k","enabled":true,"rollout_percent":101}`)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestCreateOversizeBody(t *testing.T) {
	big := strings.Repeat("a", maxBodyBytes+1)
	rr, _ := postFlags(`{"key":"` + big + `","enabled":true}`)
	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", rr.Code)
	}
}
