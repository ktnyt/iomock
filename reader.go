package iomock

import (
	"io"
)

// Reader represents a mocked Read call.
type Reader func(p []byte) (int, error)

// Read simulates a Read call by using the provided mock function.
func (read Reader) Read(p []byte) (int, error) {
	return read(p)
}

// CallCountReader simulates a read call with the number of calls to the
// reader recorded.
func CallCountReader(read func(i int, p []byte) (int, error)) Reader {
	i := 0
	return func(p []byte) (int, error) {
		i++
		return read(i, p)
	}
}

// ByteCountReader simulates a read call number of bytes read from the reader
// recorded.
func ByteCountReader(read func(i int, p []byte) (int, error)) Reader {
	i := 0
	return func(p []byte) (int, error) {
		n, err := read(i, p)
		i += n
		return n, err
	}
}

// ReadMocker provides decorators for augmenting an existing io.Reader with
// mocking functionalities.
type ReadMocker struct {
	r io.Reader
}

// NewReadMocker returns a new ReadMocker
func NewReadMocker(r io.Reader) ReadMocker {
	return ReadMocker{r}
}

// ErrOnCall returns a Reader which returns the number of bytes read from the
// underlying io.Reader and the error provided when the number of calls to the
// Reader is count, and returns the result of the underlying io.Reader
// otherwise.
func (mock ReadMocker) ErrOnCall(count int, err error) Reader {
	return CallCountReader(func(i int, p []byte) (int, error) {
		n, rerr := mock.r.Read(p)
		if i == count {
			return n, err
		}
		return n, rerr
	})
}

// ErrOnByte returns a Reader which returns the number of bytes read from the
// underlying io.Reader and the error provided when the total number of bytes
// read exceeds count, and returns the result of the underlying io.Reader
// otherwise.
func (mock ReadMocker) ErrOnByte(count int, err error) Reader {
	return ByteCountReader(func(i int, p []byte) (int, error) {
		n := count - i
		switch {
		case len(p) < n:
			return mock.r.Read(p)
		case 0 < n:
			if m, rerr := mock.r.Read(p[:n]); rerr != nil {
				return m, rerr
			}
			return n, err
		default:
			return 0, err
		}
	})
}
