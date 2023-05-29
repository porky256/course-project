package render

import (
	"github.com/alexedwards/scs/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/porky256/course-project/internal/config"
	"github.com/porky256/course-project/internal/helpers"
	"github.com/porky256/course-project/internal/models"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const testString = "test string"

var _ = Describe("Render", func() {
	var session *scs.SessionManager
	var request *http.Request
	var render *Render
	var app config.AppConfig
	BeforeEach(func() {
		session = scs.New()
		r, err := http.NewRequest("GET", "/some-url", nil)
		Expect(err).ToNot(HaveOccurred())
		ctx := r.Context()
		ctx, err = session.Load(ctx, r.Header.Get("X-Session"))
		Expect(err).ToNot(HaveOccurred())
		request = r.WithContext(ctx)
		app = config.AppConfig{
			Session:  session,
			RootPath: "./../..",
		}
		app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
		app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
		helpers.NewHelpers(&app)
		render = NewRender(&app)
	})

	Context("addDefaultData", func() {
		It("without any additional info", func() {
			result := render.addDefaultData(&models.TemplateData{}, request)
			Expect(len(result.StringMap)).To(Equal(0))
			Expect(len(result.IntMap)).To(Equal(0))
			Expect(len(result.Float32Map)).To(Equal(0))
			Expect(len(result.Data)).To(Equal(0))
			Expect(result.CSRFToken).To(Equal(""))
			Expect(result.Flash).To(Equal(""))
			Expect(result.Warning).To(Equal(""))
			Expect(result.Error).To(Equal(""))
			Expect(result.Form).To(BeNil())
		})

		It("with flash, error and warning", func() {
			session.Put(request.Context(), "flash", testString)
			session.Put(request.Context(), "warning", testString)
			session.Put(request.Context(), "error", testString)
			result := render.addDefaultData(&models.TemplateData{}, request)
			Expect(len(result.StringMap)).To(Equal(0))
			Expect(len(result.IntMap)).To(Equal(0))
			Expect(len(result.Float32Map)).To(Equal(0))
			Expect(len(result.Data)).To(Equal(0))
			Expect(result.CSRFToken).To(Equal(""))
			Expect(result.Flash).To(Equal(testString))
			Expect(result.Warning).To(Equal(testString))
			Expect(result.Error).To(Equal(testString))
			Expect(result.Form).To(BeNil())
		})
	})

	Context("CreateTemplateCacheMap", func() {
		It("Check if function works correctly", func() {
			cache, err := CreateTemplateCacheMap(render.app)
			Expect(err).ToNot(HaveOccurred())
			files, err := filepath.Glob(app.RootPath + "/static/templates/*page.tmpl")
			var cachedPages []string
			for k := range cache {
				cachedPages = append(cachedPages, k)
			}
			var templateNames []string
			for _, f := range files {
				templateNames = append(templateNames, filepath.Base(f))
			}

			By("check actual and cached pages bijection")
			Expect(cachedPages).To(ContainElements(templateNames))
			Expect(templateNames).To(ContainElements(cachedPages))
		})
	})

	Context("Template", func() {
		var ww mockWriter
		BeforeEach(func() {
			tc, err := CreateTemplateCacheMap(&app)
			Expect(err).ToNot(HaveOccurred())
			app.TemplateCache = tc
			ww = mockWriter{}
		})
		It("test with real page", func() {
			err := render.Template(&ww, request, "home.page.tmpl", &models.TemplateData{})
			Expect(err).ToNot(HaveOccurred())
		})

		It("test with non-existing page", func() {
			err := render.Template(&ww, request, "non-existent.page.tmpl", &models.TemplateData{})
			Expect(err).To(HaveOccurred())
		})

	})
})
