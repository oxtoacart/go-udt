package udt

// Structure of packets and functions for writing/reading them

import (
	"bytes"
	"io"
	"net"
)

type handshakePacket struct {
	controlPacket
	udtVer         uint32 // UDT version
	sockType       uint32 // Socket Type (0 = STREAM or 1 = DGRAM)
	initPktSeq     uint32 // initial packet sequence number
	maxPktSize     uint32 // maximum packet size (including UDP/IP headers)
	maxFlowWinSize uint32 // maximum flow window size
	connType       uint32 // connection type (regular or rendezvous)
	sockId         uint32 // socket ID
	synCookie      uint32 // SYN cookie
	sockAddr       net.IP // the IP address of the peer's UDP socket
}

func (p *handshakePacket) writeTo(w io.Writer) (err error) {
	if err := p.writeHeaderTo(w, handshake, noinfo); err != nil {
		return err
	}
	if err := writeBinary(w, p.udtVer); err != nil {
		return err
	}
	if err := writeBinary(w, p.sockType); err != nil {
		return err
	}
	if err := writeBinary(w, p.initPktSeq); err != nil {
		return err
	}
	if err := writeBinary(w, p.maxPktSize); err != nil {
		return err
	}
	if err := writeBinary(w, p.maxFlowWinSize); err != nil {
		return err
	}
	if err := writeBinary(w, p.connType); err != nil {
		return err
	}
	if err := writeBinary(w, p.sockId); err != nil {
		return err
	}
	if err := writeBinary(w, p.synCookie); err != nil {
		return err
	}
	if _, err := w.Write(p.sockAddr); err != nil {
		return err
	}
	return
}

func (p *handshakePacket) readFrom(b []byte, r *bytes.Reader) (err error) {
	if _, err = p.readHeaderFrom(r); err != nil {
		return
	}
	if err = readBinary(r, &p.udtVer); err != nil {
		return
	}
	if err = readBinary(r, &p.sockType); err != nil {
		return
	}
	if err = readBinary(r, &p.initPktSeq); err != nil {
		return
	}
	if err = readBinary(r, &p.maxPktSize); err != nil {
		return
	}
	if err = readBinary(r, &p.maxFlowWinSize); err != nil {
		return
	}
	if err = readBinary(r, &p.connType); err != nil {
		return
	}
	if err = readBinary(r, &p.sockId); err != nil {
		return
	}
	if err = readBinary(r, &p.synCookie); err != nil {
		return
	}
	p.sockAddr = net.IP(b[len(b)-r.Len():])
	return
}
