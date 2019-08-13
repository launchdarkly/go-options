package test

import (
	"net/url"
	"time"
)

//go:generate go-options -type=config -imports=time,url=net/url
type config struct {
	myInt                   int
	myIntWithDefault        int `options:",1"`
	myRenamedInt            int `options:"yourInt"`
	myRenamedIntWithDefault int `options:"yourIntWithDefault,1"`

	myFloat            float64
	myFloatWithDefault float64 `options:",1.23"`

	myString              string
	myStringWithDefault   string `options:",default string"`
	myStringWithoutOption string `options:"-"` // nolint:structcheck,unused // not expected to be used

	// types requiring imports
	myURL url.URL
	myDuration time.Duration
}

//go:generate go-options -type=configWithDifferentApply -func applyDifferent -option DifferentOption -new=false
type configWithDifferentApply struct {
}

//go:generate go-options -type=configWithDifferentPrefix -prefix Opt -option MyOpt
type configWithDifferentPrefix struct {
	myFloat float64
}
