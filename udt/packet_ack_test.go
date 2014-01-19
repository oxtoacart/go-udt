package udt

import (
	"testing"
)

func TestACKPacket(t *testing.T) {
	testPacket(
		&ackPacket{
			h: header{
				ts:        100,
				dstSockId: 59,
			},
			ackSeqNo:    90,
			pktSeqHi:    91,
			rtt:         92,
			rttVar:      93,
			buffAvail:   94,
			pktRecvRate: 95,
			estLinkCap:  96,
		}, t)
}
