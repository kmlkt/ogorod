package postgres

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/kmlkt/ogorod/abs"
	"github.com/kmlkt/ogorod/utils"
)

type Postgres struct {
	Version string
	files   struct {
		index        string
		src          string
		srcConfigure string
		bin          string
		binInitDb    string
		binPgCtl     string
		data         string
	}
}

const postgresPath = "https://ftp.postgresql.org/pub/source/v%[1]v/postgresql-%[1]v.tar.gz"

func (p *Postgres) Key() string {
	return "postgres"
}

func (p *Postgres) Install() error {
	p.prepare()
	if p.State() != abs.Missing {
		return nil
	}
	return utils.RunAll(p.download, p.configure, p.makeInstall, p.initDb)
}

func (p *Postgres) download() error {
	return utils.Then(p.fetchSource, p.decompress)()
}

func (p *Postgres) fetchSource() (io.Reader, error) {
	path := fmt.Sprintf(postgresPath, p.Version)
	resp, err := http.Get(path)
	if resp == nil {
		return nil, err
	}
	return resp.Body, err
}

func (p *Postgres) decompress(source io.Reader) error {
	mkdir := func(path string) error {
		fmt.Println(path)
		return os.MkdirAll(p.convertArchivePath(path), 0777)
	}
	openFile := func(path string) (io.Writer, error) {
		convPath := p.convertArchivePath(path)
		if convPath == "" {
			return nil, nil
		}
		return os.OpenFile(convPath, os.O_RDWR|os.O_CREATE, 0777)
	}
	writeToFile := func(file io.Writer, data io.Reader) error {
		_, err := io.Copy(file, data)
		return err
	}
	return utils.DecompressTarGz(source, mkdir, utils.ThenAA(openFile, writeToFile))
}

func (p *Postgres) convertArchivePath(path string) string {
	firstDirIndex := strings.IndexRune(path, os.PathSeparator)
	if firstDirIndex == -1 {
		return ""
	}
	pathWithoutFirstDir := path[firstDirIndex+1:]
	return utils.MakePath(p, "src/"+pathWithoutFirstDir)
}

func (p *Postgres) configure() error {
	return utils.RunCommand(p.files.srcConfigure, p.files.src,
		"--prefix", p.files.index)
}

func (p *Postgres) makeInstall() error {
	return utils.RunCommand("make", p.files.src, "install")
}

func (p *Postgres) initDb() error {
	return utils.RunCommand(p.files.binInitDb, p.files.bin,
		p.files.data)
}

func (p *Postgres) State() abs.ServiceState {
	p.prepare()
	exists, err := utils.FileExists(p.files.srcConfigure)
	if !exists || err != nil {
		return abs.Missing
	}
	output, err := utils.RunCommandSilent(p.files.binPgCtl, p.files.bin,
		"-D", p.files.data,
		"status")
	if strings.Contains(output, "no server running") || err != nil {
		return abs.Idle
	}
	return abs.Running
}

func (p *Postgres) Start() error {
	if p.State() != abs.Idle {
		return nil
	}
	p.prepare()
	return utils.RunCommand(p.files.binPgCtl, p.files.bin,
		"-D", p.files.data,
		"start")
}

func (p *Postgres) Stop() error {
	if p.State() != abs.Running {
		return nil
	}
	p.prepare()
	return utils.RunCommand(p.files.binPgCtl, p.files.bin,
		"-D", p.files.data,
		"stop")
}

func (p *Postgres) prepare() {
	if p.files.src == "" {
		p.files.index = utils.MakePath(p, "")
		p.files.src = utils.MakePath(p, "src")
		p.files.srcConfigure = utils.MakePath(p, "src/configure")
		p.files.binInitDb = utils.MakePath(p, "bin/initdb")
		p.files.bin = utils.MakePath(p, "bin")
		p.files.binPgCtl = utils.MakePath(p, "bin/pg_ctl")
		p.files.data = utils.MakePath(p, "data")
	}
}
