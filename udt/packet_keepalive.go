package udt

// Structure of packets and functions for writing/reading them

import (
	"io"
)

type keepAlivePacket struct {
	h header
}

func (p *keepAlivePacket) socketId() (sockId uint32) {
	return p.h.dstSockId
}

func (p *keepAlivePacket) sendTime() (ts uint32) {
	return p.h.ts
}

func (p *keepAlivePacket) writeTo(w io.Writer) (err error) {
	return p.h.writeTo(w, keepalive, noinfo)
}

func (p *keepAlivePacket) readFrom(r io.Reader) (err error) {
	_, err = p.h.readFrom(r)
	return
}
