package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialise(t *testing.T) {

	assert.Error(t, Initialise())

	// create test file
	os.WriteFile("test.yaml", []byte("test: test"), 0644)
	os.WriteFile("test.json", []byte("{\"test\": \"test\"}"), 0644)
	defer os.Remove("test.yaml")
	defer os.Remove("test.json")

	assert.NoError(t, Initialise(WithFiles("test.yaml")))
	assert.NoError(t, Initialise(WithFiles("test.yaml", "test.json")))
	assert.Error(t, Initialise(WithFiles("test.txt")))
	assert.Error(t, Initialise(WithFiles()))

	// create test env variable
	os.Setenv("test", "test")
	defer os.Unsetenv("test")

	assert.NoError(t, Initialise(WithEnvVariables("test")))
	assert.Error(t, Initialise(WithEnvVariables()))
}

func TestWithFiles(t *testing.T) {

	// create test file
	os.WriteFile("test.yaml", []byte("test: test"), 0644)
	defer os.Remove("test.yaml")

	assert.NoError(t, Initialise(WithFiles("test.yaml")))
	assert.Error(t, Initialise(WithFiles("test.txt")))

	// create test.txt
	os.WriteFile("test.txt", []byte("test: test"), 0644)
	defer os.Remove("test.txt")
	assert.Error(t, Initialise(WithFiles("test.txt")))
}

func TestWithEnvVariables(t *testing.T) {

	// create test env variable
	os.Setenv("test", "test")
	defer os.Unsetenv("test")

	assert.NoError(t, Initialise(WithEnvVariables("test")))
	assert.Error(t, Initialise(WithEnvVariables()))
}

func TestClose(t *testing.T) {

	assert.NotPanics(t, Close)
}

func TestGet(t *testing.T) {

	// set env variable
	os.Setenv("test", "test")
	defer os.Unsetenv("test")

	// initialise config
	assert.NoError(t, Initialise(WithEnvVariables("test")))
	val, err := Get("test")
	assert.NoError(t, err)
	assert.Equal(t, "test", val.(string))
}

func TestGetBool(t *testing.T) {

	// create test file
	os.WriteFile("test.yaml", []byte("test: true"), 0644)
	// initialise config
	assert.NoError(t, Initialise(WithFiles("test.yaml")))
	val, err := GetBool("test")
	assert.NoError(t, err)
	assert.Equal(t, true, val)
}

func TestGetInt(t *testing.T) {

	// create test file
	os.WriteFile("test.yaml", []byte("test: 1"), 0644)
	defer os.Remove("test.yaml")

	// initialise config
	assert.NoError(t, Initialise(WithFiles("test.yaml")))
	val, err := GetInt("test")
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
}

func TestGetString(t *testing.T) {

	// set env variable
	os.Setenv("test", "test")
	defer os.Unsetenv("test")

	// initialise config
	assert.NoError(t, Initialise(WithEnvVariables("test")))
	val, err := GetString("test")
	assert.NoError(t, err)
	assert.Equal(t, "test", val)
}

func TestGetFloat(t *testing.T) {

	// create test file
	os.WriteFile("test.yaml", []byte("test: 1.1"), 0644)
	defer os.Remove("test.yaml")

	// initialise config
	assert.NoError(t, Initialise(WithFiles("test.yaml")))
	val, err := GetFloat("test")
	assert.NoError(t, err)
	assert.Equal(t, 1.1, val)
}

func TestGetSlice(t *testing.T) {

	// create test file
	os.WriteFile("test.yaml", []byte("slice: [\"test\"]"), 0644)
	defer os.Remove("test.yaml")

	// initialise config
	assert.NoError(t, Initialise(WithFiles("test.yaml")))

	val, err := GetSlice("slice")
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{"test"}, val)

	valString, err := GetString("slice.0")
	assert.NoError(t, err)
	assert.Equal(t, "test", valString)
}

func TestGetMap(t *testing.T) {

	// create test file
	os.WriteFile("test.yaml", []byte("test:\n  test: test"), 0644)
	defer os.Remove("test.yaml")

	// initialise config
	assert.NoError(t, Initialise(WithFiles("test.yaml")))
	val, err := GetMap("test")
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"test": "test"}, val)
}

func TestGetStringSlice(t *testing.T) {

	// create test file
	os.WriteFile("test.yaml", []byte("slice: [\"test\"]"), 0644)
	defer os.Remove("test.yaml")

	// initialise config
	assert.NoError(t, Initialise(WithFiles("test.yaml")))
	val, err := GetStringSlice("slice")
	assert.NoError(t, err)
	assert.Equal(t, []string{"test"}, val)
}

func TestGetStringMap(t *testing.T) {

	// create test file
	os.WriteFile("test.yaml", []byte("test:\n  test: test"), 0644)
	defer os.Remove("test.yaml")

	// initialise config
	assert.NoError(t, Initialise(WithFiles("test.yaml")))
	val, err := GetStringMap("test")
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"test": "test"}, val)
}

func TestGetAll(t *testing.T) {

	// set env variable
	os.Setenv("test", "test")
	defer os.Unsetenv("test")

	// create test file
	os.WriteFile("test.yaml", []byte("test: test"), 0644)
	defer os.Remove("test.yaml")

	// initialise config
	assert.NoError(t, Initialise(WithFiles("test.yaml"), WithEnvVariables("test")))
	val := GetAll()
	assert.Equal(t, map[string]interface{}{"test": "test"}, val)
}

func TestReload(t *testing.T) {

	// create test file
	os.WriteFile("test.yaml", []byte("test: test"), 0644)

	// initialise config
	assert.NoError(t, Initialise(WithFiles("test.yaml")))
	os.Remove("test.yaml")

	os.WriteFile("test.yaml", []byte("test: test-reloaded"), 0644)
	defer os.Remove("test.yaml")

	assert.NoError(t, Reload())
	val, err := GetString("test")
	assert.NoError(t, err)
	assert.Equal(t, "test-reloaded", val)
}
