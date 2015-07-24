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
	Add(String("config", "config.json", "The filename of the config file to use"))
	Add(Bool("config-export", false, "Export the as-run configuration to a file"))
	Add(Bool("config-generate", false, "Export the as-run configuration to a file, then exit"))
}

// Add adds an Option to the config's OptionSet
func Add(o *Option) {
	baseOptionSet.Add(o)
}

// Build builds the configuration object. Starts by setting the default values as defined in code, then parses the config file,
// then loads the overridden options from flag. If set, this also exports the as-run configuration to the the filename
// set in the "config" option.
func Build() error {
	// parse flags
	fs := NewFlagSet(os.Args[0], os.Args[1:])
	perr := fs.ParseBuiltIn()
	if perr != nil {
		os.Exit(2)
	}

	// determine location of config file, import it
	file := FileIO{Filename: Require("config").String()}
	err := file.Read()
	if err != nil {
		if _, ok := err.(jsonConfigMapParseErrorList); ok {
			return err.(jsonConfigMapParseErrorList)
		}

		return fmt.Errorf("Error building config file: %s", err)
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
	if Require("config-export").Bool() || Require("config-generate").Bool() {
		file.Write()
	}

	if Require("config-generate").Bool() {
		os.Exit(0)
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
