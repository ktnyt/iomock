package iomock

import "io"

// WriteMock represents a mocked Write call.
type WriteMock func(p []byte) (int, error)

// NewWriter returns a new mocked io.Writer.
func NewWriter(mock WriteMock) io.Writer {
	return mock
}

// Write simulates a Write call by using the provided mock function.
func (write WriteMock) Write(p []byte) (int, error) {
	return write(p)
}
