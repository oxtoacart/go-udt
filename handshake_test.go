package udt

import (
	"net"
	"testing"
)

func TestHandshakePacket(t *testing.T) {
	testPacket(
		&handshakePacket{
			controlPacket: controlPacket{
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
