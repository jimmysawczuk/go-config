# go-config

[ ![travis-ci status for jimmysawczuk/go-config](https://travis-ci.org/jimmysawczuk/go-config.svg)](https://travis-ci.org/jimmysawczuk/go-config) [![GoDoc](https://godoc.org/github.com/jimmysawczuk/go-config?status.svg)](https://godoc.org/github.com/jimmysawczuk/go-config) [![Go Report Card](https://goreportcard.com/badge/github.com/jimmysawczuk/go-config)](https://goreportcard.com/report/github.com/jimmysawczuk/go-config)

**go-config** is a configuration library for Go that can process configuration from the source code itself, config files, and the command line.

## Usage

### Example

Here's some simple usage of **go-config** in action:

```go
package main

import (
	"fmt"
	"github.com/jimmysawczuk/go-config"
	"math"
)

func main() {
	// Set up the app information for the --help output
	config.Name = "config-test"
	config.Version = "1.0.0"
	config.Description = "Takes two arguments and does an operation on them"

	config.Examples = []config.Example{
		{
			Cmd:         `config-tester -addend.a=1 -addend-b=2`,
			Description: "Adds 1 and 2, returns 3",
		},

		{
			Cmd:         `config-tester -addend.a=3 -addend-b=2 -subtract`,
			Description: "Subtracts 2 from 3, returns 1",
		},
	}

	var res float64

	// Add the variables, store the options as variables for later
	a := config.Add(config.Int("addend.a", 10, "The first addend").Exportable(true))
	b := config.Add(config.Float("addend.b", math.Pi, "The second addend").Exportable(true))
	sub := config.Add(config.Bool("subtract", false, "Subtract instead of add").Exportable(true))

	// Build the config
	config.Build()

	// Calculate the result
	res = op(float64(a.Int()), b.Float(), sub.Bool())

	fmt.Println("From the stored variables")
	fmt.Println(a.Int(), b.Float(), sub.Bool())
	fmt.Println(res)

	// You can also get the config values on demand
	addend_a := float64(config.Require("addend.a").Int())
	addend_b := float64(config.Require("addend.b").Float())
	subtract := config.Require("subtract").Bool()

	res = op(addend_a, addend_b, subtract)

	fmt.Println("------")
	fmt.Println("From the on demand lookup")
	fmt.Println(addend_a, addend_b, sub)
	fmt.Println(res)

}

func op(a, b float64, subtract bool) (r float64) {
	if !subtract {
		r = a + b
	} else {
		r = a - b
	}
	return r
}
```

The above program produces the following output:
```bash
$ config-test
From the stored variables
10 3.141592653589793 false
13.141592653589793
------
From the on demand lookup
10 3.141592653589793 false
13.141592653589793
```

You can override individual config values using flags:

```bash
$ config-test --help
config-test (ver. 1.0.0)
Takes two arguments and does an operation on them

Examples:
 # Adds 1 and 2, returns 3
 $ config-tester -addend.a=1 -addend.b=2

 # Subtracts 2 from 3, returns 1
 $ config-tester -addend.a=3 -addend.b=2 -subtract

Flags:
 -addend.a     (default: 10)
     The first addend

 -addend.b     (default: 3.141592653589793)
     The second addend

 -config-debug (default: false)
     Show the files that are parsed and where each config value comes from

 -subtract     (default: false)
     Subtract instead of add


 -config       (default: <empty>)
     The filename of the config file to use

 -config-save  (default: <empty>)
     Export the as-run configuration to one of the loaded configuration scopes

 -config-write (default: <empty>)
     Export the as-run configuration to one of the loaded configuration scopes, then exit

```

### Automatic config file generation

WIP

## More documentation

More documentation is available [via GoDoc][godoc].

## License

go-config is released under [the MIT license][license].

  [license]: https://github.com/jimmysawczuk/go-config/blob/master/LICENSE
  [godoc]: http://godoc.org/github.com/jimmysawczuk/go-config
