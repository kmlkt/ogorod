package archive

import (
	"archive/tar"
	"compress/gzip"
	"io"
)

func DecompressTarGz(targz io.Reader, open func(path string) (io.Writer, error)) error {
	ungzip, err := gzip.NewReader(targz)
	if err != nil {
		return err
	}
	untar := tar.NewReader(ungzip)
	for f, err := untar.Next(); f != nil && err == nil; f, err = untar.Next() {
		if f.FileInfo().IsDir() {
			continue
		}
		writer, err := open(f.Name)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, untar)
		if err != nil {
			return err
		}
	}
	return nil
}
