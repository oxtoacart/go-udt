package udt

import (
	"net"
)

/*
listener implements the io.Listener interface for UDT.
*/
type listener struct {
	conn *net.UDPConn
}

func (m *multiplexer) Accept() (c net.Conn, err error) {
	return m.conn, nil
}

func (m *multiplexer) Close() (err error) {
	return m.conn.Close()
}

func (m *multiplexer) Addr() (addr net.Addr) {
	return m.laddr
}

