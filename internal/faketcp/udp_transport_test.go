package faketcp

import (
	"context"
	"net/netip"
	"testing"
	"time"
)

func TestUDPTransportRoundTrip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	a, err := ListenUDP("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	b, err := ListenUDP("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()

	want := []byte("wg-payload")
	pkt := Packet{
		SrcIP:   netip.MustParseAddr("10.0.0.2"),
		DstIP:   netip.MustParseAddr("10.0.0.1"),
		SrcPort: 51820,
		DstPort: 443,
		Seq:     1,
		Flags:   FlagACK | FlagPSH,
		Payload: want,
	}
	if err := a.SendTo(ctx, b.LocalAddr().String(), pkt); err != nil {
		t.Fatal(err)
	}
	got, _, err := b.Recv(ctx, 65535)
	if err != nil {
		t.Fatal(err)
	}
	if string(got.Payload) != string(want) {
		t.Fatalf("payload mismatch: %q", got.Payload)
	}
}
