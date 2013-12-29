package udt

import (
	"testing"
)

func TestNAKPacket(t *testing.T) {
	testPacket(
		&nakPacket{
			h: header{
				ts:        100,
				dstSockId: 59,
			},
			cmpLossInfo: 90,
		}, t)
}
