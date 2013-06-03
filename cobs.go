// Package cobs implements Consistent Overhead Byte Stuffing encoding and decoding.
//
// References:
//     https://en.wikipedia.org/wiki/Consistent_Overhead_Byte_Stuffing and
//     http://conferences.sigcomm.org/sigcomm/1997/papers/p062.pdf
package cobs

// TODO(dgryski): fix api to allow passing in decode buffer
// TODO(dgryski): zero-pair elimination

func Encode(src []byte) (dst []byte) {

	// guess at how much extra space we need
	var l int
	if len(src) <= 254 {
		l = len(src) + 1
	} else {
		l = (len(src) * 104) / 100 // approx
	}

	dst = make([]byte, 1, l)

	code_ptr := 0
	code := byte(0x01)

	for _, b := range src {
		if b == 0 {
			dst[code_ptr] = code
			code_ptr = len(dst)
			dst = append(dst, 0)
			code = byte(0x01)
			continue
		}

		dst = append(dst, b)
		code++
		if code == 0xff {
			dst[code_ptr] = code
			code_ptr = len(dst)
			dst = append(dst, 0)
			code = byte(0x01)
		}
	}

	dst[code_ptr] = code

	return dst
}

func Decode(src []byte) (dst []byte) {

	dst = make([]byte, 0, len(src))

	var ptr = 0

	for ptr < len(src) {
		code := src[ptr]

		ptr++

		for i := 1; i < int(code); i++ {
			dst = append(dst, src[ptr])
			ptr++
		}
		if code < 0xff {
			dst = append(dst, 0)
		}
	}

	return dst[0 : len(dst)-1] // trim phantom zero
}
