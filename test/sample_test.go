package test

import (
	"errors"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOptions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Options Suite")
}

// Make sure we don't try to generate this
type OptionMyRenamedInt int // nolint:structcheck,unused // just exists so we would conflict with it

type OptionSetMyInt123 struct{}

func (o OptionSetMyInt123) apply(c *config) error {
	c.myInt = int(123)
	return nil
}

type OptionMakeError struct{}

func (o OptionMakeError) apply(c *config) error {
	return errors.New("bad news")
}

var _ = Describe("Generating options", func() {
	cfg := config{}

	It("generates options to set config value", func() {
		err := applyConfigOptions(&cfg,
			OptionMyInt(123),
			OptionMyFloat(4.56),
			OptionMyString("my-string"),
		)
		Ω(err).NotTo(HaveOccurred())
		Ω(cfg.myInt).To(Equal(123))
		Ω(cfg.myFloat).To(Equal(4.56))
		Ω(cfg.myString).To(Equal("my-string"))
	})

	It("generates an new function create a config", func() {
		cfg, err := newConfig(OptionMyInt(123))
		Ω(err).NotTo(HaveOccurred())
		Ω(cfg.myInt).To(Equal(123))
	})

	It("sets default values", func() {
		err := applyConfigOptions(&cfg)
		Ω(err).NotTo(HaveOccurred())
		Ω(cfg.myIntWithDefault).To(Equal(1))
		Ω(cfg.myStringWithDefault).To(Equal("default string"))
		Ω(cfg.myFloatWithDefault).To(Equal(1.23))
	})

	It("defines constants for default values", func() {
		err := applyConfigOptions(&cfg, OptionMakeError{})
		Ω(err).To(MatchError("bad news"))
	})

	Describe("custom options", func() {
		It("can be extended with custom options", func() {
			err := applyConfigOptions(&cfg, OptionSetMyInt123{})
			Ω(err).NotTo(HaveOccurred())
			Ω(cfg.myInt).To(Equal(123))
		})

		It("returns error from custom options", func() {
			err := applyConfigOptions(&cfg, OptionMakeError{})
			Ω(err).To(MatchError("bad news"))
		})
	})
})

var _ = Describe("Customizing the apply function name", func() {
	cfg := configWithDifferentApply{}

	It("uses the provided function name", func() {
		err := applyDifferent(&cfg)
		Ω(err).NotTo(HaveOccurred())
	})
})
