package udt

// Structure of packets and functions for writing/reading them

import (
	"io"
)

type nakPacket struct {
	h           header
	cmpLossInfo uint32 // integer array of compressed loss information
}

func (p *nakPacket) socketId() (sockId uint32) {
	return p.h.dstSockId
}

func (p *nakPacket) sendTime() (ts uint32) {
	return p.h.ts
}

func (p *nakPacket) writeTo(w io.Writer) (err error) {
	if err := p.h.writeTo(w, nak, noinfo); err != nil {
		return err
	}
	if err := writeBinary(w, p.cmpLossInfo); err != nil {
		return err
	}
	return
}

func (p *nakPacket) readFrom(r io.Reader) (err error) {
	if _, err = p.h.readFrom(r); err != nil {
		return
	}
	if err = readBinary(r, &p.cmpLossInfo); err != nil {
		return
	}
	return
}
