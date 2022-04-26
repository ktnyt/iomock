package iomock

// Writer represents a mocked Write call.
type Writer func(p []byte) (int, error)

// NewWriter returns a new mocked io.Writer.
func NewWriter(mock Writer) Writer {
	return mock
}

// Write simulates a Write call by using the provided mock function.
func (write Writer) Write(p []byte) (int, error) {
	return write(p)
}

// CallCountWriter simulates a write call with the number of calls to the
// writer recorded.
func CallCountWriter(write func(i int, p []byte) (int, error)) Writer {
	i := 0
	return func(p []byte) (int, error) {
		i++
		return write(i, p)
	}
}

// ErrOnCallWriter returns a Writer which returns (0, err) when the number of
// calls to the Writer is count, and returns (len(p), nil) otherwise.
func ErrOnCallWriter(count int, err error) Writer {
	return CallCountWriter(func(i int, p []byte) (int, error) {
		if count == i {
			return 0, err
		}
		return len(p), nil
	})
}

// ByteCountWriter simulates a write call with the number of bytes written to
// the writer recorded.
func ByteCountWriter(write func(i int, p []byte) (int, error)) Writer {
	i := 0
	return func(p []byte) (int, error) {
		n, err := write(i, p)
		i += len(p)
		return n, err
	}
}

// ErrOnByteWriter returns a Writer which returns (n >= 0, err) when the total
// number of bytes written exceeds count, and (len(p), nil) otherwise.
func ErrOnByteWriter(count int, err error) Writer {
	return ByteCountWriter(func(i int, p []byte) (int, error) {
		n := count - i
		switch {
		case len(p) < n:
			return len(p), nil
		case 0 < n:
			return n, err
		default:
			return 0, err
		}
	})
}
