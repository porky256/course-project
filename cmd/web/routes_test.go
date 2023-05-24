package main

import (
	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/porky256/course-project/internal/config"
	"github.com/porky256/course-project/internal/driver"
	"github.com/porky256/course-project/internal/handlers"
	"github.com/porky256/course-project/internal/render"
)

const (
	get  = "GET"
	post = "POST"
)

var _ = Describe("Routes private functions", func() {
	var app config.AppConfig
	var r *render.Render
	var h *handlers.Handlers
	Context("routes", func() {
		BeforeEach(func() {
			app = config.AppConfig{}
			r = render.NewRender(&app)
			h = handlers.NewHandlers(&app, r, &driver.DB{})
		})

		It("Check if routes added correctly", func() {
			mux := routes(&app, h).(*chi.Mux)
			routes := mux.Routes()

			Expect(routeExists(get, "/", routes)).To(Equal(true))

			Expect(routeExists(get, "/about", routes)).To(Equal(true))

			Expect(routeExists(get, "/generals-quarters", routes)).To(Equal(true))

			Expect(routeExists(get, "/majors-suite", routes)).To(Equal(true))

			Expect(routeExists(get, "/make-reservation", routes)).To(Equal(true))
			Expect(routeExists(post, "/make-reservation", routes)).To(Equal(true))

			Expect(routeExists(get, "/search-availability", routes)).To(Equal(true))
			Expect(routeExists(post, "/search-availability", routes)).To(Equal(true))
			Expect(routeExists(post, "/search-availability-json", routes)).To(Equal(true))

			Expect(routeExists(get, "/reservation-summary", routes)).To(Equal(true))

			Expect(routeExists(get, "/contact", routes)).To(Equal(true))
		})
	})
})

func routeExists(handler string, pattern string, routes []chi.Route) bool {
	for _, route := range routes {
		if route.Pattern == pattern && route.Handlers[handler] != nil {
			return true
		}
	}
	return false
}
