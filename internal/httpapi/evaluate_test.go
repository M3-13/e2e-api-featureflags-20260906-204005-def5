package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"featureflags/internal/store"
)

func evaluateRequest(s *store.Store, target string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /flags/{key}/evaluate", EvaluateHandler(s))
	req := httptest.NewRequest(http.MethodGet, target, nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func decodeResult(t *testing.T, rr *httptest.ResponseRecorder) bool {
	t.Helper()
	var body map[string]bool
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON body %q: %v", rr.Body.String(), err)
	}
	res, ok := body["result"]
	if !ok {
		t.Fatalf("response missing result field: %q", rr.Body.String())
	}
	return res
}

func newEvalStore() *store.Store {
	s := store.New()
	s.Create(store.Flag{Key: "feature", Enabled: true, RolloutPercent: 50})
	s.Create(store.Flag{Key: "off", Enabled: false, RolloutPercent: 50})
	s.Create(store.Flag{Key: "none", Enabled: true, RolloutPercent: 0})
	s.Create(store.Flag{Key: "all", Enabled: true, RolloutPercent: 100})
	return s
}

func TestEvaluateDeterministicSameUser(t *testing.T) {
	s := newEvalStore()
	first := decodeResult(t, evaluateRequest(s, "/flags/feature/evaluate?user=user-42"))
	for i := 0; i < 100; i++ {
		if got := decodeResult(t, evaluateRequest(s, "/flags/feature/evaluate?user=user-42")); got != first {
			t.Fatalf("same user got different results: %v then %v", first, got)
		}
	}
}

func TestEvaluateRolloutZeroFalse(t *testing.T) {
	s := newEvalStore()
	rr := evaluateRequest(s, "/flags/none/evaluate?user=user-42")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if decodeResult(t, rr) {
		t.Fatal("expected false for rollout_percent 0")
	}
}

func TestEvaluateRolloutHundredTrue(t *testing.T) {
	s := newEvalStore()
	rr := evaluateRequest(s, "/flags/all/evaluate?user=user-42")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if !decodeResult(t, rr) {
		t.Fatal("expected true for rollout_percent 100")
	}
}

func TestEvaluateDisabledFlagFalse(t *testing.T) {
	s := newEvalStore()
	rr := evaluateRequest(s, "/flags/off/evaluate?user=user-42")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if decodeResult(t, rr) {
		t.Fatal("expected false for disabled flag")
	}
}

func TestEvaluateUnknownKey404(t *testing.T) {
	s := newEvalStore()
	rr := evaluateRequest(s, "/flags/missing/evaluate?user=user-42")
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestEvaluateMissingUser400(t *testing.T) {
	s := newEvalStore()
	for _, target := range []string{
		"/flags/feature/evaluate",
		"/flags/feature/evaluate?user=",
	} {
		rr := evaluateRequest(s, target)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for %q, got %d", target, rr.Code)
		}
	}
}
