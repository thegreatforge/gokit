package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/thegreatforge/gokit/config/provider"
	"go.uber.org/zap"
)

type app struct {
	logger          *zap.Logger
	data            map[string]interface{}
	configProviders []provider.Provider
}

type Option func(*app) error

var initialisedApp *app

// Initialise initialises the config
func Initialise(opts ...Option) error {
	zl, _ := zap.NewDevelopment()
	initialisedApp = &app{
		logger: zl,
		data:   make(map[string]interface{}),
	}

	for _, opt := range opts {
		if err := opt(initialisedApp); err != nil {
			return err
		}
	}

	for _, provider := range initialisedApp.configProviders {
		err := provider.LoadConfig(initialisedApp.data)
		if err != nil {
			initialisedApp.logger.Error("failed to load config", zap.Error(err))
			return err
		}
	}

	return nil
}

// Close closes the config goroutines
func Close() {
	// TODO: implement if needed
}

// WithLogger sets the logger for the config
func WithLogger(logger *zap.Logger) Option {
	return func(c *app) error {
		c.logger = logger
		return nil
	}
}

// WithFiles sets the yaml / yml  / json files to load the config from
// paths is a list of paths to load the config from
func WithFiles(paths []string) Option {
	return func(c *app) error {
		if len(paths) == 0 {
			c.logger.Error("no config files provided")
			return errors.New("no config files provided")
		}

		for _, path := range paths {
			_, err := os.Stat(path)
			if err != nil {
				c.logger.Error("failed to load config file ", zap.String("path", path), zap.Error(err))
				return err
			}

			ext := filepath.Ext(path)
			if ext != ".yaml" && ext != ".yml" && ext != ".json" {
				c.logger.Error("invalid config file extension", zap.String("path", path), zap.String("extension", ext))
				return errors.New("invalid config file extension: " + ext)
			}
		}

		c.configProviders = append(c.configProviders, provider.NewFileProvider(paths))
		return nil
	}
}

// WithEnvVariables sets the env variables to load the config from
// variables is a list of env variables to load the config from
func WithEnvVariables(variables []string) Option {
	return func(c *app) error {

		if len(variables) == 0 {
			c.logger.Error("no env variables provided")
			return errors.New("no env variables provided")
		}

		c.configProviders = append(c.configProviders, provider.NewEnvProvider(variables))
		return nil
	}
}

// Get returns the config value for the given key
func Get(key string) interface{} {
	r, exists := initialisedApp.data[key]
	if !exists {
		return nil
	}
	return r
}

// GetBool returns the config value for the given key as a bool
func GetBool(key string) bool {
	r, exists := initialisedApp.data[key]
	if !exists {
		return false
	}
	return r.(bool)
}

// GetInt returns the config value for the given key as an int
func GetInt(key string) int {
	r, exists := initialisedApp.data[key]
	if !exists {
		return 0
	}
	return r.(int)
}

// GetUint returns the config value for the given key as a uint
func GetUint(key string) uint {
	r, exists := initialisedApp.data[key]
	if !exists {
		return 0
	}
	return r.(uint)
}

// GetFloat returns the config value for the given key as a float
func GetFloat(key string) float64 {
	r, exists := initialisedApp.data[key]
	if !exists {
		return 0
	}
	return r.(float64)
}

// GetString returns the config value for the given key as a string
func GetString(key string) string {
	r, exists := initialisedApp.data[key]
	if !exists {
		return ""
	}
	return r.(string)
}

// GetSlice returns the config value for the given key as a slice
func GetSlice(key string) []interface{} {
	r, exists := initialisedApp.data[key]
	if !exists {
		return nil
	}
	return r.([]interface{})
}

// GetMap returns the config value for the given key as a map
func GetMap(key string) map[string]interface{} {
	r, exists := initialisedApp.data[key]
	if !exists {
		return nil
	}
	return r.(map[string]interface{})
}

// GetAll returns all the config values
func GetAll() map[string]interface{} {
	return initialisedApp.data
}
