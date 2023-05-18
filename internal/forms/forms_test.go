package forms_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/porky256/course-project/internal/forms"
	"net/url"
)

var _ = Describe("Forms", func() {
	var testForm *forms.Form
	BeforeEach(func() {
		testForm = forms.New(url.Values{})
	})

	Context("New", func() {
		It("check function creates new form with preset data", func() {
			values := url.Values{}
			values["test"] = []string{"test", "test"}
			newForm := forms.New(values)
			Expect(len(newForm.Errors)).To(Equal(0))
			Expect(newForm.Values).To(Equal(values))
		})
	})

	Context("Valid", func() {
		var testForm *forms.Form
		BeforeEach(func() {
			testForm = forms.New(url.Values{})
		})

		It("form is valid", func() {
			Expect(testForm.Valid()).To(Equal(true))
		})

		It("form is invalid", func() {
			testForm.Errors.Add("test error", "test error")
			Expect(testForm.Valid()).To(Equal(false))
		})

	})

	Context("Required", func() {
		BeforeEach(func() {
			testForm.Values["test1"] = []string{"test1"}
			testForm.Values["test2"] = []string{"test2"}
		})

		It("form do have required fields", func() {
			testForm.Required("test1", "test2")
			Expect(testForm.Valid()).To(Equal(true))
		})

		It("form don't have required fields", func() {
			testForm.Required("test1", "test2", "test3", "test4")
			Expect(testForm.Valid()).To(Equal(false))
			Expect(testForm.Errors.Get("test1")).To(Equal(""))
			Expect(testForm.Errors.Get("test2")).To(Equal(""))
			Expect(testForm.Errors.Get("test3")).To(Equal("This field is required"))
			Expect(testForm.Errors.Get("test4")).To(Equal("This field is required"))
		})
	})

	Context("Has", func() {
		BeforeEach(func() {
			testForm.Values["test"] = []string{"test"}
		})

		It("form do have required field", func() {
			Expect(testForm.Has("test")).To(Equal(true))
		})
		It("form don't have required field", func() {
			Expect(testForm.Has("test2")).To(Equal(false))
		})
	})

	Context("MinLength", func() {
		BeforeEach(func() {
			testForm.Values["test"] = []string{"test"}
		})

		It("field exists and longer", func() {
			Expect(testForm.MinLength("test", 3)).To(Equal(true))
			Expect(testForm.Valid()).To(Equal(true))
		})

		It("field exists and shorter", func() {
			Expect(testForm.MinLength("test", 5)).To(Equal(false))
			Expect(testForm.Valid()).To(Equal(false))
			Expect(testForm.Errors.Get("test")).To(Equal("This field must be at least 5 symbols long"))
		})

		It("field don't exists", func() {
			Expect(testForm.MinLength("test2", 5)).To(Equal(false))
			Expect(testForm.Valid()).To(Equal(false))
			Expect(testForm.Errors.Get("test2")).To(Equal("This field must be at least 5 symbols long"))
		})
	})

	Context("Has", func() {

		It("field is correct email", func() {
			testForm.Values["email"] = []string{"email@here.com"}
			Expect(testForm.IsEmail("email")).To(Equal(true))
			Expect(testForm.Valid()).To(Equal(true))
		})

		It("field isn't correct email", func() {
			testForm.Values["email"] = []string{"email@here"}
			Expect(testForm.IsEmail("email")).To(Equal(false))
			Expect(testForm.Valid()).To(Equal(false))
			Expect(testForm.Errors.Get("email")).To(Equal("This field is not a valid email"))
		})
	})

	Context("Errors order", func() {
		It("field required and shorter", func() {
			testForm.Required("test")
			testForm.MinLength("test", 3)
			Expect(testForm.Valid()).To(Equal(false))
			Expect(testForm.Errors.Get("test")).To(Equal("This field is required"))
		})

		It("field shorter and required", func() {
			testForm.MinLength("test", 3)
			testForm.Required("test")
			Expect(testForm.Valid()).To(Equal(false))
			Expect(testForm.Errors.Get("test")).To(Equal("This field must be at least 3 symbols long"))
		})
	})
})
