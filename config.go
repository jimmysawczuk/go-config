package config

import (
	"errors"
	"flag"
	"reflect"
)

var optionsDict OptionsSet

func init() {
	optionsDict = OptionsSet{
		options: make(map[string]*Option),
	}

	String("config", "config.json", "The filename of the config file to use", false)
	Bool("config-export", false, "Export the as-run configuration to a file", false)
}

func String(name string, default_value string, description string, exportable bool) {

	opt := Option{
		Name:          name,
		Description:   description,
		default_value: default_value,

		exportable: exportable,
		flag_val:   flag.String(name, default_value, description),
		val:        reflect.ValueOf(default_value),
		opt_type:   reflect.TypeOf(default_value),
	}

	optionsDict.Add(opt)
}

func Bool(name string, default_value bool, description string, exportable bool) {
	opt := Option{
		Name:          name,
		Description:   description,
		default_value: default_value,

		exportable: exportable,
		flag_val:   flag.Bool(name, default_value, description),
		val:        reflect.ValueOf(default_value),
		opt_type:   reflect.TypeOf(default_value),
	}

	optionsDict.Add(opt)
}

func Int(name string, default_value int, description string, exportable bool) {
	opt := Option{
		Name:          name,
		Description:   description,
		default_value: default_value,

		exportable: exportable,
		flag_val:   flag.Int(name, default_value, description),
		val:        reflect.ValueOf(default_value),
		opt_type:   reflect.TypeOf(default_value),
	}

	optionsDict.Add(opt)
}

func Int64(name string, default_value int64, description string, exportable bool) {
	opt := Option{
		Name:          name,
		Description:   description,
		default_value: default_value,

		exportable: exportable,
		flag_val:   flag.Int64(name, default_value, description),
		val:        reflect.ValueOf(default_value),
		opt_type:   reflect.TypeOf(default_value),
	}

	optionsDict.Add(opt)
}

func Float64(name string, default_value float64, description string, exportable bool) {
	opt := Option{
		Name:          name,
		Description:   description,
		default_value: default_value,

		exportable: exportable,
		flag_val:   flag.Float64(name, default_value, description),
		val:        reflect.ValueOf(default_value),
		opt_type:   reflect.TypeOf(default_value),
	}

	optionsDict.Add(opt)
}

func Build() {

	importFlags(true)
	config_filename := Require("config").String()
	importConfigFile(config_filename)
	flag.Parse()
	importFlags(false)

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

	s, exists := optionsDict.Get(key)

	if !exists {
		return nil, errors.New("Config with key " + key + " not found")
	}

	return s, nil
}

func Set(key string, value interface{}) {
	s, exists := optionsDict.Get(key)

	if !exists {
		// TODO: return error
		return
	}

	s.Set(value)
}

func importFlags(visitall bool) {

	setter := func(f *flag.Flag) {
		if v, exists := optionsDict.Get(f.Name); exists {

			switch v.opt_type.Name() {
			case "string":
				v.Set(*(v.flag_val.(*string)))
			case "bool":
				v.Set(*(v.flag_val.(*bool)))
			case "int":
				v.Set(*(v.flag_val.(*int)))
			case "int64":
				v.Set(*(v.flag_val.(*int64)))
			case "float64":
				v.Set(*(v.flag_val.(*float64)))
			}
		}
	}

	if visitall {
		flag.VisitAll(setter)
	} else {
		flag.Visit(setter)
	}
}
