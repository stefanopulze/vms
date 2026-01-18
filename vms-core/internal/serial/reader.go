package serial

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/tarm/serial"
)

func readBufferedResponse(port *serial.Port) ([]byte, error) {
	var buffer bytes.Buffer
	tempBuf := make([]byte, 1)
	startTime := time.Now()
	timeout := 1 * time.Second

	for {
		if time.Since(startTime) > timeout {
			return nil, fmt.Errorf("timeout ricezione risposta")
		}

		n, err := port.Read(tempBuf)
		if err != nil {
			if err == io.EOF {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			return nil, fmt.Errorf("cannot read: %w", err)
		}

		if n == 0 {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		buffer.WriteByte(tempBuf[0])

		if tempBuf[0] == 0x0D {
			break
		}
	}

	return buffer.Bytes(), nil
}
