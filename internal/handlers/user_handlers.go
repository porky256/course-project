package handlers

import (
	"github.com/porky256/course-project/internal/forms"
	"github.com/porky256/course-project/internal/models"
	"net/http"
)

// Login handles request to login
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	err := h.render.Template(w, r, "user.login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
	if err != nil {
		h.app.ErrorLog.Println(err)
	}
}

// PostLogin handles request to post login
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
		err = h.render.Template(w, r, "user.login.page.tmpl", &models.TemplateData{
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
	//h.app.Session.Put(r.Context(), "flash", "Authenticated successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout handles request to logout
func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	err := h.app.Session.Destroy(r.Context())
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "Error with logging out")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err = h.app.Session.RenewToken(r.Context())
	if err != nil {
		h.app.ErrorLog.Println(err)
		h.app.Session.Put(r.Context(), "error", "Error with logging out")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
