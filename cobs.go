// Package cobs implements Consistent Overhead Byte Stuffing encoding and decoding.
//
// References:
//     https://en.wikipedia.org/wiki/Consistent_Overhead_Byte_Stuffing and
//     http://conferences.sigcomm.org/sigcomm/1997/papers/p062.pdf
//     https://tools.ietf.org/html/draft-ietf-pppext-cobs-00
package cobs

// TODO(dgryski): fix api to allow passing in decode buffer

func Encode(src []byte) (dst []byte) {

	// guess at how much extra space we need
	var l int
	l = int(float64(len(src)) * 1.05)

	if len(src) == 0 {
		return []byte{}
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
		if code == 0xFF {
			dst[code_ptr] = code
			code_ptr = len(dst)
			dst = append(dst, 0)
			code = byte(0x01)
		}
	}

	dst[code_ptr] = code

	return dst
}

func EncodeZPE(src []byte) (dst []byte) {

	// guess at how much extra space we need
	l := int(float64(len(src)) * 1.05)

	if len(src) == 0 {
		return []byte{}
	}

	dst = make([]byte, 1, l)

	code_ptr := 0
	code := byte(0x01)

	wantPair := false
	for _, b := range src {

		if wantPair {
			wantPair = false // only valid for next byte
			if b == 0 {
				// assert code < 31
				dst[code_ptr] = code | 0xE0
				code_ptr = len(dst)
				dst = append(dst, 0)
				code = byte(0x01)
				continue
			}

			// was looking for a pair of zeros but didn't find it -- encode as normal
			dst[code_ptr] = code
			code_ptr = len(dst)
			dst = append(dst, 0)
			code = byte(0x01)

			dst = append(dst, b)
			code++

			continue
		}

		if b == 0 {
			if code < 31 {
				wantPair = true
				continue
			}

			// too long to encode with ZPE -- encode as normal
			dst[code_ptr] = code
			code_ptr = len(dst)
			dst = append(dst, 0)
			code = byte(0x01)
			continue
		}

		dst = append(dst, b)
		code++
		if code == 0xE0 {
			dst[code_ptr] = code
			code_ptr = len(dst)
			dst = append(dst, 0)
			code = byte(0x01)
		}
	}

	if wantPair {
		// assert(code < 31)
		code = 0xE0 | code
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

	if len(dst) == 0 {
		return dst
	}

	return dst[0 : len(dst)-1] // trim phantom zero
}

func DecodeZPE(src []byte) (dst []byte) {

	dst = make([]byte, 0, len(src))

	var ptr = 0

	for ptr < len(src) {
		code := src[ptr]

		ptr++

		l := int(code)

		if code > 0xE0 {
			l = int(code & 0x1F)
		}

		for i := 1; i < l; i++ {
			dst = append(dst, src[ptr])
			ptr++
		}

		switch {
		case code > 0xE0:
			dst = append(dst, 0)
			dst = append(dst, 0)
		case code < 0xE0:
			dst = append(dst, 0)
		case code == 0xE0:
			// nothing
		}

	}

	if len(dst) == 0 {
		return dst
	}

	return dst[0 : len(dst)-1] // trim phantom zero
}
