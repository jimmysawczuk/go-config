package config

import (
	"flag"
	"fmt"
	"os"
	"reflect"
)

var baseOptionSet OptionSet

func init() {
	resetBaseOptionSet()
	flag.Usage = func() {}
}

func resetBaseOptionSet() {
	baseOptionSet = make(OptionSet)
	Add(String("config", "config.json", "The filename of the config file to use", false))
	Add(Bool("config-export", false, "Export the as-run configuration to a file", false))
	Add(Bool("config-generate", false, "Export the as-run configuration to a file, then exit", false))
}

// Create an Option with the parameters given of type string
func String(name string, default_value string, description string, exportable bool) *Option {

	opt := Option{
		Name:        name,
		Description: description,

		DefaultValue: reflect.ValueOf(default_value),
		Value:        default_value,
		Type:         reflect.TypeOf(default_value),

		Exportable: exportable,
		// flag:       configFlags.String(name, default_value, description),
	}

	return &opt
}

// Create an Option with the parameters given of type bool
func Bool(name string, default_value bool, description string, exportable bool) *Option {

	opt := Option{
		Name:        name,
		Description: description,

		DefaultValue: reflect.ValueOf(default_value),
		Value:        default_value,
		Type:         reflect.TypeOf(default_value),

		Exportable: exportable,
		// flag:       configFlags.Bool(name, default_value, description),
	}

	return &opt
}

// Create an Option with the parameters given of type int64
func Int(name string, default_value int64, description string, exportable bool) *Option {

	opt := Option{
		Name:        name,
		Description: description,

		DefaultValue: reflect.ValueOf(default_value),
		Value:        default_value,
		Type:         reflect.TypeOf(default_value),

		Exportable: exportable,
		// flag:       configFlags.Int64(name, default_value, description),
	}

	return &opt
}

// Create an Option with the parameters given of type float64
func Float(name string, default_value float64, description string, exportable bool) *Option {

	opt := Option{
		Name:        name,
		Description: description,

		DefaultValue: reflect.ValueOf(default_value),
		Value:        default_value,
		Type:         reflect.TypeOf(default_value),

		Exportable: exportable,
		// flag:       configFlags.Float64(name, default_value, description),
	}

	return &opt
}

// Adds an Option to the config
func Add(o *Option) {
	baseOptionSet[o.Name] = o
}

// Builds the configuration object. Starts by setting the default values as defined in code, then parses the config file,
// then loads the overridden options from flag. If set, this also exports the as-run configuration to the the filename
// set in the "config" option.
func Build() error {
	// parse flags
	fs := NewFlagSet(os.Args[0], os.Args[1:])
	parse_err := fs.ParseBuiltIn()
	if parse_err != nil {
		os.Exit(2)
	}

	// determine location of config file, import it
	file := FileIO{Filename: Require("config").String()}
	err := file.Read()
	if err != nil {
		if _, ok := err.(jsonConfigMapParseErrorList); ok {
			return err.(jsonConfigMapParseErrorList)
		} else {
			return fmt.Errorf("Error building config file: %s", err)
		}
	}

	fs = NewFlagSet(os.Args[0], os.Args[1:])
	parse_err = fs.Parse()
	if parse_err != nil {
		return parse_err
	}

	if fs.HasHelpFlag() {
		Usage()
		os.Exit(0)
		return nil
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

// Requires that an Option with name key be found, otherwise panics.
func Require(key string) *Option {

	s, err := Get(key)

	if err != nil {
		panic(err)
	}

	return s
}

// Looks for an Option with name key, returns it if found, otherwise an error.
func Get(key string) (*Option, error) {

	s, exists := baseOptionSet.Get(key)

	if !exists {
		return nil, fmt.Errorf("config option with key %s not found", key)
	}

	return s, nil
}
