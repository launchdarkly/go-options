//go:build testing
// +build testing

package test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("able to add a build tag", func() {
	It("it should have a build tag", func() {
		Ω(fmt.Sprintf("%v", BuildOptMyInt(1))).Should(Equal("BuildOptMyInt: 1"))
	})
})
