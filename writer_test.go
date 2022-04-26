package iomock

import (
	"errors"
	"fmt"
	"testing"
)

var errWrite = errors.New("write error")

type writeAction struct {
	in  []byte
	n   int
	err error
}

var writerTests = []struct {
	name    string
	w       Writer
	actions []writeAction
}{
	{
		name: "Writer (success)",
		w: NewWriter(func(p []byte) (int, error) {
			return len(p), nil
		}),
		actions: []writeAction{
			{in: []byte("foo"), n: 3, err: nil},
		},
	},
	{
		name: "Writer (error)",
		w: NewWriter(func(p []byte) (int, error) {
			return 0, errWrite
		}),
		actions: []writeAction{
			{in: []byte("foo"), n: 0, err: errWrite},
		},
	},
	{
		name: "ErrOnCallWriter",
		w:    ErrOnCallWriter(3, errWrite),
		actions: []writeAction{
			{in: []byte("foo"), n: 3, err: nil},
			{in: []byte("foo"), n: 3, err: nil},
			{in: []byte("foo"), n: 0, err: errWrite},
			{in: []byte("foo"), n: 3, err: nil},
		},
	},
	{
		name: "ErrOnByteWriter",
		w:    ErrOnByteWriter(8, errWrite),
		actions: []writeAction{
			{in: []byte("foo"), n: 3, err: nil},
			{in: []byte("foo"), n: 3, err: nil},
			{in: []byte("foo"), n: 2, err: errWrite},
			{in: []byte("foo"), n: 0, err: errWrite},
		},
	},
}

func TestWriter(t *testing.T) {
	for i, tt := range writerTests {
		t.Run(fmt.Sprintf("case [%d]: %s)", i+1, tt.name), func(t *testing.T) {
			for j, action := range tt.actions {
				t.Run(fmt.Sprintf("write [%d]", j+1), func(t *testing.T) {
					if n, err := tt.w.Write(action.in); n != action.n || !errors.Is(err, action.err) {
						t.Errorf("w.Write([]byte(%q)) = %d, %v: want %d, %v", string(action.in), n, err, action.n, action.err)
					}
				})
			}
		})
	}
}
