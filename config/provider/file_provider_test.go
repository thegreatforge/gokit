package provider

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileProvider(t *testing.T) {
	assert.NotNil(t, NewFileProvider([]string{"test"}))
	assert.Equal(t, &fileProvider{
		paths:     []string{"test"},
		delimiter: ".",
	}, NewFileProvider([]string{"test"}))
}

func TestFileLoadConfig(t *testing.T) {
	fp := &fileProvider{
		paths:     []string{"test"},
		delimiter: ".",
	}

	data := make(map[string]interface{})
	err := fp.LoadConfig(data)
	assert.Error(t, err)

	// create test file
	os.WriteFile("test.yaml", []byte("test: test"), 0644)
	defer os.Remove("test.yaml")

	fp.paths = []string{"test.yaml"}
	err = fp.LoadConfig(data)
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"test": "test"}, data)

}

func TestParseMap(t *testing.T) {
	// write map in yaml format to file
	os.WriteFile("test.yaml", []byte("test: test"), 0644)
	defer os.Remove("test.yaml")

	// create test file provider
	fp := &fileProvider{
		paths:     []string{"test.yaml"},
		delimiter: ".",
	}

	// create test map
	testMap := make(map[string]interface{})
	testMap["test"] = "test"

	// parse map
	parsedMap, err := fp.parseMap(testMap, "")
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"test": "test"}, parsedMap)
}

func TestParseSlice(t *testing.T) {
	// write slice in yaml format to file
	os.WriteFile("test.yaml", []byte("- test"), 0644)
	defer os.Remove("test.yaml")

	// create test file provider
	fp := &fileProvider{
		paths:     []string{"test.yaml"},
		delimiter: ".",
	}

	// create test slice
	testSlice := make([]interface{}, 1)
	testSlice[0] = "test"

	// parse slice
	parsedSlice, err := fp.parseSlice(testSlice, "")
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"0": "test"}, parsedSlice)
}
