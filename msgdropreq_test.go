package udt

import (
	"testing"
)

func TestMsgDropReqPacket(t *testing.T) {
	testPacket(
		&msgDropReqPacket{
			controlPacket: controlPacket{
				ts:        100,
				dstSockId: 59,
			},
			msgId:    90,
			firstSeq: 91,
			lastSeq:  92,
		}, t)
}
