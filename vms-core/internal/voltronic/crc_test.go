package voltronic

import (
	"testing"
)

func TestChecksum(t *testing.T) {
	crc := Checksum([]byte("QPI"))
	if crc != 0xBEAC {
		t.Error("CRC error")
	}

	crc = Checksum([]byte("QID"))
	if crc != 0xD6EA {
		t.Error("CRC error")
	}
}
