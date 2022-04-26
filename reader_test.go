package iomock

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
)

var errRead = errors.New("read error")

type readAction struct {
	out []byte
	n   int
	err error
}

var readerTests = []struct {
	name    string
	r       Reader
	actions []readAction
}{
	{
		name: "Reader (success)",
		r: NewReader(func(p []byte) (int, error) {
			return copy(p, []byte("foo")), nil
		}),
		actions: []readAction{
			{out: []byte("foo"), n: 3, err: nil},
		},
	},
	{
		name: "Reader (error)",
		r: NewReader(func(p []byte) (int, error) {
			return 0, errRead
		}),
		actions: []readAction{
			{out: []byte{}, n: 0, err: errRead},
		},
	},
	{
		name: "MockReader.ErrOnCall",
		r:    NewReadMocker(bytes.NewBuffer([]byte("foobarbaz"))).ErrOnCall(3, errRead),
		actions: []readAction{
			{out: []byte("foo"), n: 3, err: nil},
			{out: []byte("bar"), n: 3, err: nil},
			{out: []byte("baz"), n: 3, err: errRead},
			{out: []byte{0, 0, 0}, n: 0, err: io.EOF},
		},
	},
	{
		name: "MockReader.ErrOnByte",
		r:    NewReadMocker(bytes.NewBuffer([]byte("foobarbaz"))).ErrOnByte(8, errRead),
		actions: []readAction{
			{out: []byte("foo"), n: 3, err: nil},
			{out: []byte("bar"), n: 3, err: nil},
			{out: []byte{98, 97, 0}, n: 2, err: errRead},
			{out: []byte{0, 0, 0}, n: 0, err: errRead},
		},
	},
	{
		name: "MockReader.ErrOnByte",
		r: NewReadMocker(NewReader(func(p []byte) (int, error) {
			return 0, io.EOF
		})).ErrOnByte(2, errRead),
		actions: []readAction{
			{out: []byte{0, 0, 0}, n: 0, err: io.EOF},
		},
	},
}

func TestReader(t *testing.T) {
	for i, tt := range readerTests {
		t.Run(fmt.Sprintf("case [%d]: %s)", i+1, tt.name), func(t *testing.T) {
			for j, action := range tt.actions {
				t.Run(fmt.Sprintf("read [%d]", j+1), func(t *testing.T) {
					out := make([]byte, len(action.out))
					if n, err := tt.r.Read(out); n != action.n || !errors.Is(err, action.err) {
						t.Errorf("r.Read([]byte(%d)) = %d, %v: want %d, %v", len(out), n, err, action.n, action.err)
					}
					if !bytes.Equal(out, action.out) {
						t.Errorf("unexpected read output %v, want %v", out, action.out)
					}
				})
			}
		})
	}
}
