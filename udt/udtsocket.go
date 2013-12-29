package udt

import (
	"math"
	"net"
)

const (
	sock_state_new              = 0
	sock_state_handshake_init   = 1
	sock_state_handshake_finish = 2
	sock_state_connected        = 3
)

type udtSocket struct {
	sockState uint8
	ackPeriod uint32       // in microseconds
	nakPeriod uint32       // in microseconds
	expPeriod uint32       // in microseconds
	sndPeriod uint32       // in microseconds
	ctrlOut   chan packet  // outbound control packets
	dataOut   *packetQueue // queue of outbound data packets

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

func (s *udtSocket) Read(p []byte) (n int, err error) {
	return
}

func (s *udtSocket) Write(p []byte) (n int, err error) {
	return
}

func (s *udtSocket) Close() (err error) {
	return
}

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
func newServerSocket(ctrlOut chan packet, p handshakePacket) (s *udtSocket, err error) {
	s = &udtSocket{
		sockState:      sock_state_new,
		udtVer:         p.udtVer,
		initPktSeq:     p.initPktSeq,
		maxPktSize:     p.maxPktSize,
		maxFlowWinSize: p.maxFlowWinSize,
		sockType:       p.sockType,
		sockId:         p.sockId,
		synCookie:      randUint32(),
		ctrlOut:        ctrlOut,
		dataOut:        newPacketQueue(),
	}

	s.respondInitHandshake()

	return
}

func newClientSocket(ctrlOut chan packet, peerIp net.IP, sockId uint32) (s *udtSocket, err error) {
	s = &udtSocket{
		sockState:      sock_state_new,
		udtVer:         4,
		initPktSeq:     randUint32(),
		maxPktSize:     max_packet_size,
		maxFlowWinSize: 8192, // todo: figure out if/how we should calculate this and/or configure it
		sockType:       DGRAM,
		sockId:         sockId,
		sockAddr:       peerIp,
		ctrlOut:        ctrlOut,
		dataOut:        newPacketQueue(),
	}

	s.initHandshake()

	return
}

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
	s.ctrlOut <- &p
}

func (s *udtSocket) respondInitHandshake() {
	p := handshakePacket{
		h: header{
			dstSockId: s.sockId,
		},
		udtVer:   s.udtVer,
		sockType: 1,
	}
	s.sockState = sock_state_handshake_init
	s.ctrlOut <- &p
}
