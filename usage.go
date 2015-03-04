package config

import (
	"fmt"
	"io"
	"os"
	"sort"
)

var UsageWriter io.Writer = os.Stdout

func Usage() {

	uprintf := func(fmt_str string, args ...interface{}) {
		fmt.Fprintf(UsageWriter, fmt_str, args...)
	}

	uprintln := func(fmt_str string, args ...interface{}) {
		uprintf(fmt_str+"\n", args...)
	}

	max_len := 0
	opts := []Option{}
	for _, opt := range baseOptionSet {
		opts = append(opts, *opt)
		s := fmt.Sprintf("%s=%s", opt.Name, opt.DefaultValueString())
		if len(s) > max_len {
			max_len = len(s)
		}
	}

	sort.Sort(SortedOptionSlice(opts))

	uprintln(`%s`, os.Args[0])

	fmt_str := fmt.Sprintf(`  -%%-%ds  %%s`, max_len)
	for _, opt := range opts {
		uprintln(fmt_str,
			fmt.Sprintf("%s=%s", opt.Name, opt.DefaultValueString()),
			opt.Description,
		)
	}
}
