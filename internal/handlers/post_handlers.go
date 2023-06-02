package handlers

import (
	"encoding/json"
	"github.com/porky256/course-project/internal/forms"
	"github.com/porky256/course-project/internal/models"
	"net/http"
	"strconv"
	"time"
)

// PostMakeReservation handles the posting of a reservation form
func (h *Handlers) PostMakeReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "bad form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	reservation, ok := h.app.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		h.app.ErrorLog.Printf("cannot find reservation")
		h.app.Session.Put(r.Context(), "error", "cannot find reservation")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		err := h.render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		if err != nil {
			h.app.ErrorLog.Println(err)
		}
		return
	}
	h.app.InfoLog.Printf("saving to db reservation: %+v\n", reservation)
	newID, err := h.DB.InsertReservation(&reservation)
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "can't insert reservation")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	h.app.InfoLog.Println("new reservation's id is: ", newID)

	rmrs := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newID,
		RestrictionID: 1,
	}
	rmrsID, err := h.DB.InsertRoomRestriction(&rmrs)
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "can't insert room restriction")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	rmrs.ID = rmrsID
	h.app.InfoLog.Printf("room restriction %+v saved to db\n", rmrs)

	h.app.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// PostSearchAvailability handles the posting of a search availability form
func (h *Handlers) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "bad form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")
	startDate, err := time.Parse(h.app.DateLayout, start)
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "bad start time")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	endDate, err := time.Parse(h.app.DateLayout, end)
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "bad end time")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	rooms, err := h.DB.AvailabilityOfAllRooms(startDate, endDate)
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "can't get rooms")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if len(rooms) == 0 {
		h.app.InfoLog.Printf("no rooms for dates: %s to %s\n", start, end)
		h.app.Session.Put(r.Context(), "error", "sorry, no available rooms on this dates >:(")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"rooms": rooms,
	}

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	h.app.Session.Put(r.Context(), "reservation", res)

	err = h.render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
	return
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	RoomID    string `json:"room_id"`
}

// SearchAvailabilityJson handles request for availability and sends JSON response
func (h *Handlers) SearchAvailabilityJson(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "bad form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

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

	rid := r.Form.Get("room_id")
	roomID, err := strconv.Atoi(rid)
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "bad room id")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ok, err := h.DB.LookForAvailabilityOfRoom(startDate, endDate, roomID)
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "problem with searching room")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	req := jsonResponse{
		OK:        ok,
		StartDate: sd,
		EndDate:   ed,
		RoomID:    rid,
	}

	out, err := json.MarshalIndent(req, "", "\t")
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "json marshalling error")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(out)
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}
