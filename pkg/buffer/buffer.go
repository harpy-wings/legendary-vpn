package buffer

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/harpy-wings/legendary-vpn/pkg/packet"
)

type Buffer struct {
	mutex sync.Mutex

	buf       *bufferList
	rate      int32
	flushChan chan *packet.Packet
	newPack   chan struct{}
	config    struct {
		newPackBufferSize int
	}
}

func NewBuffer(flushChan chan *packet.Packet) *Buffer {
	b := new(Buffer)
	b.setDefaults()
	b.flushChan = flushChan
	b.init()
	return b
}

func (b *Buffer) init() {
	b.buf = newBufferList()
	b.newPack = make(chan struct{}, b.config.newPackBufferSize)
	go func() {
		for {
			p := b.Pop()
			if p == nil {
				continue
			}
			b.flushChan <- p
		}
	}()

}

func (b *Buffer) setDefaults() {
	b.config.newPackBufferSize = defaultNewPackBufferSize
}

func (b *Buffer) Push(p *packet.Packet) {
	atomic.AddInt32(&b.rate, 1)
	b.buf.Push(int64(p.Seq), p)
	b.newPack <- struct{}{}
}

func (b *Buffer) Pop() *packet.Packet {
	<-b.newPack
	r := int(b.rate & rateMask)
	if b.buf.count < 8+r { // what is that 8?
		// todo 20 , 50 magic words !!!!
		time.Sleep(time.Duration(r*20+50) * time.Microsecond)
		b.rate = b.rate >> 1
	}
	p := b.buf.Pop().(*packet.Packet) // todo handle type assertion
	return p
}
