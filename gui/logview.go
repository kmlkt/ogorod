package gui

import (
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/kmlkt/ogorod/sys"
)

type LogView struct {
	BaseView
	Logs string
}

func (v BaseView) logView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	id := r.PathValue("service_id")
	shortid, _ := sys.ParseShortID(id)
	slog.Debug("log", "service_id", id)
	var currentService *sys.Service
	for _, service := range v.Config.Services {
		if service.ID == shortid {
			currentService = &service
			break
		}
	}
	logFile, err := currentService.LogFileReader("run")
	if err != nil {
		slog.Error("open log file", "error", err)
		return
	}
	defer logFile.Close()
	logsBytes, err := io.ReadAll(logFile)
	if err != nil {
		slog.Error("read log file", "error", err)
		return
	}
	slog.Debug("log", "logs", string(logsBytes))
	templates.ExecuteTemplate(w, "logview.html", LogView{v, string(logsBytes)})
	rc := http.NewResponseController(w)
	rc.Flush()
	time.Sleep(time.Second * 3)
	w.Write([]byte("Hello, World!"))
	rc.Flush()
	select {
	case <-r.Context().Done():
		return
	}
}
