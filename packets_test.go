package udt

import (
	"bytes"
	"net"
	"testing"
)

func TestData(t *testing.T) {
	p := DataPacket{
		seq:       50,
		ts:        1409,
		dstSockId: 90,
		data:      []byte("Hello UDT World!"),
	}
	b := bytes.NewBuffer([]byte{})
	if err := p.writeTo(b); err != nil {
		t.Errorf("Unable to write packet: %s", err)
	}
	if _p2, err := readPacketFromBytes(b.Bytes(), 1000); err != nil {
		t.Errorf("Unable to read data packet: %s", err)
	} else {
		switch p2 := _p2.(type) {
		case DataPacket:
			if p.seq != p2.seq {
				t.Errorf("Read seq did not match written.\n\nWrote: %s\nRead:  %s", p.seq, p2.seq)
			}
			if p.ts != p2.ts {
				t.Errorf("Read ts did not match written.\n\nWrote: %s\nRead:  %s", p.ts, p2.ts)
			}
			if p.dstSockId != p2.dstSockId {
				t.Errorf("Read dstSockId did not match written.\n\nWrote: %s\nRead:  %s", p.dstSockId, p2.dstSockId)
			}
			if !bytes.Equal(p.data, p2.data) {
				t.Errorf("Read data did not match written.\n\nWrote: %x\nRead:  %x", p.data, p2.data)
			}
		default:
			t.Errorf("Read packet is of incorrect type: %s", _p2)
		}
	}
}

func TestHandshake(t *testing.T) {
	p := HandshakePacket{
		ch: ControlPacketHeader{
			msgType:   HANDSHAKE,
			ts:        100,
			dstSockId: 59,
		},
		udtVer:         4,
		sockType:       STREAM,
		initPktSeq:     50,
		maxPktSize:     1000,
		maxFlowWinSize: 500,
		connType:       1,
		sockId:         59,
		synCookie:      978,
		sockAddr:       net.IP{5, 9, 2},
	}
	b := bytes.NewBuffer([]byte{})
	if err := p.writeTo(b); err != nil {
		t.Errorf("Unable to write packet: %s", err)
	}
	if _p2, err := readPacketFromBytes(b.Bytes(), 1000); err != nil {
		t.Errorf("Unable to read data packet: %s", err)
	} else {
		switch p2 := _p2.(type) {
		case HandshakePacket:
			if p.ch != p2.ch {
				t.Errorf("Read ch did not match written.\n\nWrote: %s\nRead:  %s", p.ch, p2.ch)
			}
			if p.udtVer != p2.udtVer {
				t.Errorf("Read udtVer did not match written.\n\nWrote: %s\nRead:  %s", p.udtVer, p2.udtVer)
			}
			if p.sockType != p2.sockType {
				t.Errorf("Read sockType did not match written.\n\nWrote: %s\nRead:  %s", p.sockType, p2.sockType)
			}
			if p.initPktSeq != p2.initPktSeq {
				t.Errorf("Read initPktSeq did not match written.\n\nWrote: %s\nRead:  %s", p.initPktSeq, p2.initPktSeq)
			}
			if p.maxPktSize != p2.maxPktSize {
				t.Errorf("Read maxPktSize did not match written.\n\nWrote: %s\nRead:  %s", p.maxPktSize, p2.maxPktSize)
			}
			if p.maxFlowWinSize != p2.maxFlowWinSize {
				t.Errorf("Read maxFlowWinSize did not match written.\n\nWrote: %s\nRead:  %s", p.maxFlowWinSize, p2.maxFlowWinSize)
			}
			if p.connType != p2.connType {
				t.Errorf("Read connType did not match written.\n\nWrote: %s\nRead:  %s", p.connType, p2.connType)
			}
			if p.sockId != p2.sockId {
				t.Errorf("Read sockId did not match written.\n\nWrote: %s\nRead:  %s", p.sockId, p2.sockId)
			}
			if p.synCookie != p2.synCookie {
				t.Errorf("Read synCookie did not match written.\n\nWrote: %s\nRead:  %s", p.synCookie, p2.synCookie)
			}
			if !bytes.Equal(p.sockAddr, p2.sockAddr) {
				t.Errorf("Read sockAddr did not match written.\n\nWrote: %x\nRead:  %x", p.sockAddr, p2.sockAddr)
			}
		default:
			t.Errorf("Read packet is of incorrect type: %s", _p2)
		}
	}
}
