package logger

import (
	"context"

	"go.uber.org/zap"
)

func DefaultLogger() *Logger {
	return defaultLogger
}

func DefaultZapLogger() *zap.Logger {
	return defaultLogger.zapLogger
}

// WithContext adds an error as single field to the Entry.
func (logger *Logger) WithError(err error) *Logger {
	return &Logger{
		zapLogger: logger.zapLogger.With(zap.Error(err)),
		level:     logger.level,
	}
}

// WithField creates an entry from the standard logger and adds a field to
// it.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func (logger *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		zapLogger: logger.zapLogger.Sugar().With(key, value).Desugar(),
		level:     logger.level,
	}
}

// WithField creates an entry from the standard logger and adds a field to
// it.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *Logger {
	return defaultLogger.WithField(key, value)
}

// WithError creates an entry from the standard logger and adds an error to
// it.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithError(err error) *Logger {
	return &Logger{
		zapLogger: defaultLogger.zapLogger.With(zap.Error(err)),
		level:     defaultLogger.level,
	}
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
//
// Note: Use With instead for structured fields for better performance
func (logger *Logger) WithFields(fields map[string]interface{}) *Logger {
	kvs := make([]interface{}, 0)
	for k, v := range fields {
		kvs = append(kvs, k, v)
	}

	return &Logger{
		zapLogger: logger.zapLogger.Sugar().With(kvs...).Desugar(),
		level:     logger.level,
	}
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
//
// Note: Use With instead for structured fields for better performance
func WithFields(fields map[string]interface{}) *Logger {
	return defaultLogger.WithFields(fields)
}

// With adds multiple fields to the logger.
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Logger it returns.
func With(fields ...Field) *Logger {
	return defaultLogger.With(fields...)
}

// With adds multiple fields to the logger.
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Logger it returns.
func (logger *Logger) With(fields ...Field) *Logger {
	return &Logger{
		zapLogger: logger.zapLogger.With(fields...),
		level:     logger.level,
	}
}

// GetLevel returns the current logging level
// Logging levels can be - "debug", "info", "warning", "error", "fatal"
func GetLevel() string {
	return defaultLogger.level.String()
}

// GetLevel returns the current logging level
// Logging levels can be - "debug", "info", "warning", "error", "fatal"
func (logger *Logger) GetLevel() string {
	return logger.level.String()
}

// SetLevel sets the current logging level
func SetLevel(level string) {
	defaultLogger.SetLevel(level)
}

// SetLevel sets the current logging level
func (logger *Logger) SetLevel(level string) {
	lvl := parseLevelString(level)
	logger.level.SetLevel(lvl)
}

// Debug logs a message at DEBUG level using the default logger
func Debug(args ...interface{}) {
	defaultLogger.zapLogger.Sugar().Debug(args...)
}

// Debug logs a message at DEBUG level using the provided logger
func (logger *Logger) Debug(args ...interface{}) {
	logger.zapLogger.Sugar().Debug(args...)
}

// Debugf logs a message at DEBUG level using the default logger
func Debugf(format string, args ...interface{}) {
	defaultLogger.zapLogger.Sugar().Debugf(format, args...)
}

// Debugf logs a message at DEBUG level using the provided logger
func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.zapLogger.Sugar().Debugf(format, args...)
}

// Debugln logs a message at DEBUG level using the default logger
func Debugln(args ...interface{}) {
	defaultLogger.zapLogger.Sugar().Debug(args...)
}

// Debugln logs a message at DEBUG level using the provided logger
func (logger *Logger) Debugln(args ...interface{}) {
	logger.zapLogger.Sugar().Debug(args...)
}

// Print logs a message at INFO level using the default logger.
func Print(args ...interface{}) {
	defaultLogger.zapLogger.Sugar().Info(args...)
}

// Print logs a message at INFO level using the provided logger.
func (logger *Logger) Print(args ...interface{}) {
	logger.zapLogger.Sugar().Info(args...)
}

// Printf logs a message at INFO level using the default logger.
func Printf(format string, args ...interface{}) {
	defaultLogger.zapLogger.Sugar().Infof(format, args...)
}

// Printf logs a message at INFO level using the provided logger.
func (logger *Logger) Printf(format string, args ...interface{}) {
	logger.zapLogger.Sugar().Infof(format, args...)
}

// Println logs a message at INFO level using the default logger
func Println(args ...interface{}) {
	defaultLogger.zapLogger.Sugar().Info(args...)
}

// Println logs a message at INFO level using the provided logger
func (logger *Logger) Println(args ...interface{}) {
	logger.zapLogger.Sugar().Info(args...)
}

// Info logs a message at INFO level using the default logger
func Info(args ...interface{}) {
	defaultLogger.zapLogger.Sugar().Info(args...)
}

// Info logs a message at INFO level using the provided logger
func (logger *Logger) Info(args ...interface{}) {
	logger.zapLogger.Sugar().Info(args...)
}

// Infof logs a message at INFO level using the default logger
func Infof(format string, args ...interface{}) {
	defaultLogger.zapLogger.Sugar().Infof(format, args...)
}

// Infof logs a message at INFO level using the provided logger
func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.zapLogger.Sugar().Infof(format, args...)
}

// Infoln logs a message at INFO level using the default logger
func Infoln(args ...interface{}) {
	defaultLogger.zapLogger.Sugar().Info(args...)
}

// Infoln logs a message at INFO level using the provided logger
func (logger *Logger) Infoln(args ...interface{}) {
	logger.zapLogger.Sugar().Info(args...)
}

// Warn logs a message at WARNING level using the default logger
func Warn(args ...interface{}) {
	defaultLogger.zapLogger.Sugar().Warn(args...)
}

// Warn logs a message at WARNING level using the provided logger
func (logger *Logger) Warn(args ...interface{}) {
	logger.zapLogger.Sugar().Warn(args...)
}

// Warnf logs a message at WARNING level using the default logger
func Warnf(format string, args ...interface{}) {
	defaultLogger.zapLogger.Sugar().Warnf(format, args...)
}

// Warnf logs a message at WARNING level using the provided logger
func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.zapLogger.Sugar().Warnf(format, args...)
}

// Warnln logs a message at WARNING level using the default logger
func Warnln(args ...interface{}) {
	defaultLogger.zapLogger.Sugar().Warn(args...)
}

// Warnln logs a message at WARNING level using the provided logger
func (logger *Logger) Warnln(args ...interface{}) {
	logger.zapLogger.Sugar().Warn(args...)
}

// Error logs a message at ERROR level using the default logger
func Error(args ...interface{}) {
	defaultLogger.zapLogger.Sugar().Error(args...)
}

// Error logs a message at ERROR level using the provided logger
func (logger *Logger) Error(args ...interface{}) {
	logger.zapLogger.Sugar().Error(args...)
}

// Errorf logs a message at ERROR level using the default logger
func Errorf(format string, args ...interface{}) {
	defaultLogger.zapLogger.Sugar().Errorf(format, args...)
}

// Errorf logs a message at ERROR level using the provided logger
func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.zapLogger.Sugar().Errorf(format, args...)
}

// Errorln logs a message at ERROR level using the default logger
func Errorln(args ...interface{}) {
	defaultLogger.zapLogger.Sugar().Error(args...)
}

// Errorln logs a message at ERROR level using the provided logger
func (logger *Logger) Errorln(args ...interface{}) {
	logger.zapLogger.Sugar().Error(args...)
}

// Panic logs a message at PANIC level using the default logger
func Panic(args ...interface{}) {
	defaultLogger.zapLogger.Sugar().Panic(args...)
}

// Panic logs a message at PANIC level using the provided logger
func (logger *Logger) Panic(args ...interface{}) {
	logger.zapLogger.Sugar().Panic(args...)
}

// Panicf logs a message at PANIC level using the default logger
func Panicf(format string, args ...interface{}) {
	defaultLogger.zapLogger.Sugar().Panicf(format, args...)
}

// Panicf logs a message at PANIC level using the given logger
func (logger *Logger) Panicf(format string, args ...interface{}) {
	logger.zapLogger.Sugar().Panicf(format, args...)
}

// Panicln logs a message at PANIC level using the default logger
func Panicln(args ...interface{}) {
	defaultLogger.zapLogger.Sugar().Panic(args...)
}

// Panicln logs a message at PANIC level using the given logger
func (logger *Logger) Panicln(args ...interface{}) {
	logger.zapLogger.Sugar().Panic(args...)
}

// WithContextFields logs messages with the fields added to the context
func WithContextFields(ctx context.Context) *Logger {
	var fields []Field
	if f, ok := ctx.Value(LoggerCtxFieldsKey).([]Field); ok {
		fields = f
	}
	return defaultLogger.With(fields...)
}
