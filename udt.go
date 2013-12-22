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
	"bytes"
	"io"
)

type Packet interface {
	writeTo(io.Writer) (err error)

	readFrom(b []byte, r *bytes.Reader) (err error)
}
