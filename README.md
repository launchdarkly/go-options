# LaunchDarkly Options Generator

[![CircleCI](https://circleci.com/gh/launchdarkly/go-options.svg?style=svg)](https://circleci.com/gh/launchdarkly/go-options)

The LaunchDarkly Options Generator generates boilerplate code for setting options for a configuration struct using varargs syntax.  You write this:

```go
//go:generate go-options config
type config struct {
	howMany int
}
```

Then run go generate and you can write this:

```go
cfg, err := newConfig(OptionHowMany(100))
```

or, more interestingly, this:

```go
type Collection {
    config
}

func NewCollection(options... Option) (Foo, err) {
    cfg, err := newConfig(options...)
    return Collection{cfg}, err
}
```

You can also specify default values and override the option name as follows:

```go
//go:generate go-options config
type config struct {
	howMany int `options:"number,5"
}
```

This would create `OptionNumber` with a default value of 5.  Entering the the tag `options:",5"` would keep the default `OptionHowMany` name.

You can also specify documentation using docstrings or line strings, so:

```go
//go:generate go-options config
type config struct {
    // indicates the number of items
    howMany int // no more than 10
}
```

would generate code that looks like this:

```go
// OptionHowMany indicates the number of items
// no more than ten
func OptionHowMany(o int) applyOptionFunc {
    // ...
}
```

You can use nested structures to create multi-field options, so:


```go
type config struct {
    number struct {
        a, b int
    }
}
```

would yield:

```go
func OptionNumber(a int, b int) applyOptionFunc {
    // ...
}
```

You can use also use "..." at the end of a name in `options` tag to create variadic arguments, so:

```go
type config struct {
  numbers []int `options:"..."`
  nums []int `options:"ints..."`
}
```

would yield:

```
func OptionNumbers(numbers ...int) applyOptionFunc {
    // ...
}

func OptionInts(nums ...int) applyOptionFunc {
    // ...
}
```


You can use also use "*" at the beginning of a name in `options` tag to record whether an option was set, so:

```go
type config struct {
  value *int `options:"*"`
  v *int `options:"*myValue"`
}
```

would yield:

```
func OptionValue(o ...int) applyOptionFunc {
    // ...
}

func OptionMyValue(o ...int) applyOptionFunc {
    // ...
}
```

Generated options are interoperable with any other user-created options that support the option interface:

```
type Option interface {
    apply(config *c) error
}
```

The name `Option` can be customized along with various method names as shown under [Options](#options) below.

## Installation

Install with `go get -u github.com/launchdarkly/go-options`.

## Tag Syntax

The syntax for a tag is:

`<alternateName or blank>,[optional default value]`

## For testing and debugging

By default, generated options can be compared using `cmp.Equal` from `github.com/google/go-cmp`.  Simple options can
also be compared simply with `==` because they are structs; more complex options involving variadic slices and pointers 
require using `cmp.Equal` because pointers inside the options will not match.  To allow `cmp.Equal` compare options, the
tool generates  an `Equal` method for each option.  Generation of `Equal` methods can be disabled by setting
`-cmp=false`.

To aid with debugging and producing more meaningful error in tests, the tool generates a `String()` method for each
option.  This method fulfills the `fmt.Stringer` interface, allowing more details about the option to be included in the
`%v` format verb.  This behavior can be disabled by setting `-stringer=false`.

## Options

`go-options` can be customized with several command-line arguments:

- `-fmt=false` disable running gofmt
- `-func <string>` sets the name of function created to apply options to <type> (default is apply&lt;Type&gt;Options)
- `-new=false` controls generation of the function that returns a new config (default true)
- `-cmp=false` controls whether we generate an `Equal` method that works with `github.com/google/go-cmp` (default true)
- `-imports=[<path>|<alias>=<path>],...` add imports to generated file
- `-option <string>` sets name of the interface to use for options (default "Option")
- `-output <string>` sets the name of the output file (default is <type>_options.go)
- `-input <string>` sets the name of the input file. When set uses "go/build" and "go/parser" directly, which can result in performance improvements
- `-prefix <string>` sets prefix to be used for options (defaults to the value of `option`)
- `-quote-default-strings=false` disables default quoting of default values for string
- `-stringer=false` controls whether we generate an `String()` method that exposes option names and values.  Useful for debugging tests. (default true)
- `-suffix <string>` sets suffix to be used for options (instead of prefix, cannot be used with `prefix` option)
- `-type <string>` name of struct type to create options for (original syntax before multiple types on command-line were supported)
