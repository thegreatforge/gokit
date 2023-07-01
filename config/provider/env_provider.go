package provider

import "os"

type envProvider struct {
	variables []string
}

func NewEnvProvider(variables []string) Provider {
	return &envProvider{
		variables: variables,
	}
}

func (ep *envProvider) LoadConfig(data map[string]interface{}) error {
	for _, variable := range ep.variables {
		value, ok := os.LookupEnv(variable)
		if !ok {
			continue
		}

		data[variable] = value
	}

	return nil
}
