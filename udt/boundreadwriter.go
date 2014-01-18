package udt

import (
	"net"
)

/*
boundReadWriter is an implementation of io.ReadWriter that reads from a UDPConn
and writes packets to that UDPConn, addressed to a specific remote address.
*/
type boundReadWriter struct {
	conn  *net.UDPConn
	raddr *net.UDPAddr
}

func newBoundReadWriter(conn *net.UDPConn, raddr *net.UDPAddr) (rw *boundReadWriter) {
	return &boundReadWriter{
		conn:  conn,
		raddr: raddr,
	}
}

func (rw *boundReadWriter) Read(p []byte) (n int, err error) {
	return rw.conn.Read(p)
}

func (rw *boundReadWriter) Write(p []byte) (n int, err error) {
	return rw.conn.WriteToUDP(p, rw.raddr)
}
