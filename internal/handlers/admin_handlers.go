package handlers

import (
	"fmt"
	"github.com/porky256/course-project/internal/models"
	"net/http"
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
	fmt.Println("sd")
	err = h.render.Template(w, r, "admin.all-reservations.page.tmpl", &models.TemplateData{})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

func (h *Handlers) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	err := h.render.Template(w, r, "admin.new-reservations.page.tmpl", &models.TemplateData{})
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
