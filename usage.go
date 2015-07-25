package config

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// AppInfo describes a Go program that uses go-config. It is used primarily in Usage(), or when a user uses -help.
type AppInfo struct {
	Name        string
	Description string
	Version     string
	Examples    []Example
}

// Example describes an example of a proper way to invoke the current Go program.
type Example struct {
	Cmd      string
	Function string
}

// App is the AppInfo for the current Go program.
var App = AppInfo{
	Name:     os.Args[0],
	Examples: []Example{},
}

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
		s := fmt.Sprintf("%s", opt.Name)
		if len(s) > mlen {
			mlen = len(s)
		}
	}

	sort.Sort(sortedUsageOptionSlice(opts))

	uprintln(`Usage of %s:`, App.Name)
	if App.Description != "" {
		uprintln(`  %s`, App.Description)
	}

	uprintln("")

	if len(App.Examples) > 0 {
		uprintln("  Examples:")
		for _, v := range App.Examples {
			uprintln("     %s\n     %s\n", v.Cmd, v.Function)
		}
	}

	if len(opts) > 0 {
		uprintln("  Flags:")

		lastSort := opts[0].Options.SortOrder

		fmtStr := fmt.Sprintf("    -%%-%ds (default: %%s)\n     %%s\n", mlen)
		for _, opt := range opts {
			if opt.Options.SortOrder != lastSort {
				uprintln("")
			}

			uprintln(fmtStr,
				fmt.Sprintf("%s", opt.Name),
				opt.DefaultValueString(),
				opt.Description,
			)

			lastSort = opt.Options.SortOrder
		}
	}
}
