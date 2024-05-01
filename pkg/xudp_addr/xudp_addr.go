package xudpaddr

import "net"

type XUDPAddr struct {
	// u is the root UDP address
	u *net.UDPAddr

	// hash is the short form of the ip(4 bytes) + port(2 bytes)
	hash [6]byte
}

func NewXUDPFromUDP(a *net.UDPAddr) *XUDPAddr {
	return &XUDPAddr{a, xudpHash(a)}
}

func xudpHash(a *net.UDPAddr) [6]byte {
	var b [6]byte
	copy(b[:4], []byte(a.IP)[:4])
	p := uint16(a.Port)
	b[4] = byte((p >> 8) & 0xFF)
	b[5] = byte(p & 0xFF)
	return b
}
