package udt

// Structure of packets and functions for writing/reading them

import (
	"bytes"
	"io"
)

type ack2Packet struct {
	controlPacket
	ackSeqNo uint32 // ACK sequence number
}

func (p *ack2Packet) writeTo(w io.Writer) (err error) {
	if err := p.writeHeaderTo(w, ack2, p.ackSeqNo); err != nil {
		return err
	}
	return
}

func (p *ack2Packet) readFrom(b []byte, r *bytes.Reader) (err error) {
	p.ackSeqNo, err = p.readHeaderFrom(r)
	return
}
