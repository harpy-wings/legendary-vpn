package packet

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"

	"github.com/harpy-wings/legendary-vpn/pkg/cipher"
	"github.com/harpy-wings/legendary-vpn/pkg/xerror"
)

type IPacket interface {
	Pack() []byte
	Size() int
	String() string
}

type packetHeader struct {
	Flag       byte
	Seq        uint32
	PLen       uint16
	FragPrefix uint16
	Frag       uint8
	Sid        uint32
	DLen       uint16
}

type Packet struct {
	packetHeader

	payload []byte
	noise   []byte
	buf     []byte
}

func (p *Packet) Pack() []byte {
	var buf *bytes.Buffer
	p.DLen = uint16(len(p.payload))
	if p.buf != nil {
		// to reduce memcopy
		//todo why we do this?
		buf = bytes.NewBuffer(p.buf[:0])
		binary.Write(buf, binary.BigEndian, p.packetHeader)
		return p.buf // todo replace with cipher.Encrypt(p.buf)
	}

	buf.Write(p.payload)
	buf.Write(p.noise)
	p.buf = buf.Bytes()

	return p.buf // todo replace with cipher.Encrypt(p.buf)
}

func (p *Packet) Size() int {
	return HOP_HDR_LEN + len(p.payload) + len(p.noise)
}

func (p *Packet) String() string {
	return fmt.Sprintf(packetFormatString, p.packetHeader, p.payload, p.noise)
}

func (p *Packet) setSid(sid [4]byte) {
	p.Sid = binary.BigEndian.Uint32(sid[:])
}

func (p *Packet) addNoise(n int) {
	if p.buf != nil {
		// already packed before
		noiseOffset := HOP_HDR_LEN + len(p.payload)
		p.noise = p.buf[:noiseOffset:len(p.buf)]
	} else {
		p.noise = make([]byte, n)
	}
	rand.Read(p.noise)
}

func UnpackPacket(b []byte, c cipher.Cipher) (*Packet, error) {
	frame, err := c.Decrypt(b)
	if err != nil {
		return nil, err
	}
	if frame == nil {
		return nil, xerror.ErrNilFrame
	}

	buf := bytes.NewBuffer(frame)
	p := new(Packet)
	err = binary.Read(buf, binary.BigEndian, &p.packetHeader)
	if err != nil {
		return nil, err
	}

	// todo read length validation with header
	_, err = buf.Read(p.payload)
	if err != nil {
		return nil, err
	}
	return p, nil
}
