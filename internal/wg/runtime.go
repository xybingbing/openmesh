package wg

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xybingbing/openmesh/internal/system"
)

type RuntimeConfig struct {
	Interface string
	Config    string
	Address   string
	MTU       int
}

type Runtime struct {
	Runner system.Runner
}

func (r Runtime) Up(ctx context.Context, cfg RuntimeConfig) error {
	if cfg.Interface == "" {
		return fmt.Errorf("interface is required")
	}
	if cfg.Config == "" {
		return fmt.Errorf("WireGuard config path is required")
	}
	runner := r.Runner
	if runner == nil {
		runner = system.ExecRunner{}
	}

	if err := runner.Run(ctx, "ip", "link", "add", cfg.Interface, "type", "wireguard"); err != nil {
		return err
	}
	if cfg.MTU > 0 {
		if err := runner.Run(ctx, "ip", "link", "set", "dev", cfg.Interface, "mtu", fmt.Sprint(cfg.MTU)); err != nil {
			return err
		}
	}
	if err := runner.Run(ctx, "wg", "setconf", cfg.Interface, cfg.Config); err != nil {
		return err
	}
	if cfg.Address != "" {
		if err := runner.Run(ctx, "ip", "addr", "add", cfg.Address, "dev", cfg.Interface); err != nil {
			return err
		}
	}
	return runner.Run(ctx, "ip", "link", "set", "up", "dev", cfg.Interface)
}

func (r Runtime) Down(ctx context.Context, iface string) error {
	if iface == "" {
		return fmt.Errorf("interface is required")
	}
	runner := r.Runner
	if runner == nil {
		runner = system.ExecRunner{}
	}
	return runner.Run(ctx, "ip", "link", "del", iface)
}

func WriteConfig(path string, content string) error {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return err
		}
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(content), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
