package config

import (
	"fmt"
	"reflect"
	"strconv"
)

// Holds information for a configuration option
type Option struct {
	// The name of the option is what's used to reference the option and its value during the program
	Name string

	// What the option is for. This also shows up when invoking `program --help`.
	Description string

	// Holds the value contained by this option
	Value interface{}

	// Holds the default value for this option
	DefaultValue interface{}

	// Holds the type of this option
	Type reflect.Type

	// If true, this option is exportable to a config.json file if generated.
	Exportable bool

	// A pointer to the flag value to set this option from the command line.
	flag interface{}
}

// Returns the string value of the option. Will panic if it's not a string.
func (this Option) String() string {
	return reflect.ValueOf(this.Value).String()
}

// Returns the bool value of the option. Will panic if it's not a bool.
func (this Option) Bool() bool {
	return reflect.ValueOf(this.Value).Bool()
}

// Returns the float64 value of the option. Will panic if it's not a float64.
func (this Option) Float() float64 {
	return reflect.ValueOf(this.Value).Float()
}

// Returns the int64 value of the option. Will panic if it's not an int64.
func (this Option) Int() int64 {
	return reflect.ValueOf(this.Value).Int()
}

func (this Option) DefaultValueString() string {
	v := this.DefaultValue.(reflect.Value)

	switch this.Type.Kind() {
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

func (this *Option) SetFromString(val string) (err error) {
	switch this.Type.Kind() {
	case reflect.String:
		this.Value = val

	case reflect.Int64:
		v, err := strconv.ParseInt(val, 0, 64)
		if err == nil {
			this.Value = v
		}

	case reflect.Float64:
		v, err := strconv.ParseFloat(val, 64)
		if err == nil {
			this.Value = v
		}

	case reflect.Bool:
		switch val {
		case "1", "t", "T", "true", "TRUE", "True":
			this.Value = true
		case "0", "f", "F", "false", "FALSE", "False":
			this.Value = false
		default:
			err = fmt.Errorf("Invalid boolean value: %s", val)
		}
	}

	return
}

type SortedOptionSlice []Option

func (this SortedOptionSlice) Less(a, b int) bool { return this[a].Name < this[b].Name }
func (this SortedOptionSlice) Swap(a, b int)      { this[a], this[b] = this[b], this[a] }
func (this SortedOptionSlice) Len() int           { return len(this) }
