package gui

import (
	"embed"
	"log/slog"
	"net/http"
	"strings"
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
}

func GuiHandler(c *sys.Config, pm *sys.ProcessMap, apply func() error) http.Handler {
	mux := http.NewServeMux()
	view := BaseView{Config: c, ProcessMap: pm}
	mux.HandleFunc("/", view.index)
	mux.HandleFunc("/{service_id}", view.index)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
	return mux
}

type ServiceView struct {
	BaseView
	Service *sys.Service
}

func (v ServiceView) IframeURL() string {
	if v.Service == nil {
		return ""
	}
	result := v.Service.URL
	if strings.HasPrefix(result, "/") {
		result = "http://localhost:8080" + result // TODO rm hardcode
	}
	if !strings.HasPrefix(result, "http") {
		result = "https://" + result
	}
	return result
}

func (v BaseView) index(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("service_id")
	slog.Debug("index", "service_id", id)
	var currentService *sys.Service
	for _, service := range v.Config.Services {
		if service.ID.String() == id {
			currentService = &service
			break
		}
	}
	templates.ExecuteTemplate(w, "index.html", ServiceView{v, currentService})
}
