package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type LocalConfig struct {
	ControllerURL string `json:"controller_url"`
	Token         string `json:"token"`
	NodeID        string `json:"node_id"`
	WGConfigPath  string `json:"wg_config_path"`
	WGInterface   string `json:"wg_interface"`
	WGAddress     string `json:"wg_address"`
	WGMTU         int    `json:"wg_mtu"`
	SyncCommand   string `json:"sync_command"`
}

func LoadLocalConfig(path string) (LocalConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return LocalConfig{}, err
	}
	var cfg LocalConfig
	if err := json.Unmarshal(b, &cfg); err != nil {
		return LocalConfig{}, err
	}
	applyDefaults(&cfg)
	return cfg, nil
}

func SaveLocalConfig(path string, cfg LocalConfig) error {
	if cfg.ControllerURL == "" {
		return fmt.Errorf("controller URL is required")
	}
	if cfg.NodeID == "" {
		return fmt.Errorf("node id is required")
	}
	applyDefaults(&cfg)
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return err
		}
	}
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o600)
}

func applyDefaults(cfg *LocalConfig) {
	if cfg.WGConfigPath == "" {
		cfg.WGConfigPath = "/etc/wireguard/openmesh.conf"
	}
	if cfg.WGInterface == "" {
		cfg.WGInterface = "openmesh"
	}
	if cfg.WGMTU == 0 {
		cfg.WGMTU = 1280
	}
}
