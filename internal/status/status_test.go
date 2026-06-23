package status

import (
	"testing"
	"time"

	"github.com/xybingbing/openmesh/internal/model"
)

func TestEvaluateNodeOfflineWithoutHeartbeat(t *testing.T) {
	got := EvaluateNode(model.Node{ID: "node-1", Name: "a", Status: "registered"}, time.Now(), DefaultOfflineAfter)
	if got.Status != "offline" {
		t.Fatalf("expected offline, got %s", got.Status)
	}
}

func TestEvaluateNodeOnlineWithRecentHeartbeat(t *testing.T) {
	now := time.Now()
	got := EvaluateNode(model.Node{ID: "node-1", Name: "a", Status: "online", LastSeen: now.Add(-10 * time.Second)}, now, DefaultOfflineAfter)
	if got.Status != "online" {
		t.Fatalf("expected online, got %s", got.Status)
	}
}

func TestEvaluateNodeOfflineAfterTimeout(t *testing.T) {
	now := time.Now()
	got := EvaluateNode(model.Node{ID: "node-1", Name: "a", Status: "online", LastSeen: now.Add(-2 * time.Minute)}, now, DefaultOfflineAfter)
	if got.Status != "offline" {
		t.Fatalf("expected offline, got %s", got.Status)
	}
}
