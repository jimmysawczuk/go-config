package config

import (
	"fmt"
	"strconv"
)

// Type is a string representing the type of data stored by an Option
type Type string

// These Types constants parallel their standard counterparts, and are the four elementary types that come
// when unmarshaling JSON
const (
	BoolType   Type = "bool"
	StringType      = "string"
	FloatType       = "float64"
	IntType         = "int64"
	CustomType      = "custom"
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
	Type Type

	// Extra options
	Options OptionMeta

	overridden bool
	scopes     []string
	isBuiltIn  bool
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

// Str creates an Option with the parameters given of type string
func Str(name string, defaultValue string, description string) *Option {

	v := Option{
		Name:        name,
		Description: description,

		DefaultValue: defaultValue,
		Value:        defaultValue,
		Type:         StringType,

		Options: DefaultOptionMeta,
	}

	return &v
}

// Bool creates an Option with the parameters given of type bool
func Bool(name string, defaultValue bool, description string) *Option {
	v := Option{
		Name:        name,
		Description: description,

		DefaultValue: defaultValue,
		Value:        defaultValue,
		Type:         BoolType,

		Options: DefaultOptionMeta,
	}

	return &v
}

// Int creates an Option with the parameters given of type int64
func Int(name string, defaultValue int64, description string) *Option {
	v := Option{
		Name:        name,
		Description: description,

		DefaultValue: defaultValue,
		Value:        defaultValue,
		Type:         IntType,

		Options: DefaultOptionMeta,
	}

	return &v
}

// Float creates an Option with the parameters given of type float64
func Float(name string, defaultValue float64, description string) *Option {
	v := Option{
		Name:        name,
		Description: description,

		DefaultValue: defaultValue,
		Value:        defaultValue,
		Type:         FloatType,

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

		DefaultValue: defaultValue,
		Value:        defaultValue,
		Type:         StringType,

		Options: DefaultOptionMeta,
	}

	v.
		Validate(true).
		AddFilter(IsOneOfStrings(possibleValues))

	return &v
}

// DebugString returns a string describing some attributes about the Option, including the name, value, type and what scopes it came from.
func (o Option) DebugString() string {
	return fmt.Sprintf(`name: %s, value: %v, type: %s, scopes: %s`, o.Name, o.Value, o.Type, o.scopes)
}

// String implements fmt.Stringer. This is used for printing the OptionSet if needed; you should use Str() to
// return the string value of a string Option, as it'll return what you expect all the time.
func (o Option) String() string {
	return fmt.Sprintf(`%v`, o.Value)
}

// Str returns the string value of the option. Will panic if the Option's type is not a string.
func (o Option) Str() string {
	return o.Value.(string)
}

// Bool returns the bool value of the option. Will panic if the Option's type is not a bool.
func (o Option) Bool() bool {
	return o.Value.(bool)
}

// Float returns the float64 value of the option. Will panic if the Option's type is not a float64.
func (o Option) Float() float64 {
	return o.Value.(float64)
}

// Int returns the int64 value of the option. Will panic if the Option's type not an int64.
func (o Option) Int() int64 {
	return o.Value.(int64)
}

// defaultValueString returns the Option's default value as a string. If that value resolves to "", it'll return the
// emptyReplacement argument instead.
func (o Option) defaultValueString(emptyReplacement string) string {
	ret := fmt.Sprintf(`%v`, o.DefaultValue)

	if ret == "" {
		ret = emptyReplacement
	}

	return ret
}

// AddScope adds a scope to an Option indicating that it was parsed in a file with the given scope.
func (o *Option) AddScope(s string) {
	if o.scopes == nil {
		o.scopes = make([]string, 0)
	}

	o.scopes = append(o.scopes, s)
}

// HasScope returns true if the Option has the specified scope.
func (o *Option) HasScope(s string) bool {
	for _, v := range o.scopes {
		if v == s {
			return true
		}
	}
	return false
}

// DefaultValueString returns the Option's default value as a string.
func (o Option) DefaultValueString() string {
	return o.defaultValueString("")
}

// SetFromFlagValue attempts to set the Option's value as its proper type by parsing the string argument, and also
// sets a hidden value on the Option indicating it was overridden by a flag argument.
func (o *Option) SetFromFlagValue(val string) (err error) {
	err = o.SetFromString(val)
	if err != nil {
		return err
	}

	o.overridden = true
	o.AddScope("flag")
	return nil
}

// SetFromString attempts to set the Option's value as its proper type by parsing the string argument
func (o *Option) SetFromString(val string) (err error) {
	switch o.Type {
	case StringType:
		o.Value = val

	case IntType:
		v, err := strconv.ParseInt(val, 0, 64)
		if err == nil {
			o.Value = v
		}

	case FloatType:
		v, err := strconv.ParseFloat(val, 64)
		if err == nil {
			o.Value = v
		}

	case BoolType:
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

func (o *Option) builtIn() *Option {
	o.isBuiltIn = true
	return o
}
