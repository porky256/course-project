package handlers

import (
	"fmt"
	"github.com/porky256/course-project/internal/forms"
	"github.com/porky256/course-project/internal/models"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	now := time.Now()
	if r.URL.Query().Get("y") != "" {
		year, err := strconv.Atoi(r.URL.Query().Get("y"))
		if err != nil {
			h.app.ErrorLog.Println(err)
			h.app.Session.Put(r.Context(), "error", fmt.Sprintf("can't get year from url: %s", r.RequestURI))
			http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
			return
		}
		month, err := strconv.Atoi(r.URL.Query().Get("m"))
		if err != nil {
			h.app.ErrorLog.Println(err)
			h.app.Session.Put(r.Context(), "error", fmt.Sprintf("can't get month from url: %s", r.RequestURI))
			http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
			return
		}
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	last := now.AddDate(0, -1, 0)
	next := now.AddDate(0, 1, 0)
	stringMap := make(map[string]string)
	stringMap["last_month"] = last.Format("01")
	stringMap["last_year"] = last.Format("2006")
	stringMap["next_month"] = next.Format("01")
	stringMap["next_year"] = next.Format("2006")
	stringMap["now_month"] = now.Format("01")
	stringMap["now_year"] = now.Format("2006")

	y, m, _ := now.Date()
	intMap := make(map[string]int)
	firstDayOfMonth := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)
	intMap["number_of_days"] = lastDayOfMonth.Day()

	data := make(map[string]interface{})
	data["now"] = now
	data["last"] = last
	data["next"] = next

	rooms, err := h.DB.GetAllRooms()

	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "can't get rooms")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}
	data["rooms"] = rooms

	for _, room := range rooms {
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for current := firstDayOfMonth; !current.After(lastDayOfMonth); current = current.AddDate(0, 0, 1) {
			reservationMap[current.Format("2006-01-2")] = 0
			blockMap[current.Format("2006-01-2")] = 0
		}
		roomRestrictions, err := h.DB.GetRoomRestrictionsByRoomIdWithinDates(room.ID, firstDayOfMonth, lastDayOfMonth.AddDate(0, 0, 1))
		if err != nil {
			h.app.ErrorLog.Println(err)
			h.app.Session.Put(r.Context(), "error", "can't get room restrictions")
			http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
			return
		}
		for _, rr := range roomRestrictions {
			if rr.Reservation != nil {
				for current := rr.Reservation.StartDate; !current.Equal(rr.Reservation.EndDate); current = current.AddDate(0, 0, 1) {
					reservationMap[current.Format("2006-01-2")] = rr.ReservationID
				}
			} else {
				for current := rr.StartDate; !current.Equal(rr.EndDate); current = current.AddDate(0, 0, 1) {
					blockMap[current.Format("2006-01-2")] = rr.ID
				}
			}
		}

		data[fmt.Sprintf("reservation_map_%d", room.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", room.ID)] = blockMap

		h.app.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", room.ID), blockMap)
	}

	err = h.render.Template(w, r, "admin.reservation-calendar.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		IntMap:    intMap,
		Data:      data,
	})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

func (h *Handlers) AdminSingleReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	if len(exploded) < 5 {
		h.app.ErrorLog.Printf("incorrect request url: %s", r.RequestURI)
		h.app.Session.Put(r.Context(), "error", "incorrect request url")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	}

	stringMap := make(map[string]string)
	stringMap["src"] = exploded[3]

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	stringMap["year"] = year
	stringMap["month"] = month

	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "wrong id")
		http.Redirect(w, r, fmt.Sprintf("/admin/%s-reservations", exploded[3]), http.StatusSeeOther)
		return
	}
	reservation, err := h.DB.GetReservationByID(id)
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "can't find reservation")
		http.Redirect(w, r, fmt.Sprintf("/admin/%s-reservations", exploded[3]), http.StatusSeeOther)
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

func (h *Handlers) AdminPostSingleReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")

	if len(exploded) < 5 {
		h.app.ErrorLog.Printf("incorrect request url: %s", r.RequestURI)
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	}

	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "wrong id")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "bad form")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}
	redirectString := fmt.Sprintf("/admin/%s-reservations", exploded[3])
	month := r.Form.Get("month")
	year := r.Form.Get("year")
	if month != "" && year != "" {
		redirectString = fmt.Sprintf("/admin/reservation-calendar?y=%s&m=%s", year, month)
	}

	reservation := models.Reservation{}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")
	reservation.ID = id

	err = h.DB.UpdateReservation(reservation)
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "can't update reservation")
		http.Redirect(w, r, redirectString, http.StatusSeeOther)
		return
	}

	h.app.Session.Put(r.Context(), "flash", "reservation updated")
	http.Redirect(w, r, redirectString, http.StatusSeeOther)
}

func (h *Handlers) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")

	if len(exploded) != 5 {
		h.app.ErrorLog.Printf("incorrect request url: %s", r.RequestURI)
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	}

	redirectString := fmt.Sprintf("/admin/%s-reservations", exploded[3])
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")
	if month != "" && year != "" {
		redirectString = fmt.Sprintf("/admin/reservation-calendar?y=%s&m=%s", year, month)
	}

	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "wrong id")
		http.Redirect(w, r, redirectString, http.StatusSeeOther)
		return
	}

	err = h.DB.UpdateReservationProcessed(id, 1)
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "can't update reservation")
		http.Redirect(w, r, redirectString, http.StatusSeeOther)
		return
	}

	h.app.Session.Put(r.Context(), "flash", "reservation is marked as processed")
	http.Redirect(w, r, redirectString, http.StatusSeeOther)
}

func (h *Handlers) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")

	if len(exploded) < 5 {
		h.app.ErrorLog.Printf("incorrect request url: %s", r.RequestURI)
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	}

	redirectString := fmt.Sprintf("/admin/%s-reservations", exploded[3])
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")
	if month != "" && year != "" {
		redirectString = fmt.Sprintf("/admin/reservation-calendar?y=%s&m=%s", year, month)
	}

	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "wrong id")
		http.Redirect(w, r, redirectString, http.StatusSeeOther)
		return
	}

	err = h.DB.DeleteReservationByID(id)
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "can't update reservation")
		http.Redirect(w, r, redirectString, http.StatusSeeOther)
		return
	}

	h.app.Session.Put(r.Context(), "flash", "reservation is deleted")
	http.Redirect(w, r, redirectString, http.StatusSeeOther)
}

func (h *Handlers) AdminPostReservationCalendar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "bad form")
		http.Redirect(w, r, "/admin/reservation-calendar", http.StatusSeeOther)
		return
	}

	year, err := strconv.Atoi(r.Form.Get("y"))
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "wrong year")
		http.Redirect(w, r, "/admin/reservation-calendar", http.StatusSeeOther)
		return
	}
	month, err := strconv.Atoi(r.Form.Get("m"))
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "wrong month")
		http.Redirect(w, r, "/admin/reservation-calendar", http.StatusSeeOther)
		return
	}

	rooms, err := h.DB.GetAllRooms()
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "can't get rooms")
		http.Redirect(w, r, fmt.Sprintf("/admin/reservation-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
		return
	}

	form := forms.New(r.Form)

	for _, room := range rooms {
		curMap, ok := h.app.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", room.ID)).(map[string]int)
		if !ok {
			h.app.ErrorLog.Println("can't get block map for room: ", room.Name)
			h.app.Session.Put(r.Context(), "error", fmt.Sprintf("can't get block map for room: %s", room.Name))
			http.Redirect(w, r, fmt.Sprintf("/admin/reservation-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
			return
		}
		for name, value := range curMap {
			if value > 0 {
				if !form.Has(fmt.Sprintf("remove_block_%d_%s", room.ID, name)) {
					err = h.DB.DeleteRoomRestrictionByID(value)
					if err != nil {
						h.app.ErrorLog.Println(err)
						h.app.Session.Put(r.Context(), "error", fmt.Sprintf("can't delete restriction %d", value))
						http.Redirect(w, r, fmt.Sprintf("/admin/reservation-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
						return
					}
				}
			}
		}
	}

	for name, _ := range r.Form {
		if strings.HasPrefix(name, "add_block_") {
			exploded := strings.Split(name, "_")
			roomId, err := strconv.Atoi(exploded[2])
			if err != nil {
				h.app.ErrorLog.Println(err)
				h.app.Session.Put(r.Context(), "error", fmt.Sprintf("can't parse room id %s", name))
				http.Redirect(w, r, fmt.Sprintf("/admin/reservation-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
				return
			}
			date, err := time.Parse("2006-01-02", exploded[3])
			if err != nil {
				h.app.ErrorLog.Println(err)
				h.app.Session.Put(r.Context(), "error", fmt.Sprintf("can't parse day %s", name))
				http.Redirect(w, r, fmt.Sprintf("/admin/reservation-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
				return
			}

			_, err = h.DB.AddSingleDayRoomRestriction(roomId, 2, date)
			if err != nil {
				h.app.ErrorLog.Println(err)
				h.app.Session.Put(r.Context(), "error", fmt.Sprintf("can't save this restriction %s", name))
				http.Redirect(w, r, fmt.Sprintf("/admin/reservation-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
				return
			}
		}
	}

	h.app.Session.Put(r.Context(), "flash", "changes saved!")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservation-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
}
