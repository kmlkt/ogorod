package ogorod

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

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
	if p.State() != Missing {
		return nil
	}
	path := fmt.Sprintf(postgresPath, p.Version)
	resp, err := http.Get(path)
	if err != nil {
		return err
	}
	err = archive.DecompressTarGz(resp.Body, func(path string) (io.Writer, error) {
		firstDirIndex := strings.IndexRune(path, os.PathSeparator)
		if firstDirIndex == -1 {
			return nil, nil
		}
		pathWithoutFirstDir := path[firstDirIndex+1:]
		return OpenServiceFile(p, "src", pathWithoutFirstDir)
	})
	if err != nil {
		return err
	}
	configure := RunServiceFile(p, "src", "configure")
	configure.Args = append(configure.Args, "--prefix="+ServiceFilePath(p, []string{}))
	configure.Stdout = os.Stdout
	configure.Stderr = os.Stdout
	configure.Dir = ServiceFilePath(p, []string{"src"})
	err = configure.Start()
	if err != nil {
		return err
	}
	err = configure.Wait()
	if err != nil {
		return err
	}

	makeInstall := exec.Command("make", "install")
	makeInstall.Stdout = os.Stdout
	makeInstall.Stderr = os.Stdout
	makeInstall.Dir = ServiceFilePath(p, []string{"src"})
	makeInstall.Env = configure.Env
	err = makeInstall.Start()
	if err != nil {
		return err
	}
	err = makeInstall.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (p Postgres) State() ServiceState {
	exists, err := ServiceFileExists(p, "postgres")
	if !exists || err != nil {
		return Missing
	} else {
		return Idle
	}
}

func (p Postgres) Start() error {
	return nil
}

func (p Postgres) Stop() error {
	return nil
}
