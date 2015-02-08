package config

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"
)

type jsonConfigMap struct {
	config map[string]interface{}
	err    error
}

func (this *jsonConfigMap) UnmarshalJSON(in []byte) (err error) {
	return json.Unmarshal(in, &this.config)
}

func (this *jsonConfigMap) Parse() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%#v\n", r)
			this.err = fmt.Errorf("%s", r)
		}
	}()

	return parse(this.config, "")
}

type jsonConfigMapError interface {
	error
}

type jsonConfigMapParseError struct {
	key      string
	got      reflect.Kind
	expected reflect.Kind
}

func (this jsonConfigMapParseError) Error() string {
	return fmt.Sprintf("unexpected type: %q: expected %s, got %s", this.key, this.expected, this.got)
}

type jsonConfigMapTruncateError struct {
	key        string
	got        reflect.Kind
	expected   reflect.Kind
	difference float64
}

func (this jsonConfigMapTruncateError) Error() string {
	return fmt.Sprintf("possible truncate: %q: expected %s, got %s; difference: %e", this.key, this.expected, this.got, this.difference)
}

type jsonConfigMapParseErrorList []jsonConfigMapError

func (this *jsonConfigMapParseErrorList) Merge(inc jsonConfigMapParseErrorList) {
	*this = append(*this, inc...)
}

func (this jsonConfigMapParseErrorList) Error() string {
	strs := make([]string, len(this))
	for i, err := range this {
		strs[i] = "  " + err.Error()
	}

	return fmt.Sprintf("json parse error(s):\n%s", strings.Join(strs, ", "))
}

func (this jsonConfigMapParseErrorList) Len() int {
	return len(this)
}

func parse(configMap map[string]interface{}, prefix string) (err error) {

	errs := make(jsonConfigMapParseErrorList, 0)

	for k, v := range configMap {
		s, exists := baseOptionSet.Get(prefix + k)
		if exists {
			val := reflect.ValueOf(v)
			switch val.Kind() {
			case reflect.Float64:
				if s.Type.Kind() == reflect.Float64 {
					s.Value = val.Float()
				} else if s.Type.Kind() == reflect.Int64 {
					s.Value = int64(val.Float())
					rounded := math.Floor(val.Float() + 0.5)
					diff := math.Abs(rounded - val.Float())
					if diff > 1e-32 {
						errs = append(errs, jsonConfigMapTruncateError{
							key:        prefix + k,
							got:        val.Kind(),
							expected:   s.Type.Kind(),
							difference: diff,
						})
					}
				} else {
					errs = append(errs, jsonConfigMapParseError{
						key:      prefix + k,
						got:      val.Kind(),
						expected: s.Type.Kind(),
					})
				}
			case reflect.Bool:
				if s.Type.Kind() == reflect.Bool {
					s.Value = val.Bool()
				} else {
					errs = append(errs, jsonConfigMapParseError{
						key:      prefix + k,
						got:      val.Kind(),
						expected: s.Type.Kind(),
					})
				}
			case reflect.String:
				if s.Type.Kind() == reflect.String {
					s.Value = val.String()
				} else {
					errs = append(errs, jsonConfigMapParseError{
						key:      prefix + k,
						got:      val.Kind(),
						expected: s.Type.Kind(),
					})
				}
			}
		} else {
			switch v.(type) {
			case map[string]interface{}:
				cerr := parse(v.(map[string]interface{}), k+".")
				if cerr != nil {
					if child_errs, ok := cerr.(jsonConfigMapParseErrorList); ok {
						errs.Merge(child_errs)
					}
				}
			}
		}
	}

	if errs.Len() > 0 {
		return errs
	}

	return nil
}