package gui

import (
	"embed"
	"net/http"
	"text/template"

	"github.com/kmlkt/ogorod/sys"
)

//go:embed *.html
var templateFS embed.FS
var templates = template.Must(template.ParseFS(templateFS, "*.html"))

//go:embed *.css
var staticFS embed.FS

type BaseView struct {
	Config     *sys.Config
	ProcessMap *sys.ProcessMap
	Apply      Applier
}

type Applier func(changed []sys.ShortID) error

func GuiHandler(c *sys.Config, pm *sys.ProcessMap, apply Applier) http.Handler {
	mux := http.NewServeMux()
	view := BaseView{Config: c, ProcessMap: pm, Apply: apply}
	mux.HandleFunc("/", view.serviceView)
	mux.HandleFunc("/{service_id}", view.serviceView)
	mux.HandleFunc("/logs/{service_id}", view.logView)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
	return mux
}
