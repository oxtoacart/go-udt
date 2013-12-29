package udt

import (
	"container/heap"
)

/*
A packetQueue is a priority queue of packets sorted by ts.
*/
type packetQueue struct {
	h packetHeap
	l uint32
}

func (q *packetQueue) push(p packet) {
	heap.Push(&q.h, p)
	q.l += 1
}

func (q *packetQueue) peek() (p packet) {
	if q.l == 0 {
		return nil
	} else {
		return q.h[0]
	}
}

func (q *packetQueue) pop() (p packet) {
	if q.l == 0 {
		return nil
	} else {
		q.l -= 1
		return heap.Pop(&q.h).(packet)
	}
}

func newPacketQueue() (q *packetQueue) {
	q = &packetQueue{}
	heap.Init(&q.h)
	return
}

/*
A packetHeap is the internal implementation of a Heap used by packetQueue.
*/
type packetHeap []packet

func (h packetHeap) Len() int           { return len(h) }
func (h packetHeap) Less(i, j int) bool { return h[i].sendTime() < h[j].sendTime() }
func (h packetHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *packetHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(packet))
}

func (h *packetHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
