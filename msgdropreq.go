package udt

// Structure of packets and functions for writing/reading them

import (
	"bytes"
	"io"
)

type msgDropReqPacket struct {
	controlPacket
	msgId    uint32 // Message ID
	firstSeq uint32 // First sequence number in the message
	lastSeq  uint32 // Last sequence number in the message
}

func (p *msgDropReqPacket) writeTo(w io.Writer) (err error) {
	if err := p.writeHeaderTo(w, msg_drop_req, p.msgId); err != nil {
		return err
	}
	if err := writeBinary(w, p.firstSeq); err != nil {
		return err
	}
	if err := writeBinary(w, p.lastSeq); err != nil {
		return err
	}
	return
}

func (p *msgDropReqPacket) readFrom(b []byte, r *bytes.Reader) (err error) {
	if p.msgId, err = p.readHeaderFrom(r); err != nil {
		return
	}
	if err = readBinary(r, &p.firstSeq); err != nil {
		return
	}
	if err = readBinary(r, &p.lastSeq); err != nil {
		return
	}
	return
}
