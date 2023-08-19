package provider

type IProvider interface {
	LoadConfig(data map[string]interface{}) error
}
