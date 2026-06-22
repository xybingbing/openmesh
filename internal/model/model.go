package model

import "time"

type Node struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	PublicKey string    `json:"public_key"`
	MeshIP    string    `json:"mesh_ip"`
	Endpoint  string    `json:"endpoint,omitempty"`
	Status    string    `json:"status,omitempty"`
	Version   string    `json:"version,omitempty"`
	LastSeen  time.Time `json:"last_seen,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
	Endpoint  string `json:"endpoint,omitempty"`
}

type RegisterResponse struct {
	Node Node `json:"node"`
}

type HeartbeatRequest struct {
	Version  string `json:"version"`
	Endpoint string `json:"endpoint,omitempty"`
	Status   string `json:"status,omitempty"`
}

type HeartbeatResponse struct {
	Node Node `json:"node"`
}

type ConfigResponse struct {
	Node      Node   `json:"node"`
	WGConfig  string `json:"wg_config"`
	Generated string `json:"generated"`
}
