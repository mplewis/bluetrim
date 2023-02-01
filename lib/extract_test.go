package lib_test

import (
	. "github.com/mplewis/bluetrim/lib"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("extract", func() {
	Describe("TsToFrame", func() {
		It("converts timestamps to frames", func() {
			Expect(TsToFrame("00:00:00.000", 60)).To(Equal(int64(1)))
			Expect(TsToFrame("00:00:01.000", 60)).To(Equal(int64(61)))
			Expect(TsToFrame("00:00:15.000", 60)).To(Equal(int64(901)))
			Expect(TsToFrame("00:01:00.000", 60)).To(Equal(int64(3601)))
			Expect(TsToFrame("01:00:00.000", 60)).To(Equal(int64(216001)))
			Expect(TsToFrame("01:00:00.000", 29.97)).To(Equal(int64(107893)))
		})
	})
})
