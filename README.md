# OpenMesh

OpenMesh is a self-hosted, OpenWrt-first mesh VPN project.

Goal:

```text
Tailscale + Phantun + OpenWrt
```

The project is built around:

- WireGuard as the encrypted data plane
- a small controller for node registration and peer config generation
- an agent for OpenWrt/Linux nodes
- a FakeTCP transport engine for UDP-restricted networks

## Current status

This repository is now initialized as the real engineering project.

The first milestone focuses on a minimal working system:

- Controller HTTP API
- Node registration
- Mesh IP allocation
- WireGuard config generation
- Agent registration/config pull
- FakeTCP package boundary

## Build

```bash
go test ./...
go build -o openmesh ./cmd/openmesh
```

## Run controller

```bash
./openmesh controller --listen :8080 --data ./openmesh.json --token dev-token
```

## Register node

```bash
./openmesh agent register \
  --controller http://127.0.0.1:8080 \
  --token dev-token \
  --name node-a \
  --public-key example-public-key
```

## Pull WireGuard config

```bash
./openmesh agent config \
  --controller http://127.0.0.1:8080 \
  --token dev-token \
  --node-id <node-id>
```

## License

TBD.
