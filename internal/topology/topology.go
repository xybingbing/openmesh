package topology

import (
	"time"

	"github.com/xybingbing/openmesh/internal/model"
	"github.com/xybingbing/openmesh/internal/status"
)

type Peer struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	MeshIP   string `json:"mesh_ip"`
	Endpoint string `json:"endpoint,omitempty"`
	Status   string `json:"status"`
	Version  string `json:"version,omitempty"`
}

type NodePeers struct {
	Node  Peer   `json:"node"`
	Peers []Peer `json:"peers"`
}

func BuildNodePeers(nodes []model.Node, nodeID string, now time.Time) (NodePeers, bool) {
	statuses := make(map[string]string, len(nodes))
	for _, ns := range status.EvaluateNodes(nodes, now, status.DefaultOfflineAfter) {
		statuses[ns.NodeID] = ns.Status
	}

	var self model.Node
	found := false
	for _, n := range nodes {
		if n.ID == nodeID {
			self = n
			found = true
			break
		}
	}
	if !found {
		return NodePeers{}, false
	}

	out := NodePeers{Node: toPeer(self, statuses[self.ID])}
	for _, n := range nodes {
		if n.ID == nodeID {
			continue
		}
		out.Peers = append(out.Peers, toPeer(n, statuses[n.ID]))
	}
	return out, true
}

func BuildAllPeers(nodes []model.Node, now time.Time) []NodePeers {
	out := make([]NodePeers, 0, len(nodes))
	for _, n := range nodes {
		peers, ok := BuildNodePeers(nodes, n.ID, now)
		if ok {
			out = append(out, peers)
		}
	}
	return out
}

func toPeer(n model.Node, st string) Peer {
	if st == "" {
		st = "unknown"
	}
	return Peer{ID: n.ID, Name: n.Name, MeshIP: n.MeshIP, Endpoint: n.Endpoint, Status: st, Version: n.Version}
}
