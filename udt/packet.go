package udt

// Structure of packets and functions for writing/reading them

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	flag_bit_32 = 1 << 31 // leading bit for distinguishing control from data packets (32 bit version)
	flag_bit_16 = 1 << 15 // leading bit for distinguishing control from data packets (16 bit version)

	// Control packet types
	handshake    = 0x0
	keepalive    = 0x1
	ack          = 0x2
	nak          = 0x3
	unused       = 0x4
	shutdown     = 0x5
	ack2         = 0x6
	msg_drop_req = 0x7

	// Socket types
	// Note - these are reversed from the spec, as the C++ implementation of UDT
	// seems to have them reversed and we want to stay interoperable
	STREAM = 0
	DGRAM  = 1 // not supported!

	// No info for info section of header
	noinfo = 0
)

var (
	endianness = binary.BigEndian
)

type packet interface {
	// sendTime retrieves the timesamp of the packet
	sendTime() (ts uint32)

	writeTo(io.Writer) (err error)

	/*
		readFrom reads the packet from a Reader
	*/
	readFrom(r io.Reader) (err error)
}

type dataPacket struct {
	seq       uint32
	ts        uint32
	dstSockId uint32
	data      []byte
}

type header struct {
	ts        uint32
	dstSockId uint32
}

func (p *dataPacket) sendTime() (ts uint32) {
	return p.ts
}

func (dp *dataPacket) writeTo(w io.Writer) (err error) {
	if err := writeBinary(w, dp.seq); err != nil {
		return err
	}
	if err := writeBinary(w, dp.ts); err != nil {
		return err
	}
	if err := writeBinary(w, dp.dstSockId); err != nil {
		return err
	}
	if _, err := w.Write(dp.data); err != nil {
		return err
	}
	return
}

func (p *dataPacket) readFrom(r io.Reader) (err error) {
	if err = readBinary(r, &p.ts); err != nil {
		return
	}
	if err = readBinary(r, &p.dstSockId); err != nil {
		return
	}
	// The data is whatever is what comes after the 12 bytes of header
	var buf *bytes.Buffer
	switch b := r.(type) {
	case *bytes.Buffer:
		buf = b
	default:
		buf := bytes.NewBuffer([]byte{})
		buf.ReadFrom(r)
	}
	p.data = buf.Bytes()
	return
}

func (h *header) writeTo(w io.Writer, msgType uint16, info uint32) (err error) {
	// Sets the flag bit to indicate this is a control packet
	if err := writeBinary(w, msgType|flag_bit_16); err != nil {
		return err
	}
	// Write 16 bit reserved data
	if err := writeBinary(w, uint16(0)); err != nil {
		return err
	}
	if err := writeBinary(w, info); err != nil {
		return err
	}
	if err := writeBinary(w, h.ts); err != nil {
		return err
	}
	if err := writeBinary(w, h.dstSockId); err != nil {
		return err
	}
	return
}

func (p *header) readFrom(r io.Reader) (addtlInfo uint32, err error) {
	if err = readBinary(r, &addtlInfo); err != nil {
		return
	}
	if err = readBinary(r, &p.ts); err != nil {
		return
	}
	if err = readBinary(r, &p.dstSockId); err != nil {
		return
	}
	return
}

func readPacketFrom(r io.Reader) (p packet, err error) {
	var h uint32
	if err = readBinary(r, &h); err != nil {
		return
	}
	if h&flag_bit_32 == flag_bit_32 {
		// this is a control packet
		// Remove flag bit
		h = h &^ flag_bit_32
		// Message type is leading 16 bits
		msgType := h >> 16
		switch msgType {
		case handshake:
			p = &handshakePacket{}
		case keepalive:
			p = &keepAlivePacket{}
		case ack:
			p = &ackPacket{}
		case nak:
			p = &nakPacket{}
		case shutdown:
			p = &shutdownPacket{}
		case ack2:
			p = &ack2Packet{}
		case msg_drop_req:
			p = &msgDropReqPacket{}
		default:
			err = fmt.Errorf("Unkown control packet type: %X", msgType)
			return nil, err
		}
		err = p.readFrom(r)
		return
	} else {
		// this is a data packet
		p = &dataPacket{
			seq: h,
		}
		err = p.readFrom(r)
	}
	return
}

func writeBinary(w io.Writer, n interface{}) (err error) {
	return binary.Write(w, endianness, n)
}

func readBinary(r io.Reader, n interface{}) (err error) {
	return binary.Read(r, endianness, n)
}
