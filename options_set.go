package config

type OptionsSet struct {
	options map[string]*Option
}

func (os OptionsSet) Export() map[string]interface{} {
	tbr := make(map[string]interface{})
	for _, v := range os.options {
		if v.exportable {
			tbr[v.Name] = v.val
		}
	}
	return tbr
}

func (os *OptionsSet) Add(o Option) {
	os.options[o.Name] = &o
}

func (os OptionsSet) Get(key string) (*Option, bool) {
	result, exists := os.options[key]
	return result, exists
}
