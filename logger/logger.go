package logger

import (
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zapLogger *zap.Logger
	level     *zap.AtomicLevel
}

type config struct {
	zap.Config
	callerSkip int
}

type option func(*config) error

var defaultLogger *Logger

func init() {
	var err error
	defaultLogger, err = NewLogger()
	if err != nil {
		panic(err)
	}
}

// Initialize the logger with configuration options
// Follows the "Functional Options Pattern" for configuration. Described at:
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func Initialize(options ...option) (err error) {
	defaultLogger, err = NewLogger(options...)
	if err != nil {
		return err
	}

	return nil
}

func NewLogger(opts ...option) (*Logger, error) {
	conf := &config{zap.NewProductionConfig(), 1} // default skip of 1
	opts = append(opts, timestamp())
	for _, opt := range opts {
		err := opt(conf)

		if err != nil {
			return nil, err
		}
	}

	zapLogger, err := conf.Build(
		zap.AddStacktrace(zapcore.DPanicLevel), // Automatically add stack trace above this level
		zap.AddCallerSkip(conf.callerSkip),     // skip configured level to print the correct caller
	)
	if err != nil {
		return nil, err
	}

	l := &Logger{
		zapLogger: zapLogger,
		level:     &conf.Level,
	}

	return l, nil

}

// Formatter sets the log format to development friendly text or production friendly JSON
// based on the config value
// Allowed values: "text", "json"
func Formatter(format string) option {
	return func(conf *config) error {
		switch format {
		case "text", "ascii", "terminal":
			conf.Encoding = "console"
		default:
			conf.Encoding = "json"

		}

		return nil
	}
}

// Level sets the logging level.
// Allowed values in decreasing order of verbosity - "debug", "info", "warning", "error", "fatal", "panic"
// Setting a higher logging level will ignore logs from lower levels in output
// Default: "info"
func Level(level string) option {
	return func(conf *config) error {
		lvl := parseLevelString(level)
		conf.Level.SetLevel(lvl)

		return nil
	}
}

// CallerSkip add specified skip frames while reporting caller
// should be used to make caller to appropriate frame if this pkg
// is not used in vanilla form. Multiple calls are not additive.
func CallerSkip(skip int) option {
	return func(conf *config) error {
		conf.callerSkip = skip

		return nil
	}
}

// timestamp sets the time key name and format
func timestamp() option {
	return func(conf *config) error {
		conf.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339Nano)
		conf.EncoderConfig.TimeKey = "@timestamp"
		return nil
	}
}

// Sync flushes any buffered log entries
func Sync() error {
	return defaultLogger.zapLogger.Sync()
}

// Sync flushes any buffered log entries
func (logger *Logger) Sync() error {
	return logger.zapLogger.Sync()
}

func parseLevelString(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}
