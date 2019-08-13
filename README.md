# LaunchDarkly Options Generator

The LaunchDarkly Options Generator generates boilerplate code for setting options for a configuration struct using varargs syntax.  You write this:

```
//go:generate go-options config
type config struct {
	howMany int
}
```

Then run go generate and you can write this:

```
cfg, err := newConfig(OptionHowMany(100))

```

or, more interestingly, this:

```
type Collection {
    config
}

func NewCollection(options... Option) (Foo, err) {
    cfg, err := newConfig(options...)
    return Collection{cfg}, nil
}
```

You can also specify default values and override the option name as follows:

```
//go:generate go-options config
type config struct {
	howMany int `options:"number,5"
}
```

This would create `OptionNumber` with a default value of 5.  Entering the the tag `options:",5"` would keep the default `OptionHowMany` name.

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

## Options

`go-options` can be customized with several command-line arguments:

- `-fmt=false` disable running gofmt
- `-func <string>` sets the name of function created to apply options to <type> (default is apply&lt;Type&gt;Options)
- `-new=false` controls generation of the function that returns a new config (default true)
- `-imports=[<path>|<alias>=<path>],...` add imports to generated file
- `-option <string>` sets name of the interface to use for options (default "Option")
- `-output <string>` sets the name of the output file (default is <type>_options.go)
- `-prefix <string>` sets prefix to be used for options (defaults to the value of `option`)
- `-type <string>` name of struct type to create options for (original syntax before multiple types on command-line were supported)
