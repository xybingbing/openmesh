package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/xybingbing/openmesh/internal/model"
)

type RegisterConfig struct {
	ControllerURL string
	Token         string
	Name          string
	PublicKey     string
	Endpoint      string
}

type ConfigConfig struct {
	ControllerURL string
	Token         string
	NodeID        string
}

func Register(ctx context.Context, cfg RegisterConfig, out io.Writer) error {
	if cfg.ControllerURL == "" {
		return fmt.Errorf("controller URL is required")
	}
	body, err := json.Marshal(model.RegisterRequest{Name: cfg.Name, PublicKey: cfg.PublicKey, Endpoint: cfg.Endpoint})
	if err != nil {
		return err
	}

	url := strings.TrimRight(cfg.ControllerURL, "/") + "/api/v1/nodes"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.Token)
	}

	var resp model.RegisterResponse
	if err := do(req, &resp); err != nil {
		return err
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(resp)
}

func Config(ctx context.Context, cfg ConfigConfig, out io.Writer) error {
	if cfg.ControllerURL == "" || cfg.NodeID == "" {
		return fmt.Errorf("controller URL and node id are required")
	}
	url := strings.TrimRight(cfg.ControllerURL, "/") + "/api/v1/nodes/" + cfg.NodeID + "/config"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.Token)
	}

	var resp model.ConfigResponse
	if err := do(req, &resp); err != nil {
		return err
	}
	_, err = fmt.Fprint(out, resp.WGConfig)
	return err
}

func do(req *http.Request, out any) error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("request failed: %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
