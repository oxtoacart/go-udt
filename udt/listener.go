package udt

import (
	"net"
	"io"
)

/*
listener implements the io.Listener interface for UDT.
*/
type listener struct {
	conn *net.UDPConn
}

func (m *multiplexer) Accept() (c io.ReadWriteCloser, err error) {
	return
}

func (m *multiplexer) Close() (err error) {
	return
}

func (m *multiplexer) Addr() (addr net.Addr) {
	return
}

