package faketcp

import (
	"bytes"
	"net/netip"
	"testing"
)

func TestEncodeDecodePacket(t *testing.T) {
	payload := []byte("hello-wireguard")
	in := Packet{
		SrcIP:   netip.MustParseAddr("10.0.0.2"),
		DstIP:   netip.MustParseAddr("10.0.0.1"),
		SrcPort: 40000,
		DstPort: 443,
		Seq:     100,
		Ack:     42,
		Flags:   FlagACK | FlagPSH,
		Payload: payload,
	}
	b, err := Encode(in)
	if err != nil {
		t.Fatal(err)
	}
	out, err := Decode(b)
	if err != nil {
		t.Fatal(err)
	}
	if out.SrcIP != in.SrcIP || out.DstIP != in.DstIP || out.SrcPort != in.SrcPort || out.DstPort != in.DstPort {
		t.Fatalf("decoded tuple mismatch: %#v", out)
	}
	if out.Seq != in.Seq || out.Ack != in.Ack || out.Flags != in.Flags {
		t.Fatalf("decoded TCP state mismatch: %#v", out)
	}
	if !bytes.Equal(out.Payload, payload) {
		t.Fatalf("payload mismatch: %q", string(out.Payload))
	}
}

func TestEncodeRejectsIPv6(t *testing.T) {
	_, err := Encode(Packet{SrcIP: netip.MustParseAddr("2001:db8::1"), DstIP: netip.MustParseAddr("10.0.0.1")})
	if err == nil {
		t.Fatal("expected IPv6 rejection")
	}
}
