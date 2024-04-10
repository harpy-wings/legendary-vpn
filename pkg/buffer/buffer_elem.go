package buffer

type bufferElem struct {
	p    interface{} //todo might be replaced with Packet
	key  int64
	next *bufferElem
	prev *bufferElem
}
