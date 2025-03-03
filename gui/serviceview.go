package gui

import (
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/kmlkt/ogorod/sys"
)

type ServiceView struct {
	BaseView
	Service *sys.Service
	Random  sys.ShortID
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

func (v BaseView) serviceView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("service_id")
	shortid, _ := sys.ParseShortID(id)
	slog.Debug("index", "service_id", id)
	var currentService *sys.Service
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			slog.Error("parse form", "error", err)
			return
		}
		currentService = &sys.Service{
			ID:         shortid,
			Label:      r.FormValue("Label"),
			Repository: r.FormValue("Repository"),
			Branch:     r.FormValue("Branch"),
			Hosting:    sys.HostingMethod(r.FormValue("Hosting")),
			URL:        r.FormValue("URL"),
			BuildCmd:   r.FormValue("BuildCmd"),
			RunCmd:     r.FormValue("RunCmd"),
		}
	}
	for i, service := range v.Config.Services {
		if service.ID == shortid {
			if currentService == nil {
				currentService = &service
			}
			v.Config.Services[i] = *currentService
			break
		}
	}
	if r.Method == "POST" {
		slices.SortStableFunc(v.Config.Services, func(a, b sys.Service) int {
			return strings.Compare(a.Label, b.Label)
		})
		v.Apply([]sys.ShortID{shortid})
	}
	templates.ExecuteTemplate(w, "serviceview.html", ServiceView{v, currentService, sys.NewShortID()})
}
