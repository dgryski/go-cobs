package cobs

import (
	"bytes"
	"testing"

	"github.com/dgryski/go-tinyfuzz"
)

type inputOutput struct {
	in  []byte
	out []byte
}

func testEncodeDecode(t *testing.T, codec Encoder, tests []inputOutput) {

	for _, tst := range tests {
		o := codec.Encode(tst.in)
		if !bytes.Equal(o, tst.out) {
			t.Errorf("encode failed: got % 02x wanted % 02x\n", o, tst.out)
		}

		o, err := codec.Decode(o)
		if err != nil || !bytes.Equal(o, tst.in) {
			t.Errorf("decode failed: got % 02x wanted % 02x\n", o, tst.in)
		}
	}
}

var cobsTests = []inputOutput{
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

func TestCOBS(t *testing.T) {
	testEncodeDecode(t, New(), cobsTests)
}

var zpeTests = []inputOutput{
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

func TestCOBSZPE(t *testing.T) {
	testEncodeDecode(t, NewZPE(), zpeTests)
}

func testQuick(t *testing.T, codec Encoder) {
	f := func(s []byte) bool {
		e := codec.Encode(s)
		// make sure there are no 0 bytes in the encoded data
		for i, l := 0, len(e); i < l; i++ {
			if e[i] == 0 {
				return false
			}
		}
		o, err := codec.Decode(e)
		return err == nil && bytes.Equal(s, o)
	}

	if err := tinyfuzz.Fuzz(f, nil); err != nil {
		t.Errorf("cobs roundtrip failed: %v\n", err)
	}
}

func TestQuickCOBS(t *testing.T) {
	testQuick(t, New())
	testQuick(t, NewZPE())
}

func testLengths(t *testing.T, codec Encoder) {
	var b []byte

	for i := 0; i < 512; i++ {
		b = append(b, 0)

		e := codec.Encode(b)
		o, err := codec.Decode(e)
		if err != nil || !bytes.Equal(b, o) {
			t.Errorf("length test failed for 0x11 x %d...\n", i)
		}

		b[len(b)-1] = 0x11
	}

}

func TestLengths(t *testing.T) {
	testLengths(t, New())
	testLengths(t, NewZPE())
}
