package config

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

type jsonConfigMap struct {
	scope  string
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

	err := parse(j.scope, j.config, "")
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
	got      interface{}
	expected Type
}

func (j jsonConfigMapParseError) Error() string {
	return fmt.Sprintf("unexpected type: %q: expected %s, got %T", j.key, j.expected, j.got)
}

type jsonConfigMapTruncateError struct {
	key        string
	got        interface{}
	expected   Type
	difference float64
}

func (j jsonConfigMapTruncateError) Error() string {
	return fmt.Sprintf("possible truncate: %q: expected %s, got %T; difference: %e", j.key, j.expected, j.got, j.difference)
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

func parse(scope string, configMap map[string]interface{}, prefix string) (err error) {

	errs := make(jsonConfigMapParseErrorList, 0)

	for k, v := range configMap {
		s, exists := baseOptionSet.Get(prefix + k)
		if exists {
			err := parseElem(scope, s, prefix+k, v)
			if err != nil {
				errs = append(errs, err)
			}
		} else {
			switch v.(type) {
			case map[string]interface{}:
				cerr := parse(scope, v.(map[string]interface{}), prefix+k+".")
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

func parseElem(scope string, opt *Option, key string, v interface{}) error {
	switch v.(type) {

	case float64:
		if opt.Type == FloatType {
			opt.Value = v.(float64)
		} else if opt.Type == IntType {
			opt.Value = int64(v.(float64))
			rounded := math.Floor(v.(float64) + 0.5)
			diff := math.Abs(rounded - v.(float64))
			if diff > 1e-32 {
				return jsonConfigMapTruncateError{
					key:        key,
					got:        v,
					expected:   opt.Type,
					difference: diff,
				}
			}
		} else {
			return jsonConfigMapParseError{
				key:      key,
				got:      v,
				expected: opt.Type,
			}
		}
	case bool:
		if opt.Type == BoolType {
			opt.Value = v.(bool)
		} else {
			return jsonConfigMapParseError{
				key:      key,
				got:      v,
				expected: opt.Type,
			}
		}
	case string:
		if opt.Type == StringType {
			opt.Value = v.(string)
		} else {
			return jsonConfigMapParseError{
				key:      key,
				got:      v,
				expected: opt.Type,
			}
		}
	}

	opt.AddScope(scope)

	return nil
}
