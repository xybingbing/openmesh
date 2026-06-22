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

type HeartbeatConfig struct {
	ControllerURL string
	Token         string
	NodeID        string
	Version       string
	Endpoint      string
	Status        string
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
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	setToken(request, cfg.Token)

	var resp model.RegisterResponse
	if err := do(request, &resp); err != nil {
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
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	setToken(request, cfg.Token)

	var resp model.ConfigResponse
	if err := do(request, &resp); err != nil {
		return err
	}
	_, err = fmt.Fprint(out, resp.WGConfig)
	return err
}

func Heartbeat(ctx context.Context, cfg HeartbeatConfig, out io.Writer) error {
	if cfg.ControllerURL == "" || cfg.NodeID == "" {
		return fmt.Errorf("controller URL and node id are required")
	}
	if cfg.Version == "" {
		cfg.Version = "dev"
	}
	body, err := json.Marshal(model.HeartbeatRequest{Version: cfg.Version, Endpoint: cfg.Endpoint, Status: cfg.Status})
	if err != nil {
		return err
	}
	url := strings.TrimRight(cfg.ControllerURL, "/") + "/api/v1/nodes/" + cfg.NodeID + "/heartbeat"
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	setToken(request, cfg.Token)

	var resp model.HeartbeatResponse
	if err := do(request, &resp); err != nil {
		return err
	}
	if out != nil {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(resp)
	}
	return nil
}

func setToken(req *http.Request, token string) {
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
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
