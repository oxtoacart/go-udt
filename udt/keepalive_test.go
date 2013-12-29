package udt

import (
	"testing"
)

func TestKeepAlivePacket(t *testing.T) {
	testPacket(
		&keepAlivePacket{
			h: header{
				ts:        100,
				dstSockId: 59,
			},
		}, t)
}
