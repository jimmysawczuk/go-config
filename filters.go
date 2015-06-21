package config

import (
	"fmt"
)

func validEnum(possibleValues []string) func(*Option) (bool, error) {
	return func(v *Option) (bool, error) {
		val := v.String()
		for _, s := range possibleValues {
			if val == s {
				return true, nil
			}
		}
		return false, fmt.Errorf("invalid value for enum: %s", val)
	}
}

func validString() func(*Option) (bool, error) {
	return func(v *Option) (bool, error) {
		if !v.Options.Required {
			return true, nil
		}

		s := v.String()
		if s != "" {
			return true, nil
		} else {
			return false, fmt.Errorf("empty string")
		}
	}
}
