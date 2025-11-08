package ogorod

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

var basePath string

func init() {
	basePath, _ = os.Getwd()
}

func OpenFile(path ...string) (*os.File, error) {
	fullPath := filepath.Join(path...)
	os.MkdirAll(filepath.Dir(fullPath), 0777)
	return os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0777)
}

func FileExists(path ...string) (bool, error) {
	fullPath := filepath.Join(path...)
	os.MkdirAll(filepath.Dir(fullPath), 0777)
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func ServiceFilePath(s Service, path []string) string {
	return filepath.Join(append([]string{basePath, "services", s.Key()}, path...)...)
}

func OpenServiceFile(s Service, path ...string) (io.ReadWriter, error) {
	return OpenFile(ServiceFilePath(s, path))
}

func ServiceFileExists(s Service, path ...string) (bool, error) {
	return FileExists(ServiceFilePath(s, path))
}

func RunServiceFile(s Service, path ...string) *exec.Cmd {
	return exec.Command(ServiceFilePath(s, path))
}
