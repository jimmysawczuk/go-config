package config

import (
	"flag"
	"fmt"
	"os"
	"reflect"
)

var baseOptionSet OptionSet
var configFlags *flag.FlagSet

func init() {
	resetBaseOptionSet(true)
	flag.Usage = func() {}
}

func resetBaseOptionSet(add_defaults bool) {
	baseOptionSet = make(OptionSet)

	configFlags = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	if add_defaults {
		Add(String("config", "config.json", "The filename of the config file to use", false))
		Add(Bool("config-export", false, "Export the as-run configuration to a file", false))
		Add(Bool("config-generate", false, "Export the as-run configuration to a file, then exit", false))
	}
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
		flag:       configFlags.String(name, default_value, description),
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
		flag:       configFlags.Bool(name, default_value, description),
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
		flag:       configFlags.Int64(name, default_value, description),
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
		flag:       configFlags.Float64(name, default_value, description),
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
	parse_err := configFlags.Parse(os.Args[1:])
	if parse_err != flag.ErrHelp && parse_err != nil {
		os.Exit(2)
	}

	// set default values
	importFlags(true)

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

	// overwrite with flag
	importFlags(false)

	if parse_err == flag.ErrHelp {
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

	return nil
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

func importFlags(visitall bool) {
	setter := func(f *flag.Flag) {
		if v, exists := baseOptionSet.Get(f.Name); exists {
			var val interface{}

			switch v.flag.(type) {
			case *string:
				val = *(v.flag.(*string))
			case *int64:
				val = *(v.flag.(*int64))
			case *float64:
				val = *(v.flag.(*float64))
			case *bool:
				val = *(v.flag.(*bool))
			}

			v.Value = val
		}
	}

	if visitall {
		configFlags.VisitAll(setter)
	} else {
		configFlags.Visit(setter)
	}
}
