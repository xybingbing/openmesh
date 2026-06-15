package faketcp

import (
	"context"
	"fmt"
	"net"
	"time"
)

type UDPTransport struct {
	conn *net.UDPConn
}

func ListenUDP(addr string) (*UDPTransport, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	return &UDPTransport{conn: conn}, nil
}

func (t *UDPTransport) Close() error {
	return t.conn.Close()
}

func (t *UDPTransport) LocalAddr() net.Addr {
	return t.conn.LocalAddr()
}

func (t *UDPTransport) SendTo(ctx context.Context, remote string, p Packet) error {
	addr, err := net.ResolveUDPAddr("udp", remote)
	if err != nil {
		return err
	}
	b, err := Encode(p)
	if err != nil {
		return err
	}
	if deadline, ok := ctx.Deadline(); ok {
		_ = t.conn.SetWriteDeadline(deadline)
	} else {
		_ = t.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	}
	n, err := t.conn.WriteToUDP(b, addr)
	if err != nil {
		return err
	}
	if n != len(b) {
		return fmt.Errorf("short write: %d/%d", n, len(b))
	}
	return nil
}

func (t *UDPTransport) Recv(ctx context.Context, maxSize int) (Packet, net.Addr, error) {
	if maxSize <= 0 {
		maxSize = 65535
	}
	if deadline, ok := ctx.Deadline(); ok {
		_ = t.conn.SetReadDeadline(deadline)
	} else {
		_ = t.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	}
	buf := make([]byte, maxSize)
	n, addr, err := t.conn.ReadFrom(buf)
	if err != nil {
		return Packet{}, nil, err
	}
	p, err := Decode(buf[:n])
	if err != nil {
		return Packet{}, nil, err
	}
	return p, addr, nil
}
