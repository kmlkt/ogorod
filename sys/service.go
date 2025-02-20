package sys

type Service struct {
	ID         ShortID
	Label      string
	Repository string
	Branch     string
	URL        string
	Hosting    HostingMethod
	BuildCmd   string
	RunCmd     string
}

type HostingMethod string

const (
	StaticHosting HostingMethod = "static"
	Executable    HostingMethod = "executable"
)

func (s Service) DoStuff(runNecessary bool) (Process, error) {
	changed, path, err := s.download()
	if err != nil {
		return Process{Path: path}, err
	}
	if s.BuildCmd != "" && changed {
		err = s.build(path)
		if err != nil {
			return Process{}, err
		}
	}
	if s.RunCmd != "" && (changed || runNecessary) {
		port, err := FreePort()
		err, pid := s.run(path, port)
		if err != nil {
			return Process{}, err
		}
		return Process{Port: port, Pid: pid, Path: path}, nil
	}
	return Process{Path: path}, nil
}
