package sys

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

func UseConfig(config Config, pm *ProcessMap, server *http.Server) error {
	mux := http.NewServeMux()
	err := config.DoStuff(pm, func(newPm ProcessMap) {
		for _, service := range config.Services {
			mux.Handle(service.URL, http.StripPrefix(service.URL, service.Handler(newPm[service.ID])))
		}
		time.Sleep(time.Second)
	})
	server.Handler = mux
	if err != nil {
		return err
	}
	return nil
}

func (s Service) Handler(p Process) http.Handler {
	if s.Hosting == StaticHosting {
		return http.FileServerFS(os.DirFS(p.Path))
	}
	return httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("127.0.0.1:%d", p.Port),
	})
}

func FreePort() (port int, err error) {
	address, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return
	}
	listener, err := net.ListenTCP("tcp", address)
	defer listener.Close()
	if err != nil {
		return
	}
	port = listener.Addr().(*net.TCPAddr).Port
	return
}
