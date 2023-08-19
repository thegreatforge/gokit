package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {

	assert.Equal(t, "config: no config files", ErrNoConfigFiles.Error())
	assert.Equal(t, "config: failed to load file", ErrFailedToLoadFile.Error())
	assert.Equal(t, "config: invalid file type", ErrInvalidFileType.Error())
	assert.Equal(t, "config: no env variables", ErrNoEnvVariables.Error())
	assert.Equal(t, "config: config not exists", ErrConfigNotExists.Error())
	assert.Equal(t, "config: invalid type", ErrConfigInvalidType.Error())
	assert.Equal(t, "config: config file data type not supported", ErrConfigFileDataTypeNotSupported.Error())

}

func TestString(t *testing.T) {

	assert.Equal(t, "config: no config files", ErrNoConfigFiles.String())
	assert.Equal(t, "config: failed to load file", ErrFailedToLoadFile.String())
	assert.Equal(t, "config: invalid file type", ErrInvalidFileType.String())
	assert.Equal(t, "config: no env variables", ErrNoEnvVariables.String())
	assert.Equal(t, "config: config not exists", ErrConfigNotExists.String())
	assert.Equal(t, "config: invalid type", ErrConfigInvalidType.String())
	assert.Equal(t, "config: config file data type not supported", ErrConfigFileDataTypeNotSupported.String())

}
