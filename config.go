package config

import (
	"errors"
	"flag"
	"reflect"

	"fmt"
)

type Option struct {
	Name        string
	Description string

	optType    reflect.Type
	exportable bool

	default_value interface{}

	val      interface{}
	flag_val interface{}
}

func (o Option) isExportable() bool {
	return o.exportable
}

func (o *Option) Set(v interface{}) error {
	switch o.optType.Name() {
	case "stringOption":
		switch v.(type) {
		case string:
			o.val = v.(string)
		case fmt.Stringer:
			o.val = v.(fmt.Stringer).String()
		default:
			return fmt.Errorf("Invalid data passed to *Option.Set(): %v (type %T), expected string or fmt.Stringer")
		}

	case "boolOption":
		switch v.(type) {
		case bool:
			o.val = v.(bool)
		default:
			return fmt.Errorf("Invalid data passed to *Option.Set(): %v (type %T), expected bool")
		}
		o.val = v.(bool)

	case "intOption":
		switch v.(type) {
		case int:
			o.val = v.(int)
		case float64:
			o.val = int(v.(float64))
		default:
			return fmt.Errorf("Invalid data passed to *Option.Set(): %v (type %T), expected int or float64")
		}

	case "int64Option":
		switch v.(type) {
		case int:
			o.val = int64(v.(int))
		case int64:
			o.val = v.(int64)
		case float64:
			o.val = int64(v.(float64))
		default:
			return fmt.Errorf("Invalid data passed to *Option.Set(): %v (type %T), expected int64 or float64")
		}

	case "float64Option":
		switch v.(type) {
		case float64:
			o.val = v.(float64)
		case float32:
			o.val = float64(v.(float32))
		default:
			return fmt.Errorf("Invalid data passed to *Option.Set(): %v (type %T), expected float or float64")
		}
	}

	return nil
}

type OptionsSet struct {
	options map[string]*Option
}

func (os OptionsSet) Export() map[string]interface{} {
	tbr := make(map[string]interface{})
	for _, v := range os.options {
		if v.isExportable() {
			tbr[v.Name] = v.val
		}
	}
	return tbr
}

func (os *OptionsSet) Add(o interface{}) {

	t := reflect.TypeOf(o)
	var opt Option

	switch o.(type) {
	case stringOption:
		opt = Option(o.(stringOption))
	case boolOption:
		opt = Option(o.(boolOption))
	case intOption:
		opt = Option(o.(intOption))
	case int64Option:
		opt = Option(o.(int64Option))
	case float64Option:
		opt = Option(o.(float64Option))
	}

	opt.optType = t

	os.options[opt.Name] = &opt
}

func (os OptionsSet) Get(key string) (*Option, bool) {
	result, exists := os.options[key]
	return result, exists
}

var optionsDict OptionsSet

func init() {
	optionsDict = OptionsSet{
		options: make(map[string]*Option),
	}

	String("config", "config.json", "The filename of the config file to use", false)
	Bool("config-export", false, "Export the as-run configuration to a file", false)
}

func String(name string, default_value string, description string, exportable bool) {

	opt := stringOption{
		Name:          name,
		Description:   description,
		default_value: default_value,

		exportable: exportable,
		flag_val:   flag.String(name, default_value, description),
		val:        default_value,
	}

	optionsDict.Add(opt)
}

func Bool(name string, default_value bool, description string, exportable bool) {
	opt := boolOption{
		Name:          name,
		Description:   description,
		default_value: default_value,

		exportable: exportable,
		flag_val:   flag.Bool(name, default_value, description),
		val:        default_value,
	}

	optionsDict.Add(opt)
}

func Int(name string, default_value int, description string, exportable bool) {
	opt := intOption{
		Name:          name,
		Description:   description,
		default_value: default_value,

		exportable: exportable,
		flag_val:   flag.Int(name, default_value, description),
		val:        default_value,
	}

	optionsDict.Add(opt)
}

func Int64(name string, default_value int64, description string, exportable bool) {
	opt := int64Option{
		Name:          name,
		Description:   description,
		default_value: default_value,

		exportable: exportable,
		flag_val:   flag.Int64(name, default_value, description),
		val:        default_value,
	}

	optionsDict.Add(opt)
}

func Float64(name string, default_value float64, description string, exportable bool) {
	opt := float64Option{
		Name:          name,
		Description:   description,
		default_value: default_value,

		exportable: exportable,
		flag_val:   flag.Float64(name, default_value, description),
		val:        default_value,
	}

	optionsDict.Add(opt)
}

func Build() {

	importFlags(true)
	config_filename := Require("config").(string)
	importConfigFile(config_filename)
	flag.Parse()
	importFlags(false)

	if Require("config-export").(bool) {
		exportConfigToFile(config_filename)
	}
}

func Require(key string) (s interface{}) {

	s, err := Get(key)

	if err != nil {
		panic(err)
	}

	return s
}

func Get(key string) (val interface{}, err error) {

	s, exists := optionsDict.Get(key)

	if !exists {
		return nil, errors.New("Config with key " + key + " not found")
	}

	return s.val, nil
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

			switch v.optType.Name() {
			case "stringOption":
				v.Set(*(v.flag_val.(*string)))
			case "boolOption":
				v.Set(*(v.flag_val.(*bool)))
			case "intOption":
				v.Set(*(v.flag_val.(*int)))
			case "int64Option":
				v.Set(*(v.flag_val.(*int64)))
			case "float64Option":
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
