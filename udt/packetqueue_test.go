package udt

import (
	"testing"
)

func TestPacketQueue(t *testing.T) {
	q := newPacketQueue()
	p1 := &dataPacket{
		ts: 1,
	}
	p2 := &dataPacket{
		ts: 2,
	}
	p3 := &dataPacket{
		ts: 3,
	}
	q.push(p1)
	q.push(p3)
	q.push(p2)
	el := q.peek()
	if el != p1 {
		t.Errorf("Wrong packet peeked.  Expected: %s, Got: %s", p1, el)
	}
	el = q.pop()
	if el != p1 {
		t.Errorf("Wrong packet popped.  Expected: %s, Got: %s", p1, el)
	}
	
	el = q.peek()
	if el != p2 {
		t.Errorf("Wrong packet peeked.  Expected: %s, Got: %s", p2, el)
	}
	el = q.pop()
	if el != p2 {
		t.Errorf("Wrong packet popped.  Expected: %s, Got: %s", p2, el)
	}
	
	el = q.peek()
	if el != p3 {
		t.Errorf("Wrong packet peeked.  Expected: %s, Got: %s", p3, el)
	}
	el = q.pop()
	if el != p3 {
		t.Errorf("Wrong packet popped.  Expected: %s, Got: %s", p3, el)
	}
	el = q.pop()
	if el != nil {
		t.Errorf("Expected nil packet, got %s", el)
	}
}
