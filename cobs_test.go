package cobs

import (
	"bytes"
	"testing"
)

var encodeTests = []struct {
	in  []byte
	out []byte
}{
	{
		[]byte{0x00},
		[]byte{0x01, 0x01},
	},
	{
		[]byte{0x11, 0x22, 0x00, 0x33},
		[]byte{0x03, 0x11, 0x22, 0x02, 0x33},
	},
	{
		[]byte{0x11, 0x00, 0x00, 0x00},
		[]byte{0x02, 0x11, 0x01, 0x01, 0x01},
	},
}

func TestEncode(t *testing.T) {

    for _, tst := range encodeTests {
        o := Encode(tst.in)
        if !bytes.Equal(o, tst.out) {
            t.Errorf("encode failed: got % 02x wanted % 02x\n", o, tst.out)
        }

        o = Decode(o)
        if !bytes.Equal(o, tst.in) {
            t.Errorf("decode failed: got % 02x wanted % 02x\n", o, tst.in)
        }
    }
}
