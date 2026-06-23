package topology

import (
	"testing"
	"time"

	"github.com/xybingbing/openmesh/internal/model"
)

func TestBuildNodePeers(t *testing.T) {
	now := time.Now()
	nodes := []model.Node{
		{ID: "node-1", Name: "a", MeshIP: "100.64.0.2/32", Status: "online", LastSeen: now},
		{ID: "node-2", Name: "b", MeshIP: "100.64.0.3/32", Status: "online", LastSeen: now},
	}
	got, ok := BuildNodePeers(nodes, "node-1", now)
	if !ok {
		t.Fatal("expected node peers")
	}
	if got.Node.ID != "node-1" {
		t.Fatalf("unexpected node: %#v", got.Node)
	}
	if len(got.Peers) != 1 || got.Peers[0].ID != "node-2" {
		t.Fatalf("unexpected peers: %#v", got.Peers)
	}
}

func TestBuildNodePeersMissing(t *testing.T) {
	_, ok := BuildNodePeers(nil, "missing", time.Now())
	if ok {
		t.Fatal("expected missing node")
	}
}
