package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
)

func DecompressTarGz(targz io.Reader,
	mkdir func(path string) error,
	write func(path string, data io.Reader) error) error {
	decompressGz := gzip.NewReader
	decompressTar := func(ungz *gzip.Reader) error {
		untar := tar.NewReader(ungz)
		files := NextIter(untar.Next)
		save := Apply(files, func(f *tar.Header) error {
			fmt.Println(f.Name)
			if f.FileInfo().IsDir() {
				return mkdir(f.Name)
			} else {
				return write(f.Name, untar)
			}
		})
		return RunAllIter(save)
	}
	return ThenA(decompressGz, decompressTar)(targz)
}
