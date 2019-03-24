# LaunchDarkly Options Generator

The LaunchDarkly Options Generator generates boilerplate code for setting options for a configuration struct using varargs syntax.

The idea is to start with this:

```
//go:generate go-options -type=config
type config struct {
	howMany int `options:",5"`
}
```

and be able to write this:

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

## Installation

Install with `go get -u github.com/launchdarkly/go-options`.

## Tag Syntax

The syntax for a tag is:

`<alternateName or blank>,[optional default value]`

## Options

`go-options` can be customized with several command-line arguments:

- `-type <string>` name of struct type to create options for (required)
- `-func <string>` sets the name of function created to apply options to <type> (default is apply<Type>Options)
- `-new=false` controls generation of the function that returns a new config (default true)
- `-option <string>` sets name of the interface to use for options (default "Option")
- `-output <string>` sets the name of the output file (default is <type>_options.go)
