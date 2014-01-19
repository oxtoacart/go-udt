package udt

// Structure of packets and functions for writing/reading them

import (
	"io"
)

type msgDropReqPacket struct {
	h        header
	msgId    uint32 // Message ID
	firstSeq uint32 // First sequence number in the message
	lastSeq  uint32 // Last sequence number in the message
}

func (p *msgDropReqPacket) socketId() (sockId uint32) {
	return p.h.dstSockId
}

func (p *msgDropReqPacket) sendTime() (ts uint32) {
	return p.h.ts
}

func (p *msgDropReqPacket) writeTo(w io.Writer) (err error) {
	if err := p.h.writeTo(w, msg_drop_req, p.msgId); err != nil {
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

func (p *msgDropReqPacket) readFrom(r io.Reader) (err error) {
	if p.msgId, err = p.h.readFrom(r); err != nil {
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
