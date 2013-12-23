package udt

import (
	"testing"
)

func TestUDTSocketQueue(t *testing.T) {
	// Set up some packets
	p1 := &dataPacket{
		ts: 1,
	}
	p2 := &dataPacket{
		ts: 2,
	}
	p3 := &dataPacket{
		ts: 3,
	}
	p4 := &dataPacket{
		ts: 4,
	}

	// Set up some packetQueues
	pq1 := newPacketQueue()
	pq2 := newPacketQueue()
	pq3 := newPacketQueue()
	pq1.push(p1)
	pq1.push(p2)
	pq2.push(p3)
	pq3.push(p4)

	// Add packetQueues to a udtSocketQueue
	q := newUdtSocketQueue()
	q.push(pq1)
	q.push(pq3)
	q.push(pq2)

	el := q.peek()
	if el != pq1 {
		t.Errorf("Wrong queue peeked.  Expected: %s, Got: %s", pq1, el)
	}
	el = q.pop()
	if el != pq1 {
		t.Errorf("Wrong queue popped.  Expected: %s, Got: %s", pq1, el)
	}

	el = q.peek()
	if el != pq2 {
		t.Errorf("Wrong queue peeked.  Expected: %s, Got: %s", pq2, el)
	}
	el = q.pop()
	if el != pq2 {
		t.Errorf("Wrong queue popped.  Expected: %s, Got: %s", pq2, el)
	}

	el = q.peek()
	if el != pq3 {
		t.Errorf("Wrong queue peeked.  Expected: %s, Got: %s", pq3, el)
	}
	el = q.pop()
	if el != pq3 {
		t.Errorf("Wrong queue popped.  Expected: %s, Got: %s", pq3, el)
	}
	el = q.pop()
	if el != nil {
		t.Errorf("Expected nil queue, got %s", el)
	}
}
