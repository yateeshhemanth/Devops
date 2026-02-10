package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"devops/backend/internal/config"
	"devops/backend/internal/store"
)

type Server struct {
	store     *store.Store
	staticDir string
	security  config.SecurityConfig
	limiter   *ipRateLimiter
}

func New(s *store.Store, staticDir string, security config.SecurityConfig) *Server {
	return &Server{store: s, staticDir: staticDir, security: security, limiter: newIPRateLimiter(80, time.Minute)}
}

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.healthz)
	mux.HandleFunc("/readyz", s.readyz)
	mux.HandleFunc("/api/v1/dashboard", s.dashboard)
	mux.HandleFunc("/api/v1/vms/", s.vmRoutes)
	mux.HandleFunc("/api/v1/storage/volumes/", s.volumeRoutes)
	mux.HandleFunc("/api/v1/hosts/", s.hostRoutes)
	if s.staticDir != "" {
		mux.HandleFunc("/", s.serveSPA)
	}

	return withRecovery(withRateLimit(s.limiter, withCORS(logging(mux))))
}

func (s *Server) healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) readyz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	if !s.authorize(w, r, "viewer") {
		return
	}
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
	case "power", "migrate", "snapshot", "console-ticket":
		if !s.authorize(w, r, "operator") {
			return
		}
	default:
		http.NotFound(w, r)
		return
	}

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
		audit(r, "vm.power", vmID)
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
		audit(r, "vm.migrate", vmID)
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
		audit(r, "vm.snapshot", vmID)
		writeJSON(w, http.StatusOK, vm)
	case "console-ticket":
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}
		audit(r, "vm.console", vmID)
		writeJSON(w, http.StatusOK, map[string]any{
			"ticket":     "nvnc-" + vmID + "-" + time.Now().Format("150405"),
			"expiresInS": 120,
			"wsURL":      "wss://console.platform.local/ws",
		})
	}
}

func (s *Server) volumeRoutes(w http.ResponseWriter, r *http.Request) {
	if !s.authorize(w, r, "operator") {
		return
	}
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
	audit(r, "volume.expand", parts[0])
	writeJSON(w, http.StatusOK, volume)
}

func (s *Server) hostRoutes(w http.ResponseWriter, r *http.Request) {
	if !s.authorize(w, r, "admin") {
		return
	}
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
	audit(r, "host.maintenance", parts[0])
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

func (s *Server) authorize(w http.ResponseWriter, r *http.Request, required string) bool {
	if len(s.security.APIKeys) == 0 {
		return true
	}
	token := strings.TrimSpace(r.Header.Get("X-API-Key"))
	role, ok := s.security.APIKeys[token]
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return false
	}
	if !roleAllowed(role, required) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
		return false
	}
	return true
}

func roleAllowed(role, required string) bool {
	levels := map[string]int{"viewer": 1, "operator": 2, "admin": 3}
	return levels[role] >= levels[required]
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
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,X-API-Key")
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
		requestID := newRequestID()
		w.Header().Set("X-Request-ID", requestID)
		log.Printf("request_id=%s method=%s path=%s", requestID, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic recovered: %v", rec)
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type ipRateLimiter struct {
	mu        sync.Mutex
	window    time.Duration
	limit     int
	requests  map[string]int
	lastSweep time.Time
}

func newIPRateLimiter(limit int, window time.Duration) *ipRateLimiter {
	return &ipRateLimiter{window: window, limit: limit, requests: map[string]int{}, lastSweep: time.Now()}
}

func withRateLimit(l *ipRateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r.RemoteAddr)
		if !l.allow(ip) {
			writeJSON(w, http.StatusTooManyRequests, map[string]string{"error": "rate limit exceeded"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (l *ipRateLimiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if time.Since(l.lastSweep) > l.window {
		l.requests = map[string]int{}
		l.lastSweep = time.Now()
	}
	l.requests[key]++
	return l.requests[key] <= l.limit
}

func clientIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return host
}

func newRequestID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return time.Now().Format("150405.000")
	}
	return hex.EncodeToString(buf)
}

func audit(r *http.Request, action, target string) {
	actor := r.Header.Get("X-API-Key")
	if actor == "" {
		actor = "anonymous"
	}
	log.Printf("audit actor=%s action=%s target=%s", actor, action, target)
}
