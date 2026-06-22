package agent

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/xybingbing/openmesh/internal/system"
	"github.com/xybingbing/openmesh/internal/wg"
)

type DaemonConfig struct {
	ConfigPath string
	Interval   time.Duration
	Once       bool
	Version    string
}

func Apply(ctx context.Context, cfg LocalConfig) error {
	_, err := ApplyAndReturnConfig(ctx, cfg)
	return err
}

func ApplyAndReturnConfig(ctx context.Context, cfg LocalConfig) (string, error) {
	applyDefaults(&cfg)
	var buf bytes.Buffer
	if err := Config(ctx, ConfigConfig{ControllerURL: cfg.ControllerURL, Token: cfg.Token, NodeID: cfg.NodeID}, &buf); err != nil {
		return "", err
	}
	wgConfig := buf.String()
	if err := writeWGConfig(cfg.WGConfigPath, []byte(wgConfig)); err != nil {
		return "", err
	}
	if cfg.SyncCommand != "" {
		cmd := exec.CommandContext(ctx, "sh", "-c", cfg.SyncCommand)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("sync command failed: %w", err)
		}
	}
	return wgConfig, nil
}

func Up(ctx context.Context, cfg LocalConfig, runner system.Runner) error {
	wgConfig, err := ApplyAndReturnConfig(ctx, cfg)
	if err != nil {
		return err
	}
	if cfg.WGAddress == "" {
		cfg.WGAddress = extractAddress(wgConfig)
	}
	return wg.Runtime{Runner: runner}.Up(ctx, wg.RuntimeConfig{Interface: cfg.WGInterface, Config: cfg.WGConfigPath, Address: cfg.WGAddress, MTU: cfg.WGMTU})
}

func Down(ctx context.Context, cfg LocalConfig, runner system.Runner) error {
	applyDefaults(&cfg)
	return wg.Runtime{Runner: runner}.Down(ctx, cfg.WGInterface)
}

func Daemon(ctx context.Context, dc DaemonConfig) error {
	if dc.ConfigPath == "" {
		dc.ConfigPath = "/etc/openmesh/agent.json"
	}
	if dc.Interval == 0 {
		dc.Interval = 30 * time.Second
	}
	if dc.Version == "" {
		dc.Version = "dev"
	}

	cfg, err := LoadLocalConfig(dc.ConfigPath)
	if err != nil {
		return err
	}
	for {
		if err := Apply(ctx, cfg); err != nil {
			fmt.Fprintln(os.Stderr, "agent apply failed:", err)
		}
		if err := Heartbeat(ctx, HeartbeatConfig{ControllerURL: cfg.ControllerURL, Token: cfg.Token, NodeID: cfg.NodeID, Version: dc.Version, Endpoint: "", Status: "online"}, nil); err != nil {
			fmt.Fprintln(os.Stderr, "agent heartbeat failed:", err)
		}
		if dc.Once {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dc.Interval):
		}
	}
}

func extractAddress(wgConfig string) string {
	for _, line := range strings.Split(wgConfig, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Address =") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Address ="))
		}
	}
	return ""
}

func writeWGConfig(path string, b []byte) error {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return err
		}
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
