package config

func validEnum(possibleValues []string) func(*Option) bool {
	return func(v *Option) bool {
		val := v.String()
		for _, s := range possibleValues {
			if val == s {
				return true
			}
		}
		return false
	}
}

func validString() func(*Option) bool {
	return func(v *Option) bool {
		if v.Options.Required {
			s := v.String()
			return s != ""
		} else {
			return true
		}
	}
}
