package forms

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Errors", func() {
	var err errors
	BeforeEach(func() {
		err = errors{}
	})

	Context("Add", func() {
		It("add test value", func() {
			err.Add("test", "test message")
			Expect(len(err)).To(Equal(1))
			Expect(err["test"]).To(Equal([]string{"test message"}))
		})
	})

	Context("Get", func() {
		It("get real value", func() {
			err.Add("test", "test message")
			Expect(err.Get("test")).To(Equal("test message"))
		})

		It("try to get non-existing value", func() {
			Expect(err.Get("test")).To(Equal(""))
		})

		It("multiple errors, should return first", func() {
			err.Add("test", "test message number one")
			err.Add("test", "test message number two")
			Expect(err.Get("test")).To(Equal("test message number one"))
		})
	})
})
