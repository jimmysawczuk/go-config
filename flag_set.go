package config

import (
	"fmt"
	"reflect"
)

type errUndefinedFlag struct {
	name string
}

type FlagSet struct {
	name string

	args     []string
	unparsed []string
	notset   []string

	triggerErrorUndefined bool

	help_flag bool
}

func NewFlagSet(name string, args []string) (f FlagSet) {
	f.name = name
	f.unparsed = args
	f.triggerErrorUndefined = true
	return
}

func (f *FlagSet) Parse() error {
	for {
		seen, err := f.parseOne()
		if seen {
			continue
		}
		if err == nil {
			break
		}

		if _, ok := err.(errUndefinedFlag); ok && f.triggerErrorUndefined {
			return err
		}
	}

	// fmt.Println("args", f.args)
	// fmt.Println("unparsed", f.unparsed)
	// fmt.Println("notset", f.notset)

	// pr, _ := json.MarshalIndent(baseOptionSet.Export(), "", "  ")
	// fmt.Println(string(pr))

	return nil
}

func (f *FlagSet) parseOne() (seen bool, err error) {
	if len(f.unparsed) == 0 {
		return false, nil
	}

	arg := f.unparsed[0]
	if len(arg) == 0 || arg[0] != '-' || len(arg) == 1 {
		if len(f.unparsed) > 0 {
			f.args = f.unparsed[0:]
			f.unparsed = []string{}
		}
		return false, nil
	}

	num_minuses := 1
	if arg[1] == '-' {
		num_minuses++
		if len(arg) == 2 {
			if len(f.unparsed) > 0 {
				f.args = f.unparsed[1:]
				f.unparsed = []string{}
			}

			return false, nil
		}
	}

	name := arg[num_minuses:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return false, fmt.Errorf("bad flag syntax: %s", arg)
	}

	has_value := false
	value := ""
	for i := 1; i < len(name); i++ { // equals cannot be first
		if name[i] == '=' {
			value = name[i+1:]
			has_value = true
			name = name[0:i]
			break
		}
	}

	// the help flag is special, so we'll stop parsing flags
	if name == "help" || name == "h" {
		f.help_flag = true
		return false, nil
	}

	if option, exists := baseOptionSet[name]; exists {
		if has_value {
			f.unparsed = f.unparsed[1:]
			// the option exists, and we have a value, so we can set it
			err := option.SetFromString(value)
			if err != nil {
				return true, fmt.Errorf("Error setting option %s to %s: %s", name, value, err)
			}

		} else if option.Type.Kind() == reflect.Bool {
			// don't need a value, and we're not allowed to use two args, so we can set the value to true normally and continue
			f.unparsed = f.unparsed[1:]
			option.Value = true
			return true, nil
		} else {
			// we need a value and don't have one yet, so we need to check the next argument
			if !has_value && len(f.unparsed) > 0 {
				// value is the next arg
				has_value = true
				value = f.unparsed[0]
				f.unparsed = f.unparsed[1:]

				err := option.SetFromString(value)
				if err != nil {
					return true, fmt.Errorf("Error setting option %s to %s: %s", name, value, err)
				}
			}

			if !has_value {
				return false, fmt.Errorf("flag needs an argument: -%s", name)
			}
		}
	} else {
		f.unparsed = f.unparsed[1:]

		if has_value {
			f.notset = append(f.notset, fmt.Sprintf("-%s=%s", name, value))
		} else if len(f.unparsed) > 0 {
			value = f.unparsed[0]
			f.unparsed = f.unparsed[1:]
			f.notset = append(f.notset, fmt.Sprintf("-%s", name), fmt.Sprintf("%s", value))
		}

		return true, errUndefinedFlag{name: name}
	}

	return true, nil
}

func (f FlagSet) ShowHelp() bool {
	return f.help_flag
}

func (f FlagSet) Release() []string {
	args := []string{f.name}
	args = append(args, f.notset...)
	args = append(args, f.args...)
	return args
}

func (e errUndefinedFlag) Error() string {
	return fmt.Sprintf("Undefined flag: -%s", e.name)
}
