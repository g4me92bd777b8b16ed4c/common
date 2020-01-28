package updater

import (
	//"bytes"
	//"encoding/gob"
	//"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/g4me92bd777b8b16ed4c/common/types"
)

var Dir = "dlfiles"

type File struct {
	Name  string
	Bytes []byte
	Dir   string
}

func (f *File) WriteToDisk() error {
	return ioutil.WriteFile(filepath.Join(Dir, f.Dir, f.Name), f.Bytes, 0755)
}

func (f File) Type() types.Type {
	return types.File
}
