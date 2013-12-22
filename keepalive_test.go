package udt

import (
	"testing"
)

func TestKeepAlivePacket(t *testing.T) {
	testPacket(
		&keepAlivePacket{
			controlPacket: controlPacket{
				ts:        100,
				dstSockId: 59,
			},
		}, t)
}
