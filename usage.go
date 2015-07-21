package config

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// UsageWriter is the io.Writer to use for outputting Usage(). Defaults to stdout.
var UsageWriter io.Writer = os.Stdout

type sortedUsageOptionSlice []Option

func (s sortedUsageOptionSlice) Less(a, b int) bool {
	if s[a].Options.SortOrder < s[b].Options.SortOrder {
		return true
	} else if s[a].Options.SortOrder > s[b].Options.SortOrder {
		return false
	}

	return s[a].Name < s[b].Name
}

func (s sortedUsageOptionSlice) Swap(a, b int) { s[a], s[b] = s[b], s[a] }
func (s sortedUsageOptionSlice) Len() int      { return len(s) }

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

	sort.Sort(sortedUsageOptionSlice(opts))

	uprintln(`%s`, os.Args[0])

	fmtStr := fmt.Sprintf(`  -%%-%ds  %%s`, mlen)
	for _, opt := range opts {
		uprintln(fmtStr,
			fmt.Sprintf("%s=%s", opt.Name, opt.DefaultValueString()),
			opt.Description,
		)
	}
}
