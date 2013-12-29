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

	// Set up some udtSockets
	s1 := &udtSocket{dataOut: newPacketQueue()}
	s2 := &udtSocket{dataOut: newPacketQueue()}
	s3 := &udtSocket{dataOut: newPacketQueue()}
	s1.dataOut.push(p1)
	s1.dataOut.push(p2)
	s2.dataOut.push(p3)
	s3.dataOut.push(p4)

	// Add sockets to a udtSocketQueue
	q := newUdtSocketQueue()
	q.push(s1)
	q.push(s2)
	q.push(s3)

	el := q.peek()
	if el != s1 {
		t.Errorf("Wrong queue peeked.  Expected: %s, Got: %s", s1, el)
	}
	el = q.pop()
	if el != s1 {
		t.Errorf("Wrong queue popped.  Expected: %s, Got: %s", s1, el)
	}

	el = q.peek()
	if el != s2 {
		t.Errorf("Wrong queue peeked.  Expected: %s, Got: %s", s2, el)
	}
	el = q.pop()
	if el != s2 {
		t.Errorf("Wrong queue popped.  Expected: %s, Got: %s", s2, el)
	}

	el = q.peek()
	if el != s3 {
		t.Errorf("Wrong queue peeked.  Expected: %s, Got: %s", s3, el)
	}
	el = q.pop()
	if el != s3 {
		t.Errorf("Wrong queue popped.  Expected: %s, Got: %s", s3, el)
	}
	el = q.pop()
	if el != nil {
		t.Errorf("Expected nil queue, got %s", el)
	}
}
