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
	err := config.Build()
	if err != nil {
		fmt.Println(err.Error())
	}

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
 $ config-test -addend.a=1 -addend-b=2

 # Subtracts 2 from 3, returns 1
 $ config-test -addend.a=3 -addend-b=2 -subtract

Flags:
 -addend.a       (default: 10)
     The first addend

 -addend.b       (default: 3.141592653589793)
     The second addend

 -subtract       (default: false)
     Subtract instead of add


 -config-debug   (default: false)
     Show the files/scopes that are parsed and which scope each config value comes from

 -config-file    (default: <empty>)
     A filename of an additional config file to use


 -config-partial (default: false)
     Export a partial copy of the configuration, only what is explicitly passed in via flags

 -config-save    (default: false)
     Export the configuration to the specified scope

 -config-scope   (default: <empty>)
     The scope that'll be written to

 -config-write   (default: false)
     Export the configuration to the specified scope, then exit

```

### Automatic config file generation

Config files can be saved as JSON files. go-config supports parsing multiple config files and in the event of two files having different values for one option, takes the most recently parsed option. The default order in which config files and arguments are parsed is:

1. `$HOME/.<program-name>/config.json` (the user's home directory) (scope: `"user"`)
2. `./config.json` (the working directory) (scope: `"app"`)
3. A config file specified via the `-config-file` flag (optional) (scope: `"custom"`)
4. Any flags specified on the command line at runtime (scope: `"flag"`)

You can automatically write a config file by specifying `-config-scope` (see the list above), a `-config-file` if necessary, and either `-config-save` (which continues execution of the program after saving the config file) or `-config-write` (which terminates the program after writing). By default, this will write all of the exportable options to the specified file, but you can specify `-config-partial` to only write the config values specified by flag (and not the rest of the exportable options).

## More documentation

More documentation is available [via GoDoc][godoc].

## License

go-config is released under [the MIT license][license].

  [license]: https://github.com/jimmysawczuk/go-config/blob/master/LICENSE
  [godoc]: http://godoc.org/github.com/jimmysawczuk/go-config
