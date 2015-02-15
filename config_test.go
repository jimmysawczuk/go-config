package config

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"fmt"
	"math"
	"os"
	"strings"
)

var our_args, test_args []string

func init() {
	our_args = []string{""}
	test_args = []string{""}

	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-test.") {
			our_args = append(our_args, arg)
		} else {
			test_args = append(test_args, arg)
		}
	}

	os.Args = test_args

	_ = fmt.Printf
}

func TestBasicConfigLoad(t *testing.T) {
	// writing the config.json to a temporary file
	config_json := []byte(`{
    "addend": {
        "a": 10,
        "b": 3.8
    },
    "subtract": false,
    "name": "Basic Example"
}`)

	filepath := os.TempDir() + "/go-config-basic-config.json"

	fp, err := os.OpenFile(filepath, os.O_RDWR+os.O_CREATE+os.O_TRUNC, 0644)
	if err != nil {
		t.Errorf("Couldn't open temporary config file")
		t.FailNow()
	}

	fp.Write(config_json)

	// rigging the test to use our temporary config file
	resetBaseOptionSet(false)
	Add(String("config", filepath, "The filename of the config file to use", false))
	Add(Bool("config-export", false, "Export the as-run configuration to a file", false))
	Add(Bool("config-generate", false, "Export the as-run configuration to a file, then exit", false))

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 10, "The first addend", true))
	Add(Float("addend.b", math.Pi, "The second addend", true))
	Add(Bool("subtract", false, "Subtract instead of add", true))
	Add(String("name", "Basic Example", "Name of the example", true))

	// and here we go!
	os.Args = our_args
	Build()

	a := Require("addend.a").Int()
	assert.Equal(t, a, 10, "addend.a should be 10")

	b := Require("addend.b").Float()
	assert.Equal(t, b, 3.8, "addend.b should be 3.8")

	sub := Require("subtract").Bool()
	assert.Equal(t, sub, false, "subtract should be false")

	c := float64(a) + float64(b)
	assert.Equal(t, c, 10+3.8, "The operation on addend.a + addend.b should be 13.8")

	name := Require("name").String()
	assert.Equal(t, name, "Basic Example", "Name should be \"Basic Example\"")

	_, err = Get("invalid-parameter")
	assert.NotEqual(t, err, nil, "Get(\"invalid-parameter\") should return an error")

	assert.Panics(t, func() {
		_ = Require("invalid-parameter")
	}, "Calling Require(\"invalid-parameter\") should panic")

	assert.NotPanics(t, func() {
		Usage()
	}, "Calling Usage() shouldn't panic")
}
