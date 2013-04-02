package config

import (
	"encoding/json"
	"os"
)

func exportConfigToFile(filename string) error {
	json, _ := json.MarshalIndent(optionsDict.Export(), "", "    ")

	fp, err := os.OpenFile(filename, os.O_RDWR+os.O_CREATE+os.O_TRUNC, 0666)

	if err == nil {
		defer fp.Close()

		fp.Write(json)
	} else {
		return err
	}

	return nil
}

func importConfigFile(filename string) error {
	fp, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer fp.Close()

	var n int64 = 0

	if fi, err := fp.Stat(); err == nil {
		if size := fi.Size(); size < 1e9 {
			n = size
		}
	}

	bytes := make([]byte, n)
	fp.Read(bytes)

	configMap := make(map[string]interface{})

	err = json.Unmarshal(bytes, &configMap)

	if err != nil {
		return err
	}

	for k, v := range configMap {
		s, exists := optionsDict.Get(k)

		if exists {
			s.Set(v)
		}
	}

	return nil
}
