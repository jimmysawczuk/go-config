package config

import (
	"reflect"
)

type Option struct {
	Name        string
	Description string

	Value        interface{}
	DefaultValue interface{}
	Type         reflect.Type

	Exportable bool

	flag interface{}
}

func (this Option) String() string {
	return reflect.ValueOf(this.Value).String()
}

func (this Option) Bool() bool {
	return reflect.ValueOf(this.Value).Bool()
}

func (this Option) Float() float64 {
	return reflect.ValueOf(this.Value).Float()
}

func (this Option) Int() int64 {
	return reflect.ValueOf(this.Value).Int()
}
