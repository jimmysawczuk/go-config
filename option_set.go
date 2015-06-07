package config

import (
	"fmt"
	"strings"
)

// An OptionSet is map of Options, keyed by the Options' Names.
type OptionSet map[string]*Option

// Export returns a map that's suitable for pushing into a config.json file.
func (os OptionSet) Export() map[string]interface{} {
	return os.export(false)
}

func (os OptionSet) export(includeAll bool) map[string]interface{} {
	tbr := make(map[string]interface{})
	for _, v := range os {
		if v.Options.Exportable || includeAll {
			parts := strings.Split(v.Name, ".")
			var i int
			var cursor = &tbr

			for i < len(parts)-1 {

				if _, exists := (*cursor)[parts[i]]; !exists {
					(*cursor)[parts[i]] = make(map[string]interface{})
				}

				v := (*cursor)[parts[i]].(map[string]interface{})
				cursor = &v

				i++
			}
			(*cursor)[parts[i]] = v.Value
		}
	}
	return tbr
}

// Add adds an Option to an OptionSet with a key of the Option's name.
func (os OptionSet) Add(o Option) {
	os[o.Name] = &o
}

// Get retrieves an Option with the Name of key, and a boolean to determine if it was found or not.
func (os OptionSet) Get(key string) (*Option, bool) {
	result, exists := os[key]
	return result, exists
}

// Require retrieves an Option with the Name of key. Panics if there was no key found.
func (os OptionSet) Require(key string) *Option {
	result, exists := os.Get(key)
	if !exists {
		panic("Option with name " + key + " doesn't exist")
	}
	return result
}

func (os OptionSet) Validate() error {
	hasError := false
	invalidOpts := []string{}
	for _, v := range os {
		if v.Options.Required {
			validOption := true
			for _, f := range v.Options.Filters {
				validOption = validOption && f(v)
			}
			if !validOption {
				invalidOpts = append(invalidOpts, v.Name)
			}
			hasError = hasError && validOption
		}
	}

	if hasError {
		return fmt.Errorf("Some options were empty or invalid: %s", strings.Join(invalidOpts, ", "))
	}
	return nil
}
