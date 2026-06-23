package status

import (
	"time"

	"github.com/xybingbing/openmesh/internal/model"
)

const DefaultOfflineAfter = 90 * time.Second

type NodeStatus struct {
	NodeID       string    `json:"node_id"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	Version      string    `json:"version,omitempty"`
	Endpoint     string    `json:"endpoint,omitempty"`
	MeshIP       string    `json:"mesh_ip"`
	LastSeen     time.Time `json:"last_seen,omitempty"`
	OfflineAfter string    `json:"offline_after"`
}

func EvaluateNode(n model.Node, now time.Time, offlineAfter time.Duration) NodeStatus {
	if offlineAfter <= 0 {
		offlineAfter = DefaultOfflineAfter
	}
	state := n.Status
	if state == "" {
		state = "unknown"
	}
	if n.LastSeen.IsZero() {
		if state == "registered" || state == "unknown" {
			state = "offline"
		}
	} else if now.Sub(n.LastSeen) > offlineAfter {
		state = "offline"
	} else if state == "" || state == "registered" {
		state = "online"
	}
	return NodeStatus{
		NodeID:       n.ID,
		Name:         n.Name,
		Status:       state,
		Version:      n.Version,
		Endpoint:     n.Endpoint,
		MeshIP:       n.MeshIP,
		LastSeen:     n.LastSeen,
		OfflineAfter: offlineAfter.String(),
	}
}

func EvaluateNodes(nodes []model.Node, now time.Time, offlineAfter time.Duration) []NodeStatus {
	out := make([]NodeStatus, 0, len(nodes))
	for _, n := range nodes {
		out = append(out, EvaluateNode(n, now, offlineAfter))
	}
	return out
}
