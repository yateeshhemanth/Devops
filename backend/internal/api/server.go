package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"devops/backend/internal/store"
)

type Server struct {
	store     *store.Store
	staticDir string
}

func New(s *store.Store, staticDir string) *Server { return &Server{store: s, staticDir: staticDir} }

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.healthz)
	mux.HandleFunc("/api/v1/dashboard", s.dashboard)
	mux.HandleFunc("/api/v1/vms/", s.vmRoutes)
	mux.HandleFunc("/api/v1/storage/volumes/", s.volumeRoutes)
	mux.HandleFunc("/api/v1/hosts/", s.hostRoutes)
	if s.staticDir != "" {
		mux.HandleFunc("/", s.serveSPA)
	}
	return withCORS(logging(mux))
}

func (s *Server) healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) dashboard(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.store.Dashboard())
}

func (s *Server) vmRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/vms/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		http.NotFound(w, r)
		return
	}
	vmID := parts[0]
	action := parts[1]

	switch action {
	case "power":
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}
		var req struct {
			Action string `json:"action"`
		}
		if !decodeJSON(w, r, &req) {
			return
		}
		vm, err := s.store.SetPower(vmID, req.Action)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, vm)
	case "migrate":
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}
		var req struct {
			Host string `json:"host"`
		}
		if !decodeJSON(w, r, &req) {
			return
		}
		vm, err := s.store.MigrateVM(vmID, req.Host)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, vm)
	case "snapshot":
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}
		vm, err := s.store.SnapshotVM(vmID)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, vm)
	case "console-ticket":
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"ticket":     "nvnc-" + vmID + "-" + time.Now().Format("150405"),
			"expiresInS": 120,
			"wsURL":      "wss://console.platform.local/ws",
		})
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) volumeRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/storage/volumes/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[1] != "expand" || r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	var req struct {
		SizeGB int `json:"sizeGb"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	volume, err := s.store.ExpandVolume(parts[0], req.SizeGB)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, volume)
}

func (s *Server) hostRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/hosts/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[1] != "maintenance" || r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	host, err := s.store.SetHostMaintenance(parts[0], req.Enabled)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, host)
}

func (s *Server) serveSPA(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		http.NotFound(w, r)
		return
	}

	relPath := filepath.Clean(strings.TrimPrefix(r.URL.Path, "/"))
	if relPath == "." {
		relPath = "index.html"
	}

	candidate := filepath.Join(s.staticDir, relPath)
	if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
		http.ServeFile(w, r, candidate)
		return
	}

	indexPath := filepath.Join(s.staticDir, "index.html")
	if _, err := os.Stat(indexPath); err == nil {
		http.ServeFile(w, r, indexPath)
		return
	}

	http.NotFound(w, r)
}

func writeJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, out any) bool {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return false
	}
	return true
}

func methodNotAllowed(w http.ResponseWriter) {
	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
