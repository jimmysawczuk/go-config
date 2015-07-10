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

func (j *jsonConfigMap) UnmarshalJSON(in []byte) (err error) {
	return json.Unmarshal(in, &j.config)
}

func (j *jsonConfigMap) Parse() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %#v\n", r)
			j.err = fmt.Errorf("%s", r)
		}
	}()

	err := parse(j.config, "")
	if jerr, ok := err.(jsonConfigMapParseErrorList); ok {
		return jerr
	}
	return err
}

type jsonConfigMapError interface {
	error
}

type jsonConfigMapParseError struct {
	key      string
	got      reflect.Kind
	expected reflect.Kind
}

func (j jsonConfigMapParseError) Error() string {
	return fmt.Sprintf("unexpected type: %q: expected %s, got %s", j.key, j.expected, j.got)
}

type jsonConfigMapTruncateError struct {
	key        string
	got        reflect.Kind
	expected   reflect.Kind
	difference float64
}

func (j jsonConfigMapTruncateError) Error() string {
	return fmt.Sprintf("possible truncate: %q: expected %s, got %s; difference: %e", j.key, j.expected, j.got, j.difference)
}

type jsonConfigMapParseErrorList []jsonConfigMapError

func (j *jsonConfigMapParseErrorList) Merge(inc jsonConfigMapParseErrorList) {
	*j = append(*j, inc...)
}

func (j jsonConfigMapParseErrorList) Error() string {
	strs := make([]string, len(j))
	for i, err := range j {
		strs[i] = "  " + err.Error()
	}

	return fmt.Sprintf("json parse error(s): %s", strings.Join(strs, ", "))
}

func (j jsonConfigMapParseErrorList) Len() int {
	return len(j)
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
				cerr := parse(v.(map[string]interface{}), prefix+k+".")
				if cerr != nil {
					if childerrs, ok := cerr.(jsonConfigMapParseErrorList); ok {
						errs.Merge(childerrs)
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
