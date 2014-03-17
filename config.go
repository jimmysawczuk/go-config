package config

import (
	"flag"
	"reflect"

	"fmt"
)

var baseOptionSet OptionSet

func init() {
	baseOptionSet = make(OptionSet)

	Add(String("config", "config.json", "The filename of the config file to use", false))
	Add(Bool("config-export", false, "Export the as-run configuration to a file", false))
}

func String(name string, default_value string, description string, exportable bool) *Option {

	opt := Option{
		Name:        name,
		Description: description,

		DefaultValue: reflect.ValueOf(default_value),
		Value:        reflect.ValueOf(default_value),
		Type:         reflect.TypeOf(default_value),

		Exportable: exportable,
		flag:       flag.String(name, default_value, description),
	}

	return &opt
}

func Bool(name string, default_value bool, description string, exportable bool) *Option {

	opt := Option{
		Name:        name,
		Description: description,

		DefaultValue: reflect.ValueOf(default_value),
		Value:        reflect.ValueOf(default_value),
		Type:         reflect.TypeOf(default_value),

		Exportable: exportable,
		flag:       flag.Bool(name, default_value, description),
	}

	return &opt
}

func Int(name string, default_value int64, description string, exportable bool) *Option {

	opt := Option{
		Name:        name,
		Description: description,

		DefaultValue: reflect.ValueOf(default_value),
		Value:        reflect.ValueOf(default_value),
		Type:         reflect.TypeOf(default_value),

		Exportable: exportable,
		flag:       flag.Int64(name, default_value, description),
	}

	return &opt
}

func Float(name string, default_value float64, description string, exportable bool) *Option {

	opt := Option{
		Name:        name,
		Description: description,

		DefaultValue: reflect.ValueOf(default_value),
		Value:        reflect.ValueOf(default_value),
		Type:         reflect.TypeOf(default_value),

		Exportable: exportable,
		flag:       flag.Float64(name, default_value, description),
	}

	return &opt
}

func Add(o *Option) {
	baseOptionSet[o.Name] = o
}

func Build() {
	// parse flags
	flag.Parse()

	// set default values
	importFlags(true)

	// determine location of config file, import it
	config_filename := Require("config").String()
	importConfigFile(config_filename)

	// overwrite with flag
	importFlags(false)

	// export new config to file if necessary
	if Require("config-export").Bool() {
		exportConfigToFile(config_filename)
	}
}

func Require(key string) *Option {

	s, err := Get(key)

	if err != nil {
		panic(err)
	}

	return s
}

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
			var target, val reflect.Value
			target = reflect.ValueOf(v).Elem().FieldByName("Value")

			switch v.flag.(type) {
			case *string:
				val = reflect.ValueOf(*(v.flag.(*string)))
			case *int64:
				val = reflect.ValueOf(*(v.flag.(*int64)))
			case *float64:
				val = reflect.ValueOf(*(v.flag.(*float64)))
			case *bool:
				val = reflect.ValueOf(*(v.flag.(*bool)))
			}

			fmt.Println(f.Name, val, visitall, f)

			target.Set(val)
		}
	}

	if visitall {
		flag.VisitAll(setter)
	} else {
		flag.Visit(setter)
	}
}
