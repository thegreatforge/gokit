package provider

type fileProvider struct {
	paths []string
}

func (fp *fileProvider) LoadConfig(data map[string]interface{}) error {
	return nil
}

func NewFileProvider(paths []string) Provider {
	return &fileProvider{
		paths: paths,
	}
}
