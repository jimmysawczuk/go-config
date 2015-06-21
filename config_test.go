package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
)

var originalOSArgs []string
var originalFlagSet *flag.FlagSet

func init() {
	// defined in usage.go, just tells Usage() where to direct output.
	// for our tests, we'll just direct it to /dev/null because we don't
	// care about the output.
	UsageWriter, _ = os.OpenFile(os.DevNull, os.O_RDWR, 0700)

	originalOSArgs = os.Args
	originalFlagSet = flag.CommandLine

	_ = fmt.Printf
}

func resetArgs() {
	os.Args = originalOSArgs
	flag.CommandLine = originalFlagSet
}

func TestBasicConfigLoad(t *testing.T) {
	// writing the config.json to a temporary file
	configJSON := []byte(`{
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

	fp.Write(configJSON)
	fp.Close()

	// rigging the test to use our temporary config file
	resetBaseOptionSet()
	baseOptionSet.Require("config").SetFromString(filepath)

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 10, "The first addend", OptionMeta{Exportable: true}))
	Add(Float("addend.b", math.Pi, "The second addend", OptionMeta{Exportable: true}))
	Add(Bool("subtract", false, "Subtract instead of add", OptionMeta{Exportable: true}))
	Add(String("name", "Basic Example", "Name of the example", OptionMeta{Exportable: true}))

	resetArgs()

	// and here we go!
	Build()

	a := Require("addend.a").Int()
	assert.Equal(t, a, int64(10), "addend.a should be 10")

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

func TestRequiredConfigLoad(t *testing.T) {
	// writing the config.json to a temporary file
	configJSON := []byte(`{
    "addend": {
        "a": 10,
        "b": 3.8
    },
    "subtract": false,
    "param-1": "",
    "param-2": "provided"
}`)

	filepath := os.TempDir() + "/go-config-required-config.json"

	fp, err := os.OpenFile(filepath, os.O_RDWR+os.O_CREATE+os.O_TRUNC, 0644)
	if err != nil {
		t.Errorf("Couldn't open temporary config file")
		t.FailNow()
	}

	fp.Write(configJSON)
	fp.Close()

	// rigging the test to use our temporary config file
	resetBaseOptionSet()
	baseOptionSet.Require("config").SetFromString(filepath)

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 10, "The first addend", OptionMeta{Exportable: true}))
	Add(Float("addend.b", math.Pi, "The second addend", OptionMeta{Exportable: true}))
	Add(Bool("subtract", false, "Subtract instead of add", OptionMeta{Exportable: true}))
	Add(String("param-1", "Value 1", "Name of the example", OptionMeta{Exportable: true, Required: true}))
	Add(String("param-2", "Value 2", "", OptionMeta{Exportable: true, Required: true}))

	resetArgs()

	// and here we go!
	err = Build()
	assert.NotNil(t, err, "There should be an error")

	a := Require("addend.a").Int()
	assert.Equal(t, a, int64(10), "addend.a should be 10")

	b := Require("addend.b").Float()
	assert.Equal(t, b, 3.8, "addend.b should be 3.8")

	sub := Require("subtract").Bool()
	assert.Equal(t, sub, false, "subtract should be false")

	c := float64(a) + float64(b)
	assert.Equal(t, c, 10+3.8, "The operation on addend.a + addend.b should be 13.8")

	param1 := Require("param-1").String()
	assert.Equal(t, param1, "", "Name should be \"\" (invalid, but it should still parse)")

	param2 := Require("param-2").String()
	assert.Equal(t, param2, "provided", "Name should be \"provided\"")

	_, err = Get("invalid-parameter")
	assert.NotEqual(t, err, nil, "Get(\"invalid-parameter\") should return an error")

	assert.Panics(t, func() {
		_ = Require("invalid-parameter")
	}, "Calling Require(\"invalid-parameter\") should panic")

	assert.NotPanics(t, func() {
		Usage()
	}, "Calling Usage() shouldn't panic")
}

func TestErroredConfigLoad(t *testing.T) {
	// writing the config.json to a temporary file
	// all of these values aren't properly set
	configJSON := []byte(`{
    "addend": {
        "a": 3.33333,
        "b": true
    },
    "subtract": "false",
    "name": false
}`)

	filepath := os.TempDir() + "/go-config-basic-config.json"

	fp, err := os.OpenFile(filepath, os.O_RDWR+os.O_CREATE+os.O_TRUNC, 0644)
	if err != nil {
		t.Errorf("Couldn't open temporary config file")
		t.FailNow()
	}

	fp.Write(configJSON)
	fp.Close()

	// rigging the test to use our temporary config file
	resetBaseOptionSet()
	baseOptionSet.Require("config").SetFromString(filepath)

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 10, "The first addend", OptionMeta{Exportable: true}))
	Add(Float("addend.b", math.Pi, "The second addend", OptionMeta{Exportable: true}))
	Add(Bool("subtract", false, "Subtract instead of add", OptionMeta{Exportable: true}))
	Add(String("name", "Basic Example", "Name of the example", OptionMeta{Exportable: true}))

	resetArgs()

	// and here we go!
	genericErr := Build()
	buildErr, ok := genericErr.(jsonConfigMapParseErrorList)
	_ = ok

	require.IsType(t, jsonConfigMapParseErrorList{}, buildErr, true, "Build() should return a jsonConfigMapParseErrorList, instead %T", genericErr)
	assert.Equal(t, 4, buildErr.Len(), "There should be 4 build errors")
}

func TestBasicConfigLoadWithFlags(t *testing.T) {
	// writing the config.json to a temporary file
	configJSON := []byte(`{
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

	fp.Write(configJSON)
	fp.Close()

	// rigging the test to use our temporary config file
	resetBaseOptionSet()
	baseOptionSet.Require("config").SetFromString(filepath)

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 10, "The first addend", OptionMeta{Exportable: true}))
	Add(Float("addend.b", math.Pi, "The second addend", OptionMeta{Exportable: true}))
	Add(Bool("subtract", false, "Subtract instead of add", OptionMeta{Exportable: true}))
	Add(String("name", "Basic Example", "Name of the example", OptionMeta{Exportable: true}))

	// and here we go!
	os.Args = []string{
		`go-config`,
		`-addend.a=4`,
		`-addend.b=4`,
		`-subtract=false`,
		`-name=Flag override`,
	}
	Build()

	a := Require("addend.a").Int()
	assert.Equal(t, a, int64(4), "addend.a should be 4")

	b := Require("addend.b").Float()
	assert.Equal(t, b, float64(4), "addend.b should be 4")

	sub := Require("subtract").Bool()
	assert.Equal(t, sub, false, "subtract should be false")

	c := float64(a) + float64(b)
	assert.Equal(t, c, float64(4+4), "The operation on addend.a + addend.b should be 8")

	name := Require("name").String()
	assert.Equal(t, name, "Flag override", "Name should be \"Flag override\"")

	_, err = Get("invalid-parameter")
	assert.NotEqual(t, err, nil, "Get(\"invalid-parameter\") should return an error")

	assert.Panics(t, func() {
		_ = Require("invalid-parameter")
	}, "Calling Require(\"invalid-parameter\") should panic")

	assert.NotPanics(t, func() {
		Usage()
	}, "Calling Usage() shouldn't panic")
}

func TestBasicConfigLoadWithOtherFlags(t *testing.T) {
	// writing the config.json to a temporary file
	configJSON := []byte(`{
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

	fp.Write(configJSON)
	fp.Close()

	// rigging the test to use our temporary config file
	resetBaseOptionSet()
	baseOptionSet.Require("config").SetFromString(filepath)

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 10, "The first addend", OptionMeta{Exportable: true}))
	Add(Float("addend.b", math.Pi, "The second addend", OptionMeta{Exportable: true}))
	Add(Bool("subtract", false, "Subtract instead of add", OptionMeta{Exportable: true}))
	Add(String("name", "Basic Example", "Name of the example", OptionMeta{Exportable: true}))

	// and here we go!
	os.Args = []string{
		`go-config`,
		`-addend.a`,
		`4`,

		`-addend.b`,
		`2`,

		`-subtract`,

		`--name`,
		`Test`,
	}
	Build()

	a := Require("addend.a").Int()
	assert.Equal(t, a, int64(4), "addend.a should be 4")

	b := Require("addend.b").Float()
	assert.Equal(t, b, float64(2), "addend.b should be 2")

	sub := Require("subtract").Bool()
	assert.Equal(t, sub, true, "subtract should be true")

	c := float64(a) - float64(b)
	assert.Equal(t, c, float64(4-2), "The operation on addend.a - addend.b should be 2")

	name := Require("name").String()
	assert.Equal(t, name, "Test", "Name should be \"Test\"")

	_, err = Get("invalid-parameter")
	assert.NotEqual(t, err, nil, "Get(\"invalid-parameter\") should return an error")

	assert.Panics(t, func() {
		_ = Require("invalid-parameter")
	}, "Calling Require(\"invalid-parameter\") should panic")

	assert.NotPanics(t, func() {
		Usage()
	}, "Calling Usage() shouldn't panic")
}

func TestBasicConfigLoadWithFinalBooleanFlag(t *testing.T) {
	// writing the config.json to a temporary file
	configJSON := []byte(`{
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

	fp.Write(configJSON)
	fp.Close()

	// rigging the test to use our temporary config file
	resetBaseOptionSet()
	baseOptionSet.Require("config").SetFromString(filepath)

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 10, "The first addend", OptionMeta{Exportable: true}))
	Add(Float("addend.b", math.Pi, "The second addend", OptionMeta{Exportable: true}))
	Add(Bool("subtract", false, "Subtract instead of add", OptionMeta{Exportable: true}))
	Add(String("name", "Basic Example", "Name of the example", OptionMeta{Exportable: true}))

	// and here we go!
	os.Args = []string{
		`go-config`,
		`-addend.a`,
		`4`,

		`-addend.b`,
		`2`,

		`--name=Test`,

		`-subtract`,
	}
	Build()

	a := Require("addend.a").Int()
	assert.Equal(t, a, int64(4), "addend.a should be 4")

	b := Require("addend.b").Float()
	assert.Equal(t, b, float64(2), "addend.b should be 2")

	sub := Require("subtract").Bool()
	assert.Equal(t, sub, true, "subtract should be true")

	c := float64(a) - float64(b)
	assert.Equal(t, c, float64(4-2), "The operation on addend.a - addend.b should be 2")

	name := Require("name").String()
	assert.Equal(t, name, "Test", "Name should be \"Test\"")

	_, err = Get("invalid-parameter")
	assert.NotEqual(t, err, nil, "Get(\"invalid-parameter\") should return an error")

	assert.Panics(t, func() {
		_ = Require("invalid-parameter")
	}, "Calling Require(\"invalid-parameter\") should panic")

	assert.NotPanics(t, func() {
		Usage()
	}, "Calling Usage() shouldn't panic")
}

func TestBasicConfigLoadWithUndefinedFlags(t *testing.T) {
	// writing the config.json to a temporary file
	configJSON := []byte(`{
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

	fp.Write(configJSON)
	fp.Close()

	// rigging the test to use our temporary config file
	resetBaseOptionSet()
	baseOptionSet.Require("config").SetFromString(filepath)

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 10, "The first addend", OptionMeta{Exportable: true}))
	Add(Float("addend.b", math.Pi, "The second addend", OptionMeta{Exportable: true}))
	Add(Bool("subtract", false, "Subtract instead of add", OptionMeta{Exportable: true}))
	Add(String("name", "Basic Example", "Name of the example", OptionMeta{Exportable: true}))

	// and here we go!
	os.Args = []string{
		`go-config`,
		`-addend.a`,
		`4`,

		`-addend.b`,
		`2`,

		`-subtract`,

		`--name`,
		`Test`,

		`-addend.c`,
		`4`,
	}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	err = Build()

	assert.NotNil(t, err, "The error should not be nil")

	a := Require("addend.a").Int()
	assert.Equal(t, a, int64(4), "addend.a should be 4")

	b := Require("addend.b").Float()
	assert.Equal(t, b, float64(2), "addend.b should be 2")

	sub := Require("subtract").Bool()
	assert.Equal(t, sub, true, "subtract should be true")

	c := float64(a) - float64(b)
	assert.Equal(t, c, float64(4-2), "The operation on addend.a - addend.b should be 2")

	name := Require("name").String()
	assert.Equal(t, name, "Test", "Name should be \"Test\"")
}

func TestBasicConfigWrite(t *testing.T) {
	// writing the config.json to a temporary file
	configJSON := []byte(`{
    "addend": {
        "a": 10,
        "b": 3.8
    },
    "subtract": false,
    "name": "Basic Example"
}`)

	newConfigJSON := []byte(`{
    "addend": {
        "a": 15,
        "b": 15
    },
    "name": "Test",
    "subtract": true
}`)

	filepath := os.TempDir() + "/go-config-basic-config-write.json"

	fp, err := os.OpenFile(filepath, os.O_RDWR+os.O_CREATE+os.O_TRUNC, 0644)
	if err != nil {
		t.Errorf("Couldn't open temporary config file")
		t.FailNow()
	}

	fp.Write(configJSON)
	fp.Close()

	// rigging the test to use our temporary config file
	resetBaseOptionSet()
	baseOptionSet.Require("config").SetFromString(filepath)

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 10, "The first addend", OptionMeta{Exportable: true}))
	Add(Float("addend.b", math.Pi, "The second addend", OptionMeta{Exportable: true}))
	Add(Bool("subtract", false, "Subtract instead of add", OptionMeta{Exportable: true}))
	Add(String("name", "Basic Example", "Name of the example", OptionMeta{Exportable: true}))

	os.Args = []string{
		`go-config`,
		`-config-export`,

		`-addend.a`,
		`15`,

		`-addend.b`,
		`15`,

		`-subtract`,

		`--name`,
		`Test`,
	}

	// and here we go!
	Build()

	resetBaseOptionSet()
	baseOptionSet.Require("config").SetFromString(filepath)

	Add(Int("addend.a", 10, "The first addend", OptionMeta{Exportable: true}))
	Add(Float("addend.b", math.Pi, "The second addend", OptionMeta{Exportable: true}))
	Add(Bool("subtract", false, "Subtract instead of add", OptionMeta{Exportable: true}))
	Add(String("name", "Basic Example", "Name of the example", OptionMeta{Exportable: true}))

	resetArgs()

	Build()

	builtJSON, _ := json.Marshal(baseOptionSet.Export())
	buf := bytes.Buffer{}
	json.Compact(&buf, newConfigJSON)

	assert.Equal(t, buf.String(), string(builtJSON), "Written config file should match expected")
}

func TestEnumConfig(t *testing.T) {
	// writing the config.json to a temporary file
	configJSON := []byte(`{
    "addend": {
        "a": 4,
        "b": 2
    },
    "mode": "subtract",
    "name": "Test"
}`)

	filepath := os.TempDir() + "/go-config-enum-config.json"

	fp, err := os.OpenFile(filepath, os.O_RDWR+os.O_CREATE+os.O_TRUNC, 0644)
	if err != nil {
		t.Errorf("Couldn't open temporary config file")
		t.FailNow()
	}

	fp.Write(configJSON)
	fp.Close()

	// rigging the test to use our temporary config file
	resetBaseOptionSet()
	baseOptionSet.Require("config").SetFromString(filepath)

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 10, "The first addend", OptionMeta{Exportable: true}))
	Add(Float("addend.b", math.Pi, "The second addend", OptionMeta{Exportable: true}))
	Add(Enum("mode", []string{"subtract", "add"}, "add", "subtract or add", OptionMeta{Exportable: true}))
	Add(String("name", "Basic Example", "Name of the example", OptionMeta{Exportable: true}))

	err = Build()

	a := Require("addend.a").Int()
	assert.Equal(t, a, int64(4), "addend.a should be 4")

	b := Require("addend.b").Float()
	assert.Equal(t, b, float64(2), "addend.b should be 2")

	mode := Require("mode").String()
	assert.Equal(t, mode, "subtract", "mode should be subtract")

	c := float64(a) - float64(b)
	assert.Equal(t, c, float64(4-2), "The operation on addend.a - addend.b should be 2")

	name := Require("name").String()
	assert.Equal(t, name, "Test", "Name should be \"Test\"")
}

func TestInvalidEnumConfig(t *testing.T) {
	// writing the config.json to a temporary file
	configJSON := []byte(`{
    "addend": {
        "a": 4,
        "b": 2
    },
    "mode": "invalid",
    "name": "Test"
}`)

	filepath := os.TempDir() + "/go-config-enum-config.json"

	fp, err := os.OpenFile(filepath, os.O_RDWR+os.O_CREATE+os.O_TRUNC, 0644)
	if err != nil {
		t.Errorf("Couldn't open temporary config file")
		t.FailNow()
	}

	fp.Write(configJSON)
	fp.Close()

	// rigging the test to use our temporary config file
	resetBaseOptionSet()
	baseOptionSet.Require("config").SetFromString(filepath)

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 10, "The first addend", OptionMeta{Exportable: true}))
	Add(Float("addend.b", math.Pi, "The second addend", OptionMeta{Exportable: true}))
	Add(Enum("mode", []string{"subtract", "add"}, "subtract", "subtract or add", OptionMeta{Exportable: true}))
	Add(String("name", "Basic Example", "Name of the example", OptionMeta{Exportable: true}))

	err = Build()
	t.Logf("%s", err)

	assert.NotNil(t, err, "There should be an error because the enum is wrong")
}
