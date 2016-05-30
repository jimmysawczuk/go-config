package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
)

var originalOSArgs []string
var originalFlagSet *flag.FlagSet
var tempDir, tempAppDir, tempUserDir, tempCustomDir string

func init() {
	// defined in usage.go, just tells Usage() where to direct output.
	// for our tests, we'll just direct it to /dev/null because we don't
	// care about the output.
	UsageWriter, _ = os.OpenFile(os.DevNull, os.O_RDWR, 0700)

	originalOSArgs = os.Args
	originalFlagSet = flag.CommandLine

	tempDir = path.Clean(os.TempDir() + "/go-config")
	err := os.MkdirAll(tempDir, 0775)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating temporary directory: %s\n", err)
		os.Exit(2)
	}

	SearchFiles = []SearchFile{
		{
			Scope: "app",
			Path:  tempDir + "/app/config.json",
		},
		{
			Scope: "user",
			Path:  tempDir + "/user/config.json",
		},
	}

	for _, s := range SearchFiles {
		err := os.MkdirAll(tempDir+"/"+s.Scope, 0775)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating temporary directory: %s\n", err)
			os.Exit(2)
		}
	}

	for _, s := range []SearchFile{SearchFile{Scope: "custom", Path: tempDir + "/custom/config.json"}} {
		err := os.MkdirAll(tempDir+"/"+s.Scope, 0775)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating temporary directory: %s\n", err)
			os.Exit(2)
		}
	}

	tempAppDir = path.Clean(tempDir + "/app")
	tempUserDir = path.Clean(tempDir + "/user")
	tempCustomDir = path.Clean(tempDir + "/custom")
}

func resetArgs() {
	os.Args = originalOSArgs
	flag.CommandLine = originalFlagSet
}

func writeToTemporaryFile(t *testing.T, out []byte, filepath string) {
	fp, err := os.OpenFile(filepath, os.O_RDWR+os.O_CREATE+os.O_TRUNC, 0644)
	if err != nil {
		t.Errorf("Couldn't open temporary config file at %s: %s", filepath, err)
		t.FailNow()
	}

	fp.Write(out)
	fp.Close()
}

func readFromTemporaryFile(t *testing.T, filepath string) []byte {
	by, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Errorf("Couldn't read temporary config file at %s: %s", filepath, err)
		t.FailNow()
	}

	return by
}

func TestBasicConfigLoad(t *testing.T) {
	var err error
	var filepath = tempAppDir + "/config.json"

	// writing the config.json to a temporary file
	configJSON := []byte(`{
    "addend": {
        "a": 10,
        "b": 3.8
    },
    "bad_string": 8.5,
    "subtract": false,
    "name": "Basic Example"
}`)

	writeToTemporaryFile(t, configJSON, filepath)
	resetBaseOptionSet()

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 10, "The first addend").Exportable(true).SortOrder(-1))
	Add(Float("addend.b", math.Pi, "The second addend").Exportable(true).SortOrder(-1))
	Add(Bool("subtract", false, "Subtract instead of add").Exportable(true))
	Add(Str("name", "Basic Example", "Name of the example").Exportable(true).SortOrder(+1))
	Add(Str("bad_string", "test", "Defined as a string, but isn't").Exportable(true).SortOrder(+1))

	resetArgs()

	// and here we go!
	err = Build()
	assert.EqualError(t, err, jsonConfigMapParseErrorList([]jsonConfigMapError{
		jsonConfigMapParseError{
			key:      "bad_string",
			got:      float64(8.5),
			expected: StringType,
		},
	}).Error())

	a := Require("addend.a").Int()
	assert.Equal(t, a, int64(10), "addend.a should be 10")

	b := Require("addend.b").Float()
	assert.Equal(t, b, 3.8, "addend.b should be 3.8")

	sub := Require("subtract").Bool()
	assert.Equal(t, sub, false, "subtract should be false")

	c := float64(a) + float64(b)
	assert.Equal(t, c, 10+3.8, "The operation on addend.a + addend.b should be 13.8")

	name := Require("name").Str()
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
	var err error
	var filepath = tempAppDir + "/config.json"

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

	writeToTemporaryFile(t, configJSON, filepath)
	resetBaseOptionSet()

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 10, "The first addend").Exportable(true))
	Add(Float("addend.b", math.Pi, "The second addend").Exportable(true))
	Add(Bool("subtract", false, "Subtract instead of add").Exportable(true))
	Add(Str("param-1", "Value 1", "Name of the example").Exportable(true).AddFilter(NonEmptyString()))
	Add(Str("param-2", "Value 2", "").Exportable(true).AddFilter(NonEmptyString()))

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

	param1 := Require("param-1").Str()
	assert.Equal(t, param1, "", "Name should be \"\" (invalid, but it should still parse)")

	param2 := Require("param-2").Str()
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
	var filepath = tempAppDir + "/config.json"

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

	writeToTemporaryFile(t, configJSON, filepath)
	resetBaseOptionSet()

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 10, "The first addend").Exportable(true))
	Add(Float("addend.b", math.Pi, "The second addend").Exportable(true))
	Add(Bool("subtract", false, "Subtract instead of add").Exportable(true))
	Add(Str("name", "Basic Example", "Name of the example").Exportable(true))

	resetArgs()

	// and here we go!
	genericErr := Build()
	buildErr, ok := genericErr.(jsonConfigMapParseErrorList)
	_ = ok

	require.IsType(t, jsonConfigMapParseErrorList{}, buildErr, true, "Build() should return a jsonConfigMapParseErrorList, instead %T", genericErr)
	assert.Equal(t, 4, buildErr.Len(), "There should be 4 build errors")
}

func TestConfigLoadWithFlags(t *testing.T) {
	var err error
	var filepath = tempAppDir + "/config.json"

	// writing the config.json to a temporary file
	configJSON := []byte(`{
    "addend": {
        "a": 10,
        "b": 3.8
    },
    "subtract": false,
    "name": "Basic Example"
}`)

	writeToTemporaryFile(t, configJSON, filepath)
	resetBaseOptionSet()

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 10, "The first addend").Exportable(true))
	Add(Float("addend.b", math.Pi, "The second addend").Exportable(true))
	Add(Bool("subtract", false, "Subtract instead of add").Exportable(true))
	Add(Str("name", "Basic Example", "Name of the example").Exportable(true))

	os.Args = []string{
		`go-config`,
		`-addend.a=4`,
		`-addend.b=4`,
		`-name=Overridden by flag`,
		`-subtract`,
	}

	Build()

	a := Require("addend.a").Int()
	assert.Equal(t, a, int64(4), "addend.a should be 4")

	b := Require("addend.b").Float()
	assert.Equal(t, b, float64(4), "addend.b should be 4")

	sub := Require("subtract").Bool()
	assert.Equal(t, sub, true, "subtract should be false")

	c := float64(a) + float64(b)
	assert.Equal(t, c, float64(4+4), "The operation on addend.a + addend.b should be 8")

	name := Require("name").Str()
	assert.Equal(t, name, "Overridden by flag", "Name should be \"Overridden by flag\"")

	_, err = Get("invalid-parameter")
	assert.NotEqual(t, err, nil, "Get(\"invalid-parameter\") should return an error")

	assert.Panics(t, func() {
		_ = Require("invalid-parameter")
	}, "Calling Require(\"invalid-parameter\") should panic")

	assert.NotPanics(t, func() {
		Usage()
	}, "Calling Usage() shouldn't panic")
}

func TestConfigLoadWithAlternateFlags(t *testing.T) {
	var err error
	var filepath = tempAppDir + "/config.json"

	// writing the config.json to a temporary file
	configJSON := []byte(`{
    "addend": {
        "a": 10,
        "b": 3.8
    },
    "subtract": false,
    "name": "Basic Example"
}`)

	writeToTemporaryFile(t, configJSON, filepath)
	resetBaseOptionSet()

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 10, "The first addend").Exportable(true))
	Add(Float("addend.b", math.Pi, "The second addend").Exportable(true))
	Add(Bool("subtract", false, "Subtract instead of add").Exportable(true))
	Add(Str("name", "Basic Example", "Name of the example").Exportable(true))

	os.Args = []string{
		`go-config`,

		`-addend.a`,
		`4`,

		`-addend.b`,
		`4`,

		`-name`,
		`Overridden by flag`,

		`-subtract`,
	}

	Build()

	a := Require("addend.a").Int()
	assert.Equal(t, a, int64(4), "addend.a should be 4")

	b := Require("addend.b").Float()
	assert.Equal(t, b, float64(4), "addend.b should be 4")

	sub := Require("subtract").Bool()
	assert.Equal(t, sub, true, "subtract should be false")

	c := float64(a) + float64(b)
	assert.Equal(t, c, float64(4+4), "The operation on addend.a + addend.b should be 8")

	name := Require("name").Str()
	assert.Equal(t, name, "Overridden by flag", "Name should be \"Overridden by flag\"")

	_, err = Get("invalid-parameter")
	assert.NotEqual(t, err, nil, "Get(\"invalid-parameter\") should return an error")

	assert.Panics(t, func() {
		_ = Require("invalid-parameter")
	}, "Calling Require(\"invalid-parameter\") should panic")

	assert.NotPanics(t, func() {
		Usage()
	}, "Calling Usage() shouldn't panic")
}

func TestConfigWrite(t *testing.T) {
	var filepath = tempAppDir + "/config.json"

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

	writeToTemporaryFile(t, configJSON, filepath)
	resetBaseOptionSet()

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 10, "The first addend").Exportable(true))
	Add(Float("addend.b", math.Pi, "The second addend").Exportable(true))
	Add(Bool("subtract", false, "Subtract instead of add").Exportable(true))
	Add(Str("name", "Basic Example", "Name of the example").Exportable(true))

	os.Args = []string{
		`go-config`,

		`-config-save`,
		`-config-scope=app`,

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
	resetArgs()

	Add(Int("addend.a", 10, "The first addend").Exportable(true))
	Add(Float("addend.b", math.Pi, "The second addend").Exportable(true))
	Add(Bool("subtract", false, "Subtract instead of add").Exportable(true))
	Add(Str("name", "Basic Example", "Name of the example").Exportable(true))

	Build()

	builtJSON, _ := json.Marshal(baseOptionSet.Export(false, true))
	buf := bytes.Buffer{}
	json.Compact(&buf, newConfigJSON)

	assert.Equal(t, buf.String(), string(builtJSON), "Written config file should match expected")
}

func TestEnumConfig(t *testing.T) {
	var filepath = tempAppDir + "/config.json"

	// writing the config.json to a temporary file
	configJSON := []byte(`{
    "addend": {
        "a": 4,
        "b": 2
    },
    "mode": "subtract",
    "name": "Test"
}`)

	writeToTemporaryFile(t, configJSON, filepath)
	resetBaseOptionSet()

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 10, "The first addend").Exportable(true))
	Add(Float("addend.b", math.Pi, "The second addend").Exportable(true))
	Add(Enum("mode", []string{"subtract", "add"}, "add", "subtract or add").Exportable(true))
	Add(Str("name", "Basic Example", "Name of the example").Exportable(true))

	Build()

	a := Require("addend.a").Int()
	assert.Equal(t, a, int64(4), "addend.a should be 4")

	b := Require("addend.b").Float()
	assert.Equal(t, b, float64(2), "addend.b should be 2")

	mode := Require("mode").Str()
	assert.Equal(t, mode, "subtract", "mode should be subtract")

	c := float64(a) - float64(b)
	assert.Equal(t, c, float64(4-2), "The operation on addend.a - addend.b should be 2")

	name := Require("name").Str()
	assert.Equal(t, name, "Test", "Name should be \"Test\"")
}

func TestInvalidEnumConfig(t *testing.T) {
	var err error
	var filepath = tempAppDir + "/config.json"

	// writing the config.json to a temporary file
	configJSON := []byte(`{
    "addend": {
        "a": 4,
        "b": 2
    },
    "mode": "invalid",
    "name": "Test"
}`)

	writeToTemporaryFile(t, configJSON, filepath)
	resetBaseOptionSet()

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 10, "The first addend").Exportable(true))
	Add(Float("addend.b", math.Pi, "The second addend").Exportable(true))
	Add(Enum("mode", []string{"subtract", "add"}, "subtract", "subtract or add").Exportable(true))
	Add(Str("name", "Basic Example", "Name of the example").Exportable(true))

	err = Build()

	assert.NotNil(t, err, "There should be an error because the enum is wrong")
}

func TestTripleNestedOption(t *testing.T) {
	var err error
	var filepath = tempAppDir + "/config.json"

	// writing the config.json to a temporary file
	configJSON := []byte(`{
	"equation": {
		"addend": {
			"a": 4,
			"b": 2
		}
	},
	"mode": "subtract",
	"name": "Test"
}`)

	writeToTemporaryFile(t, configJSON, filepath)
	resetBaseOptionSet()

	// setting up our config options to read the temporary config.json properly
	Add(Int("equation.addend.a", 10, "The first addend").Exportable(true))
	Add(Int("equation.addend.b", 5, "The second addend").Exportable(true))
	Add(Enum("mode", []string{"subtract", "add"}, "subtract", "subtract or add").Exportable(true))
	Add(Str("name", "Basic Example", "Name of the example").Exportable(true))

	err = Build()

	assert.Nil(t, err, "There should be no error here")

	a := Require("equation.addend.a").Int()
	b := Require("equation.addend.b").Int()

	assert.Equal(t, int64(4), a, "addend_a should be 4")
	assert.Equal(t, int64(2), b, "addend_b should be 2")

	assert.Equal(t, int64(2), a-b, "%d minus %d != 2", a, b)
}

func TestDifferentFlagFormats(t *testing.T) {
	var err error
	var filepath = tempAppDir + "/config.json"

	// writing the config.json to a temporary file
	configJSON := []byte(`{
	"addend": {
		"a": 1,
		"b": 1
	},
	"subtract": true,
	"name": "Config file"
}`)

	writeToTemporaryFile(t, configJSON, filepath)
	resetBaseOptionSet()

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 0, "The first addend").Exportable(true))
	Add(Int("addend.b", 0, "The second addend").Exportable(true))
	Add(Bool("subtract", true, "Subtract the arguments").Exportable(true))
	Add(Str("name", "", "Name of the example").Exportable(true))

	os.Args = []string{
		`go-config`,

		`--addend.a`,
		`2`,

		`--addend.b=2`,

		`-subtract=f`,

		`-name`,
		`Test`,
	}

	// and here we go!
	err = Build()

	assert.Nil(t, err, "There is no error here")

	a := Require("addend.a").Int()
	b := Require("addend.b").Int()
	subtract := Require("subtract").Bool()
	name := Require("name").Str()

	assert.Equal(t, int64(2), a, "addend_a should be 2")
	assert.Equal(t, int64(2), b, "addend_b should be 2")
	assert.Equal(t, false, subtract, "subtract should be false")
	assert.Equal(t, "Test", name, "name should be Test")
}

func TestCustomConfigFile(t *testing.T) {
	var err error
	var filepath = tempAppDir + "/config.json"
	var custompath = tempCustomDir + "/config.json"

	// writing the config.json to a temporary file
	configJSON := []byte(`{
	"addend": {
		"a": 1,
		"b": 1
	},
	"subtract": true,
	"name": "Config file"
}`)

	customConfigJSON := []byte(`{
	"addend": {
		"a": 2,
		"b": 2
	},
	"name": "Config file",
	"subtract": true
}`)

	writeToTemporaryFile(t, configJSON, filepath)
	resetBaseOptionSet()

	// setting up our config options to read the temporary config.json properly
	Add(Int("addend.a", 0, "The first addend").Exportable(true))
	Add(Int("addend.b", 0, "The second addend").Exportable(true))
	Add(Bool("subtract", true, "Subtract the arguments").Exportable(true))
	Add(Str("name", "", "Name of the example").Exportable(true))

	os.Args = []string{
		`go-config`,

		"-config-file",
		custompath,

		"-config-scope=custom",
		"-config-save",

		`--addend.a`,
		`2`,

		`--addend.b=2`,
	}

	// and here we go!
	err = Build()

	assert.Nil(t, err, "There is no error here")

	a := Require("addend.a").Int()
	b := Require("addend.b").Int()
	subtract := Require("subtract").Bool()
	name := Require("name").Str()

	assert.Equal(t, int64(2), a, "addend_a should be 2")
	assert.Equal(t, int64(2), b, "addend_b should be 2")
	assert.Equal(t, true, subtract, "subtract should be true")
	assert.Equal(t, "Config file", name, "name should be Config file")

	written := readFromTemporaryFile(t, custompath)

	fmt.Println(string(written), string(customConfigJSON))

	assert.Equal(t, string(customConfigJSON), string(written), "Written config file should match expected")
}
