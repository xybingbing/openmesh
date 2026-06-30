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
	"github.com/xybingbing/openmesh/internal/topology"
	"github.com/xybingbing/openmesh/internal/wg"
)

// ... unchanged above omitted for brevity in patch context ...

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

	// NAT endpoint auto-discovery (use source IP:port if agent did not provide one)
	if req.Endpoint == "" {
		// r.RemoteAddr is usually ip:port
		req.Endpoint = r.RemoteAddr
	}

	n, err := s.store.Heartbeat(nodeID, req)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	writeJSON(w, http.StatusOK, model.HeartbeatResponse{Node: n})
}
