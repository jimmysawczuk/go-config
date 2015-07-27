package config

import (
	"fmt"
	"strings"
)

// IsOneOfStrings returns an OptionFilterFunc that checks the Option value against a list of
// string values and returns true if the Option value matches one of the possible values
func IsOneOfStrings(possibleValues []string) OptionFilterFunc {
	return func(v *Option) (bool, error) {
		val := v.String()
		for _, s := range possibleValues {
			if val == s {
				return true, nil
			}
		}

		return false, fmt.Errorf("%s is not a possible value (try one of %s)", val, strings.Join(possibleValues, ", "))
	}
}

// NonEmptyString returns an OptionFilterFunc that returns true if the Option value is a non-empty
// string. It will also return false if the Option is not a string.
func NonEmptyString() OptionFilterFunc {
	return func(v *Option) (bool, error) {
		s := v.String()
		if s == "" {
			return false, fmt.Errorf("value cannot be an empty string")
		}

		return true, nil
	}
}
