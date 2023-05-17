package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main private functions", func() {
	Context("run()", func() {
		It("test if function completes correctly", func() {
			err := run()
			Expect(err).To(BeNil())
		})
	})
})
