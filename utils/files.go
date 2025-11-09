package utils

import (
	"os"
	"path/filepath"

	"github.com/kmlkt/ogorod/abs"
)

var basePath string
var servicesPath string

func init() {
	basePath, _ = os.Getwd()
	servicesPath = filepath.Join(basePath, "services")
}

func FileExists(fullPath string) (bool, error) {
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func MakePath(s abs.Service, path string) string {
	innerPath := filepath.FromSlash(path)
	fullPath := filepath.Join(servicesPath, s.Key(), innerPath)
	return fullPath
}
