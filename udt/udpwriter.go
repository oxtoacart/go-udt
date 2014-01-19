package udt

import (
	"net"
)

/*
boundUDPWriter is an io.Writer that writes everything to a *net.UDPConn with
a specific destination *net.UDPAddr.
*/
type boundUDPWriter struct {
	conn *net.UDPConn
	addr *net.UDPAddr
}

func (w *boundUDPWriter) Write(b []byte) (int, error) {
	return w.conn.WriteToUDP(b, w.addr)
}
