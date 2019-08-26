package test

import (
	"net/url"
	"time"
	time2 "time"
)

//go:generate go-options -imports=time,net/url,time2=time config
type config struct {
	myInt int // takes an integer
	// has documentation
	myIntWithDefault        int `options:",1"`
	myRenamedInt            int `options:"yourInt"` // does something
	myRenamedIntWithDefault int `options:"yourIntWithDefault,1"`

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
}

//go:generate go-options -type=configWithDifferentApply -func applyDifferent -option DifferentOption -new=false
type configWithDifferentApply struct {
}

//go:generate go-options -type=configWithDifferentPrefix -prefix Opt -option MyOpt
type configWithDifferentPrefix struct {
	myFloat float64
}
