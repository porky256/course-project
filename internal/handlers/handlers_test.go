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
	"github.com/porky256/course-project/internal/driver"
	"github.com/porky256/course-project/internal/handlers"
	"github.com/porky256/course-project/internal/helpers"
	"github.com/porky256/course-project/internal/models"
	"github.com/porky256/course-project/internal/render"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
)

type params struct {
	key   string
	value string
}

var app config.AppConfig

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
		url                string
		params             []params
		expectedStatusCode int
	}{
		{url: "/make-reservation", params: []params{
			{key: "start", value: "2023-01-02"},
			{key: "end", value: "2023-01-04"},
		}, expectedStatusCode: http.StatusOK},
		{url: "/search-availability", params: []params{
			{key: "first_name", value: "John"},
			{key: "last_name", value: "Black"},
			{key: "email", value: "test@test.com"},
			{key: "phone", value: "555-555-5555"},
		}, expectedStatusCode: http.StatusOK},
		{url: "/search-availability-json", params: []params{
			{key: "first_name", value: "John"},
			{key: "last_name", value: "Black"},
			{key: "email", value: "test@test.com"},
			{key: "phone", value: "555-555-5555"},
		}, expectedStatusCode: http.StatusOK},
	}

	var server *httptest.Server
	BeforeAll(func() {
		gob.Register(models.Reservation{})
		app = config.AppConfig{Session: scs.New(), UseCache: false, RootPath: "./../.."}
		app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
		app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
		helpers.NewHelpers(&app)
		r := render.NewRender(&app)
		h := handlers.NewHandlers(&app, r, &driver.DB{})
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

	It("check POST handlers", func() {
		for _, testCase := range postTestCases {
			By(fmt.Sprintf("testing '%s' page", testCase.url))
			values := url.Values{}
			for _, value := range testCase.params {
				values.Add(value.key, value.value)
			}
			resp, err := server.Client().PostForm(server.URL+testCase.url, values)
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

	mux.Get("/reservation-summary", http.HandlerFunc(handler.ReservationSummary))

	mux.Get("/contact", http.HandlerFunc(handler.Contact))

	return mux
}

func SessionLoad(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
}
