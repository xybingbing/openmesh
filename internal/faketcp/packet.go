package faketcp

import (
	"encoding/binary"
	"fmt"
	"net/netip"
)

const (
	ipHeaderLen  = 20
	tcpHeaderLen = 20
	protoTCP     = 6
)

const (
	FlagFIN uint8 = 0x01
	FlagSYN uint8 = 0x02
	FlagRST uint8 = 0x04
	FlagPSH uint8 = 0x08
	FlagACK uint8 = 0x10
)

type Packet struct {
	SrcIP   netip.Addr
	DstIP   netip.Addr
	SrcPort uint16
	DstPort uint16
	Seq     uint32
	Ack     uint32
	Flags   uint8
	TTL     uint8
	Payload []byte
}

func Encode(p Packet) ([]byte, error) {
	if !p.SrcIP.Is4() || !p.DstIP.Is4() {
		return nil, fmt.Errorf("only IPv4 is supported")
	}
	if p.TTL == 0 {
		p.TTL = 64
	}
	totalLen := ipHeaderLen + tcpHeaderLen + len(p.Payload)
	b := make([]byte, totalLen)

	b[0] = 0x45
	b[1] = 0
	binary.BigEndian.PutUint16(b[2:4], uint16(totalLen))
	binary.BigEndian.PutUint16(b[4:6], 0)
	binary.BigEndian.PutUint16(b[6:8], 0x4000)
	b[8] = p.TTL
	b[9] = protoTCP
	copy(b[12:16], p.SrcIP.AsSlice())
	copy(b[16:20], p.DstIP.AsSlice())
	binary.BigEndian.PutUint16(b[10:12], checksum(b[:ipHeaderLen]))

	tcp := b[ipHeaderLen:]
	binary.BigEndian.PutUint16(tcp[0:2], p.SrcPort)
	binary.BigEndian.PutUint16(tcp[2:4], p.DstPort)
	binary.BigEndian.PutUint32(tcp[4:8], p.Seq)
	binary.BigEndian.PutUint32(tcp[8:12], p.Ack)
	tcp[12] = byte(tcpHeaderLen/4) << 4
	tcp[13] = p.Flags
	binary.BigEndian.PutUint16(tcp[14:16], 65535)
	copy(tcp[tcpHeaderLen:], p.Payload)
	binary.BigEndian.PutUint16(tcp[16:18], tcpChecksum(p.SrcIP, p.DstIP, tcp))

	return b, nil
}

func Decode(b []byte) (Packet, error) {
	if len(b) < ipHeaderLen+tcpHeaderLen {
		return Packet{}, fmt.Errorf("packet too short")
	}
	if b[0]>>4 != 4 {
		return Packet{}, fmt.Errorf("not IPv4")
	}
	ihl := int(b[0]&0x0f) * 4
	if ihl < ipHeaderLen || len(b) < ihl+tcpHeaderLen {
		return Packet{}, fmt.Errorf("invalid IPv4 header length")
	}
	if b[9] != protoTCP {
		return Packet{}, fmt.Errorf("not TCP")
	}
	totalLen := int(binary.BigEndian.Uint16(b[2:4]))
	if totalLen == 0 || totalLen > len(b) {
		totalLen = len(b)
	}

	srcIP := netip.AddrFrom4([4]byte{b[12], b[13], b[14], b[15]})
	dstIP := netip.AddrFrom4([4]byte{b[16], b[17], b[18], b[19]})
	tcp := b[ihl:totalLen]
	dataOffset := int(tcp[12]>>4) * 4
	if dataOffset < tcpHeaderLen || len(tcp) < dataOffset {
		return Packet{}, fmt.Errorf("invalid TCP data offset")
	}

	payload := make([]byte, len(tcp[dataOffset:]))
	copy(payload, tcp[dataOffset:])
	return Packet{
		SrcIP:   srcIP,
		DstIP:   dstIP,
		SrcPort: binary.BigEndian.Uint16(tcp[0:2]),
		DstPort: binary.BigEndian.Uint16(tcp[2:4]),
		Seq:     binary.BigEndian.Uint32(tcp[4:8]),
		Ack:     binary.BigEndian.Uint32(tcp[8:12]),
		Flags:   tcp[13],
		TTL:     b[8],
		Payload: payload,
	}, nil
}

func checksum(b []byte) uint16 {
	var sum uint32
	for len(b) >= 2 {
		sum += uint32(binary.BigEndian.Uint16(b[:2]))
		b = b[2:]
	}
	if len(b) == 1 {
		sum += uint32(b[0]) << 8
	}
	for sum>>16 != 0 {
		sum = (sum & 0xffff) + (sum >> 16)
	}
	return ^uint16(sum)
}

func tcpChecksum(src, dst netip.Addr, tcp []byte) uint16 {
	pseudo := make([]byte, 12+len(tcp))
	copy(pseudo[0:4], src.AsSlice())
	copy(pseudo[4:8], dst.AsSlice())
	pseudo[8] = 0
	pseudo[9] = protoTCP
	binary.BigEndian.PutUint16(pseudo[10:12], uint16(len(tcp)))
	copy(pseudo[12:], tcp)
	return checksum(pseudo)
}
