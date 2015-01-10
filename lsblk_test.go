package main_test

import (
	. "github.com/JamesClonk/cryptmount"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lsblk", func() {

	var (
		LSBLK_CMD = `cat fixtures/lsblk.output`
	)

	BeforeEach(func() {
		LSBLK_CMD = `cat fixtures/lsblk.output`
	})

	Describe("Calling Lsdsk()", func() {
		Context("With all possible device combinations", func() {
			It("should return the expected data", func() {
				Expect(Lsdsk()).To(Equal(nil))
			})
		})
	})

})
