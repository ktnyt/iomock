package iomock

import "io"

// ReadMock represents a mocked Read call.
type ReadMock func(p []byte) (int, error)

// NewReader returns a new mocked io.Reader.
func NewReader(mock ReadMock) io.Reader {
	return mock
}

// Read simulates a Read call by using the provided mock function.
func (read ReadMock) Read(p []byte) (int, error) {
	return read(p)
}
