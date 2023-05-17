package handlers_test

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/porky256/course-project/internal/config"
	"github.com/porky256/course-project/internal/handlers"
	"github.com/porky256/course-project/internal/models"
	"github.com/porky256/course-project/internal/render"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Handlers", Ordered, func() {

	var getTestCases = []struct {
		url                string
		expectedStatusCode int
	}{
		{url: "/", expectedStatusCode: http.StatusOK},
		{url: "/about", expectedStatusCode: http.StatusOK},
		{url: "/generals-quarters", expectedStatusCode: http.StatusOK},
		{url: "/majors-suite", expectedStatusCode: http.StatusOK},
		{url: "/make-reservation", expectedStatusCode: http.StatusOK},
		{url: "/search-availability", expectedStatusCode: http.StatusOK},
		{url: "/contact", expectedStatusCode: http.StatusOK},
	}

	var postTestCases = []struct {
		url    string
		params []struct {
			key   string
			value string
		}
		expectedStatusCode int
	}{
		{url: "/", params},
	}

	var server *httptest.Server
	var app config.AppConfig
	BeforeAll(func() {
		gob.Register(models.Reservation{})
		app = config.AppConfig{Session: scs.New()}
		r := render.NewRender(&app)
		h := handlers.NewHandlers(&app, r)
		server = httptest.NewTLSServer(routes(h))
	})

	It("check GET handlers", func() {

		for _, testCase := range getTestCases {
			By(fmt.Sprintf("testing '%s' page", testCase.url))
			resp, err := server.Client().Get(server.URL + testCase.url)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(testCase.expectedStatusCode))
		}
	})

	//It(" check reservation summary Get handler", func() {
	//	testReservation := models.Reservation{
	//		FirstName: "John",
	//		LastName:  "Black",
	//		Email:     "test@test.com",
	//		Phone:     "555-555-5555",
	//	}
	//	app.Session.Put(context.Background(), "reservation", testReservation)
	//	resp, err := server.Client().Get(server.URL + "/reservation-summary")
	//	Expect(err).ToNot(HaveOccurred())
	//	Expect(resp.StatusCode).To(Equal(http.StatusOK))
	//})

	AfterAll(func() {
		server.Close()
	})
})

func routes(handler *handlers.Handlers) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

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

	mux.Get("/reservation-summary", http.HandlerFunc(handler.ReservationSummary))

	mux.Get("/contact", http.HandlerFunc(handler.Contact))

	return mux
}
