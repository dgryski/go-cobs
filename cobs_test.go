package cobs

import (
	"bytes"
	"testing"
	"testing/quick"
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

var encodeZPETests = []struct {
	in  []byte
	out []byte
}{
	{
		[]byte{0x45, 0x00, 0x00, 0x2C, 0x4C, 0x79, 0x00, 0x00, 0x40, 0x06, 0x4F, 0x37},
		[]byte{0xE2, 0x45, 0xE4, 0x2C, 0x4C, 0x79, 0x05, 0x40, 0x06, 0x4F, 0x37},
	},
	{
		[]byte{0x11, 0x00, 0x00, 0x00},
		[]byte{0xE2, 0x11, 0xE1},
	},
	{
		[]byte{0x11, 0x22, 0x00, 0x33},
		[]byte{0x03, 0x11, 0x22, 0x02, 0x33},
	},
}

func TestEncodeZPE(t *testing.T) {

	for _, tst := range encodeZPETests {
		o := EncodeZPE(tst.in)
		if !bytes.Equal(o, tst.out) {
			t.Errorf("encode zpe failed: got % 02x wanted % 02x\n", o, tst.out)
		}

		o = DecodeZPE(o)
		if !bytes.Equal(o, tst.in) {
			t.Errorf("decode zpe failed: got % 02x wanted % 02x\n", o, tst.in)
		}
	}
}

func TestQuick(t *testing.T) {

	f := func(s []byte) bool {
		e := Encode(s)
		o := Decode(e)
		return bytes.Equal(s, o)
	}

	quick.Check(f, nil)
}

func TestLengths(t *testing.T) {
	var b []byte

	for i := 0; i < 512; i++ {
		b = append(b, 0)

		e := Encode(b)
		o := Decode(e)
		if !bytes.Equal(b, o) {
			t.Errorf("length test failed for 0x11 x %d...\n", i)
		}

		e = EncodeZPE(b)
		o = DecodeZPE(e)
		if !bytes.Equal(b, o) {
			t.Errorf("length test failed for 0x11 x %d...\n", i)
		}

		b[len(b)-1] = 0x11
	}

}

func TestZPEQuick(t *testing.T) {

	f := func(s []byte) bool {
		e := EncodeZPE(s)
		o := DecodeZPE(e)
		return bytes.Equal(s, o)
	}

	quick.Check(f, nil)
}
