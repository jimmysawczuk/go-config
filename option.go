package config

import (
	"fmt"
	"reflect"
	"strconv"
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

	// Validate is true if the option is required
	Validate bool

	// Filters is a set of boolean functions that are tested with the given value. If Validate is true, all of these must succeed.
	Filters []OptionFilterFunc

	// SortOrder controls the sort order of Options when displayed in Usage(). Defaults to 0; ties are resolved alphabetically.
	SortOrder int
}

// OptionFilterFunc is a function type that takes an *Option as a parameter. It returns true, nil if the *Option passes the filter, and false, error with a reason why if it didn't.
type OptionFilterFunc func(*Option) (bool, error)

// DefaultOptionMeta returns the default OptionMeta object
var DefaultOptionMeta = OptionMeta{
	Exportable: false,
	Validate:   true,
	Filters:    []OptionFilterFunc{},
	SortOrder:  0,
}

// String creates an Option with the parameters given of type string
func String(name string, defaultValue string, description string) *Option {

	v := Option{
		Name:        name,
		Description: description,

		DefaultValue: reflect.ValueOf(defaultValue),
		Value:        defaultValue,
		Type:         reflect.TypeOf(defaultValue),

		Options: DefaultOptionMeta,
	}

	return &v
}

// Bool creates an Option with the parameters given of type bool
func Bool(name string, defaultValue bool, description string) *Option {
	v := Option{
		Name:        name,
		Description: description,

		DefaultValue: reflect.ValueOf(defaultValue),
		Value:        defaultValue,
		Type:         reflect.TypeOf(defaultValue),

		Options: DefaultOptionMeta,
	}

	return &v
}

// Int creates an Option with the parameters given of type int64
func Int(name string, defaultValue int64, description string) *Option {
	v := Option{
		Name:        name,
		Description: description,

		DefaultValue: reflect.ValueOf(defaultValue),
		Value:        defaultValue,
		Type:         reflect.TypeOf(defaultValue),

		Options: DefaultOptionMeta,
	}

	return &v
}

// Float creates an Option with the parameters given of type float64
func Float(name string, defaultValue float64, description string) *Option {
	v := Option{
		Name:        name,
		Description: description,

		DefaultValue: reflect.ValueOf(defaultValue),
		Value:        defaultValue,
		Type:         reflect.TypeOf(defaultValue),

		Options: DefaultOptionMeta,
	}

	return &v
}

// Enum creates an Option with the parameters given of type string and a built-in validation to make sure
// that the parsed Option value is contained within the possibleValues argument.
func Enum(name string, possibleValues []string, defaultValue string, description string) *Option {

	v := Option{
		Name:        name,
		Description: description,

		DefaultValue: reflect.ValueOf(defaultValue),
		Value:        defaultValue,
		Type:         reflect.TypeOf(defaultValue),

		Options: DefaultOptionMeta,
	}

	v.
		Validate(true).
		AddFilter(IsOneOfStrings(possibleValues))

	return &v
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

// defaultValueString returns the Option's default value as a string. If that value resolves to "", it'll return the
// emptyReplacement argument instead.
func (o Option) defaultValueString(emptyReplacement string) string {
	v := o.DefaultValue.(reflect.Value)

	ret := ""

	switch o.Type.Kind() {
	case reflect.String:
		ret = fmt.Sprintf(`%v`, v.String())
	case reflect.Int64:
		ret = fmt.Sprintf(`%v`, v.Int())
	case reflect.Float64:
		ret = fmt.Sprintf(`%v`, v.Float())
	case reflect.Bool:
		ret = fmt.Sprintf(`%v`, v.Bool())
	}

	if ret == "" {
		ret = emptyReplacement
	}

	return ret
}

// DefaultValueString returns the Option's default value as a string.
func (o Option) DefaultValueString() string {
	return o.defaultValueString("")
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

// Exportable sets whether or not the Option is exportable to a config file.
func (o *Option) Exportable(v bool) *Option {
	o.Options.Exportable = v
	return o
}

// Validate sets whether or not the Filters on the Option will be tested for validity before being accepted.
func (o *Option) Validate(v bool) *Option {
	o.Options.Validate = v
	return o
}

// AddFilter adds an OptionFilterFunc to the Option's filter set. It also sets Validate to true.
func (o *Option) AddFilter(v OptionFilterFunc) *Option {
	o.Options.Filters = append(o.Options.Filters, v)
	o.Options.Validate = true
	return o
}

// SortOrder sets the sort order on the Option used in Usage().
func (o *Option) SortOrder(i int) *Option {
	o.Options.SortOrder = i
	return o
}
