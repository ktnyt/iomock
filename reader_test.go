package iomock

import (
	"bytes"
	"errors"
	"testing"
)

var errRead = errors.New("read error")

func TestReader(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		in := "foo"

		r := NewReader(func(p []byte) (int, error) {
			copy(p, []byte(in))
			return len(p), nil
		})

		p := make([]byte, len([]byte(in)))
		if n, err := r.Read(p); n != len([]byte(in)) || err != nil {
			t.Errorf("r.Read(p) = %d, %v: want %d, nil", n, err, len(in))
		}

		if string(p) != in {
			t.Errorf("string(p) = %q, want %q", string(p), in)
		}
	})

	t.Run("error", func(t *testing.T) {
		r := NewReader(func(p []byte) (int, error) {
			return 0, errRead
		})

		p := make([]byte, 3)
		if n, err := r.Read(p); n != 0 || !errors.Is(err, errRead) {
			t.Errorf("r.Read(p) = %d, %v: want 0, %v", n, err, errRead)
		}

		exp := make([]byte, 3)
		if !bytes.Equal(p, exp) {
			t.Errorf("p = %v, want %v", p, exp)
		}
	})
}
