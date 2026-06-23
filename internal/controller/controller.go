package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/xybingbing/openmesh/internal/model"
	"github.com/xybingbing/openmesh/internal/status"
	"github.com/xybingbing/openmesh/internal/store"
	"github.com/xybingbing/openmesh/internal/wg"
)

type Config struct {
	Listen   string
	DataPath string
	Token    string
}

type Server struct {
	cfg   Config
	store *store.Store
}

func Run(ctx context.Context, cfg Config) error {
	if cfg.Listen == "" {
		cfg.Listen = ":8080"
	}
	st, err := store.Open(cfg.DataPath)
	if err != nil {
		return err
	}
	s := &Server{cfg: cfg, store: st}

	h := http.NewServeMux()
	h.HandleFunc("/healthz", s.healthz)
	h.HandleFunc("/api/v1/status", s.auth(s.statusAll))
	h.HandleFunc("/api/v1/nodes", s.auth(s.nodes))
	h.HandleFunc("/api/v1/nodes/", s.auth(s.nodeRoute))

	srv := &http.Server{Addr: cfg.Listen, Handler: h}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()
	fmt.Println("openmesh controller listening on", cfg.Listen)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.cfg.Token != "" {
			got := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if got != s.cfg.Token {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}

func (s *Server) healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) statusAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"nodes": status.EvaluateNodes(s.store.Nodes(), time.Now().UTC(), status.DefaultOfflineAfter)})
}

func (s *Server) nodes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, map[string]any{"nodes": s.store.Nodes()})
	case http.MethodPost:
		var req model.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		n, err := s.store.Register(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, model.RegisterResponse{Node: n})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) nodeRoute(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/nodes/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}
	switch parts[1] {
	case "config":
		s.nodeConfig(w, r, parts[0])
	case "heartbeat":
		s.nodeHeartbeat(w, r, parts[0])
	case "status":
		s.nodeStatus(w, r, parts[0])
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) nodeConfig(w http.ResponseWriter, r *http.Request, nodeID string) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	n, ok := s.store.Node(nodeID)
	if !ok {
		http.NotFound(w, r)
		return
	}
	cfg := wg.RenderNodeConfig(n, s.store.Nodes(), wg.RenderOptions{})
	writeJSON(w, http.StatusOK, model.ConfigResponse{Node: n, WGConfig: cfg, Generated: time.Now().UTC().Format(time.RFC3339)})
}

func (s *Server) nodeHeartbeat(w http.ResponseWriter, r *http.Request, nodeID string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req model.HeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	n, err := s.store.Heartbeat(nodeID, req)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, model.HeartbeatResponse{Node: n})
}

func (s *Server) nodeStatus(w http.ResponseWriter, r *http.Request, nodeID string) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	n, ok := s.store.Node(nodeID)
	if !ok {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, status.EvaluateNode(n, time.Now().UTC(), status.DefaultOfflineAfter))
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
