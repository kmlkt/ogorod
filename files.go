package ogorod

import (
	"io"
	"os"
	"path/filepath"
)

var basePath = "./"

func OpenFile(path ...string) (*os.File, error) {
	fullPath := filepath.Join(path...)
	os.MkdirAll(filepath.Dir(fullPath), 0777)
	return os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0777)
}

func OpenServiceFile(s Service, path ...string) (io.ReadWriter, error) {
	return OpenFile(basePath, "services", s.Key(), filepath.Join(path...))
}
