package udt

// Structure of packets and functions for writing/reading them

import (
	"io"
)

func (p *shutdownPacket) socketId() (sockId uint32) {
	return p.h.dstSockId
}

type shutdownPacket struct {
	h header
}

func (p *shutdownPacket) sendTime() (ts uint32) {
	return p.h.ts
}

func (p *shutdownPacket) writeTo(w io.Writer) (err error) {
	return p.h.writeTo(w, shutdown, noinfo)
}

func (p *shutdownPacket) readFrom(r io.Reader) (err error) {
	_, err = p.h.readFrom(r)
	return
}
