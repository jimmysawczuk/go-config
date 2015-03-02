package config

import (
	"fmt"
	"io"
	"os"
	"sort"
)

var output_writer io.Writer = os.Stdout

func Usage() {
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

	fmt.Fprintf(output_writer, "%s\n", os.Args[0])

	fmt_str := fmt.Sprintf("  -%%-%ds  %%s\n", max_len)
	for _, opt := range opts {
		fmt.Fprintf(output_writer, fmt_str,
			fmt.Sprintf("%s=%s", opt.Name, opt.DefaultValueString()),
			opt.Description,
		)
	}
}
