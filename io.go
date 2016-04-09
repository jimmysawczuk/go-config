package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// IO defines an interface that allows reading and writing of OptionSets to external storage
type IO interface {
	Read() error
	Write() error
}

// FileIO implements IO and writes to the filesystem
type FileIO struct {
	Filename string
}

func (f FileIO) Write() error {
	json, err := json.MarshalIndent(baseOptionSet.Export(), "", "    ")
	if err != nil {
		return fmt.Errorf("go-config: error marshaling config: %s", err)
	}

	fp, err := os.OpenFile(f.Filename, os.O_RDWR+os.O_CREATE+os.O_TRUNC, 0666)
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

func (f FileIO) Read() (err error) {
	fp, err := os.Open(f.Filename)
	if err != nil {
		if os.IsNotExist(err) {
			return IOError{
				Type: "exist",
				Path: f.Filename,
				err:  err,
			}
		}

		return IOError{
			Type: "open",
			Path: f.Filename,
			err:  err,
		}
	}
	defer fp.Close()

	fi, err := fp.Stat()
	if err != nil {
		return IOError{
			Type: "stat",
			Path: f.Filename,
			err:  err,
		}
	}

	n := fi.Size()

	by := make([]byte, n)
	read, err := fp.Read(by)
	if err != nil || int64(read) < n {
		return IOError{
			Type: "read",
			Path: f.Filename,
			err:  err,
		}
	}

	jmap := jsonConfigMap{}
	err = json.Unmarshal(by, &jmap)
	if err != nil {
		return IOError{
			Type: "unmarshal",
			Path: f.Filename,
			err:  err,
		}
	}

	err = jmap.Parse()
	if err != nil {
		return err
	}

	return nil
}

// IOError describes an error related to loading a config file.
type IOError struct {
	Type string
	Path string
	err  error
}

func (e IOError) Error() string {
	return fmt.Sprintf("go-config: file i/o %s error on %s: %s", e.Type, e.Path, e.err)
}
