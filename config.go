package config

import (
	"flag"
	"fmt"
	"os"
)

var baseOptionSet OptionSet

func init() {
	resetBaseOptionSet()
	flag.Usage = func() {}
}

func resetBaseOptionSet() {
	baseOptionSet = make(OptionSet)
	Add(Str("config", "", "A filename of an additional config file to use").SortOrder(998))
	Add(Bool("config-debug", false, "Show the files/scopes that are parsed and which scope each config value comes from").SortOrder(998))
	Add(Str("config-scope", "", "The scope that'll be written to").SortOrder(999))
	Add(Bool("config-partial", false, "Export a partial copy of the configuration, only what is explicitly passed in via flags").SortOrder(999))
	Add(Bool("config-save", false, "Export the configuration to the specified scope").SortOrder(999))
	Add(Bool("config-write", false, "Export the configuration to the specified scope, then exit").SortOrder(999))
}

// Add adds an Option to the config's OptionSet
func Add(o *Option) *Option {
	baseOptionSet.Add(o)
	return o
}

// Build builds the configuration object. Starts by setting the default values as defined in code, then parses the config file,
// then loads the overridden options from flag. If set, this also exports the as-run configuration to the the filename
// set in the "config" option.
func Build() error {
	var err error

	// parse flags
	fs := NewFlagSet(os.Args[0], os.Args[1:])
	perr := fs.ParseBuiltIn()
	if perr != nil {
		os.Exit(2)
	}

	searchFiles := make([]SearchFile, len(SearchFiles))
	copy(searchFiles, SearchFiles)

	overrideName := Require("config").String()
	if overrideName != "" {
		searchFiles = append([]SearchFile{{
			Scope: "flag",
			Path:  overrideName,
		}}, searchFiles...)
	}

	// find all the config files, import them
	for i := len(searchFiles) - 1; i >= 0; i-- {
		// fmt.Println("Parsing", os.ExpandEnv(searchFiles[i]))
		file := FileIO{
			filename: searchFiles[i].ExpandedPath(),
			scope:    searchFiles[i].Scope,
		}
		err = file.Read()
		if err != nil {
			if ioerr, ok := err.(IOError); ok {
				if ioerr.Type == "exist" {
					continue
				}

				fmt.Fprintf(os.Stderr, "go-config: error parsing config file: %s\n", ioerr.err)
				continue
			}

			if _, ok := err.(jsonConfigMapParseErrorList); ok {
				fmt.Println("Error:", err.Error())
				return err.(jsonConfigMapParseErrorList)
			}

			fmt.Println("Error:", err.Error())
			return fmt.Errorf("Error building config file: %s", err)
		}

		// fmt.Println(baseOptionSet)
	}

	fs = NewFlagSet(os.Args[0], os.Args[1:])
	perr = fs.Parse()
	if perr != nil {
		return perr
	}

	if fs.HasHelpFlag() {
		Usage()
		os.Exit(0)
		return nil
	}

	// validate all options that are required
	err = baseOptionSet.Validate()
	if err != nil {
		return err
	}

	for _, v := range baseOptionSet {
		fmt.Println(v.DebugString())
	}

	// export new config to file if necessary
	if Require("config-save").Bool() || Require("config-write").Bool() {

		scope := Require("config-scope").Str()

		found := false
		for _, v := range searchFiles {
			if v.Scope == scope {
				found = true
				file := FileIO{
					filename: v.ExpandedPath(),
					scope:    scope,
				}
				err := file.Write()
				if err != nil {
					return fmt.Errorf("go-config: can't write to file: %s", err)
				}
				if Require("config-write").Bool() {
					os.Exit(0)
				}
			}
		}

		if !found {
			return fmt.Errorf("go-config: can't find a config file with the scope %s", scope)
		}
	}

	os.Args = fs.Release()
	if len(os.Args) > 0 {
		err = flag.CommandLine.Parse(os.Args[1:])
	}

	return err
}

// Require looks for an Option with the name of `key`. If no Option is found, this function panics.
func Require(key string) *Option {

	s, err := Get(key)

	if err != nil {
		panic(err)
	}

	return s
}

// Get looks for an Option with the name of `key`. If no Option is found, this function returns an error.
func Get(key string) (*Option, error) {

	s, exists := baseOptionSet.Get(key)

	if !exists {
		return nil, fmt.Errorf("config option with key %s not found", key)
	}

	return s, nil
}
