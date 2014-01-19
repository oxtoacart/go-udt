package udt

// Structure of packets and functions for writing/reading them

import (
	"io"
)

type ackPacket struct {
	h        header
	ackSeqNo uint32 // ACK sequence number
	pktSeqHi uint32 // The packet sequence number to which all the previous packets have been received (excluding)

	// The below are optional
	rtt         uint32 // RTT (in microseconds)
	rttVar      uint32 // RTT variance
	buffAvail   uint32 // Available buffer size (in bytes)
	pktRecvRate uint32 // Packets receiving rate (in number of packets per second)
	estLinkCap  uint32 // Estimated link capacity (in number of packets per second)
}

func (p *ackPacket) socketId() (sockId uint32) {
	return p.h.dstSockId
}

func (p *ackPacket) sendTime() (ts uint32) {
	return p.h.ts
}

func (p *ackPacket) writeTo(w io.Writer) (err error) {
	if err := p.h.writeTo(w, ack, p.ackSeqNo); err != nil {
		return err
	}
	if err := writeBinary(w, p.pktSeqHi); err != nil {
		return err
	}
	if err := writeBinary(w, p.rtt); err != nil {
		return err
	}
	if err := writeBinary(w, p.rttVar); err != nil {
		return err
	}
	if err := writeBinary(w, p.buffAvail); err != nil {
		return err
	}
	if err := writeBinary(w, p.pktRecvRate); err != nil {
		return err
	}
	if err := writeBinary(w, p.estLinkCap); err != nil {
		return err
	}
	return
}

func (p *ackPacket) readFrom(r io.Reader) (err error) {
	if p.ackSeqNo, err = p.h.readFrom(r); err != nil {
		return
	}
	if err = readBinary(r, &p.pktSeqHi); err != nil {
		return
	}
	if err = readBinary(r, &p.rtt); err != nil {
		return
	}
	if err = readBinary(r, &p.rttVar); err != nil {
		return
	}
	if err = readBinary(r, &p.buffAvail); err != nil {
		return
	}
	if err = readBinary(r, &p.pktRecvRate); err != nil {
		return
	}
	if err = readBinary(r, &p.estLinkCap); err != nil {
		return
	}
	return
}
