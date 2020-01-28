package updater

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io"
	"github.com/g4me92bd777b8b16ed4c/common/types"
)

type File struct {
	Name  string
	Bytes []byte
	Dir   string
}

func (f *File) Decode(b []byte) error {
	if err := gob.NewDecoder(bytes.NewBuffer(b)).Decode(f); err != nil && err != io.EOF {
		return err
	}
	return nil
}
func (f *File) Encode(b []byte) (int, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(f); err != nil {
		return 0, err
	}
	n := copy(b, f.Bytes)
	if n < len(f.Bytes) {
		return n, errors.New("short write")
	}
	return n, nil
}

func (f File) Type() byte {
	return types.File
}