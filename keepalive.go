package udt

// Structure of packets and functions for writing/reading them

import (
	"bytes"
	"io"
)

type keepAlivePacket struct {
	controlPacket
}

func (p *keepAlivePacket) writeTo(w io.Writer) (err error) {
	return p.writeHeaderTo(w, keepalive, noinfo)
}

func (p *keepAlivePacket) readFrom(b []byte, r *bytes.Reader) (err error) {
	_, err = p.readHeaderFrom(r)
	return
}
