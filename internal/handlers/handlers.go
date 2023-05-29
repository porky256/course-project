package handlers

import (
	"encoding/json"
	"github.com/porky256/course-project/internal/config"
	"github.com/porky256/course-project/internal/driver"
	"github.com/porky256/course-project/internal/forms"
	"github.com/porky256/course-project/internal/models"
	"github.com/porky256/course-project/internal/render"
	"github.com/porky256/course-project/internal/repository"
	"github.com/porky256/course-project/internal/repository/dbrepo"
	mock_dbrepo "github.com/porky256/course-project/internal/repository/mock"
	"net/http"
	"strconv"
	"strings"
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

func NewTestHandlers(app *config.AppConfig, render *render.Render, db *mock_dbrepo.MockDatabaseRepo) *Handlers {
	return &Handlers{
		app:    app,
		render: render,
		DB:     db,
	}
}

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
	room, err := h.DB.GetRoom(res.RoomID)
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

// SearchAvailability renders search availability page
func (h *Handlers) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	err := h.render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
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
	room, err := h.DB.GetRoom(roomID)
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

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	err := h.render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

func (h *Handlers) PostLogin(w http.ResponseWriter, r *http.Request) {
	err := h.app.Session.RenewToken(r.Context())
	if err != nil {
		h.app.ErrorLog.Println(err)
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "bad form")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)

	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		err = h.render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		if err != nil {
			h.app.ErrorLog.Println(err)
		}
		return
	}

	id, _, err := h.DB.Authenticate(email, password)
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	h.app.Session.Put(r.Context(), "user_id", id)
	h.app.Session.Put(r.Context(), "flash", "Authenticated successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
