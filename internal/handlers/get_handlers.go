package handlers

import (
	"github.com/porky256/course-project/internal/forms"
	"github.com/porky256/course-project/internal/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Home renders home page
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	err := h.render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

// About renders about page
func (h *Handlers) About(w http.ResponseWriter, r *http.Request) {
	err := h.render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

// Contact renders contact page
func (h *Handlers) Contact(w http.ResponseWriter, r *http.Request) {
	err := h.render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

// GeneralsQuarters renders room page
func (h *Handlers) GeneralsQuarters(w http.ResponseWriter, r *http.Request) {
	err := h.render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

// MajorsSuite renders room page
func (h *Handlers) MajorsSuite(w http.ResponseWriter, r *http.Request) {
	err := h.render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

// MakeReservation renders make reservation page
func (h *Handlers) MakeReservation(w http.ResponseWriter, r *http.Request) {
	res, ok := h.app.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		h.app.ErrorLog.Printf("can't find reservation")
		h.app.Session.Put(r.Context(), "error", "can't find reservation")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	room, err := h.DB.GetRoomByID(res.RoomID)
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "no such room")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	res.Room = room

	h.app.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format(h.app.DateLayout)
	ed := res.EndDate.Format(h.app.DateLayout)

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	err = h.render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

// SearchAvailability renders search availability page
func (h *Handlers) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	err := h.render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

// ReservationSummary handles request for reservation summary
func (h *Handlers) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	res, ok := h.app.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		h.app.ErrorLog.Println("Can't find reservation >:(")
		h.app.Session.Put(r.Context(), "error", "can't find reservation")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	h.app.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = res

	sd := res.StartDate.Format(h.app.DateLayout)
	ed := res.EndDate.Format(h.app.DateLayout)

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	err := h.render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

// ChooseRoom handles request to choose room
func (h *Handlers) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.URL.RequestURI(), "/")
	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "can't find such room")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	res, ok := h.app.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		h.app.ErrorLog.Printf("cannot find reservation")
		h.app.Session.Put(r.Context(), "error", "can't find reservation")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	res.RoomID = roomID
	h.app.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoom handles request to book room
func (h *Handlers) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "can't find such room")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")
	startDate, err := time.Parse(h.app.DateLayout, sd)
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "bad start time")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	endDate, err := time.Parse(h.app.DateLayout, ed)
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "bad end time")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	room, err := h.DB.GetRoomByID(roomID)
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "no such room")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
		Room:      room,
	}

	h.app.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
