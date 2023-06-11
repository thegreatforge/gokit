# logger

logger package is a wrapper over uber/zap logger package.


### Usage

```go
 // initialization
err = logger.Initialize()

if err != nil {
	fmt.Println("failed to initialise logg")
	os.Exit(1)
}

logger.Info("starting service: magicnotes")
```