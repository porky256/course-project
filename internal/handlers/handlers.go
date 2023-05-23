package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/porky256/course-project/internal/config"
	"github.com/porky256/course-project/internal/driver"
	"github.com/porky256/course-project/internal/forms"
	"github.com/porky256/course-project/internal/helpers"
	"github.com/porky256/course-project/internal/models"
	"github.com/porky256/course-project/internal/render"
	"github.com/porky256/course-project/internal/repository"
	"github.com/porky256/course-project/internal/repository/dbrepo"
	"net/http"
	"strconv"
	"time"
)

type Handlers struct {
	app    *config.AppConfig
	render *render.Render
	DB     repository.DatabaseRepo
}

func NewHandlers(app *config.AppConfig, render *render.Render, db *driver.DB) *Handlers {
	return &Handlers{
		app:    app,
		render: render,
		DB:     dbrepo.NewPostgressDB(db.DB, app),
	}
}

// Home renders home page
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	err := h.render.RenderTemplateV3(w, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

// About renders about page
func (h *Handlers) About(w http.ResponseWriter, r *http.Request) {
	err := h.render.RenderTemplateV3(w, r, "about.page.tmpl", &models.TemplateData{})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

// Contact renders contact page
func (h *Handlers) Contact(w http.ResponseWriter, r *http.Request) {

	err := h.render.RenderTemplateV3(w, r, "contact.page.tmpl", &models.TemplateData{})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

// GeneralsQuarters renders room page
func (h *Handlers) GeneralsQuarters(w http.ResponseWriter, r *http.Request) {

	err := h.render.RenderTemplateV3(w, r, "generals.page.tmpl", &models.TemplateData{})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

// MajorsSuite renders room page
func (h *Handlers) MajorsSuite(w http.ResponseWriter, r *http.Request) {

	err := h.render.RenderTemplateV3(w, r, "majors.page.tmpl", &models.TemplateData{})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

// MakeReservation renders make reservation page
func (h *Handlers) MakeReservation(w http.ResponseWriter, r *http.Request) {
	res, ok := h.app.Session.Pop(r.Context(), "reservation").(models.Reservation)
	if !ok {
		h.app.ErrorLog.Printf("cannot find reservation")
		return
	}
	room, err := h.DB.GetRoom(res.RoomID)
	if err != nil {
		h.app.ErrorLog.Println(err)
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

	err = h.render.RenderTemplateV3(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

// PostMakeReservation handles the posting of a reservation form
func (h *Handlers) PostMakeReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation, ok := h.app.Session.Pop(r.Context(), "reservation").(models.Reservation)
	if !ok {
		h.app.ErrorLog.Printf("cannot find reservation")
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

		err := h.render.RenderTemplateV3(w, r, "make-reservation.page.tmpl", &models.TemplateData{
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
	h.app.InfoLog.Println("new reservation's id is: ", newID)
	if err != nil {
		h.app.ErrorLog.Println(err)
		return
	}

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
		return
	}
	rmrs.ID = rmrsID
	h.app.InfoLog.Printf("room restriction %+v saved to db\n", rmrs)

	h.app.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// SearchAvailability renders search availability page
func (h *Handlers) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	err := h.render.RenderTemplateV3(w, r, "search-availability.page.tmpl", &models.TemplateData{})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

// PostSearchAvailability handles the posting of a search availability form
func (h *Handlers) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	startDate, err := time.Parse(h.app.DateLayout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(h.app.DateLayout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	rooms, err := h.DB.AvailabilityOfAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if len(rooms) == 0 {
		h.app.Session.Put(r.Context(), "error", "sorry, no available rooms on this dates >:(")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
	}

	data := map[string]interface{}{
		"rooms": rooms,
	}

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	h.app.Session.Put(r.Context(), "reservation", res)

	err = h.render.RenderTemplateV3(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
	return
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
	_, err = w.Write(out)
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

func (h *Handlers) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	res, ok := h.app.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		h.app.ErrorLog.Println("Can't find reservation >:(")
		h.app.Session.Put(r.Context(), "error", "can't find reservation")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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

	err := h.render.RenderTemplateV3(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

func (h *Handlers) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.app.ErrorLog.Println(err)
		return
	}
	res, ok := h.app.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		h.app.ErrorLog.Printf("cannot find reservation")
		return
	}
	res.RoomID = roomID
	h.app.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
