# Logger Package

The Logger package provides a flexible and configurable logging solution based on [Zap](https://pkg.go.dev/go.uber.org/zap). It allows you to easily integrate structured logging into your Go applications.

## Installation

To use this library in your Go project, you can simply import it:

```go
import "github.com/thegreatforge/gokit/logger"
```

## Getting Started

The Logger package provides a simple and convenient way to set up a logger with various configuration options. It follows the "Functional Options Pattern" for configuration, allowing you to customize the logger's behavior.

### Initializing the Logger

You can initialize the logger with custom configuration options using the `Initialize` function:

```go
err := logger.Initialize(
    logger.Level("debug"),    // Set log level to debug
)
if err != nil {
    panic(err)
}
```

By default, the library initializes a logger with production settings, but you can override these settings as shown above.

### Logging

Once the logger is initialized, you can use it to log messages at different log levels: `Debug`, `Info`, `Warn`, `Error`, `Fatal`, and `Panic`. For example:

```go
logger.Debug("This is a debug message")
logger.Info("This is an info message")
logger.Error("This is an error message")
```

You can also log messages with formatting:

```go
logger.Infof("Formatted message with a value: %d", 42)
```

### Adding Fields

The library allows you to add structured fields to log entries using the `WithField` and `WithFields` methods. Fields provide additional context to your log messages:

```go
logger.WithField("user_id", 123).Info("User logged in")
```

### Using WithContextFields

The `WithContextFields` function helps you retrieve fields from contexts and add them as fields to log entries:

```go
logger.WithContextFields(ctx).Info("Processing request")
```

### Customizing Logging Level

You can change the logging level dynamically using the `SetLevel` function:

```go
logger.SetLevel("info") // Set the logging level to "info"
```

## Contributing

Contributions are welcome! If you find any issues, have suggestions, or want to add new features, feel free to open an issue or submit a pull request on the [GitHub repository](https://github.com/thegreatforge/gokit).