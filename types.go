package config

type stringOption Option

func (s stringOption) DefaultValue() string {
	return s.default_value.(string)
}

func (s stringOption) Get() string {
	return s.val.(string)
}

//////

type boolOption Option

func (b boolOption) DefaultValue() bool {
	return b.default_value.(bool)
}

func (b boolOption) Get() bool {
	return b.val.(bool)
}

//////

type intOption Option

func (i intOption) DefaultValue() int {
	return i.default_value.(int)
}

func (i intOption) Get() int {
	return i.val.(int)
}

//////

type int64Option Option

func (i64 int64Option) DefaultValue() int64 {
	return i64.default_value.(int64)
}

func (i64 int64Option) Get() int64 {
	return i64.val.(int64)
}

//////

type float64Option Option

func (f float64Option) DefaultValue() float64 {
	return f.default_value.(float64)
}

func (f float64Option) Get() float64 {
	return f.val.(float64)
}
