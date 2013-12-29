package udt

import (
	"github.com/oxtoacart/bpool"
	"log"
	"net"
	"sync"
	"bytes"
	"time"
)

/*
A multiplexer multiplexes multiple UDT sockets over a single PacketConn.
*/
type multiplexer struct {
	conn            *net.UDPConn          // the PacketConn from which we read/write
	mode            uint8                 // client or server
	sockets         map[uint32]*udtSocket // the server udtSockets handled by this multiplexer, by sockId
	socketsMutex    *sync.Mutex
	sendQ           *udtSocketQueue   // priority queue of udtSockets awaiting a send (actually includes ones with no packets waiting too)
	ctrlOut         chan packet       // control packets meant queued for sending
	in              chan packet       // packets inbound from the PacketConn
	out             chan packet       // packets outbound to the PacketConn
	writeBufferPool *bpool.BufferPool // leaky buffer pool for writing to conn
	readBytePool    *bpool.BytePool   // leaky byte pool for reading from conn
}

func (m *multiplexer) Accept() (c *net.UDPConn, err error) {
	return
}

func (m *multiplexer) Close() (err error) {
	return
}

func (m *multiplexer) Addr() (addr net.Addr) {
	return
}

/*
multiplexerFor gets or creates a multiplexer for the given address.  If a new
multiplexer is created, the given init function is run to obtain a net.UDPConn.
*/
func multiplexerFor(_net string, laddr *net.UDPAddr, init func() (*net.UDPConn, error)) (m *multiplexer, err error) {
	multiplexersMutex.Lock()
	defer multiplexersMutex.Unlock()
	if laddr != nil {
		// Explicit local address given, check if we have a multiplexer for it
		m = multiplexers[laddr.String()]
	}
	if m == nil {
		// No multiplexer, need to create connection
		var conn *net.UDPConn
		if conn, err = init(); err == nil {
			m = newMultiplexer(conn)
			multiplexers[conn.LocalAddr().String()] = m
		}
	}
	return
}

func newMultiplexer(conn *net.UDPConn) (m *multiplexer) {
	m = &multiplexer{
		conn:            conn,
		sockets:         make(map[uint32]*udtSocket),
		socketsMutex:    new(sync.Mutex),
		sendQ:           newUdtSocketQueue(),
		ctrlOut:         make(chan packet, 100),                    // todo: figure out how to size this
		in:              make(chan packet, 100),                    // todo: make this tunable
		out:             make(chan packet, 100),                    // todo: make this tunable
		writeBufferPool: bpool.NewBufferPool(25600),                // todo: make this tunable
		readBytePool:    bpool.NewBytePool(25600, max_packet_size), // todo: make this tunable
	}

	go m.coordinate()
	go m.read()
	go m.write()

	return
}

func (m *multiplexer) newClientSocket() (s *udtSocket, err error) {
	m.socketsMutex.Lock()
	defer m.socketsMutex.Unlock()
	sids -= 1
	sid := sids
	if host, _, err := net.SplitHostPort(m.conn.RemoteAddr().String()); err != nil {
		err = err
	} else {
		if s, err = newClientSocket(m.ctrlOut, net.ParseIP(host), sid); err == nil {
			m.sockets[sid] = s
		}
	}

	for {
		s.initHandshake()
		time.Sleep(200 * time.Millisecond)
	}

	return
}

/*
read runs in a goroutine and reads packets from conn using a buffer from the
readBufferPool, or a new buffer.
*/
func (m *multiplexer) read() {
	for {
		b := m.readBytePool.Get()
		defer m.readBytePool.Put(b)
		if _, err := m.conn.Read(b); err != nil {
			log.Printf("Unable to read into buffer: %s", err)
		} else {
			r := bytes.NewReader(b)
			if p, err := readPacketFrom(r); err != nil {
				log.Printf("Unable to read packet: %s", err)
			} else {
				m.in <- p
			}
		}
	}
}

/*
write runs in a goroutine and writes packets to conn using a buffer from the
writeBufferPool, or a new buffer.
*/
func (m *multiplexer) write() {
	for {
		select {
		case p := <-m.ctrlOut:
			b := m.writeBufferPool.Get()
			defer m.writeBufferPool.Put(b)
			if err := p.writeTo(b); err != nil {
				// TODO: handle write error
				log.Fatalf("Unable to write out: %s", err)
			} else {
				if _, err := b.WriteTo(m.conn); err != nil {
					// TODO: handle write error
					log.Fatalf("Unable to write out: %s", err)
				}
			}
		}
	}
}

// coordinate runs in a goroutine and coordinates all of the multiplexer's work
func (m *multiplexer) coordinate() {
	for {
		select {
		case p := <-m.in:
			m.handleInbound(p)
		}
	}
}

func (m *multiplexer) handleInbound(_p interface{}) {
	switch p := _p.(type) {
	case handshakePacket:
		// Only process packet if version and type are supported
		if p.udtVer == 4 && p.sockType == STREAM {
			s := m.sockets[p.sockId]
			if p.sockType == init_client_handshake {
				if s == nil {
					// create a new udt socket and remember it
					var err error
					if s, err = newServerSocket(m.ctrlOut, p); err == nil {
						m.sockets[p.sockId] = s
					}
				}
			}

			if s == nil {
				// s may still be nil if we couldn't create a new socket
				// in this case, we ignore the error
			} else {
				// Okay, we have a socket

			}
		}
	}
}
