package config

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// UsageWriter is the io.Writer to use for outputting Usage(). Defaults to stdout.
var UsageWriter io.Writer = os.Stdout

// Usage prints the help information to UsageWriter (defaults to stdout).
func Usage() {

	uprintf := func(strFmt string, args ...interface{}) {
		fmt.Fprintf(UsageWriter, strFmt, args...)
	}

	uprintln := func(strFmt string, args ...interface{}) {
		uprintf(strFmt+"\n", args...)
	}

	mlen := 0
	opts := []Option{}
	for _, opt := range baseOptionSet {
		opts = append(opts, *opt)
		s := fmt.Sprintf("%s=%s", opt.Name, opt.DefaultValueString())
		if len(s) > mlen {
			mlen = len(s)
		}
	}

	sort.Sort(sortedOptionSlice(opts))

	uprintln(`%s`, os.Args[0])

	fmtStr := fmt.Sprintf(`  -%%-%ds  %%s`, mlen)
	for _, opt := range opts {
		uprintln(fmtStr,
			fmt.Sprintf("%s=%s", opt.Name, opt.DefaultValueString()),
			opt.Description,
		)
	}
}
