package udt

import (
	"bytes"
	"github.com/oxtoacart/bpool"
	"log"
	"net"
	"sync"
	"time"
)

/*
A multiplexer multiplexes multiple UDT sockets over a single PacketConn.
*/
type multiplexer struct {
	laddr           *net.UDPAddr          // the local address handled by this multiplexer
	conn            *net.UDPConn          // the UDPConn from which we read/write
	mode            uint8                 // client or server
	sockets         map[uint32]*udtSocket // the server udtSockets handled by this multiplexer, by sockId
	socketsMutex    *sync.Mutex
	sendQ           *udtSocketQueue   // priority queue of udtSockets awaiting a send (actually includes ones with no packets waiting too)
	ctrlOut         chan packet       // control packets meant queued for sending
	in              chan packetHolder // packets inbound from the PacketConn
	out             chan packet       // packets outbound to the PacketConn
	writeBufferPool *bpool.BufferPool // leaky buffer pool for writing to conn
	readBytePool    *bpool.BytePool   // leaky byte pool for reading from conn
}

type packetHolder struct {
	packet packet
	from   *net.UDPAddr
}

/*
multiplexerFor gets or creates a multiplexer for the given local address.  If a
new multiplexer is created, the given init function is run to obtain an
io.ReadWriter.
*/
func multiplexerFor(laddr *net.UDPAddr, init func() (*net.UDPConn, error)) (m *multiplexer, err error) {
	multiplexersMutex.Lock()
	defer multiplexersMutex.Unlock()
	key := laddr.String()
	m = multiplexers[key]
	if m == nil {
		// No multiplexer, need to create connection
		var conn *net.UDPConn
		if conn, err = init(); err == nil {
			m = newMultiplexer(laddr, conn)
			multiplexers[key] = m
		}
	}
	return
}

func newMultiplexer(laddr *net.UDPAddr, conn *net.UDPConn) (m *multiplexer) {
	m = &multiplexer{
		laddr:           laddr,
		conn:            conn,
		sockets:         make(map[uint32]*udtSocket),
		socketsMutex:    new(sync.Mutex),
		sendQ:           newUdtSocketQueue(),
		ctrlOut:         make(chan packet, 100),                    // todo: figure out how to size this
		in:              make(chan packetHolder, 100),              // todo: make this tunable
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
	if s, err = newClientSocket(m, sid); err == nil {
		m.sockets[sid] = s
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
		if _, from, err := m.conn.ReadFromUDP(b); err != nil {
			log.Printf("Unable to read into buffer: %s", err)
		} else {
			r := bytes.NewReader(b)
			if p, err := readPacketFrom(r); err != nil {
				log.Printf("Unable to read packet: %s", err)
			} else {
				m.in <- packetHolder{p, from}
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
				log.Fatalf("Unable to buffer out: %s", err)
			} else {
				log.Printf("Writing to %s", p.socketId())
				if _, err := b.WriteTo(m.sockets[p.socketId()].boundWriter); err != nil {
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

func (m *multiplexer) handleInbound(ph packetHolder) {
	switch p := ph.packet.(type) {
	case *handshakePacket:
		// Only process packet if version and type are supported
		log.Println("Got handshake packet")
		if p.udtVer == 4 && p.sockType == STREAM {
			log.Println("Right version and type")
			s := m.sockets[p.sockId]
			if p.sockType == init_client_handshake {
				if s == nil {
					// create a new udt socket and remember it
					var err error
					if s, err = newServerSocket(m, ph.from, p); err == nil {
						m.sockets[p.sockId] = s
						log.Println("Responding to handshake")
						s.respondInitHandshake()
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
