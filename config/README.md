# Config Package

The `config` package is a Go library that provides a simple and flexible way to manage configuration settings for your applications. It allows you to load configuration data from various sources such as files (YAML, JSON), environment variables, and more. This README provides an overview of the features and usage of the `config` package.

## Installation

To use the `config` package in your Go project, you need to import it and install it using the `go get` command:

```bash
go get github.com/thegreatforge/gokit/config
```

## Usage

### Initializing Configuration

To start using the `config` package, you need to initialize it with configuration options. The `Initialise` function is used to initialize the configuration with the desired options. You can specify configuration sources using options like `WithFiles` and `WithEnvVariables`. Here's how you can initialize the configuration:

```go
import (
	"github.com/yourusername/yourproject/config"
)

func main() {
	err := config.Initialise(
		config.WithFiles("config.yaml", "config.json"),
		config.WithEnvVariables("APP_PORT", "DB_HOST"),
	)
	if err != nil {
		// Handle error
	}
	defer config.Close() // Close configuration if needed

	// Your application logic
}
```

### Getting Configuration Values

Once the configuration is initialized, you can easily retrieve configuration values using various getters provided by the package. Here are some examples:

```go
value, err := config.Get("key")
boolValue, err := config.GetBool("bool_key")
intValue, err := config.GetInt("int_key")
floatValue, err := config.GetFloat("float_key")
stringValue, err := config.GetString("string_key")
sliceValue, err := config.GetSlice("slice_key")
stringSliceValue, err := config.GetStringSlice("string_slice_key")
mapValue, err := config.GetMap("map_key")
stringMapValue, err := config.GetStringMap("string_map_key")

allConfigValues := config.GetAll()
```

### Configuration Sources

The `config` package supports loading configuration from various sources:

- **File Sources:** Load configuration from YAML, YML, or JSON files using the `WithFiles` option.

- **Environment Variables:** Load configuration from environment variables using the `WithEnvVariables` option.

## Contributing

Contributions are welcome! If you find any issues, have suggestions, or want to add new features, feel free to open an issue or submit a pull request on the [GitHub repository](https://github.com/thegreatforge/gokit).

## License

This project is licensed under the [MIT License](LICENSE). Feel free to use, modify, and distribute it according to the terms of the license.

---

If you encounter any problems or need assistance, please don't hesitate to reach out by opening an issue on the [GitHub repository](https://github.com/thegreatforge/gokit).