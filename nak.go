package udt

// Structure of packets and functions for writing/reading them

import (
	"bytes"
	"io"
)

type nakPacket struct {
	controlPacket
	cmpLossInfo uint32 // integer array of compressed loss information
}

func (p *nakPacket) writeTo(w io.Writer) (err error) {
	if err := p.writeHeaderTo(w, nak, noinfo); err != nil {
		return err
	}
	if err := writeBinary(w, p.cmpLossInfo); err != nil {
		return err
	}
	return
}

func (p *nakPacket) readFrom(b []byte, r *bytes.Reader) (err error) {
	if _, err = p.readHeaderFrom(r); err != nil {
		return
	}
	if err = readBinary(r, &p.cmpLossInfo); err != nil {
		return
	}
	return
}
