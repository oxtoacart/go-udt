package udt

// Structure of packets and functions for writing/reading them

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const (
	FLAG_BIT = 1 << 31 // leading bit for distinguishing control from data packets

	// Control packet types
	HANDSHAKE            = 0x0
	KEEPALIVE            = 0x1
	ACK                  = 0x2
	NAK                  = 0x3
	UNUSED               = 0x4
	SHUTDOWN             = 0x5
	ACK2                 = 0x6
	MESSAGE_DROP_REQUEST = 0x7

	// Socket types
	STREAM = 0
	DGRAM  = 1 // not supported!
)

var (
	endianness = binary.BigEndian
)

type DataPacket struct {
	seq       uint32
	ts        uint32
	dstSockId uint32
	data      []byte
}

// Control Packets

type ControlPacketHeader struct {
	msgType   uint32
	info      uint32
	ts        uint32
	dstSockId uint32
}

type HandshakePacket struct {
	ch             ControlPacketHeader
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

type KeepAlivePacket struct {
	ch ControlPacketHeader
}

type ACKPacket struct {
	ch       ControlPacketHeader
	pktSeqHi uint32 // The packet sequence number to which all the previous packets have been received (excluding)

	// The below are optional
	rtt         uint32 // RTT (in microseconds)
	rttVar      uint32 // RTT variance
	buffAvail   uint32 // Available buffer size (in bytes)
	pktRecvRate uint32 // Packets receiving rate (in number of packets per second)
	estLinkCap  uint32 // Estimated link capacity (in number of packets per second)
}

type NAKPacket struct {
	ch          ControlPacketHeader
	cmpLossInfo uint32 // integer array of compressed loss information
}

type ShutdownPacket struct {
	ch ControlPacketHeader
}

type ACK2Packet struct {
	ch ControlPacketHeader
}

type MessageDropRequestPacket struct {
	ch       ControlPacketHeader
	firstSeq uint32 // First sequence number in the message
	lastSeq  uint32 // Last sequence number in the message
}

func (dp *DataPacket) writeTo(w io.Writer) (err error) {
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

func readDataPacket(b []byte, r *bytes.Reader, h uint32, maxPacketSize uint16) (p *DataPacket, err error) {
	p = &DataPacket{
		seq:  h,
		data: make([]byte, maxPacketSize),
	}
	if err = readBinary(r, &p.ts); err != nil {
		return
	}
	if err = readBinary(r, &p.dstSockId); err != nil {
		return
	}
	// The data is whatever is left over after reading
	p.data = b[len(b)-r.Len():]
	return
}

func (h *ControlPacketHeader) writeTo(w io.Writer) (err error) {
	// Set the flag bit to indicate this is a control packet
	msgType := h.msgType | FLAG_BIT
	if err := writeBinary(w, msgType); err != nil {
		return err
	}
	if err := writeBinary(w, h.info); err != nil {
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

func readControlPacketHeader(h uint32, r io.Reader) (ch ControlPacketHeader, err error) {
	ch = ControlPacketHeader{msgType: h}
	if err = readBinary(r, &ch.info); err != nil {
		return
	}
	if err = readBinary(r, &ch.ts); err != nil {
		return
	}
	if err = readBinary(r, &ch.dstSockId); err != nil {
		return
	}
	return
}

func (p *HandshakePacket) writeTo(w io.Writer) (err error) {
	if err := p.ch.writeTo(w); err != nil {
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

func readHandshakePacket(b []byte, r *bytes.Reader, ch ControlPacketHeader, maxPacketSize uint16) (p *HandshakePacket, err error) {
	p = &HandshakePacket{ch: ch}
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

func readPacketFromBytes(b []byte, maxPacketSize uint16) (p Packet, err error) {
	// Wrap the byte slice with a reader so that we can use binary.Read() for the metadata
	r := bytes.NewReader(b)
	var h uint32
	if err = readBinary(r, &h); err != nil {
		return
	}
	if h&FLAG_BIT == FLAG_BIT {
		// this is a control packet
		// Remove flag bit
		h = h &^ FLAG_BIT
		if ch, err := readControlPacketHeader(h, r); err != nil {
			return nil, err
		} else {
			switch ch.msgType {
			case HANDSHAKE:
				p, err = readHandshakePacket(b, r, ch, maxPacketSize)
			default:
				err = fmt.Errorf("Unkown control packet type: %X", ch.msgType)
				return nil, err
			}
		}
	} else {
		// this is a data packet
		p, err = readDataPacket(b, r, h, maxPacketSize)
	}
	return
}

func readPacketFromConn(n net.PacketConn, maxPacketSize uint16) (packet Packet, err error) {
	b := make([]byte, maxPacketSize)
	if n, _, err := n.ReadFrom(b); err != nil {
		return nil, err
	} else {
		return readPacketFromBytes(b[:n], maxPacketSize)
	}
}

func writeBinary(w io.Writer, n interface{}) (err error) {
	return binary.Write(w, endianness, n)
}

func readBinary(r io.Reader, n interface{}) (err error) {
	return binary.Read(r, endianness, n)
}
