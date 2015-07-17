// +build gofuzz

package cobs

func Fuzz(data []byte) int {
	// _, err := Decode(data)
	_, err := DecodeZPE(data)
	if err != nil {
		return 0
	}
	return 1
}
