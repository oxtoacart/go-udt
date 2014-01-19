/*
Package udt provides a pure Go implementation of the UDT protocol per
http://udt.sourceforge.net/doc/draft-gg-udt-03.txt.

udt does not implement all of the spec.  In particular, the following are not
implemented:

- Rendezvous mode
- DGRAM mode (only streaming is supported)

*/
package udt

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	_net "net"
	"sync"
)

type Conn interface {
	io.ReadWriteCloser
}

type Listener interface {
	// Accept waits for and returns the next connection to the Listener.
	Accept() (c io.ReadWriteCloser, err error)

	// Close closes the Listener.
	// Any blocked Accept operations will be unblocked and return errors.
	Close() (err error)

	// Addr returns the Listener's network address.
	Addr() (addr _net.Addr)
}

/*
DialUDT establishes an outbound UDT connection using the supplied net, laddr and
raddr.  See function net.DialUDP for a description of net, laddr and raddr.
*/
func DialUDT(net string, laddr, raddr *_net.UDPAddr) (conn Conn, err error) {
	var m *multiplexer

	dial := func() (*_net.UDPConn, error) {
		return _net.DialUDP(net, laddr, raddr)
	}

	if m, err = multiplexerFor(laddr, dial); err == nil {
		if m.mode == mode_server {
			err = fmt.Errorf("Attempted to dial out from a server socket")
		} else {
			m.mode = mode_client
			conn, err = m.newClientSocket()
		}
	}

	return
}

/*
ListenUDT listens for incoming UDT connections addressed to the local address
laddr. See function net.ListenUDP for a description of net and laddr.
*/
func ListenUDT(net string, laddr *_net.UDPAddr) (l Listener, err error) {
	var m *multiplexer

	listen := func() (*_net.UDPConn, error) {
		return _net.ListenUDP(net, laddr)
	}

	if m, err = multiplexerFor(laddr, listen); err == nil {
		if m.mode == mode_client {
			err = fmt.Errorf("Attempted to listen on a client socket")
		} else {
			m.mode = mode_server
			l = m
		}
	}

	return
}

const (
	syn_time              = 10000 // in microseconds
	max_packet_size       = 576   // todo: make this tunable
	init_client_handshake = 1

	// Multiplexer modes
	mode_client = 1
	mode_server = 2
)

var (
	multiplexers      = map[string]*multiplexer{}
	multiplexersMutex sync.Mutex
	sids              uint32 // socketId sequence
)

func init() {
	sids = randUint32()
}

/*
randInt32 generates a secure random value between 0 and the max possible uint32
*/
func randUint32() (r uint32) {
	if _r, err := rand.Int(rand.Reader, big.NewInt(math.MaxUint32)); err != nil {
		log.Fatalf("Unable to generate random uint32: %s", err)
	} else {
		r = uint32(_r.Uint64())
	}
	return
}
