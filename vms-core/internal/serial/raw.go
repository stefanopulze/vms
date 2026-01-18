package serial

import (
	"fmt"
	"log/slog"

	"github.com/tarm/serial"
)

func NewRawSerial(port *serial.Port) Serial {
	return &rawSerial{
		port: port,
	}
}

type rawSerial struct {
	port *serial.Port
}

func (r rawSerial) Start() {
}

func (r rawSerial) Close() error {
	if r.port != nil {
		return r.port.Close()
	}

	return nil
}

func (r rawSerial) Write(data []byte) ([]byte, error) {
	_, err := r.port.Write(data)
	if err != nil {
		return nil, err
	}
	slog.Info(fmt.Sprintf("Write: %02X", data))

	return readBufferedResponse(r.port)
}
