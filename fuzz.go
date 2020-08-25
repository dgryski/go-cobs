// +build gofuzz

package cobs

import (
	"bytes"
	"fmt"
)

func Fuzz(data []byte) int {

	enc := Encode(data)

	if bytes.IndexByte(enc, 0) != -1 {
		panic("encoded has a 0")
	}

	orig, err := Decode(enc)
	if err != nil {
		panic("failed decode encoded data")
	}

	if !bytes.Equal(data, orig) {
		panic("faied to roundtrip" + fmt.Sprintf("%q != %q", data, orig))
	}

	zenc := EncodeZPE(data)
	if bytes.IndexByte(zenc, 0) != -1 {
		panic("zpe encoded has a 0")
	}
	zorig, err := DecodeZPE(zenc)
	if err != nil {
		panic("faied to zpe decode")
	}

	if !bytes.Equal(data, zorig) {
		panic("faied to zpe roundtrip" + fmt.Sprintf("%q != %q", data, zorig))
	}

	Decode(data)
	DecodeZPE(data)

	return 0
}
