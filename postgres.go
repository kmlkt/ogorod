package ogorod

import (
	"fmt"
	"io"
	"net/http"

	"github.com/kmlkt/ogorod/archive"
)

type Postgres struct {
	Version string
}

const postgresPath = "https://ftp.postgresql.org/pub/source/v%[1]v/postgresql-%[1]v.tar.gz"

func (p Postgres) Key() string {
	return "postgres"
}

func (p Postgres) Install() error {
	path := fmt.Sprintf(postgresPath, p.Version)
	resp, err := http.Get(path)
	if err != nil {
		return err
	}
	err = archive.DecompressTarGz(resp.Body, func(path string) (io.Writer, error) {
		return OpenServiceFile(p, "src", path)
	})
	if err != nil {
		return err
	}

	return nil
}

func (p Postgres) State() ServiceState {
	return Idle
}

func (p Postgres) Start() error {
	return nil
}

func (p Postgres) Stop() error {
	return nil
}
