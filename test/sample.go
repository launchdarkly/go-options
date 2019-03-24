package test

//go:generate go-options -type=config
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
}

//go:generate go-options -type=configWithDifferentApply -func applyDifferent -option DifferentOption -new=false
type configWithDifferentApply struct {
}
