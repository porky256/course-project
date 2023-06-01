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
		app.MailChan = make(chan models.MailData, 100)
		helpers.NewHelpers(&app)
		r := render.NewRender(&app)
		mockDB = mock_dbrepo.NewMockDatabaseRepo(ctrl)
		h = handlers.NewTestHandlers(&app, r, mockDB)
		server = httptest.NewTLSServer(routes(h))
	})

	AfterAll(func() {
		server.Close()
		ctrl.Finish()
		close(app.MailChan)
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
			{"adminDashboard", "/admin/dashboard", "GET", http.StatusOK},
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

	type testData struct {
		val            *url.Values
		reservation    *models.Reservation
		statusCode     int
		errorString    string
		url            string
		redirectURL    string
		dataForSession map[string]interface{}
	}

	var handler http.HandlerFunc
	var rr *httptest.ResponseRecorder
	var method string
	doall := func(data testData) {
		rr = httptest.NewRecorder()

		var req *http.Request
		var err error
		var ctx context.Context
		if data.val != nil {
			req, err = http.NewRequest(method, data.url, strings.NewReader(data.val.Encode()))
		} else {
			req, err = http.NewRequest(method, data.url, nil)
		}
		Expect(err).ToNot(HaveOccurred())
		ctx, err = getCtx(req, &app)
		Expect(err).ToNot(HaveOccurred())
		req = req.WithContext(ctx)
		req.RequestURI = data.url
		req.URL.RequestURI()
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if data.reservation != nil {
			app.Session.Put(ctx, "reservation", *data.reservation)
		}
		for key, value := range data.dataForSession {
			app.Session.Put(ctx, key, value)
		}
		handler.ServeHTTP(rr, req)

		text, ok := app.Session.Get(ctx, "error").(string)
		if data.errorString != "" {
			Expect(text).To(Equal(data.errorString))
			Expect(ok).To(Equal(true))
		} else {
			Expect(text).To(Equal(""))
			Expect(ok).To(Equal(false))
		}
		Expect(rr.Code).To(Equal(data.statusCode))

		if data.redirectURL != "" {
			actLoc, _ := rr.Result().Location()
			Expect(actLoc.String()).To(Equal(data.redirectURL))
		}
	}

	Context("MakeReservation", func() {
		var basicRes models.Reservation

		BeforeEach(func() {
			basicRes = models.Reservation{
				RoomID: 1,
			}
			handler = h.MakeReservation
			method = "GET"
		})

		It("test with right data", func() {
			mockDB.EXPECT().GetRoomByID(gomock.Eq(1)).Return(&models.Room{
				ID:   1,
				Name: "room name",
			}, nil)
			data := testData{
				val:         nil,
				reservation: &basicRes,
				statusCode:  http.StatusOK,
				errorString: "",
				url:         "/some-url",
			}
			doall(data)
		})

		It("test with incorrect room", func() {
			mockDB.EXPECT().GetRoomByID(gomock.Eq(100)).Return(nil, errors.New("no such room"))
			basicRes.RoomID = 100
			data := testData{
				reservation: &basicRes,
				statusCode:  http.StatusSeeOther,
				errorString: "no such room",
				url:         "/some-url",
				redirectURL: "/",
			}
			doall(data)
		})

		It("test with insufficient reservation", func() {
			data := testData{
				statusCode:  http.StatusSeeOther,
				errorString: "can't find reservation",
				url:         "/some-url",
				redirectURL: "/",
			}
			doall(data)
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
			method = "POST"
		})

		It("normal", func() {
			mockDB.EXPECT().InsertReservation(gomock.Any()).Return(1, nil)
			mockDB.EXPECT().InsertRoomRestriction(gomock.Any()).Return(1, nil)
			mockDB.EXPECT().GetRoomByID(gomock.Any()).Return(&models.Room{
				ID:   1,
				Name: "room name",
			}, nil)
			data := testData{
				val:         &basicVal,
				reservation: &basicRes,
				statusCode:  http.StatusSeeOther,
				url:         "/some-url",
				redirectURL: "/reservation-summary",
			}
			doall(data)
		})

		It("bad form", func() {
			data := testData{
				reservation: &basicRes,
				statusCode:  http.StatusSeeOther,
				errorString: "bad form",
				url:         "/some-url",
				redirectURL: "/",
			}
			doall(data)
		})

		It("no reservation", func() {
			data := testData{
				val:         &basicVal,
				statusCode:  http.StatusSeeOther,
				errorString: "cannot find reservation",
				url:         "/some-url",
				redirectURL: "/",
			}
			doall(data)
		})

		It("form is invalid", func() {
			basicVal.Set("first_name", "")
			data := testData{
				val:         &basicVal,
				reservation: &basicRes,
				statusCode:  http.StatusOK,
				errorString: "",
				url:         "/some-url",
			}
			doall(data)
		})

		It("can't insert reservation", func() {
			mockDB.EXPECT().InsertReservation(gomock.Any()).Return(0, errors.New("can't insert reservation"))
			data := testData{
				val:         &basicVal,
				reservation: &basicRes,
				statusCode:  http.StatusSeeOther,
				errorString: "can't insert reservation",
				url:         "/some-url",
				redirectURL: "/",
			}
			doall(data)
		})

		It("can't insert room restriction", func() {
			mockDB.EXPECT().InsertReservation(gomock.Any()).Return(1, nil)
			mockDB.EXPECT().InsertRoomRestriction(gomock.Any()).Return(0, errors.New("can't insert room restriction"))
			data := testData{
				val:         &basicVal,
				reservation: &basicRes,
				statusCode:  http.StatusSeeOther,
				errorString: "can't insert room restriction",
				url:         "/some-url",
				redirectURL: "/",
			}
			doall(data)
		})

	})

	Context("PostSearchAvailability", func() {
		var basicVal url.Values

		BeforeEach(func() {
			basicVal = url.Values{}
			basicVal.Add("start", "2050-01-01")
			basicVal.Add("end", "2050-01-02")

			handler = h.PostSearchAvailability
			method = "POST"
		})

		It("normal", func() {
			mockDB.EXPECT().AvailabilityOfAllRooms(gomock.Any(), gomock.Any()).Return([]models.Room{
				{
					Name: "name 1",
					ID:   1,
				},
				{
					Name: "name 2",
					ID:   2,
				},
			}, nil)
			data := testData{
				val:        &basicVal,
				statusCode: http.StatusOK,
				url:        "/some-url",
			}
			doall(data)
		})

		It("bad form", func() {
			data := testData{
				statusCode:  http.StatusSeeOther,
				errorString: "bad form",
				url:         "/some-url",
				redirectURL: "/",
			}
			doall(data)
		})

		It("bad start", func() {
			basicVal.Set("start", "bad")
			data := testData{
				val:         &basicVal,
				statusCode:  http.StatusSeeOther,
				errorString: "bad start time",
				url:         "/some-url",
				redirectURL: "/",
			}
			doall(data)
		})

		It("bad end", func() {
			basicVal.Set("end", "bad")
			data := testData{
				val:         &basicVal,
				statusCode:  http.StatusSeeOther,
				errorString: "bad end time",
				url:         "/some-url",
				redirectURL: "/",
			}
			doall(data)
		})

		It("AvailabilityOfAllRooms error", func() {
			mockDB.EXPECT().AvailabilityOfAllRooms(gomock.Any(), gomock.Any()).Return(nil, errors.New("text"))
			data := testData{
				val:         &basicVal,
				reservation: nil,
				statusCode:  http.StatusSeeOther,
				errorString: "can't get rooms",
				url:         "/some-url",
				redirectURL: "/",
			}
			doall(data)
		})

		It("no available rooms", func() {
			mockDB.EXPECT().AvailabilityOfAllRooms(gomock.Any(), gomock.Any()).Return([]models.Room{}, nil)
			data := testData{
				val:         &basicVal,
				statusCode:  http.StatusSeeOther,
				errorString: "sorry, no available rooms on this dates >:(",
				url:         "/some-url",
				redirectURL: "/search-availability",
			}
			doall(data)
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
			method = "POST"
		})

		It("normal", func() {
			mockDB.EXPECT().LookForAvailabilityOfRoom(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
			data := testData{
				val:        &basicVal,
				statusCode: http.StatusOK,
				url:        "/some-url",
			}
			doall(data)
		})

		It("bad form", func() {
			data := testData{
				val:         nil,
				reservation: nil,
				statusCode:  http.StatusSeeOther,
				errorString: "bad form",
				url:         "/some-url",
				redirectURL: "/",
			}
			doall(data)
		})

		It("bad start", func() {
			basicVal.Set("start", "bad")
			data := testData{
				val:         &basicVal,
				reservation: nil,
				statusCode:  http.StatusSeeOther,
				errorString: "bad start time",
				url:         "/some-url",
				redirectURL: "/",
			}
			doall(data)
		})

		It("bad end", func() {
			basicVal.Set("end", "bad")
			data := testData{
				val:         &basicVal,
				statusCode:  http.StatusSeeOther,
				errorString: "bad end time",
				url:         "/some-url",
				redirectURL: "/",
			}
			doall(data)
		})

		It("bad room id", func() {
			basicVal.Set("room_id", "bad")
			data := testData{
				val:         &basicVal,
				statusCode:  http.StatusSeeOther,
				errorString: "bad room id",
				url:         "/some-url",
				redirectURL: "/",
			}
			doall(data)
		})

		It("AvailabilityOfAllRooms error", func() {
			mockDB.EXPECT().LookForAvailabilityOfRoom(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, errors.New("text"))
			data := testData{
				val:         &basicVal,
				statusCode:  http.StatusSeeOther,
				errorString: "problem with searching room",
				url:         "/some-url",
				redirectURL: "/",
			}
			doall(data)
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
					Name: "name",
					ID:   1,
				},
			}
			handler = h.ReservationSummary
			method = "GET"
		})

		It("test with right data", func() {
			data := testData{
				reservation: &basicRes,
				statusCode:  http.StatusOK,
				url:         "/some-url",
			}
			doall(data)
		})

		It("test with insufficient reservation", func() {
			data := testData{
				statusCode:  http.StatusSeeOther,
				errorString: "can't find reservation",
				url:         "/some-url",
				redirectURL: "/",
			}
			doall(data)
		})
	})

	Context("ChooseRoom", func() {
		var basicRes models.Reservation

		BeforeEach(func() {
			handler = h.ChooseRoom
			method = "GET"
		})

		It("test with right data", func() {
			data := testData{
				val:         nil,
				reservation: &basicRes,
				statusCode:  http.StatusSeeOther,
				errorString: "",
				url:         "/choose-room/1",
				redirectURL: "/make-reservation",
			}
			doall(data)
		})

		It("test with insufficient room id", func() {
			data := testData{
				val:         nil,
				reservation: &basicRes,
				statusCode:  http.StatusSeeOther,
				errorString: "can't find such room",
				url:         "/choose-room/f",
				redirectURL: "/",
			}
			doall(data)
		})

		It("test with insufficient reservation", func() {
			data := testData{
				statusCode:  http.StatusSeeOther,
				errorString: "can't find reservation",
				url:         "/choose-room/1",
				redirectURL: "/",
			}
			doall(data)
		})
	})

	Context("BookRoom", func() {

		BeforeEach(func() {
			handler = h.BookRoom
			method = "GET"
		})

		It("test with right data", func() {
			mockDB.EXPECT().GetRoomByID(gomock.Eq(1)).Return(&models.Room{
				ID:   1,
				Name: "room name",
			}, nil).AnyTimes()
			data := testData{
				statusCode:  http.StatusSeeOther,
				url:         "/book-room?s=2050-01-01&e=2050-01-02&id=1",
				redirectURL: "/make-reservation",
			}
			doall(data)
		})

		It("bad start", func() {
			data := testData{
				statusCode:  http.StatusSeeOther,
				errorString: "bad start time",
				url:         "/book-room?s=insufficient&e=2050-01-02&id=1",
				redirectURL: "/",
			}
			doall(data)
		})

		It("bad end", func() {
			data := testData{
				statusCode:  http.StatusSeeOther,
				errorString: "bad end time",
				url:         "/book-room?s=2050-01-01&e=insufficient&id=1",
				redirectURL: "/",
			}
			doall(data)
		})

		It("test with insufficient room id", func() {
			data := testData{
				statusCode:  http.StatusSeeOther,
				errorString: "can't find such room",
				url:         "/book-room?s=2050-01-01&e=2050-01-02&id=insufficient",
				redirectURL: "/",
			}
			doall(data)
		})

		It("test with error in getRoomByID ", func() {
			mockDB.EXPECT().GetRoomByID(gomock.Eq(2)).Return(nil, errors.New("no such room")).Times(1)
			data := testData{
				statusCode:  http.StatusSeeOther,
				errorString: "no such room",
				url:         "/book-room?s=2050-01-01&e=2050-01-01&id=2",
				redirectURL: "/",
			}
			doall(data)
		})
	})

	Context("AdminAllReservations", func() {

		BeforeEach(func() {
			handler = h.AdminAllReservations
			method = "GET"
		})

		It("test with right data", func() {
			mockDB.EXPECT().GetAllReservations().Return([]models.Reservation{}, nil).Times(1)
			data := testData{
				statusCode: http.StatusOK,
				url:        "/some-url",
			}
			doall(data)
		})

		It("bad db call", func() {
			mockDB.EXPECT().GetAllReservations().Return([]models.Reservation{}, errors.New("error text")).Times(1)
			data := testData{
				statusCode:  http.StatusSeeOther,
				errorString: "can't get all reservations",
				url:         "/some-url",
				redirectURL: "/admin/dashboard",
			}
			doall(data)
		})
	})

	Context("AdminNewReservations", func() {

		BeforeEach(func() {
			handler = h.AdminNewReservations
			method = "GET"
		})

		It("test with right data", func() {
			mockDB.EXPECT().GetNewReservations().Return([]models.Reservation{}, nil).Times(1)
			data := testData{
				statusCode: http.StatusOK,
				url:        "/some-url",
			}
			doall(data)
		})

		It("bad db call", func() {
			mockDB.EXPECT().GetNewReservations().Return([]models.Reservation{}, errors.New("error text")).Times(1)
			data := testData{
				statusCode:  http.StatusSeeOther,
				errorString: "can't get new reservations",
				url:         "/some-url",
				redirectURL: "/admin/dashboard",
			}
			doall(data)
		})
	})

	Context("AdminReservationCalendar", func() {

		BeforeEach(func() {
			handler = h.AdminReservationCalendar
			method = "GET"
		})

		It("test with right data with empty rooms", func() {
			mockDB.EXPECT().GetAllRooms().Return([]models.Room{}, nil).Times(1)
			data := testData{
				statusCode: http.StatusOK,
				url:        "/reservation-calendar",
			}
			doall(data)
		})

		It("test with right data with 1 room", func() {
			mockDB.EXPECT().GetAllRooms().Return([]models.Room{
				{
					ID:   1,
					Name: "Room #1",
				},
			}, nil).Times(1)
			mockDB.EXPECT().GetRoomRestrictionsByRoomIdWithinDates(gomock.Eq(1), gomock.Any(), gomock.Any()).
				Return([]models.RoomRestriction{}, nil).Times(1)
			data := testData{
				statusCode: http.StatusOK,
				url:        "/reservation-calendar",
			}
			doall(data)
		})

		It("test with insufficient y", func() {
			data := testData{
				statusCode:  http.StatusSeeOther,
				url:         "/reservation-calendar?y=insufficient&m=1",
				errorString: "can't get year from url: /reservation-calendar?y=insufficient&m=1",
				redirectURL: "/admin/dashboard",
			}
			doall(data)
		})

		It("test with insufficient m", func() {
			data := testData{
				statusCode:  http.StatusSeeOther,
				url:         "/reservation-calendar?y=2023&m=insufficient",
				errorString: "can't get month from url: /reservation-calendar?y=2023&m=insufficient",
				redirectURL: "/admin/dashboard",
			}
			doall(data)
		})

		It("test with error in getAllRooms", func() {
			mockDB.EXPECT().GetAllRooms().Return([]models.Room{}, errors.New("error test")).Times(1)
			data := testData{
				statusCode:  http.StatusSeeOther,
				url:         "/reservation-calendar",
				errorString: "can't get rooms",
				redirectURL: "/admin/dashboard",
			}
			doall(data)
		})

		It("test with error in GetRoomRestrictionsByRoomIdWithinDates", func() {
			mockDB.EXPECT().GetAllRooms().Return([]models.Room{
				{
					ID:   1,
					Name: "Room #1",
				},
			}, nil).Times(1)
			mockDB.EXPECT().GetRoomRestrictionsByRoomIdWithinDates(gomock.Eq(1), gomock.Any(), gomock.Any()).
				Return([]models.RoomRestriction{}, errors.New("error test")).Times(1)
			data := testData{
				statusCode:  http.StatusSeeOther,
				url:         "/reservation-calendar",
				errorString: "can't get room restrictions",
				redirectURL: "/admin/dashboard",
			}
			doall(data)
		})

	})

	Context("AdminSingleReservation", func() {

		BeforeEach(func() {
			handler = h.AdminSingleReservation
			method = "GET"
		})

		It("test with right data", func() {
			mockDB.EXPECT().GetReservationByID(gomock.Eq(1)).Return(&models.Reservation{
				ID:          1,
				FirstName:   "First",
				LastName:    "Last",
				Email:       "email@mail.com",
				Phone:       "12345",
				RoomID:      1,
				IsProcessed: 0,
			}, nil).Times(1)
			data := testData{
				statusCode: http.StatusOK,
				url:        "/admin/reservations/new/1/show",
			}
			doall(data)
		})

		It("test with wrong url", func() {
			data := testData{
				statusCode:  http.StatusSeeOther,
				url:         "/admin/reservations",
				errorString: "incorrect request url",
				redirectURL: "/admin/dashboard",
			}
			doall(data)
		})

		It("test with wrong id", func() {
			data := testData{
				statusCode:  http.StatusSeeOther,
				errorString: "wrong id",
				url:         "/admin/reservations/new/q/show",
				redirectURL: "/admin/new-reservations",
			}
			doall(data)
		})

		It("test with insufficient reservation", func() {
			mockDB.EXPECT().GetReservationByID(gomock.Eq(1)).
				Return(&models.Reservation{}, errors.New("error text")).Times(1)
			data := testData{
				statusCode:  http.StatusSeeOther,
				errorString: "can't find reservation",
				url:         "/admin/reservations/new/1/show",
				redirectURL: "/admin/new-reservations",
			}
			doall(data)
		})

	})

	Context("AdminPostSingleReservation", func() {
		var basicVal url.Values
		BeforeEach(func() {
			basicVal = url.Values{}
			basicVal.Add("first_name", "First")
			basicVal.Add("last_name", "Last")
			basicVal.Add("email", "e@e.com")
			basicVal.Add("phone", "12345")
			handler = h.AdminPostSingleReservation
			method = "POST"
		})

		It("test with right data to new", func() {
			mockDB.EXPECT().UpdateReservation(gomock.Any()).Return(nil).Times(1)
			data := testData{
				val:         &basicVal,
				statusCode:  http.StatusSeeOther,
				url:         "/admin/reservations/new/1/show",
				redirectURL: "/admin/new-reservations",
			}
			doall(data)
		})

		It("test with right data to all", func() {
			mockDB.EXPECT().UpdateReservation(gomock.Any()).Return(nil).Times(1)
			data := testData{
				val:         &basicVal,
				statusCode:  http.StatusSeeOther,
				url:         "/admin/reservations/all/1/show",
				redirectURL: "/admin/all-reservations",
			}
			doall(data)
		})

		It("test with right data to calendar", func() {
			basicVal.Add("year", "2023")
			basicVal.Add("month", "11")
			mockDB.EXPECT().UpdateReservation(gomock.Any()).Return(nil).Times(1)
			data := testData{
				val:         &basicVal,
				statusCode:  http.StatusSeeOther,
				url:         "/admin/reservations/cal/1/show?y=2023&m=11",
				redirectURL: "/admin/reservation-calendar?y=2023&m=11",
			}
			doall(data)
		})

		It("test with wrong url", func() {
			data := testData{
				val:         &basicVal,
				statusCode:  http.StatusSeeOther,
				url:         "/admin/reservations",
				errorString: "incorrect request url",
				redirectURL: "/admin/dashboard",
			}
			doall(data)
		})

		It("test with wrong id", func() {
			data := testData{
				val:         &basicVal,
				statusCode:  http.StatusSeeOther,
				errorString: "wrong id",
				url:         "/admin/reservations/new/q/show",
				redirectURL: "/admin/dashboard",
			}
			doall(data)
		})

		It("test with bad form", func() {
			data := testData{
				statusCode:  http.StatusSeeOther,
				errorString: "bad form",
				url:         "/admin/reservations/new/1/show",
				redirectURL: "/admin/dashboard",
			}
			doall(data)
		})

		It("test with error in updateReservation", func() {
			basicVal.Add("year", "2023")
			basicVal.Add("month", "11")
			mockDB.EXPECT().UpdateReservation(gomock.Any()).Return(errors.New("error text")).Times(1)
			data := testData{
				val:         &basicVal,
				statusCode:  http.StatusSeeOther,
				errorString: "can't update reservation",
				url:         "/admin/reservations/cal/1/show?y=2023&m=11",
				redirectURL: "/admin/reservation-calendar?y=2023&m=11",
			}
			doall(data)
		})

	})

	Context("AdminProcessReservation", func() {
		BeforeEach(func() {
			handler = h.AdminProcessReservation
			method = "GET"
		})

		It("test with right data to new", func() {
			mockDB.EXPECT().UpdateReservationProcessed(gomock.Any(), gomock.Eq(1)).Return(nil).Times(1)
			data := testData{
				statusCode:  http.StatusSeeOther,
				url:         "/admin/process-reservation/new/1/do",
				redirectURL: "/admin/new-reservations",
			}
			doall(data)
		})

		It("test with right data to all", func() {
			mockDB.EXPECT().UpdateReservationProcessed(gomock.Any(), gomock.Eq(1)).Return(nil).Times(1)
			data := testData{
				statusCode:  http.StatusSeeOther,
				url:         "/admin/process-reservation/all/1/do",
				redirectURL: "/admin/all-reservations",
			}
			doall(data)
		})

		It("test with right data to calendar", func() {
			mockDB.EXPECT().UpdateReservationProcessed(gomock.Any(), gomock.Eq(1)).Return(nil).Times(1)
			data := testData{
				statusCode:  http.StatusSeeOther,
				url:         "/admin/process-reservation/cal/1/do?y=2023&m=11",
				redirectURL: "/admin/reservation-calendar?y=2023&m=11",
			}
			doall(data)
		})

		It("test with wrong url", func() {
			data := testData{
				statusCode:  http.StatusSeeOther,
				url:         "/admin/process-reservation",
				errorString: "incorrect request url",
				redirectURL: "/admin/dashboard",
			}
			doall(data)
		})

		It("test with wrong id", func() {
			data := testData{
				statusCode:  http.StatusSeeOther,
				errorString: "wrong id",
				url:         "/admin/process-reservation/new/q/do",
				redirectURL: "/admin/new-reservations",
			}
			doall(data)
		})

		It("test with error in UpdateReservationProcessed", func() {

			mockDB.EXPECT().UpdateReservationProcessed(gomock.Eq(1), gomock.Eq(1)).Return(errors.New("error text")).Times(1)
			data := testData{
				statusCode:  http.StatusSeeOther,
				errorString: "can't update reservation",
				url:         "/admin/process-reservation/new/1/do?y=2023&m=11",
				redirectURL: "/admin/reservation-calendar?y=2023&m=11",
			}
			doall(data)
		})

	})

	Context("AdminDeleteReservation", func() {
		BeforeEach(func() {
			handler = h.AdminDeleteReservation
			method = "GET"
		})

		It("test with right data to new", func() {
			mockDB.EXPECT().DeleteReservationByID(gomock.Eq(1)).Return(nil).Times(1)
			data := testData{
				statusCode:  http.StatusSeeOther,
				url:         "/admin/delete-reservation/new/1/do",
				redirectURL: "/admin/new-reservations",
			}
			doall(data)
		})

		It("test with right data to all", func() {
			mockDB.EXPECT().DeleteReservationByID(gomock.Eq(1)).Return(nil).Times(1)
			data := testData{
				statusCode:  http.StatusSeeOther,
				url:         "/admin/delete-reservation/all/1/do",
				redirectURL: "/admin/all-reservations",
			}
			doall(data)
		})

		It("test with right data to calendar", func() {
			mockDB.EXPECT().DeleteReservationByID(gomock.Eq(1)).Return(nil).Times(1)
			data := testData{
				statusCode:  http.StatusSeeOther,
				url:         "/admin/delete-reservation/cal/1/do?y=2023&m=11",
				redirectURL: "/admin/reservation-calendar?y=2023&m=11",
			}
			doall(data)
		})

		It("test with wrong url", func() {
			data := testData{
				statusCode:  http.StatusSeeOther,
				url:         "/admin/delete-reservation",
				errorString: "incorrect request url",
				redirectURL: "/admin/dashboard",
			}
			doall(data)
		})

		It("test with wrong id", func() {
			data := testData{
				statusCode:  http.StatusSeeOther,
				errorString: "wrong id",
				url:         "/admin/delete-reservation/new/q/do",
				redirectURL: "/admin/new-reservations",
			}
			doall(data)
		})

		It("test with error in UpdateReservationProcessed", func() {

			mockDB.EXPECT().DeleteReservationByID(gomock.Eq(1)).Return(errors.New("error text")).Times(1)
			data := testData{
				statusCode:  http.StatusSeeOther,
				errorString: "can't delete reservation",
				url:         "/admin/delete-reservation/new/1/do?y=2023&m=11",
				redirectURL: "/admin/reservation-calendar?y=2023&m=11",
			}
			doall(data)
		})

	})

	Context("AdminPostReservationCalendar", func() {
		var basicVal url.Values
		BeforeEach(func() {
			basicVal = url.Values{}
			basicVal.Add("y", "2023")
			basicVal.Add("m", "2")
			handler = h.AdminPostReservationCalendar
			method = "POST"
		})

		It("test with right data", func() {
			mockDB.EXPECT().GetAllRooms().Return([]models.Room{
				{
					ID:   1,
					Name: "Room",
				},
			}, nil).Times(1)
			mockDB.EXPECT().AddSingleDayRoomRestriction(gomock.Eq(1), gomock.Eq(2), gomock.Any()).
				Return(0, nil).Times(1)
			mockDB.EXPECT().DeleteRoomRestrictionByID(gomock.Eq(1)).Return(nil).Times(1)
			basicVal.Add("add_block_1_2023-02-02", "1")
			//basicVal.Add("remove_block_1_2023-02-05", "1")
			data := testData{
				val: &basicVal,
				dataForSession: map[string]interface{}{
					"block_map_1": map[string]int{
						"2023-02-05": 1,
					},
				},
				statusCode:  http.StatusSeeOther,
				url:         "/admin/reservation-calendar?y=2023&m=2",
				redirectURL: "/admin/reservation-calendar?y=2023&m=2",
			}
			doall(data)
		})

		It("test with bad form", func() {
			data := testData{
				statusCode:  http.StatusSeeOther,
				errorString: "bad form",
				url:         "/admin/reservation-calendar?y=2023&m=2",
				redirectURL: "/admin/reservation-calendar",
			}
			doall(data)
		})

		It("test with insufficient year", func() {
			basicVal.Set("y", "insufficient")
			data := testData{
				val:         &basicVal,
				statusCode:  http.StatusSeeOther,
				errorString: "wrong year",
				url:         "/admin/reservation-calendar?y=2023&m=2",
				redirectURL: "/admin/reservation-calendar",
			}
			doall(data)
		})

		It("test with insufficient month", func() {
			basicVal.Set("m", "insufficient")
			data := testData{
				val:         &basicVal,
				statusCode:  http.StatusSeeOther,
				errorString: "wrong month",
				url:         "/admin/reservation-calendar?y=2023&m=2",
				redirectURL: "/admin/reservation-calendar",
			}
			doall(data)
		})

		It("test with error in rooms", func() {
			mockDB.EXPECT().GetAllRooms().Return([]models.Room{}, errors.New("error text")).Times(1)
			data := testData{
				val:         &basicVal,
				statusCode:  http.StatusSeeOther,
				errorString: "can't get rooms",
				url:         "/admin/reservation-calendar?y=2023&m=2",
				redirectURL: "/admin/reservation-calendar?y=2023&m=2",
			}
			doall(data)
		})

		It("test without block_map_1 in session", func() {
			mockDB.EXPECT().GetAllRooms().Return([]models.Room{
				{
					ID:   1,
					Name: "Room",
				},
			}, nil).Times(1)

			data := testData{
				val:         &basicVal,
				statusCode:  http.StatusSeeOther,
				errorString: "can't get block map for room: Room",
				url:         "/admin/reservation-calendar?y=2023&m=2",
				redirectURL: "/admin/reservation-calendar?y=2023&m=2",
			}
			doall(data)
		})

		It("test with error in delete", func() {
			mockDB.EXPECT().GetAllRooms().Return([]models.Room{
				{
					ID:   1,
					Name: "Room",
				},
			}, nil).Times(1)
			mockDB.EXPECT().DeleteRoomRestrictionByID(gomock.Eq(1)).Return(errors.New("error text")).Times(1)
			data := testData{
				val: &basicVal,
				dataForSession: map[string]interface{}{
					"block_map_1": map[string]int{
						"2023-02-05": 1,
					},
				},
				statusCode:  http.StatusSeeOther,
				errorString: "can't delete restriction 1",
				url:         "/admin/reservation-calendar?y=2023&m=2",
				redirectURL: "/admin/reservation-calendar?y=2023&m=2",
			}
			doall(data)
		})

		It("test with insufficient room id in add block", func() {
			mockDB.EXPECT().GetAllRooms().Return([]models.Room{
				{
					ID:   1,
					Name: "Room",
				},
			}, nil).Times(1)
			basicVal.Add("add_block_insufficient_2023-02-02", "1")
			data := testData{
				val: &basicVal,
				dataForSession: map[string]interface{}{
					"block_map_1": map[string]int{},
				},
				errorString: "can't parse room id add_block_insufficient_2023-02-02",
				statusCode:  http.StatusSeeOther,
				url:         "/admin/reservation-calendar?y=2023&m=2",
				redirectURL: "/admin/reservation-calendar?y=2023&m=2",
			}
			doall(data)
		})

		It("test with insufficient date in add block", func() {
			mockDB.EXPECT().GetAllRooms().Return([]models.Room{
				{
					ID:   1,
					Name: "Room",
				},
			}, nil).Times(1)
			basicVal.Add("add_block_1_insufficient", "1")
			data := testData{
				val: &basicVal,
				dataForSession: map[string]interface{}{
					"block_map_1": map[string]int{},
				},
				errorString: "can't parse date add_block_1_insufficient",
				statusCode:  http.StatusSeeOther,
				url:         "/admin/reservation-calendar?y=2023&m=2",
				redirectURL: "/admin/reservation-calendar?y=2023&m=2",
			}
			doall(data)
		})

		It("test with error in AddSingleDayRoomRestriction", func() {
			mockDB.EXPECT().GetAllRooms().Return([]models.Room{
				{
					ID:   1,
					Name: "Room",
				},
			}, nil).Times(1)
			mockDB.EXPECT().AddSingleDayRoomRestriction(gomock.Eq(1), gomock.Eq(2), gomock.Any()).
				Return(0, errors.New("error text")).Times(1)
			basicVal.Add("add_block_1_2023-02-02", "1")
			data := testData{
				val: &basicVal,
				dataForSession: map[string]interface{}{
					"block_map_1": map[string]int{},
				},
				statusCode:  http.StatusSeeOther,
				errorString: "can't save this restriction add_block_1_2023-02-02",
				url:         "/admin/reservation-calendar?y=2023&m=2",
				redirectURL: "/admin/reservation-calendar?y=2023&m=2",
			}
			doall(data)
		})

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
		r.Post("/reservation-calendar", http.HandlerFunc(handler.AdminPostReservationCalendar))

		r.Get("/reservations/{src}/{id}/show", http.HandlerFunc(handler.AdminSingleReservation))
		r.Post("/reservations/{src}/{id}/show", http.HandlerFunc(handler.AdminPostSingleReservation))

		r.Get("/process-reservation/{src}/{id}/do", http.HandlerFunc(handler.AdminProcessReservation))
		r.Get("/delete-reservation/{src}/{id}/do", http.HandlerFunc(handler.AdminDeleteReservation))
	})
	return mux
}

func SessionLoad(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
}

func getCtx(req *http.Request, app *config.AppConfig) (context.Context, error) {
	return app.Session.Load(req.Context(), req.Header.Get("X-Session"))
}
