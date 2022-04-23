package iomock

import (
	"errors"
	"testing"
)

var errWrite = errors.New("write error")

func TestWriter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		w := NewWriter(func(p []byte) (int, error) {
			return len(p), nil
		})

		in := "foo"
		if n, err := w.Write([]byte(in)); n != len([]byte(in)) || err != nil {
			t.Errorf("w.Write([]byte(%q)) = %d, %v: want %d, nil", in, n, err, len(in))
		}
	})

	t.Run("error", func(t *testing.T) {
		w := NewWriter(func(p []byte) (int, error) {
			return 0, errWrite
		})

		in := "foo"
		if n, err := w.Write([]byte(in)); n != 0 || !errors.Is(err, errWrite) {
			t.Errorf("w.Write([]byte(%q)) = %d, %v: want 0, %v", in, n, err, errWrite)
		}
	})
}
