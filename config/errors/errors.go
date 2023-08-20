package errors

type Error string

const (
	ErrNoConfigFiles                  Error = "config: no config files"
	ErrFailedToLoadFile               Error = "config: failed to load file"
	ErrInvalidFileType                Error = "config: invalid file type"
	ErrNoEnvVariables                 Error = "config: no env variables"
	ErrConfigNotExists                Error = "config: config not exists"
	ErrConfigInvalidType              Error = "config: invalid type"
	ErrConfigFileDataTypeNotSupported Error = "config: config file data type not supported"
	ErrNoConfigProviders              Error = "config: no config providers"
	ErrConfigNotInitialised           Error = "config: config not initialised"
)

func (e Error) Error() string {
	return string(e)
}

func (e Error) String() string {
	return e.Error()
}
