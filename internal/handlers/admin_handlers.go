package handlers

import (
	"github.com/porky256/course-project/internal/forms"
	"github.com/porky256/course-project/internal/models"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handlers) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	err := h.render.Template(w, r, "admin.dashboard.page.tmpl", &models.TemplateData{})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

func (h *Handlers) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	reservations, err := h.DB.GetAllReservations()
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "can't get all reservations")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}
	data["reservations"] = reservations
	err = h.render.Template(w, r, "admin.all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

func (h *Handlers) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	reservations, err := h.DB.GetNewReservations()
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "can't get new reservations")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}
	data["reservations"] = reservations
	err = h.render.Template(w, r, "admin.new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

func (h *Handlers) AdminReservationCalendar(w http.ResponseWriter, r *http.Request) {
	err := h.render.Template(w, r, "admin.reservation-calendar.page.tmpl", &models.TemplateData{})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

func (h *Handlers) AdminSingleReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	stringMap := make(map[string]string)
	stringMap["src"] = exploded[3]
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "wrong id")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}
	reservation, err := h.DB.GetReservationByID(id)
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "can't find reservation")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}
	data := make(map[string]interface{})
	data["reservation"] = reservation
	err = h.render.Template(w, r, "admin.single-reservation.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}
