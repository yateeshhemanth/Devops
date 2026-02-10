package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"devops/backend/internal/config"
	"devops/backend/internal/store"
)

func TestDashboardAndOps(t *testing.T) {
	srv := New(store.New(), "", config.SecurityConfig{})
	h := srv.Router()

	t.Run("dashboard", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 got %d", w.Code)
		}
	})

	t.Run("power vm", func(t *testing.T) {
		payload := []byte(`{"action":"stop"}`)
		r := httptest.NewRequest(http.MethodPost, "/api/v1/vms/vm-1/power", bytes.NewReader(payload))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 got %d: %s", w.Code, w.Body.String())
		}
		var resp map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		if resp["status"] != "stopped" {
			t.Fatalf("expected stopped status got %v", resp["status"])
		}
	})

	t.Run("host maintenance", func(t *testing.T) {
		payload := []byte(`{"enabled":true}`)
		r := httptest.NewRequest(http.MethodPost, "/api/v1/hosts/host-1/maintenance", bytes.NewReader(payload))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 got %d", w.Code)
		}
	})
}

func TestAuthEnabled(t *testing.T) {
	sec := config.SecurityConfig{APIKeys: map[string]string{"viewer-key": "viewer", "admin-key": "admin"}}
	srv := New(store.New(), "", sec)
	h := srv.Router()

	r := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}

	r2 := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard", nil)
	r2.Header.Set("X-API-Key", "viewer-key")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w2.Code)
	}

	r3 := httptest.NewRequest(http.MethodPost, "/api/v1/hosts/host-1/maintenance", bytes.NewReader([]byte(`{"enabled":true}`)))
	r3.Header.Set("X-API-Key", "viewer-key")
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, r3)
	if w3.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", w3.Code)
	}

	r4 := httptest.NewRequest(http.MethodPost, "/api/v1/hosts/host-1/maintenance", bytes.NewReader([]byte(`{"enabled":true}`)))
	r4.Header.Set("X-API-Key", "admin-key")
	w4 := httptest.NewRecorder()
	h.ServeHTTP(w4, r4)
	if w4.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w4.Code)
	}
}
