package provider

type Provider interface {
	LoadConfig(data map[string]interface{}) error
}
