package handlers

import (
	"github.com/porky256/course-project/pkg/config"
	"github.com/porky256/course-project/pkg/models"
	"github.com/porky256/course-project/pkg/render"
	"net/http"
)

type Handlers struct {
	app    *config.AppConfig
	render *render.Render
}

func NewHandlers(app *config.AppConfig, render *render.Render) *Handlers {
	return &Handlers{
		app:    app,
		render: render,
	}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	h.app.Session.Put(r.Context(), "remote_ip", remoteIP)
	h.render.RenderTemplateV3(w, "home.page.html", &models.TemplateData{})
}

func (h *Handlers) About(w http.ResponseWriter, r *http.Request) {
	remoteIP := h.app.Session.GetString(r.Context(), "remote_ip")
	stringMap := map[string]string{
		"test":      "testdata",
		"remote_ip": remoteIP,
	}

	h.render.RenderTemplateV3(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}
