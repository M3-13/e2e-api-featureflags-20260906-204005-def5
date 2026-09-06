package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	WriteJSON(rr, http.StatusOK, map[string]string{"status": "ok"})

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %q", ct)
	}
	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("unexpected body: %v", body)
	}
}

func TestWriteError(t *testing.T) {
	rr := httptest.NewRecorder()
	WriteError(rr, http.StatusBadRequest, "boom")

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if body["error"] != "boom" {
		t.Fatalf("unexpected body: %v", body)
	}
}

func TestDecodeJSON(t *testing.T) {
	var v map[string]any
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"key":"x"}`))
	rr := httptest.NewRecorder()

	if err := DecodeJSON(rr, req, &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v["key"] != "x" {
		t.Fatalf("unexpected value: %v", v)
	}
}

func TestDecodeJSONBroken(t *testing.T) {
	var v map[string]any
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{broken`))
	rr := httptest.NewRecorder()

	if err := DecodeJSON(rr, req, &v); err == nil {
		t.Fatal("expected error for broken JSON")
	}
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestDecodeJSONTooLarge(t *testing.T) {
	var v map[string]any
	big := strings.Repeat("a", maxBodyBytes+1)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"x":"`+big+`"}`))
	rr := httptest.NewRecorder()

	if err := DecodeJSON(rr, req, &v); err == nil {
		t.Fatal("expected error for oversized body")
	}
	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", rr.Code)
	}
}
