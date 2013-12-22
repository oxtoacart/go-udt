package udt

import (
	"testing"
)

func TestShutdownPacket(t *testing.T) {
	testPacket(
		&shutdownPacket{
			controlPacket: controlPacket{
				ts:        100,
				dstSockId: 59,
			},
		}, t)
}
