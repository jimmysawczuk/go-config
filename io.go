package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// IO defines an interface that allows reading and writing of OptionSets to external storage
type IO interface {
	Read() error
	Write() error
	Scope() string
}

// FileIO implements IO and writes to the filesystem
type FileIO struct {
	filename string
	scope    string
}

func (f FileIO) Write() (err error) {
	partialExport := Require("config-partial").Bool()

	json, err := json.MarshalIndent(baseOptionSet.Export(false, !partialExport), "", "\t")
	if err != nil {
		return fmt.Errorf("go-config: error marshaling config: %s", err)
	}

	err = os.MkdirAll(filepath.Dir(f.filename), 0755)
	if err != nil {
		return fmt.Errorf("go-config: file i/o directory error: %s", err)
	}

	fp, err := os.OpenFile(f.filename, os.O_RDWR+os.O_CREATE+os.O_TRUNC, 0644)
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
	fp, err := os.Open(f.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return IOError{
				Type: "exist",
				Path: f.filename,
				err:  err,
			}
		}

		return IOError{
			Type: "open",
			Path: f.filename,
			err:  err,
		}
	}
	defer fp.Close()

	fi, err := fp.Stat()
	if err != nil {
		return IOError{
			Type: "stat",
			Path: f.filename,
			err:  err,
		}
	}

	n := fi.Size()

	by := make([]byte, n)
	read, err := fp.Read(by)
	if err != nil || int64(read) < n {
		return IOError{
			Type: "read",
			Path: f.filename,
			err:  err,
		}
	}

	jmap := jsonConfigMap{
		scope: f.scope,
	}
	err = json.Unmarshal(by, &jmap)
	if err != nil {
		return IOError{
			Type: "unmarshal",
			Path: f.filename,
			err:  err,
		}
	}

	err = jmap.Parse()
	if err != nil {
		return err
	}

	return nil
}

func (f FileIO) Scope() string {
	return f.scope
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
