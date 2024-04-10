package peer

import (
	"net"
	"sync"
	"time"

	"github.com/harpy-wings/legendary-vpn/pkg/buffer"
	xudpaddr "github.com/harpy-wings/legendary-vpn/pkg/xudp_addr"
)

// Peer is a record of a peer's available UDP addrs
type Peer struct {
	lock sync.RWMutex

	server Server

	id            uint64
	ip            net.IP
	addrs         map[[6]byte]int
	xAddrsList    []*xudpaddr.XUDPAddr
	seq           uint32
	state         PeerState
	handShakeDone chan struct{} // handshake Done signal
	lastSeenTime  time.Time

	receiveBuffer *buffer.Buffer
}

func newPeer(id uint64, server Server, addr *net.UDPAddr, idx int) *Peer {
	p := new(Peer)
	p.id = id
	p.xAddrsList = make([]*xudpaddr.XUDPAddr, 0)
	p.addrs = make(map[[6]byte]int)
	p.state = PeerStateInit
	p.seq = 0
	p.server = server
	p.receiveBuffer = buffer.NewBuffer(server.ToIFace())
	// TODO WIP
	return p
}

type PeerState int32

const (
	PeerStateUnknown PeerState = iota
	PeerStateInit
	PeerStateHandshake
	PeerStateWorking
	PeerStateFinishing
)
