package agent

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type DaemonConfig struct {
	ConfigPath string
	Interval   time.Duration
	Once       bool
}

func Apply(ctx context.Context, cfg LocalConfig) error {
	var buf bytes.Buffer
	if err := Config(ctx, ConfigConfig{ControllerURL: cfg.ControllerURL, Token: cfg.Token, NodeID: cfg.NodeID}, &buf); err != nil {
		return err
	}
	if err := writeWGConfig(cfg.WGConfigPath, buf.Bytes()); err != nil {
		return err
	}
	if cfg.SyncCommand != "" {
		cmd := exec.CommandContext(ctx, "sh", "-c", cfg.SyncCommand)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("sync command failed: %w", err)
		}
	}
	return nil
}

func Daemon(ctx context.Context, dc DaemonConfig) error {
	if dc.ConfigPath == "" {
		dc.ConfigPath = "/etc/openmesh/agent.json"
	}
	if dc.Interval == 0 {
		dc.Interval = 30 * time.Second
	}

	cfg, err := LoadLocalConfig(dc.ConfigPath)
	if err != nil {
		return err
	}
	for {
		if err := Apply(ctx, cfg); err != nil {
			fmt.Fprintln(os.Stderr, "agent apply failed:", err)
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
