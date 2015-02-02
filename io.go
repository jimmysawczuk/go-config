package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type IO interface {
	Read() (map[string]interface{}, error)
	Write() error
}

type FileIO struct {
	Filename string
}

func (this FileIO) Write() error {
	json, err := json.MarshalIndent(baseOptionSet.Export(), "", "    ")
	if err != nil {
		return fmt.Errorf("go-config: error marshaling config: %s", err)
	}

	fp, err := os.OpenFile(this.Filename, os.O_RDWR+os.O_CREATE+os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("go-config: file i/o open error: %s", err)
	}
	defer fp.Close()

	n, err := fp.Write(json)
	if err != nil || n < len(json) {
		return fmt.Errorf("go-config: file i/o write error: %s", err)
	}

	return nil
}

func (this *FileIO) Read() (err error) {
	fp, err := os.Open(this.Filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("go-config: file i/o open error: %s", err)
	}
	defer fp.Close()

	fi, err := fp.Stat()
	if err != nil {
		return fmt.Errorf("go-config: file i/o stat error: %s", err)
	}

	n := fi.Size()

	by := make([]byte, n)
	read, err := fp.Read(by)
	if err != nil || int64(read) < n {
		return fmt.Errorf("go-config: file i/o read error: %s", err)
	}

	j_map := jsonConfigMap{}
	err = json.Unmarshal(by, &j_map)
	if err != nil {
		return fmt.Errorf("go-config: json unmarshal error: %s", err)
	}

	err = j_map.Parse()
	if err != nil {
		return fmt.Errorf("go-config: config parse error: %s", err)
	}

	return nil
}
