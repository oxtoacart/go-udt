package udt

import (
	"testing"
)

func TestACK2Packet(t *testing.T) {
	testPacket(
		&ack2Packet{
			h: header{
				ts:        100,
				dstSockId: 59,
			},
			ackSeqNo:    90,
		}, t)
}
