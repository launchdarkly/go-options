package test

import (
	"net/url"
	"time"
	time2 "time"
)

//go:generate go-options -imports=time,net/url,time2=time config
type config struct {
	myInt            int
	myIntWithDefault int `options:",1"`
	myRenamedInt     int `options:"yourInt"`

	// does something
	myDocumentedInt int
	myCommentedInt  int // for some reason

	// does something else
	myDocAndCommentInt int // for some other reason

	// takes a float
	myFloat            float64 // really a float
	myFloatWithDefault float64 `options:",1.23"`

	myString              string
	myStringWithDefault   string `options:",default string"`
	myStringWithoutOption string `options:"-"` // nolint:structcheck,unused // not expected to be used

	myFunc func() int

	myIntPointer *int

	myInterface interface{}

	// types requiring imports
	myURL       url.URL
	myDuration  time.Duration
	myDuration2 time2.Duration

	myStruct            struct{ a, b int }
	myStructWithDefault struct {
		a int `options:",1"`
	}
	myPointerToStruct         *struct{ a, b int }
	myStructWithVariadicSlice struct {
		a int
		b []int `options:"..."`
	}

	mySlice          []int  `options:"..."`
	myPointerToSlice *[]int `options:"..."`
	myRenamedSlice   []int  `options:"yourSlice..."`
}

//go:generate go-options -func applyDifferent -option DifferentOption -new=false configWithDifferentApply
type configWithDifferentApply struct {
}

//go:generate go-options -prefix Opt -option MyOpt configWithDifferentPrefix
type configWithDifferentPrefix struct {
	myFloat float64
}

//go:generate go-options -suffix Option -option SuffixOption configWithSuffix
type configWithSuffix struct {
	myFloat float64
}

//go:generate go-options -quote-default-strings=false -option UnquotedOption configWithUnquotedString
type configWithUnquotedString struct {
	myString string `options:",\"quoted\""`
}
