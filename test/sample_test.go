package test

import (
	"errors"
	"net/url"
	"testing"
	"time"

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
		myInt := 456
		err := applyConfigOptions(&cfg,
			OptionMyInt(123),
			OptionMyFloat(4.56),
			OptionMyString("my-string"),
			OptionMyIntPointer(&myInt),
			OptionMyInterface(789),
			OptionMyFunc(func() int { return 0 }),
		)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(cfg.myInt).Should(Equal(123))
		Ω(cfg.myFloat).Should(Equal(4.56))
		Ω(cfg.myString).Should(Equal("my-string"))
		Ω(cfg.myIntPointer).Should(Equal(&myInt))
		Ω(cfg.myInterface).Should(Equal(789))
	})

	It("generates an new function create a config", func() {
		cfg, err := newConfig(OptionMyInt(123))
		Ω(err).ShouldNot(HaveOccurred())
		Ω(cfg.myInt).Should(Equal(123))
	})

	It("sets default values", func() {
		err := applyConfigOptions(&cfg)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(cfg.myIntWithDefault).Should(Equal(1))
		Ω(cfg.myStringWithDefault).Should(Equal("default string"))
		Ω(cfg.myFloatWithDefault).Should(Equal(1.23))
	})

	It("defines constants for default values", func() {
		err := applyConfigOptions(&cfg, OptionMakeError{})
		Ω(err).Should(MatchError("bad news"))
	})

	Describe("custom options", func() {
		It("can be extended with custom options", func() {
			err := applyConfigOptions(&cfg, OptionSetMyInt123{})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myInt).Should(Equal(123))
		})

		It("returns error from custom options", func() {
			err := applyConfigOptions(&cfg, OptionMakeError{})
			Ω(err).Should(MatchError("bad news"))
		})
	})

	Describe("imports", func() {
		It("works with imported types", func() {
			err := applyConfigOptions(&cfg, OptionMyDuration(time.Second))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myDuration).Should(Equal(time.Second))
		})

		It("works with aliased imports", func() {
			err := applyConfigOptions(&cfg, OptionMyDuration2(time.Second))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myDuration).Should(Equal(time.Second))
		})

		It("works with nested packages", func() {
			myURL, err := url.Parse("http://example.com")
			Ω(err).ShouldNot(HaveOccurred())
			err = applyConfigOptions(&cfg, OptionMyURL(*myURL))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myURL).Should(Equal(*myURL))
		})

		It("works with nested structs", func() {
			err := applyConfigOptions(&cfg, OptionMyStruct(1, 2))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myStruct.a).Should(Equal(1))
			Ω(cfg.myStruct.b).Should(Equal(2))
		})

		It("works with nested optional structs", func() {
			err := applyConfigOptions(&cfg)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myOptionalStruct).Should(BeNil())

			err = applyConfigOptions(&cfg, OptionMyOptionalStruct(1, 2))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myOptionalStruct).ShouldNot(BeNil())
			Ω(cfg.myOptionalStruct.a).Should(Equal(1))
			Ω(cfg.myOptionalStruct.b).Should(Equal(2))
		})

		Describe("variadic slices", func() {
			It("creates a variadic constructor", func() {
				err := applyConfigOptions(&cfg, OptionMySlice(1, 2))
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cfg.mySlice).Should(ConsistOf(1, 2))
			})

			It("allows them to be optional", func() {
				err := applyConfigOptions(&cfg)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cfg.myOptionalSlice).Should(BeNil())

				err = applyConfigOptions(&cfg, OptionMyOptionalSlice(1, 2))
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cfg.myOptionalSlice).ShouldNot(BeNil())
				Ω(*cfg.myOptionalSlice).Should(ConsistOf(1, 2))
			})

			It("allows them to be renamed", func() {
				err := applyConfigOptions(&cfg, OptionYourSlice(1, 2))
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cfg.myRenamedSlice).ShouldNot(BeNil())
				Ω(cfg.myRenamedSlice).Should(ConsistOf(1, 2))
			})
		})
	})

})

var _ = Describe("Customizing the apply function name", func() {
	cfg := configWithDifferentApply{}

	It("uses the provided function name", func() {
		err := applyDifferent(&cfg)
		Ω(err).ShouldNot(HaveOccurred())
	})
})

var _ = Describe("Customizing the option prefix", func() {
	It("creates options with the custom prefix", func() {
		_, err := newConfigWithDifferentPrefix(OptMyFloat(1.23))
		Ω(err).ShouldNot(HaveOccurred())
	})
})

var _ = Describe("Customizing the option suffix", func() {
	It("creates options with the custom prefix", func() {
		_, err := newConfigWithSuffix(MyFloatOption(1.23))
		Ω(err).ShouldNot(HaveOccurred())
	})
})
