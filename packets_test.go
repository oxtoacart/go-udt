package udt

import (
	"bytes"
	"net"
	"reflect"
	"testing"
)

func TestData(t *testing.T) {
	testPacket(
		&DataPacket{
			seq:       50,
			ts:        1409,
			dstSockId: 90,
			data:      []byte("Hello UDT World!"),
		}, t)
}

func TestHandshake(t *testing.T) {
	testPacket(
		&HandshakePacket{
			ch: ControlPacketHeader{
				msgType:   HANDSHAKE,
				ts:        100,
				dstSockId: 59,
			},
			udtVer:         4,
			sockType:       STREAM,
			initPktSeq:     50,
			maxPktSize:     1000,
			maxFlowWinSize: 500,
			connType:       1,
			sockId:         59,
			synCookie:      978,
			sockAddr:       net.IP{5, 9, 2},
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
