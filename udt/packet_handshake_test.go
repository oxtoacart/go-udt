package udt

import (
	"net"
	"testing"
	"log"
)

func TestHandshakePacket(t *testing.T) {
	read := testPacket(
		&handshakePacket{
			h: header{
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
			sockAddr:       net.ParseIP("127.0.0.1"),
		}, t)
		
	log.Println((read.(*handshakePacket)).sockAddr)
}
