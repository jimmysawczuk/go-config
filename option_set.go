package config

import (
	"strings"
)

// A map of Options, keyed by the Options' Names.
type OptionSet map[string]*Option

// Exports the OptionSet into a map that's suitable for pushing into a config.json file.
func (os OptionSet) Export() map[string]interface{} {
	tbr := make(map[string]interface{})
	for _, v := range os {
		if v.Exportable {
			parts := strings.Split(v.Name, ".")
			var i int = 0
			var cursor *map[string]interface{} = &tbr

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

func (os OptionSet) Add(o Option) {
	os[o.Name] = &o
}

// Retrieves an Option with the Name of key, and a boolean to determine if it was found or not.
func (os OptionSet) Get(key string) (*Option, bool) {
	result, exists := os[key]
	return result, exists
}
