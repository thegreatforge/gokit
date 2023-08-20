package config

import (
	"os"
	"path/filepath"

	"github.com/thegreatforge/gokit/config/errors"
	"github.com/thegreatforge/gokit/config/provider"
)

type app struct {
	data            map[string]interface{}
	configProviders []provider.IProvider
}

type Option func(*app) error

var initialisedApp *app

// Initialise initialises the config
func Initialise(opts ...Option) error {
	if len(opts) == 0 {
		return errors.ErrNoConfigProviders
	}

	initialisedApp = &app{
		data: make(map[string]interface{}),
	}

	for _, opt := range opts {
		if err := opt(initialisedApp); err != nil {
			return err
		}
	}

	for _, provider := range initialisedApp.configProviders {
		err := provider.LoadConfig(initialisedApp.data)
		if err != nil {
			return err
		}
	}

	return nil
}

// Close closes the config goroutines
func Close() {
	// TODO: implement if needed
}

// WithFiles sets the yaml / yml  / json files to load the config from
// paths is a list of paths to load the config from
func WithFiles(paths ...string) Option {
	return func(c *app) error {
		if len(paths) == 0 {
			return errors.ErrNoConfigFiles
		}

		for _, path := range paths {
			_, err := os.Stat(path)
			if err != nil {
				return err
			}

			ext := filepath.Ext(path)
			if ext != ".yaml" && ext != ".yml" && ext != ".json" {
				return errors.ErrInvalidFileType
			}
		}

		c.configProviders = append(c.configProviders, provider.NewFileProvider(paths))
		return nil
	}
}

// WithEnvVariables sets the env variables to load the config from
// variables is a list of env variables to load the config from
func WithEnvVariables(variables ...string) Option {
	return func(c *app) error {

		if len(variables) == 0 {
			return errors.ErrNoEnvVariables
		}

		c.configProviders = append(c.configProviders, provider.NewEnvProvider(variables))
		return nil
	}
}

// Reload reloads the config from the config providers
func Reload() error {
	if initialisedApp == nil {
		return errors.ErrConfigNotInitialised
	}

	newConfig := make(map[string]interface{})

	for _, provider := range initialisedApp.configProviders {
		err := provider.LoadConfig(newConfig)
		if err != nil {
			return err
		}
	}

	initialisedApp.data = newConfig

	return nil
}

// Get returns the config value for the given key
func Get(key string) (interface{}, error) {
	r, exists := initialisedApp.data[key]
	if !exists {
		return nil, errors.ErrConfigNotExists
	}
	return r, nil
}

// GetBool returns the config value for the given key as a bool
func GetBool(key string) (bool, error) {
	r, exists := initialisedApp.data[key]
	if !exists {
		return false, errors.ErrConfigNotExists
	}
	val, ok := r.(bool)
	if !ok {
		return false, errors.ErrConfigInvalidType
	}
	return val, nil
}

// GetInt returns the config value for the given key as an int
func GetInt(key string) (int, error) {
	r, exists := initialisedApp.data[key]
	if !exists {
		return 0, errors.ErrConfigNotExists
	}
	val, ok := r.(int)
	if !ok {
		return 0, errors.ErrConfigInvalidType
	}
	return val, nil
}

// GetFloat returns the config value for the given key as a float
func GetFloat(key string) (float64, error) {
	r, exists := initialisedApp.data[key]
	if !exists {
		return 0, errors.ErrConfigNotExists
	}
	val, ok := r.(float64)
	if !ok {
		return 0, errors.ErrConfigInvalidType
	}
	return val, nil
}

// GetString returns the config value for the given key as a string
func GetString(key string) (string, error) {
	r, exists := initialisedApp.data[key]
	if !exists {
		return "", errors.ErrConfigNotExists
	}
	val, ok := r.(string)
	if !ok {
		return "", errors.ErrConfigInvalidType
	}
	return val, nil
}

// GetSlice returns the config value for the given key as a slice
func GetSlice(key string) ([]interface{}, error) {
	r, exists := initialisedApp.data[key]
	if !exists {
		return nil, errors.ErrConfigNotExists
	}
	val, ok := r.([]interface{})
	if !ok {
		return nil, errors.ErrConfigInvalidType
	}
	return val, nil
}

// GetStringSlice returns the config value for the given key as a string slice
func GetStringSlice(key string) ([]string, error) {
	r, exists := initialisedApp.data[key]
	if !exists {
		return nil, errors.ErrConfigNotExists
	}

	val, ok := r.([]interface{})
	if !ok {
		return nil, errors.ErrConfigInvalidType
	}

	var out []string
	for _, v := range val {
		val, ok := v.(string)
		if !ok {
			return nil, errors.ErrConfigInvalidType
		}
		out = append(out, val)
	}
	return out, nil
}

// GetMap returns the config value for the given key as a map
func GetMap(key string) (map[string]interface{}, error) {
	r, exists := initialisedApp.data[key]
	if !exists {
		return nil, errors.ErrConfigNotExists
	}
	val, ok := r.(map[string]interface{})
	if !ok {
		return nil, errors.ErrConfigInvalidType
	}
	return val, nil
}

// GetStringMap returns the config value for the given key as a map[string]string
func GetStringMap(key string) (map[string]string, error) {
	r, exists := initialisedApp.data[key]
	if !exists {
		return nil, errors.ErrConfigNotExists
	}
	val, ok := r.(map[string]interface{})
	if !ok {
		return nil, errors.ErrConfigInvalidType
	}

	var out = make(map[string]string)
	for k, v := range val {
		val, ok := v.(string)
		if !ok {
			return nil, errors.ErrConfigInvalidType
		}
		out[k] = val
	}
	return out, nil
}

// GetAll returns all the config values
func GetAll() map[string]interface{} {
	return initialisedApp.data
}
