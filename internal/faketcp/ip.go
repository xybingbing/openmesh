package faketcp

import "net/netip"

func ParseIPv4(s string) netip.Addr {
	addr := netip.MustParseAddr(s)
	if !addr.Is4() {
		panic("not an IPv4 address: " + s)
	}
	return addr
}
