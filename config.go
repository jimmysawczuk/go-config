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
	Add(Str("config", "", "The filename of the config file to use").SortOrder(999))
	Add(Bool("config-debug", false, "Show the files that are parsed and where each config value comes from"))
	Add(Str("config-save", "", "Export the as-run configuration to one of the loaded configuration scopes").SortOrder(999))
	Add(Str("config-write", "", "Export the as-run configuration to one of the loaded configuration scopes, then exit").SortOrder(999))
}

// Add adds an Option to the config's OptionSet
func Add(o *Option) {
	baseOptionSet.Add(o)
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
		file := FileIO{Filename: searchFiles[i].ExpandedPath()}
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

	// export new config to file if necessary
	if scope := Require("config-save").Str(); scope != "" {
		for _, v := range SearchFiles {
			if v.Scope == scope {
				file := FileIO{Filename: v.ExpandedPath()}
				file.Write()
			}
		}

		return fmt.Errorf("go-config: can't find a config file with the scope %s", scope)
	}

	if scope := Require("config-save").Str(); scope != "" {
		for _, v := range SearchFiles {
			if v.Scope == scope {
				file := FileIO{Filename: v.ExpandedPath()}
				file.Write()
				os.Exit(0)
			}
		}

		return fmt.Errorf("go-config: can't find a config file with the scope %s", scope)
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
