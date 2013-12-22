package udt

// Structure of packets and functions for writing/reading them

import (
	"bytes"
	"io"
)

type shutdownPacket struct {
	controlPacket
}

func (p *shutdownPacket) writeTo(w io.Writer) (err error) {
	return p.writeHeaderTo(w, shutdown, noinfo)
}

func (p *shutdownPacket) readFrom(b []byte, r *bytes.Reader) (err error) {
	_, err = p.readHeaderFrom(r)
	return
}
