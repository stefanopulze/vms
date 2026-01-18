package serial

import "io"

type Serial interface {
	Write(data []byte) ([]byte, error)
	Start()
	io.Closer
}
