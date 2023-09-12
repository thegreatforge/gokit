# Splitify

## Introduction

Splitify is a Go module designed to facilitate weighted and conditional splitting of items into different categories or handlers. It provides two main implementations: `WeightedSplit` for weighted splitting and `ConditionalSplit` for conditional splitting.

## Installation

To use Splitify in your Go project, you can simply import it using the following import statement:

```go
import "github.com/thegreatforge/gokit/splitify"
```

## Usage

Here's a basic guide on how to use the Splitify module for weighted and conditional splitting:

### Weighted Splitting

1. Import the Splitify module:

```go
import "github.com/thegreatforge/gokit/splitify"
```

2. Create a new `WeightedSplit` instance:

```go
splitter := splitify.NewWeightedSplit()
```

3. Add rules with weights to the splitter:

```go
rule1 := &splitify.Rule{
    Handler: yourHandler1,
    Weight:  5,
}
splitter.AddRule(rule1)

rule2 := &splitify.Rule{
    Handler: yourHandler2,
    Weight:  3,
}
splitter.AddRule(rule2)
```

4. Retrieve the next handler based on the weights:

```go
handler, err := splitter.Next()
if err != nil {
    // Handle the error
} else {
    // Use the selected handler
}
```

### Conditional Splitting

1. Import the Splitify module:

```go
import "github.com/thegreatforge/gokit/splitify"
```

2. Create a new `ConditionalSplit` instance with a default handler:

```go
defaultHandler := yourDefaultHandler
splitter := splitify.NewConditionalSplit(defaultHandler)
```

3. Add rules with conditions to the splitter:

```go
rule1 := &splitify.Rule{
    Handler: yourHandler1,
    Conditions: []splitify.Condition{
        yourCondition1,
        yourCondition2,
    },
}
splitter.AddRule(rule1)

rule2 := &splitify.Rule{
    Handler: yourHandler2,
    Conditions: []splitify.Condition{
        yourCondition3,
    },
}
splitter.AddRule(rule2)
```

4. Retrieve the next handler based on the satisfied conditions:

```go
splitterArg := argumentOfCondition
handler, err := splitter.Next(splitterArg)
if err != nil {
    // Handle the error
} else {
    // Use the selected handler
}
```

## Contributing

Contributions are welcome! If you find any issues, have suggestions, or want to add new features, feel free to open an issue or submit a pull request on the [GitHub repository](https://github.com/thegreatforge/gokit).

---

If you encounter any problems or need assistance, please don't hesitate to reach out by opening an issue on the [GitHub repository](https://github.com/thegreatforge/gokit).