package provider

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEnvProvider(t *testing.T) {
	assert.NotNil(t, NewEnvProvider([]string{"test"}))
	assert.Equal(t, &envProvider{
		variables: []string{"test"},
	}, NewEnvProvider([]string{"test"}))
}

func TestEnvLoadConfig(t *testing.T) {
	ep := &envProvider{
		variables: []string{"testKey"},
	}

	data := make(map[string]interface{})

	err := ep.LoadConfig(data)
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{}, data)

	os.Setenv("testKey", "testVal")

	err = ep.LoadConfig(data)
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"testKey": "testVal"}, data)
}
