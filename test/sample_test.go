package test

import (
	"errors"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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

	It("compares using standard equality", func() {
		Ω(OptionMyInt(1)).Should(Equal(OptionMyInt(1)))
	})

	It("returns errors", func() {
		err := applyConfigOptions(&cfg, OptionMakeError{})
		Ω(err).Should(MatchError("bad news"))
	})

	It("allows option constructor to be renamed", func() {
		err := applyConfigOptions(&cfg, OptionYourInt(1))
		Ω(err).ShouldNot(HaveOccurred())
		Ω(cfg.myRenamedInt).To(Equal(1))
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

		It("can store a pointer to let us know if a value was set", func() {
			err := applyConfigOptions(&cfg)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myPointerToInt).Should(BeNil())

			err = applyConfigOptions(&cfg, OptionMyPointerToInt(1))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myPointerToInt).ShouldNot(BeNil())
			Ω(*cfg.myPointerToInt).Should(Equal(1))
		})

		Describe("nested structs", func() {
			It("generates a constructor", func() {
				err := applyConfigOptions(&cfg, OptionMyStruct(1, 2))
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cfg.myStruct.a).Should(Equal(1))
				Ω(cfg.myStruct.b).Should(Equal(2))
			})

			It("allows default values", func() {
				err := applyConfigOptions(&cfg)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cfg.myStructWithDefault.a).Should(Equal(1))
			})

			It("defaults pointer structs to nil", func() {
				err := applyConfigOptions(&cfg)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cfg.myPointerToStruct).Should(BeNil())

				err = applyConfigOptions(&cfg, OptionMyPointerToStruct(1, 2))
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cfg.myPointerToStruct).ShouldNot(BeNil())
				Ω(cfg.myPointerToStruct.a).Should(Equal(1))
				Ω(cfg.myPointerToStruct.b).Should(Equal(2))
			})

			It("allows variadic arguments within a slice", func() {
				err := applyConfigOptions(&cfg, OptionMyStructWithVariadicSlice(1, 1, 2))
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cfg.myStructWithDefault.a).Should(Equal(1))
			})

			It("allows variadic arguments to be compared with cmp", func() {
				cmp.Equal(
					OptionMyStructWithVariadicSlice(1, 1, 2),
					OptionMyStructWithVariadicSlice(1, 1, 2))
			})
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
				Ω(cfg.myPointerToSlice).Should(BeNil())

				err = applyConfigOptions(&cfg, OptionMyPointerToSlice(1, 2))
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cfg.myPointerToSlice).ShouldNot(BeNil())
				Ω(*cfg.myPointerToSlice).Should(ConsistOf(1, 2))
			})

			It("allows them to be renamed", func() {
				err := applyConfigOptions(&cfg, OptionYourSlice(1, 2))
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cfg.myRenamedSlice).ShouldNot(BeNil())
				Ω(cfg.myRenamedSlice).Should(ConsistOf(1, 2))
			})

			It("allows them to the to be compared with cmp", func() {
				Ω(cmp.Equal(
					OptionYourSlice(1, 2),
					OptionYourSlice(1, 2))).Should(BeTrue())
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

var _ = Describe("not quoting strings by default", func() {
	It("requires them to be quoted", func() {
		cfg, err := newConfigWithUnquotedString()
		Ω(err).ShouldNot(HaveOccurred())
		Ω(cfg.myString).Should(Equal("quoted"))
	})
})

var _ = Describe("Disabling cmp", func() {
	It("prevents options from implementing Equal", func() {
		_, equalsFound := reflect.TypeOf(NoCmpOptionMyInt(1)).MethodByName("Equal")
		Ω(equalsFound).Should(BeFalse())
	})
})
