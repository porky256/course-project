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
	h.render.RenderTemplateV3(w, "home.page.html", &models.TemplateData{})
}

func (h *Handlers) About(w http.ResponseWriter, r *http.Request) {

	h.render.RenderTemplateV3(w, "about.page.html", &models.TemplateData{})
}

func (h *Handlers) Contact(w http.ResponseWriter, r *http.Request) {

	h.render.RenderTemplateV3(w, "contact.page.html", &models.TemplateData{})
}

func (h *Handlers) GeneralsQarters(w http.ResponseWriter, r *http.Request) {

	h.render.RenderTemplateV3(w, "generals.page.html", &models.TemplateData{})
}

func (h *Handlers) MajorsSuite(w http.ResponseWriter, r *http.Request) {

	h.render.RenderTemplateV3(w, "majors.page.html", &models.TemplateData{})
}

func (h *Handlers) MakeReservation(w http.ResponseWriter, r *http.Request) {

	h.render.RenderTemplateV3(w, "make-reservation.page.html", &models.TemplateData{})
}

func (h *Handlers) SearchAvailability(w http.ResponseWriter, r *http.Request) {

	h.render.RenderTemplateV3(w, "search-availability.page.html", &models.TemplateData{})
}
