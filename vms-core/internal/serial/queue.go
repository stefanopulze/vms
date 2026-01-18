package serial

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/tarm/serial"
)

// command data to send to serial
type command struct {
	Data     []byte
	Response chan response
}

// response from serial
type response struct {
	Data  []byte
	Error error
}

type Queue struct {
	port       *serial.Port
	queue      chan command
	mu         sync.Mutex
	processing bool
	wg         sync.WaitGroup
}

type QueueOptions struct {
	PortName     string
	PortBaudRate int
	Size         int
}

var logger *slog.Logger

func NewQueue(opts *QueueOptions) (Serial, error) {
	config := &serial.Config{
		Name:        opts.PortName,
		Baud:        opts.PortBaudRate,
		ReadTimeout: 200 * time.Millisecond,
	}

	port, err := serial.OpenPort(config)
	if err != nil {
		return nil, fmt.Errorf("cannot open port: %w", err)
	}

	sq := &Queue{
		port:  port,
		queue: make(chan command, opts.Size),
	}

	logger = slog.With(slog.String("component", "serial"))

	return sq, nil
}

func (sq *Queue) Start() {
	go sq.processQueue()
}

// Write send command to queue and wait response
func (sq *Queue) Write(data []byte) ([]byte, error) {
	cmd := command{
		Data:     data,
		Response: make(chan response, 1),
	}

	sq.queue <- cmd
	resp := <-cmd.Response

	return resp.Data, resp.Error
}

// processQueue process queue commands
func (sq *Queue) processQueue() {
	for cmd := range sq.queue {
		sq.executeCommand(cmd)
	}
}

// executeCommand write command to serial and read response
func (sq *Queue) executeCommand(cmd command) {
	_, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	sq.mu.Lock()
	sq.processing = true
	sq.mu.Unlock()

	defer func() {
		sq.mu.Lock()
		sq.processing = false
		sq.mu.Unlock()
	}()

	//slog.Debug(fmt.Sprintf("Write command: %02X", cmd.Data))
	// Send command
	_ = sq.port.Flush()
	time.Sleep(10 * time.Millisecond)
	_, err := sq.port.Write(cmd.Data)
	if err != nil {
		cmd.Response <- response{
			Data:  nil,
			Error: fmt.Errorf("errore invio: %w", err),
		}
		return
	}

	//slog.Debug(fmt.Sprintf("Write: %02X", cmd.Data))

	start := time.Now()
	buff, err := readBufferedResponse(sq.port)
	logger.Debug(fmt.Sprintf("Read: %s in %v", cmd.Data[:len(cmd.Data)-3], time.Since(start)))

	cmd.Response <- response{
		Data:  buff,
		Error: err,
	}
}

// IsProcessing queue is processing
func (sq *Queue) IsProcessing() bool {
	sq.mu.Lock()
	defer sq.mu.Unlock()
	return sq.processing
}

// QueueLength queue command length
func (sq *Queue) QueueLength() int {
	return len(sq.queue)
}

// Close close serial port
func (sq *Queue) Close() error {
	slog.Debug("Closing serial queue")
	close(sq.queue)
	sq.wg.Wait()
	return sq.port.Close()
}
