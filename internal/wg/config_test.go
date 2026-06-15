package wg

import (
	"strings"
	"testing"

	"github.com/xybingbing/openmesh/internal/model"
)

func TestRenderNodeConfig(t *testing.T) {
	self := model.Node{ID: "node-1", Name: "a", PublicKey: "pub-a", MeshIP: "100.64.0.2/32"}
	peer := model.Node{ID: "node-2", Name: "b", PublicKey: "pub-b", MeshIP: "100.64.0.3/32", Endpoint: "1.2.3.4:51820"}
	cfg := RenderNodeConfig(self, []model.Node{self, peer}, RenderOptions{})
	for _, want := range []string{"[Interface]", "Address = 100.64.0.2/32", "[Peer]", "PublicKey = pub-b", "Endpoint = 1.2.3.4:51820"} {
		if !strings.Contains(cfg, want) {
			t.Fatalf("config missing %q:\n%s", want, cfg)
		}
	}
}
