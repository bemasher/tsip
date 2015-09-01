package tsip

import (
	"bytes"
	"testing"
	"testing/quick"
	"time"

	crand "crypto/rand"
	mrand "math/rand"
)

func TestIdentity(t *testing.T) {
	gen := func(buf []byte) bool {
		idBuf := make([]byte, 1)
		crand.Read(idBuf)
		for idBuf[0] == DLE || idBuf[0] == ETX {
			crand.Read(idBuf)
		}

		tsipBuffer := new(bytes.Buffer)
		srcPacket := Packet{idBuf[0], 0, buf}
		err := srcPacket.Write(tsipBuffer)
		if err != nil {
			return false
		}

		var dstPacket Packet
		err = dstPacket.Read(tsipBuffer)
		if err != nil {
			return false
		}

		t.Logf("Src: %+v\n", srcPacket)
		t.Logf("Dst: %+v\n", dstPacket)

		if srcPacket.ID != dstPacket.ID {
			return false
		}
		if bytes.Compare(buf, dstPacket.Data) != 0 {
			return false
		}
		return bytes.Compare(srcPacket.Data, dstPacket.Data) == 0
	}

	config := new(quick.Config)
	config.Rand = mrand.New(mrand.NewSource(time.Now().UnixNano()))
	if err := quick.Check(gen, config); err != nil {
		t.Error(err)
	}
}

func TestEmpty(t *testing.T) {
	buf := bytes.NewReader([]byte{DLE, 0x8F, DLE, ETX})
	var p Packet
	err := p.Read(buf)
	if err != nil {
		t.Error(err)
	}
	if l := len(p.Data); l != 0 {
		t.Fatal("Data length non-zero:", l)
	}
	t.Logf("%+v\n", p)
}
