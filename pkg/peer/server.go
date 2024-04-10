package peer

import "github.com/harpy-wings/legendary-vpn/pkg/packet"

type Server interface {
	ToIFace() chan *packet.Packet
}
