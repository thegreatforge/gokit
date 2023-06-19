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
	initialisedApp = &app{
		logger: &zap.Logger{},
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
	// TODO: implement
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
		return nil
	}
}
