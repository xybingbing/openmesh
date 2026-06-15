package wg

import (
	"fmt"
	"strings"

	"github.com/xybingbing/openmesh/internal/model"
)

type RenderOptions struct {
	ListenPort int
	DNS        string
}

func RenderNodeConfig(self model.Node, peers []model.Node, opt RenderOptions) string {
	if opt.ListenPort == 0 {
		opt.ListenPort = 51820
	}

	var b strings.Builder
	fmt.Fprintf(&b, "[Interface]\n")
	fmt.Fprintf(&b, "Address = %s\n", self.MeshIP)
	fmt.Fprintf(&b, "ListenPort = %d\n", opt.ListenPort)
	fmt.Fprintf(&b, "PrivateKey = <node-private-key>\n")
	if opt.DNS != "" {
		fmt.Fprintf(&b, "DNS = %s\n", opt.DNS)
	}

	for _, p := range peers {
		if p.ID == self.ID {
			continue
		}
		fmt.Fprintf(&b, "\n[Peer]\n")
		fmt.Fprintf(&b, "# %s (%s)\n", p.Name, p.ID)
		fmt.Fprintf(&b, "PublicKey = %s\n", p.PublicKey)
		fmt.Fprintf(&b, "AllowedIPs = %s\n", p.MeshIP)
		if p.Endpoint != "" {
			fmt.Fprintf(&b, "Endpoint = %s\n", p.Endpoint)
			fmt.Fprintf(&b, "PersistentKeepalive = 25\n")
		}
	}
	return b.String()
}
