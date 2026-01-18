package voltronic

const poly = uint16(0x1021)

// prepareCommand Calculate CRC and append it to the command
func prepareCommand(data []byte) []byte {
	crc := Checksum(data)
	// L'inverter ha un'eccezione specifica: se il byte del CRC è
	// 0x0A (LF), 0x0D (CR) o 0x28 ( ( ), deve essere incrementato di 1
	// per evitare conflitti con i caratteri di controllo del protocollo.
	high := byte(crc >> 8)
	low := byte(crc & 0xFF)

	if high == 0x0A || high == 0x0D || high == 0x28 {
		high++
	}
	if low == 0x0A || low == 0x0D || low == 0x28 {
		low++
	}

	result := make([]byte, 0, len(data)+3)
	result = append(result, data...)
	result = append(result, high, low)
	result = append(result, 0x0D) // CR

	return result
}

func Checksum(data []byte) uint16 {
	crc := uint16(0x0000)

	for _, b := range data {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if (crc & 0x8000) != 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc <<= 1
			}
		}
	}

	return crc
}
