# FakeTCP

OpenMesh uses FakeTCP as an optional transport for UDP-restricted networks.

Current implementation status:

- IPv4/TCP packet encoder
- IPv4/TCP packet decoder
- TCP checksum and IPv4 checksum generation
- Payload round-trip tests
- UDP-based transport for local development tests

The UDP transport is not the final production transport. It exists so the packet layer can be tested without root privileges.

Planned production transports:

- raw socket mode for Linux/OpenWrt
- TUN mode for Phantun-like deployment

Design rule:

WireGuard itself is not modified. OpenMesh transports WireGuard UDP payloads outside WireGuard and can later choose direct UDP, FakeTCP, or relay paths.
