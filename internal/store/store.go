package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/netip"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/xybingbing/openmesh/internal/model"
)

type Store struct {
	mu    sync.Mutex
	path  string
	state State
}

type State struct {
	NextIP uint32       `json:"next_ip"`
	Nodes  []model.Node `json:"nodes"`
}

func Open(path string) (*Store, error) {
	s := &Store{path: path, state: State{NextIP: 2}}
	b, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(b, &s.state); err != nil {
			return nil, fmt.Errorf("decode store: %w", err)
		}
		if s.state.NextIP == 0 {
			s.state.NextIP = 2
		}
		return s, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return s, nil
}

func (s *Store) Register(req model.RegisterRequest) (model.Node, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if req.Name == "" {
		return model.Node{}, fmt.Errorf("name is required")
	}
	if req.PublicKey == "" {
		return model.Node{}, fmt.Errorf("public key is required")
	}

	now := time.Now().UTC()
	for i := range s.state.Nodes {
		if s.state.Nodes[i].PublicKey == req.PublicKey {
			s.state.Nodes[i].Name = req.Name
			s.state.Nodes[i].Endpoint = req.Endpoint
			s.state.Nodes[i].UpdatedAt = now
			return s.state.Nodes[i], s.saveLocked()
		}
	}

	n := model.Node{
		ID:        fmt.Sprintf("node-%d", len(s.state.Nodes)+1),
		Name:      req.Name,
		PublicKey: req.PublicKey,
		MeshIP:    meshIP(s.state.NextIP),
		Endpoint:  req.Endpoint,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.state.NextIP++
	s.state.Nodes = append(s.state.Nodes, n)
	return n, s.saveLocked()
}

func (s *Store) Nodes() []model.Node {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]model.Node, len(s.state.Nodes))
	copy(out, s.state.Nodes)
	return out
}

func (s *Store) Node(id string) (model.Node, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, n := range s.state.Nodes {
		if n.ID == id {
			return n, true
		}
	}
	return model.Node{}, false
}

func (s *Store) saveLocked() error {
	dir := filepath.Dir(s.path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	b, err := json.MarshalIndent(s.state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o600)
}

func meshIP(offset uint32) string {
	addr := netip.AddrFrom4([4]byte{100, 64, byte(offset >> 8), byte(offset)})
	return addr.String() + "/32"
}
