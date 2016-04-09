package config

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Name is the name of the application you're configuring.
var Name = os.Args[0]

// Description describes the application you're configuring.
var Description string

// Version is the version of the application you're configuring.
var Version string

// Examples contains a list of example commands and what they do.
var Examples = []Example{}

// SearchFiles contains a list of files which may or may not exist, and if they do, contain
// configuration files. The last entry in this list is parsed first, and its values are
// overwritten by values in files further up the list.
var SearchFiles = []SearchFile{
	{
		Scope: "app",
		Path:  "./config.json",
	},
	{
		Scope: "user",
		Path:  "$HOME/." + strings.ToLower(Name) + "/config.json",
	},
}

// Example describes an example of a proper way to invoke the current Go program.
type Example struct {
	Cmd         string
	Description string
}

// SearchFile contains a potential config file path and a scope relating to where that file
// is stored.
type SearchFile struct {
	Scope string
	Path  string
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

	if Version != "" {
		uprintln(`%s (ver. %s)`, Name, Version)
	} else {
		uprintln(`%s`, Name)
	}

	if Description != "" {
		uprintln(`%s`, Description)
	}

	uprintln("")

	if len(Examples) > 0 {
		uprintln("Examples:")
		for _, v := range Examples {
			uprintln(" # %s", v.Description)
			uprintln(" $ %s\n", v.Cmd)
		}
	}

	if len(opts) > 0 {
		uprintln("Flags:")

		lastSort := opts[0].Options.SortOrder

		fmtStr := fmt.Sprintf(" -%%-%ds (default: %%s)\n     %%s\n", mlen)
		for _, opt := range opts {
			if opt.Options.SortOrder != lastSort {
				uprintln("")
			}

			uprintln(fmtStr,
				opt.Name,
				opt.defaultValueString("<empty>"),
				opt.Description,
			)

			lastSort = opt.Options.SortOrder
		}
	}
}

// ExpandedPath returns the SearchFile's path expanded with environment variables.
func (f SearchFile) ExpandedPath() string {
	return os.ExpandEnv(f.Path)
}
