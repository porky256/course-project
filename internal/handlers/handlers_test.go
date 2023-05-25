package handlers_test

import (
	"context"
	"encoding/gob"
	"errors"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/porky256/course-project/internal/config"
	"github.com/porky256/course-project/internal/handlers"
	"github.com/porky256/course-project/internal/helpers"
	"github.com/porky256/course-project/internal/models"
	"github.com/porky256/course-project/internal/render"
	mock_dbrepo "github.com/porky256/course-project/internal/repository/mock"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"time"
)

type params struct {
	key   string
	value string
}

var app config.AppConfig

var _ = Describe("Handlers", Ordered, func() {

	var server *httptest.Server
	var ctrl *gomock.Controller
	var h *handlers.Handlers
	var mockDB *mock_dbrepo.MockDatabaseRepo
	BeforeAll(func() {
		ctrl = gomock.NewController(GinkgoT())
		gob.Register(models.Reservation{})
		app = config.AppConfig{Session: scs.New(), UseCache: false, RootPath: "./../.."}
		app.DateLayout = "2006-01-02"
		app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
		app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
		helpers.NewHelpers(&app)
		r := render.NewRender(&app)
		mockDB = mock_dbrepo.NewMockDatabaseRepo(ctrl)
		h = handlers.NewTestHandlers(&app, r, mockDB)
		server = httptest.NewTLSServer(routes(h))
	})

	Context("basic handlers", func() {
		var theTests = []struct {
			name               string
			url                string
			method             string
			expectedStatusCode int
		}{
			{"home", "/", "GET", http.StatusOK},
			{"about", "/about", "GET", http.StatusOK},
			{"gq", "/generals-quarters", "GET", http.StatusOK},
			{"ms", "/majors-suite", "GET", http.StatusOK},
			{"sa", "/search-availability", "GET", http.StatusOK},
			{"contact", "/contact", "GET", http.StatusOK},
		}

		It("iterate through all basic handlers", func() {
			for _, e := range theTests {
				By(e.name)
				resp, err := server.Client().Get(server.URL + e.url)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(e.expectedStatusCode))
			}
		})
	})

	var handler http.HandlerFunc
	var rr *httptest.ResponseRecorder
	doall := func(val *url.Values, reservation *models.Reservation, statusCode int, errorString string, url string, method string) {
		rr = httptest.NewRecorder()

		var req *http.Request
		var err error
		var ctx context.Context
		if val != nil {
			req, err = http.NewRequest(method, url, strings.NewReader(val.Encode()))
		} else {
			req, err = http.NewRequest(method, url, nil)
		}
		Expect(err).ToNot(HaveOccurred())
		ctx, err = getCtx(req, &app)
		Expect(err).ToNot(HaveOccurred())
		req = req.WithContext(ctx)
		req.RequestURI = url
		req.URL.RequestURI()
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if reservation != nil {
			app.Session.Put(ctx, "reservation", *reservation)
		}
		handler.ServeHTTP(rr, req)

		text, ok := app.Session.Get(ctx, "error").(string)
		if errorString != "" {
			Expect(text).To(Equal(errorString))
			Expect(ok).To(Equal(true))
		} else {
			Expect(text).To(Equal(""))
			Expect(ok).To(Equal(false))
		}
		Expect(rr.Code).To(Equal(statusCode))
	}

	Context("MakeReservation", func() {
		var basicRes models.Reservation

		BeforeEach(func() {
			basicRes = models.Reservation{
				RoomID: 1,
			}
			handler = h.MakeReservation
		})

		It("test with right data", func() {
			mockDB.EXPECT().GetRoom(gomock.Eq(1)).Return(&models.Room{
				ID:       1,
				RoomName: "room name",
			}, nil)
			doall(nil, &basicRes, http.StatusOK, "", "/some-url", "GET")
		})

		It("test with incorrect room", func() {
			mockDB.EXPECT().GetRoom(gomock.Eq(100)).Return(nil, errors.New("no such room"))
			basicRes.RoomID = 100
			doall(nil, &basicRes, http.StatusTemporaryRedirect, "no such room", "/some-url", "GET")
		})

		It("test with insufficient reservation", func() {
			doall(nil, nil, http.StatusTemporaryRedirect, "can't find reservation", "/some-url", "GET")
		})
	})

	Context("PostMakeReservation", func() {
		var basicVal url.Values
		var basicRes models.Reservation

		BeforeEach(func() {
			basicVal = url.Values{}
			basicVal.Add("start_date", "2050-01-01")
			basicVal.Add("end_date", "2050-01-02")
			basicVal.Add("first_name", "John")
			basicVal.Add("last_name", "Black")
			basicVal.Add("email", "john@here.com")
			basicVal.Add("phone", "123456789")
			basicVal.Add("room_id", "1")
			basicRes = models.Reservation{
				RoomID: 1,
			}
			handler = h.PostMakeReservation
		})

		It("normal", func() {
			mockDB.EXPECT().InsertReservation(gomock.Any()).Return(1, nil)
			mockDB.EXPECT().InsertRoomRestriction(gomock.Any()).Return(1, nil)

			doall(&basicVal, &basicRes, http.StatusSeeOther, "", "/some-url", "POST")
		})

		It("bad form", func() {
			doall(nil, &basicRes, http.StatusTemporaryRedirect, "bad form", "/some-url", "POST")
		})

		It("no reservation", func() {
			doall(&basicVal, nil, http.StatusTemporaryRedirect, "cannot find reservation", "/some-url", "POST")
		})

		It("form is invalid", func() {
			basicVal.Set("first_name", "")
			doall(&basicVal, &basicRes, http.StatusSeeOther, "", "/some-url", "POST")
		})

		It("can't insert reservation", func() {
			mockDB.EXPECT().InsertReservation(gomock.Any()).Return(0, errors.New("can't insert reservation"))
			doall(&basicVal, &basicRes, http.StatusTemporaryRedirect, "can't insert reservation", "/some-url", "POST")
		})

		It("can't insert room restriction", func() {
			mockDB.EXPECT().InsertReservation(gomock.Any()).Return(1, nil)
			mockDB.EXPECT().InsertRoomRestriction(gomock.Any()).Return(0, errors.New("can't insert room restriction"))

			doall(&basicVal, &basicRes, http.StatusTemporaryRedirect, "can't insert room restriction", "/some-url", "POST")
		})
	})

	Context("PostSearchAvailability", func() {
		var basicVal url.Values

		BeforeEach(func() {
			basicVal = url.Values{}
			basicVal.Add("start", "2050-01-01")
			basicVal.Add("end", "2050-01-02")

			handler = h.PostSearchAvailability
		})

		It("normal", func() {
			mockDB.EXPECT().AvailabilityOfAllRooms(gomock.Any(), gomock.Any()).Return([]models.Room{
				{
					RoomName: "name 1",
					ID:       1,
				},
				{
					RoomName: "name 2",
					ID:       2,
				},
			}, nil)
			doall(&basicVal, nil, http.StatusOK, "", "/some-url", "POST")
		})

		It("bad form", func() {
			doall(nil, nil, http.StatusTemporaryRedirect, "bad form", "/some-url", "POST")
		})

		It("bad start", func() {
			basicVal.Set("start", "bad")
			doall(&basicVal, nil, http.StatusTemporaryRedirect, "bad start time", "/some-url", "POST")
		})

		It("bad end", func() {
			basicVal.Set("end", "bad")
			doall(&basicVal, nil, http.StatusTemporaryRedirect, "bad end time", "/some-url", "POST")
		})

		It("AvailabilityOfAllRooms error", func() {
			mockDB.EXPECT().AvailabilityOfAllRooms(gomock.Any(), gomock.Any()).Return(nil, errors.New("text"))
			doall(&basicVal, nil, http.StatusTemporaryRedirect, "can't get rooms", "/some-url", "POST")
		})

		It("no available rooms", func() {
			mockDB.EXPECT().AvailabilityOfAllRooms(gomock.Any(), gomock.Any()).Return([]models.Room{}, nil)
			doall(&basicVal, nil, http.StatusSeeOther, "sorry, no available rooms on this dates >:(", "/some-url", "POST")
		})
	})

	Context("SearchAvailabilityJson", func() {
		var basicVal url.Values

		BeforeEach(func() {
			basicVal = url.Values{}
			basicVal.Add("start", "2050-01-01")
			basicVal.Add("end", "2050-01-02")
			basicVal.Add("room_id", "1")

			handler = h.SearchAvailabilityJson
		})

		It("normal", func() {
			mockDB.EXPECT().LookForAvailabilityOfRoom(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
			doall(&basicVal, nil, http.StatusOK, "", "/some-url", "POST")
		})

		It("bad form", func() {
			doall(nil, nil, http.StatusTemporaryRedirect, "bad form", "/some-url", "POST")
		})

		It("bad start", func() {
			basicVal.Set("start", "bad")
			doall(&basicVal, nil, http.StatusTemporaryRedirect, "bad start time", "/some-url", "POST")
		})

		It("bad end", func() {
			basicVal.Set("end", "bad")
			doall(&basicVal, nil, http.StatusTemporaryRedirect, "bad end time", "/some-url", "POST")
		})

		It("bad room id", func() {
			basicVal.Set("room_id", "bad")
			doall(&basicVal, nil, http.StatusTemporaryRedirect, "bad room id", "/some-url", "POST")
		})

		It("AvailabilityOfAllRooms error", func() {
			mockDB.EXPECT().LookForAvailabilityOfRoom(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, errors.New("text"))
			doall(&basicVal, nil, http.StatusTemporaryRedirect, "problem with searching room", "/some-url", "POST")
		})
	})

	Context("ReservationSummary", func() {
		var basicRes models.Reservation

		BeforeEach(func() {
			basicRes = models.Reservation{
				RoomID:    1,
				StartDate: time.Date(2050, 1, 2, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2050, 1, 3, 0, 0, 0, 0, time.UTC),
				Room: &models.Room{
					RoomName: "name",
					ID:       1,
				},
			}
			handler = h.ReservationSummary
		})

		It("test with right data", func() {
			doall(nil, &basicRes, http.StatusOK, "", "/some-url", "GET")
		})

		It("test with insufficient reservation", func() {
			doall(nil, nil, http.StatusTemporaryRedirect, "can't find reservation", "/some-url", "GET")
		})
	})

	Context("ChooseRoom", func() {
		var basicRes models.Reservation

		BeforeEach(func() {
			handler = h.ChooseRoom
		})

		It("test with right data", func() {
			doall(nil, &basicRes, http.StatusSeeOther, "", "/choose-room/1", "GET")
		})

		It("test with insufficient room id", func() {
			doall(nil, &basicRes, http.StatusTemporaryRedirect, "can't find such room", "/choose-room/f", "GET")
		})

		It("test with insufficient reservation", func() {
			doall(nil, nil, http.StatusTemporaryRedirect, "can't find reservation", "/choose-room/1", "GET")
		})
	})

	Context("BookRoom", func() {
		var basicRes models.Reservation

		BeforeEach(func() {
			handler = h.BookRoom
		})

		It("test with right data", func() {
			mockDB.EXPECT().GetRoom(gomock.Eq(1)).Return(&models.Room{
				ID:       1,
				RoomName: "room name",
			}, nil)
			doall(nil, &basicRes, http.StatusSeeOther, "", "/book-room?s=2050-01-01&e=2050-01-02&id=1", "GET")
		})

		It("bad start", func() {
			doall(nil, nil, http.StatusTemporaryRedirect, "bad start time", "/book-room?s=insufficient&e=2050-01-02&id=1", "GET")
		})

		It("bad end", func() {
			doall(nil, nil, http.StatusTemporaryRedirect, "bad end time", "/book-room?s=2050-01-01&e=insufficient&id=1", "GET")
		})

		It("test with insufficient room id", func() {
			doall(nil, &basicRes, http.StatusTemporaryRedirect, "can't find such room", "/book-room?s=2050-01-01&e=insufficient&id=insufficient", "GET")
		})
	})

	AfterAll(func() {
		server.Close()
		ctrl.Finish()
	})
})

func routes(handler *handlers.Handlers) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(SessionLoad)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

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
	return mux
}

func SessionLoad(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
}

func getCtx(req *http.Request, app *config.AppConfig) (context.Context, error) {
	return app.Session.Load(req.Context(), req.Header.Get("X-Session"))
}
