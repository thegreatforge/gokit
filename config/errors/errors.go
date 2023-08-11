package errors

type Error string

const (
	ErrConfigNotFound    Error = "config: config not found"
	ErrConfigInvalidType Error = "config: config invalid type"
	ErrNoConfigFiles     Error = "config: no config files"
	ErrFailedToLoadFile  Error = "config: failed to load file"
	ErrFailedToLoadEnv   Error = "config: failed to load env"
)

func (e *Error) Error() string {
	return string(*e)
}

func (e *Error) String() string {
	return e.Error()
}
