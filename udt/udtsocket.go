package udt

import (
	"io"
	"math"
	"net"
	"time"
)

const (
	sock_state_new              = 0
	sock_state_handshake_init   = 1
	sock_state_handshake_finish = 2
	sock_state_connected        = 3
)

/*
udtSocket encapsulates a UDT socket between a local and remote address pair, as
defined by the UDT specification.  udtSocket implements the net.Conn interface
so that it can be used anywhere that anywhere a stream-oriented network
connection (like TCP) would be used.
*/
type udtSocket struct {
	m            *multiplexer // the multiplexer that handles this socket
	raddr        *net.UDPAddr // the remote address
	boundWriter  io.Writer    // a UDP writer that knows which address to send to
	created      time.Time    // the time that this socket was created
	sockState    uint8
	ackPeriod    uint32           // in microseconds
	nakPeriod    uint32           // in microseconds
	expPeriod    uint32           // in microseconds
	sndPeriod    uint32           // in microseconds
	ctrlIn       chan *dataPacket // inbound control packets
	dataIn       chan *dataPacket // inbound data packets
	dataOut      *packetQueue     // queue of outbound data packets
	pktSeq       uint32           // the current packet sequence number
	currDp       *dataPacket      // currently reading data packet (for partial reads)
	currDpOffset int              // offset in currIn (for partial reads)

	// The below fields mirror what's seen on handshakePacket
	udtVer         uint32
	initPktSeq     uint32
	maxPktSize     uint32
	maxFlowWinSize uint32
	sockType       uint32
	sockId         uint32
	synCookie      uint32
	sockAddr       net.IP
}

/*******************************************************************************
 Implementation of net.Conn interface
*******************************************************************************/

func (s *udtSocket) Read(p []byte) (n int, err error) {
	if s.currDp == nil {
		// Grab the next data packet
		s.currDp = <-s.dataIn
		s.currDpOffset = 0
	}
	n = copy(p, s.currDp.data[s.currDpOffset:])
	s.currDpOffset += n
	if s.currDpOffset >= len(s.currDp.data) {
		// we've exhausted the current data packet, reset to nil
		s.currDp = nil
	}

	return
}

// TODO: implement ReadFrom and WriteTo for performance(?)

func (s *udtSocket) Write(p []byte) (n int, err error) {
	s.pktSeq += 1
	dp := &dataPacket{
		seq:       s.pktSeq,
		ts:        uint32(time.Now().Sub(s.created) / time.Microsecond),
		dstSockId: s.sockId,
		data:      p,
	}
	s.dataOut.push(dp)
	return
}

func (s *udtSocket) Close() (err error) {
	// TODO: implement
	return
}

func (s *udtSocket) LocalAddr() net.Addr {
	return s.m.laddr
}

func (s *udtSocket) RemoteAddr() net.Addr {
	return s.raddr
}

func (s *udtSocket) SetDeadline(t time.Time) error {
	// TODO: implement
	return nil
}

func (s *udtSocket) SetReadDeadline(t time.Time) error {
	// TODO: implement
	return nil
}

func (s *udtSocket) SetWriteDeadline(t time.Time) error {
	// TODO: implement
	return nil
}

/*******************************************************************************
 Private functions
*******************************************************************************/

/*
nextSendTime returns the ts of the next data packet with the lowest ts of
queued packets, or math.MaxUint32 if no packets are queued.
*/
func (s *udtSocket) nextSendTime() (ts uint32) {
	p := s.dataOut.peek()
	if p != nil {
		return p.sendTime()
	} else {
		return math.MaxUint32
	}
}

/**
newUdtSocket creates a new UDT socket based on an initial handshakePacket.
*/
func newServerSocket(m *multiplexer, raddr *net.UDPAddr, p *handshakePacket) (s *udtSocket, err error) {
	s = &udtSocket{
		m:              m,
		raddr:          raddr,
		boundWriter:    &boundUDPWriter{m.conn, raddr},
		sockState:      sock_state_new,
		udtVer:         p.udtVer,
		initPktSeq:     p.initPktSeq,
		maxPktSize:     p.maxPktSize,
		maxFlowWinSize: p.maxFlowWinSize,
		sockType:       p.sockType,
		sockId:         p.sockId,
		sockAddr:       raddr.IP,
		synCookie:      randUint32(),
		dataOut:        newPacketQueue(),
	}

	return
}

func newClientSocket(m *multiplexer, sockId uint32) (s *udtSocket, err error) {
	raddr := (m.conn.RemoteAddr()).(*net.UDPAddr)
	s = &udtSocket{
		m:              m,
		raddr:          raddr,
		boundWriter:    m.conn,
		sockState:      sock_state_new,
		udtVer:         4,
		initPktSeq:     randUint32(),
		maxPktSize:     max_packet_size,
		maxFlowWinSize: 8192, // todo: figure out if/how we should calculate this and/or configure it
		sockType:       STREAM,
		sockId:         sockId,
		sockAddr:       raddr.IP,
		dataOut:        newPacketQueue(),
	}

	return
}

/*******************************************************************************
 Lifecycle functions
*******************************************************************************/

func (s *udtSocket) initHandshake() {
	p := handshakePacket{
		udtVer:         s.udtVer,
		sockType:       s.sockType,
		initPktSeq:     s.initPktSeq,
		maxPktSize:     s.maxPktSize,
		maxFlowWinSize: s.maxFlowWinSize,
		connType:       1,
		sockId:         s.sockId,
		sockAddr:       s.sockAddr,
	}
	s.sockState = sock_state_handshake_init
	s.m.ctrlOut <- &p
}

func (s *udtSocket) respondInitHandshake() {
	p := handshakePacket{
		h: header{
			dstSockId: s.sockId,
		},
		udtVer:   s.udtVer,
		sockType: 1,
		sockId:   s.sockId,
	}
	s.sockState = sock_state_handshake_init
	s.m.ctrlOut <- &p
}
