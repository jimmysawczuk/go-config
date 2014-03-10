package config

import (
	"encoding/json"
	"os"
	"reflect"
)

func exportConfigToFile(filename string) error {
	json, _ := json.MarshalIndent(baseOptionSet.Export(), "", "    ")
	fp, err := os.OpenFile(filename, os.O_RDWR+os.O_CREATE+os.O_TRUNC, 0666)

	if err == nil {
		defer fp.Close()

		fp.Write(json)
	} else {
		return err
	}

	return nil
}

func loadConfigFile(filename string) (configMap map[string]interface{}, err error) {
	fp, err := os.Open(filename)
	if err != nil {
		return
	}

	defer fp.Close()

	var n int64 = 0

	if fi, err := fp.Stat(); err == nil {
		if size := fi.Size(); size < 1e9 {
			n = size
		}
	}

	by := make([]byte, n)
	fp.Read(by)

	configMap = make(map[string]interface{})
	err = json.Unmarshal(by, &configMap)

	return
}

func importConfigFile(filename string) error {

	configMap, err := loadConfigFile(filename)
	if err != nil {
		return err
	}

	parseConfigFileMap(configMap, "")
	return nil
}

func parseConfigFileMap(configMap map[string]interface{}, prefix string) {
	for k, v := range configMap {
		s, exists := baseOptionSet.Get(prefix + k)
		if exists {
			switch s.Type.Kind() {
			case reflect.Int64:
				x := int64(reflect.ValueOf(v).Float())
				s.Value = x

			case reflect.Float64:
				x := reflect.ValueOf(v).Float()
				s.Value = x

			case reflect.Bool:
				x := reflect.ValueOf(v).Bool()
				s.Value = x

			case reflect.String:
				x := reflect.ValueOf(v).String()
				s.Value = x
			}
		} else {
			switch v.(type) {
			case map[string]interface{}:
				parseConfigFileMap(v.(map[string]interface{}), k+".")
			}
		}
	}
}
