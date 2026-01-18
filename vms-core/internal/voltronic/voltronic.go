package voltronic

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"vms-core/internal/serial"
)

const nak = "(NAK"

var ErrorNak = errors.New(nak)
var logger *slog.Logger

type Client struct {
	port             serial.Serial
	ratingInfo       DeviceRatingInfo
	ratingInfoMux    sync.Mutex
	generalStatus    DeviceGeneralStatus
	generalStatusMux sync.Mutex
}

func NewClient(port serial.Serial) *Client {
	logger = slog.With(slog.String("component", "voltronic"))
	return &Client{
		port:          port,
		ratingInfo:    DeviceRatingInfo{},
		generalStatus: DeviceGeneralStatus{},
	}
}

// SendCommand send command after calculating CRC, response validation from CRC
func (c *Client) SendCommand(cmd string) ([]byte, error) {
	payload := prepareCommand([]byte(cmd))

	response, err := c.port.Write(payload)
	if err != nil {
		return nil, err
	}

	return c.validateResponse(response)
}

// SendUpdateCommand send command after calculating CRC, response validation from CRC, and check NAK
func (c *Client) SendUpdateCommand(cmd string) error {
	data, err := c.SendCommand(cmd)
	if err != nil {
		return err
	}

	if nak == string(data) {
		return ErrorNak
	}

	return nil
}

func (c *Client) validateResponse(data []byte) ([]byte, error) {
	if len(data) < 3 {
		return nil, errors.New("invalid length")
	}

	if data[len(data)-1] != 0x0D {
		return nil, errors.New("invalid terminator")
	}

	// remove CR
	data = data[:len(data)-1]

	if len(data) < 2 {
		return nil, errors.New("response without CRC")
	}

	size := len(data) - 2
	payload := data[:size]
	receivedCRC := uint16(data[size])<<8 | uint16(data[size+1])
	computedCRC := Checksum(payload)

	if receivedCRC != computedCRC {
		//return nil, errors.New("CRC error")
		return nil, fmt.Errorf("invalid CRC (%X != %X)", receivedCRC, computedCRC)
	}

	return payload, nil
}
