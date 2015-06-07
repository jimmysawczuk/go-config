package config

import (
	"fmt"
	"reflect"
	"strconv"
)

// OptionMask is a type for setting ConfigOptions on Options
type OptionMask int

const (
	_ OptionMask = iota

	// Exportable means the Option will be exported to a config.json file
	Exportable

	// Required means the Option is required and all of the Filters will be tested
	Required
)

// Option holds information for a configuration option
type Option struct {
	// The name of the option is what's used to reference the option and its value during the program
	Name string

	// What the option is for. This also shows up when invoking `program --help`.
	Description string

	// Holds the actual value contained by this option
	Value interface{}

	// Holds the default value for this option
	DefaultValue interface{}

	// Holds the type of this option
	Type reflect.Type

	// Extra options
	Options OptionMeta
}

// OptionMeta holds information for configuring options on Options
type OptionMeta struct {
	// Exportable is true if the option is exportable to a config.json file
	Exportable bool

	// Required is true if the option is required
	Required bool

	// Filters is a set of boolean functions that are tested with the given value. If Required is true, all of these must succeed.
	Filters []func(option *Option) bool
}

// String returns the string value of the option. Will panic if the Option's type is not a string.
func (o Option) String() string {
	return reflect.ValueOf(o.Value).String()
}

// Bool returns the bool value of the option. Will panic if the Option's type is not a bool.
func (o Option) Bool() bool {
	return reflect.ValueOf(o.Value).Bool()
}

// Float returns the float64 value of the option. Will panic if the Option's type is not a float64.
func (o Option) Float() float64 {
	return reflect.ValueOf(o.Value).Float()
}

// Int returns the int64 value of the option. Will panic if the Option's type not an int64.
func (o Option) Int() int64 {
	return reflect.ValueOf(o.Value).Int()
}

// DefaultValueString returns the Option's default value as a string
func (o Option) DefaultValueString() string {
	v := o.DefaultValue.(reflect.Value)

	switch o.Type.Kind() {
	case reflect.String:
		return fmt.Sprintf(`%v`, v.String())
	case reflect.Int64:
		return fmt.Sprintf(`%v`, v.Int())
	case reflect.Float64:
		return fmt.Sprintf(`%v`, v.Float())
	case reflect.Bool:
		return fmt.Sprintf(`%v`, v.Bool())
	}

	return ""
}

// SetFromString attempts to set the Option's value as its proper type by parsing the string argument
func (o *Option) SetFromString(val string) (err error) {
	switch o.Type.Kind() {
	case reflect.String:
		o.Value = val

	case reflect.Int64:
		v, err := strconv.ParseInt(val, 0, 64)
		if err == nil {
			o.Value = v
		}

	case reflect.Float64:
		v, err := strconv.ParseFloat(val, 64)
		if err == nil {
			o.Value = v
		}

	case reflect.Bool:
		switch val {
		case "1", "t", "T", "true", "TRUE", "True":
			o.Value = true
		case "0", "f", "F", "false", "FALSE", "False":
			o.Value = false
		default:
			err = fmt.Errorf("Invalid boolean value: %s", val)
		}
	}

	return
}

type sortedOptionSlice []Option

func (s sortedOptionSlice) Less(a, b int) bool { return s[a].Name < s[b].Name }
func (s sortedOptionSlice) Swap(a, b int)      { s[a], s[b] = s[b], s[a] }
func (s sortedOptionSlice) Len() int           { return len(s) }

// DefaultOptionMeta returns the default OptionMeta object
func DefaultOptionMeta() OptionMeta {
	s := OptionMeta{
		Exportable: false,
		Required:   false,
		Filters:    []func(*Option) bool{},
	}
	return s
}

func getOptionMetaFromMask(optMask OptionMask) OptionMeta {
	s := DefaultOptionMeta()
	s.Exportable = optMask%2 == 1
	s.Required = (optMask>>1)%2 == 1

	// fmt.Printf("exportable %t required %t\n", s.Exportable, s.Required)

	return s
}
