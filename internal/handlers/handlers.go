package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/porky256/course-project/internal/config"
	"github.com/porky256/course-project/internal/forms"
	"github.com/porky256/course-project/internal/helpers"
	"github.com/porky256/course-project/internal/models"
	"github.com/porky256/course-project/internal/render"
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

// Home renders home page
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	h.render.RenderTemplateV3(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About renders about page
func (h *Handlers) About(w http.ResponseWriter, r *http.Request) {

	h.render.RenderTemplateV3(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Contact renders contact page
func (h *Handlers) Contact(w http.ResponseWriter, r *http.Request) {

	h.render.RenderTemplateV3(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// GeneralsQuarters renders room page
func (h *Handlers) GeneralsQuarters(w http.ResponseWriter, r *http.Request) {

	h.render.RenderTemplateV3(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// MajorsSuite renders room page
func (h *Handlers) MajorsSuite(w http.ResponseWriter, r *http.Request) {

	h.render.RenderTemplateV3(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// MakeReservation renders make reservation page
func (h *Handlers) MakeReservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := map[string]interface{}{}
	data["reservation"] = emptyReservation

	h.render.RenderTemplateV3(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostMakeReservation handles the posting of a reservation form
func (h *Handlers) PostMakeReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		h.render.RenderTemplateV3(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	h.app.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// SearchAvailability renders search availability page
func (h *Handlers) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	h.render.RenderTemplateV3(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostSearchAvailability handles the posting of a search availability form
func (h *Handlers) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("start date: %s, end date: %s", start, end)))
}

type jsonExample struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// SearchAvailabilityJson handles request for availability and sends JSON response
func (h *Handlers) SearchAvailabilityJson(w http.ResponseWriter, r *http.Request) {
	req := jsonExample{
		OK:      true,
		Message: "example message",
	}

	out, err := json.MarshalIndent(req, "", "\t")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (h *Handlers) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	res, ok := h.app.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		h.app.ErrorLog.Println("Can't find reservation")
		h.app.Session.Put(r.Context(), "error", "can't find reservation")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	h.app.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = res

	h.render.RenderTemplateV3(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
