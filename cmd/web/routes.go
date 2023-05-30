package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/porky256/course-project/internal/config"
	"github.com/porky256/course-project/internal/handlers"
	"net/http"
)

func routes(app *config.AppConfig, handler *handlers.Handlers) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", http.HandlerFunc(handler.Home))
	mux.Get("/about", http.HandlerFunc(handler.About))
	mux.Get("/generals-quarters", http.HandlerFunc(handler.GeneralsQuarters))
	mux.Get("/majors-suite", http.HandlerFunc(handler.MajorsSuite))

	mux.Get("/make-reservation", http.HandlerFunc(handler.MakeReservation))
	mux.Post("/make-reservation", http.HandlerFunc(handler.PostMakeReservation))

	mux.Get("/search-availability", http.HandlerFunc(handler.SearchAvailability))
	mux.Post("/search-availability", http.HandlerFunc(handler.PostSearchAvailability))
	mux.Post("/search-availability-json", http.HandlerFunc(handler.SearchAvailabilityJson))
	mux.Get("/choose-room/{id}", http.HandlerFunc(handler.ChooseRoom))
	mux.Get("/book-room", http.HandlerFunc(handler.BookRoom))

	mux.Get("/reservation-summary", http.HandlerFunc(handler.ReservationSummary))

	mux.Get("/contact", http.HandlerFunc(handler.Contact))

	mux.Route("/user", func(r chi.Router) {
		r.Get("/login", http.HandlerFunc(handler.Login))
		r.Post("/login", http.HandlerFunc(handler.PostLogin))
		r.Get("/logout", http.HandlerFunc(handler.Logout))
	})

	mux.Route("/admin", func(r chi.Router) {
		//r.Use(Auth)
		r.Get("/dashboard", http.HandlerFunc(handler.AdminDashboard))
		r.Get("/new-reservations", http.HandlerFunc(handler.AdminNewReservations))
		r.Get("/all-reservations", http.HandlerFunc(handler.AdminAllReservations))
		r.Get("/reservation-calendar", http.HandlerFunc(handler.AdminReservationCalendar))
		r.Get("/reservations/{src}/{id}", http.HandlerFunc(handler.AdminSingleReservation))
		r.Post("/reservations/{src}/{id}", http.HandlerFunc(handler.AdminPostSingleReservation))
		r.Get("/process-reservation/{src}/{id}", http.HandlerFunc(handler.AdminProcessReservation))
		r.Get("/delete-reservation/{src}/{id}", http.HandlerFunc(handler.AdminDeleteReservation))
	})

	return mux
}
