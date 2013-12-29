package udt

import (
	"testing"
)

func TestShutdownPacket(t *testing.T) {
	testPacket(
		&shutdownPacket{
			h: header{
				ts:        100,
				dstSockId: 59,
			},
		}, t)
}
