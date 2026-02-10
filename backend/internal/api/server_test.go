package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"devops/backend/internal/store"
)

func TestDashboardAndOps(t *testing.T) {
	srv := New(store.New())
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
