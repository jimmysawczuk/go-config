package config

import (
	"fmt"
	"reflect"
)

var builtInFlags map[string]bool = map[string]bool{
	"config":          true,
	"config-export":   true,
	"config-generate": true,
}

type errUndefinedFlag struct {
	name string
}

type FlagSet struct {
	name string

	args     []string
	unparsed []string
	notset   []string

	help_flag bool
}

func NewFlagSet(name string, args []string) (f FlagSet) {
	f.name = name
	f.unparsed = args
	return
}

func (f *FlagSet) Parse() error {
	return f.parse(false)
}

func (f *FlagSet) ParseBuiltIn() error {
	return f.parse(true)
}

func (f *FlagSet) parse(built_in_only bool) error {
	for {
		seen, err := f.parseOne(built_in_only)
		if seen && err == nil {
			continue
		}

		if err == nil {
			break
		}

		if _, ok := err.(errUndefinedFlag); err != nil && !ok {
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

func (f *FlagSet) parseOne(built_in_only bool) (seen bool, err error) {
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

	if _, exists := builtInFlags[name]; built_in_only && !exists {
		// we saw a valid flag, but it's not the one we're looking for, so we move on
		f.unparsed = f.unparsed[1:]

		if has_value {
			f.notset = append(f.notset, fmt.Sprintf("-%s=%s", name, value))
		} else if len(f.unparsed) > 0 {
			value = f.unparsed[0]
			f.unparsed = f.unparsed[1:]
			f.notset = append(f.notset, fmt.Sprintf("-%s", name), fmt.Sprintf("%s", value))
		}

		return true, nil

	} else if option, exists := baseOptionSet[name]; exists {
		// valid flag, might need to find a value still
		f.unparsed = f.unparsed[1:]
		if has_value {
			// the option exists, and we have a value, so we can set it
			err := option.SetFromString(value)
			if err != nil {
				return true, fmt.Errorf("Error setting option %s to %s: %s", name, value, err)
			}

		} else if option.Type.Kind() == reflect.Bool {
			// don't need a value, and we're not allowed to use two args, so we can set the value to true normally and continue
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
			// this isn't the last argument, but we haven't seen a value yet, so we want to check the next argument to see if it's the value
			// for this flag
			f.notset = append(f.notset, fmt.Sprintf("-%s", name))
			next_arg := f.unparsed[0]
			if next_arg[0] != '-' {
				f.unparsed = f.unparsed[1:]
				f.notset = append(f.notset, next_arg)
			}
		} else {
			// this is the last argument, and it's a flag with no value, so we'll assume it's a boolean and pass it on through.
			f.notset = append(f.notset, fmt.Sprintf("-%s", name))
		}

		return true, errUndefinedFlag{name: name}
	}

	return true, nil
}

func (f FlagSet) HasHelpFlag() bool {
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
