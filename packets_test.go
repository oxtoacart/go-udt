package udt

import (
	"bytes"
	"reflect"
	"testing"
)

func TestDataPacket(t *testing.T) {
	testPacket(
		&dataPacket{
			seq:       50,
			ts:        1409,
			dstSockId: 90,
			data:      []byte("Hello UDT World!"),
		}, t)
}

func testPacket(p Packet, t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	if err := p.writeTo(b); err != nil {
		t.Errorf("Unable to write packet: %s", err)
	}
	if p2, err := readPacketFromBytes(b.Bytes(), 1000); err != nil {
		t.Errorf("Unable to read packet: %s", err)
	} else {
		if !reflect.DeepEqual(p, p2) {
			t.Errorf("Read did not match written.\n\nWrote: %s\nRead:  %s", p, p2)
		}
	}
}
