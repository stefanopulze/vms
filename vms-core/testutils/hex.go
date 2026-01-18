package testutils

import "encoding/hex"

func FromHex(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}
