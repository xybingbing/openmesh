package faketcp

// Package faketcp will contain the UDP-to-fake-TCP transport engine.
//
// Milestone 1 keeps this as a package boundary so the controller and agent can
// be developed and tested first. The real transport implementation will add:
//
//   - IPv4/TCP packet encode/decode
//   - checksum validation
//   - WireGuard UDP payload forwarding
//   - raw socket or TUN based transport
//   - OpenWrt firewall integration

type Mode string

const (
	ModeDisabled Mode = "disabled"
	ModeClient   Mode = "client"
	ModeServer   Mode = "server"
)

type Config struct {
	Mode   Mode
	Listen string
	Remote string
	Target string
}

func Enabled(cfg Config) bool {
	return cfg.Mode == ModeClient || cfg.Mode == ModeServer
}
