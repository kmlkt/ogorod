package main

import (
	"log/slog"
	"net/http"

	"github.com/kmlkt/ogorod/gui"
	_ "github.com/kmlkt/ogorod/gui"
	"github.com/kmlkt/ogorod/sys"
)

var sitesEnabled = "/etc/nginx/sites-enabled/ogorod"

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	err := run()
	if err != nil {
		slog.Error("", "err", err)
	}
}

func run() error {
	httpServer := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: &http.ServeMux{},
	}
	config, err := sys.LoadConfig()
	if err != nil {
		return err
	}
	var pm sys.ProcessMap
	go func() {
		sys.GracefulShutdown(&pm)
		httpServer.Shutdown(nil)
	}()
	apply := func(changed []sys.ShortID) error {
		err = sys.UseConfig(config, &pm, &httpServer, changed)
		if err != nil {
			return err
		}
		err := config.Save()
		if err != nil {
			return err
		}
		return nil
	}
	guiServer := http.Server{
		Addr:    "127.0.0.1:8082",
		Handler: gui.GuiHandler(&config, &pm, apply),
	}
	go guiServer.ListenAndServe()
	err = apply([]sys.ShortID{})
	if err != nil {
		return err
	}
	err = httpServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
