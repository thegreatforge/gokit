# config package
config package reduces the boilerplate to write the structs of configuration or marshal / unmarshal them.
Currently it supports reading configurations from yaml / yml / json files and environment variables.

### Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/thegreatforge/gokit/config"
)

func main() {
  err := config.Initialise(
		config.WithFiles([]string{"example/config.yaml"}),
		config.WithEnvVariables([]string{"CONFIG_SOURCE"}),
	)
	if err != nil {
		fmt.Printf("failed to init the config with error - %s", err.Error())
		os.Exit(1)
	}
	fmt.Println(config.GetAll())
	fmt.Println(config.GetString("1.hello"), config.GetMap("1"))
}

```