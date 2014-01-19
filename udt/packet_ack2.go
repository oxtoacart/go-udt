package udt

// Structure of packets and functions for writing/reading them

import (
	"io"
)

type ack2Packet struct {
	h        header
	ackSeqNo uint32 // ACK sequence number
}

func (p *ack2Packet) socketId() (sockId uint32) {
	return p.h.dstSockId
}

func (p *ack2Packet) sendTime() (ts uint32) {
	return p.h.ts
}

func (p *ack2Packet) writeTo(w io.Writer) (err error) {
	if err := p.h.writeTo(w, ack2, p.ackSeqNo); err != nil {
		return err
	}
	return
}

func (p *ack2Packet) readFrom(r io.Reader) (err error) {
	p.ackSeqNo, err = p.h.readFrom(r)
	return
}
