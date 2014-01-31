package config

import (
	"fmt"
	"reflect"
)

type Option struct {
	Name        string
	Description string

	opt_type reflect.Type
	val      reflect.Value

	exportable    bool
	default_value interface{}
	flag_val      interface{}
}

func (o *Option) Set(v interface{}) error {
	switch o.opt_type.Name() {
	case "string":
		switch v.(type) {
		case string:
			o.val = reflect.ValueOf(v.(string))
		case fmt.Stringer:
			o.val = reflect.ValueOf(v.(fmt.Stringer).String())
		default:
			return fmt.Errorf("Invalid data passed to *Option.Set(): %v (type %T), expected string or fmt.Stringer")
		}

	case "bool":
		switch v.(type) {
		case bool:
			o.val = reflect.ValueOf(v.(bool))
		default:
			return fmt.Errorf("Invalid data passed to *Option.Set(): %v (type %T), expected bool")
		}

	case "int":
		switch v.(type) {
		case int:
			o.val = reflect.ValueOf(v.(int))
		case float64:
			o.val = reflect.ValueOf(int(v.(float64)))
		default:
			return fmt.Errorf("Invalid data passed to *Option.Set(): %v (type %T), expected int or float64")
		}

	case "int64":
		switch v.(type) {
		case int:
			o.val = reflect.ValueOf(int64(v.(int)))
		case int64:
			o.val = reflect.ValueOf(v.(int64))
		case float64:
			o.val = reflect.ValueOf(int64(v.(float64)))
		default:
			return fmt.Errorf("Invalid data passed to *Option.Set(): %v (type %T), expected int64 or float64")
		}

	case "float64":
		switch v.(type) {
		case float64:
			o.val = reflect.ValueOf(v.(float64))
		case float32:
			o.val = reflect.ValueOf(float64(v.(float32)))
		default:
			return fmt.Errorf("Invalid data passed to *Option.Set(): %v (type %T), expected float or float64")
		}
	}

	return nil
}

func (o Option) Value() reflect.Value {
	return o.val
}

func (o Option) Bool() bool {
	switch o.opt_type.Name() {
	case "bool":
		return o.val.Bool()
	default:
		panic("Invalid type: got " + o.opt_type.Name() + ", expected bool!")
	}
}

func (o Option) String() string {
	switch o.opt_type.Name() {
	case "bool":
		if o.val.Bool() == true {
			return "true"
		} else {
			return "false"
		}
	case "string":
		return o.val.String()
	default:
		panic("Invalid type: got " + o.opt_type.Name() + ", expected string!")
	}
}

func (o Option) Float64() float64 {
	switch o.opt_type.Name() {
	case "float64":
		return o.val.Float()
	case "int":
		return float64(o.val.Int())
	case "int64":
		return float64(o.val.Int())
	default:
		panic("Invalid type: got " + o.opt_type.Name() + ", expected float64!")
	}
}
