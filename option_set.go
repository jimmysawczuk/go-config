package config

import (
	"fmt"
	"strings"
)

func init() {
	_ = fmt.Printf
}

type OptionSet map[string]*Option

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

func (os OptionSet) Get(key string) (*Option, bool) {
	result, exists := os[key]
	return result, exists
}
