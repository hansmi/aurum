package aurum

import (
	"io/fs"
	"os"
	"path/filepath"
)

type WriteFileFS interface {
	WriteFile(name string, data []byte, perm os.FileMode) error
}

type writableDirFS struct {
	fs.FS
	dir string
}

var _ WriteFileFS = (*writableDirFS)(nil)

func newWritableDirFS(dir string) *writableDirFS {
	return &writableDirFS{
		FS:  os.DirFS(dir),
		dir: dir,
	}
}

func (f *writableDirFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filepath.Join(f.dir, name), data, perm)
}
