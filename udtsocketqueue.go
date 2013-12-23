package udt

import (
	"container/heap"
)

/*
A udtSocketQueue is a priority queue of udtSockets sorted by the next packet sending time.
*/
type udtSocketQueue struct {
	h socketHeap
	l uint32
}

func (q *udtSocketQueue) push(p *packetQueue) {
	heap.Push(&q.h, p)
	q.l += 1
}

func (q *udtSocketQueue) peek() (p *packetQueue) {
	if q.l == 0 {
		return nil
	} else {
		return q.h[0]
	}
}

func (q *udtSocketQueue) pop() (p *packetQueue) {
	if q.l == 0 {
		return nil
	} else {
		q.l -= 1
		return heap.Pop(&q.h).(*packetQueue)
	}
}

func newUdtSocketQueue() (q *udtSocketQueue) {
	q = &udtSocketQueue{}
	heap.Init(&q.h)
	return
}

/*
A socketHeap is the internal implementation of a Heap used by udtSocketQueue.
*/
type socketHeap []*packetQueue

func (h socketHeap) Len() int           { return len(h) }
func (h socketHeap) Less(i, j int) bool { return h[i].nextPacketSendTime() < h[j].nextPacketSendTime() }
func (h socketHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *socketHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*packetQueue))
}

func (h *socketHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
