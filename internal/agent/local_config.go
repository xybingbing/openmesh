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
	if cfg.WGConfigPath == "" {
		cfg.WGConfigPath = "/etc/wireguard/openmesh.conf"
	}
	return cfg, nil
}

func SaveLocalConfig(path string, cfg LocalConfig) error {
	if cfg.ControllerURL == "" {
		return fmt.Errorf("controller URL is required")
	}
	if cfg.NodeID == "" {
		return fmt.Errorf("node id is required")
	}
	if cfg.WGConfigPath == "" {
		cfg.WGConfigPath = "/etc/wireguard/openmesh.conf"
	}
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
