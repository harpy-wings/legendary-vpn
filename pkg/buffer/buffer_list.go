package buffer

import "sync"

type bufferList struct {
	count      int
	head       *bufferElem
	mutex      sync.Mutex
	lastPopped int64
}

func newBufferList() *bufferList {
	bl := new(bufferList)
	bl.count = 0
	bl.lastPopped = -1
	bl.head = new(bufferElem)
	bl.head.p = nil
	bl.head.next = bl.head
	bl.head.prev = bl.head
	return bl
}

func (bl *bufferList) Push(key int64, p interface{}) {
	bl.mutex.Lock()
	defer bl.mutex.Unlock()

	elem := &bufferElem{p, key, nil, nil}

	notInserted := true
	// add the data to the end of the queue
	for curser := bl.head.prev; curser != bl.head; curser = curser.prev {
		if curser.key < key {
			notInserted = false
			elem.next = curser.next
			elem.prev = curser
			curser.next = elem
			elem.next.prev = elem
			break
		}
	}
	if notInserted {
		elem.next = bl.head.next
		elem.prev = bl.head
		bl.head.next = elem
		elem.next.prev = elem
	}
	bl.count++
}

func (bl *bufferList) Pop() interface{} {
	bl.mutex.Lock()
	defer bl.mutex.Unlock()
	if bl.count == 0 {
		// Noting to return
		// unexpected case
		return nil
	}
	elem := bl.head.next
	bl.head.next = elem.next
	elem.next.prev = bl.head
	bl.count--
	return elem.p
}
