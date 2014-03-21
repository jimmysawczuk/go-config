# go-config

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
	config.Add(config.Int("addend.a", 10, "The first addend", true))
	config.Add(config.Float("addend.b", math.Pi, "The second addend", true))
	config.Add(config.Bool("subtract", false, "Subtract instead of add", true))

	config.Build()

	addend_1 := float64(config.Require("addend.a").Int())
	addend_2 := float64(config.Require("addend.b").Float())

	subtract := config.Require("subtract").Bool()

	var res float64
	if !subtract {
		res = addend_1 + addend_2
	} else {
		res = addend_1 - addend_2
	}

	fmt.Println(addend_1, addend_2, subtract)
	fmt.Println(res)
}
```

The above program produces the following output:
```bash
$ config-test
10 3.141592653589793 false
13.141592653589793
```

You can override individual config values using flags:

```bash
$ config-test --help
Usage of config-test:
  -addend.a=10: The first addend
  -addend.b=3.141592653589793: The second addend
  -config="config.json": The filename of the config file to use
  -config-export=false: Export the as-run configuration to a file
  -subtract=false: Subtract instead of add
$ config-test --addend.a=5
5 3.141592653589793 false
8.141592653589793
```

### Automatic config file generation

Every program that uses **go-config** will parse two options by default:

* `config` (string), which indicates where to look for a config file to import and parse (default: ./config.json)
* `config-export` (bool), which, if true, will trigger an export of the config file to the path provided by **config**.

You can run your program with `config-export` to generate a fresh config file or overwrite the one that's there with the as-run configuration. Running the above program with `--config-export` yields:

```bash
$ config-test --config-export
10 3.141592653589793 false
13.141592653589793
$ cat config.json
{
    "addend": {
        "a": 10,
        "b": 3.141592653589793
    },
    "subtract": false
}
```

Notice that the two addends, which had keys with similar prefixes separated by periods (`.`), are organized in the config file to be hierarchial, but maintain a flat structure in your program.

## More documentation

More documentation is available [via GoDoc](http://godoc.org/github.com/jimmysawczuk/go-config).

## License

	The MIT License (MIT)
	Copyright (C) 2013-2014 by Jimmy Sawczuk

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in
	all copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE.
