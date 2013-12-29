package udt

import (
	"testing"
)

func TestMsgDropReqPacket(t *testing.T) {
	testPacket(
		&msgDropReqPacket{
			h: header{
				ts:        100,
				dstSockId: 59,
			},
			msgId:    90,
			firstSeq: 91,
			lastSeq:  92,
		}, t)
}
