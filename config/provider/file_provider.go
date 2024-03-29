package provider

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"github.com/thegreatforge/gokit/config/errors"

	"gopkg.in/yaml.v3"
)

type fileProvider struct {
	paths     []string
	delimiter string
}

func NewFileProvider(paths []string) IProvider {
	return &fileProvider{
		paths:     paths,
		delimiter: ".",
	}
}

func (fp *fileProvider) LoadConfig(data map[string]interface{}) error {
	for _, path := range fp.paths {
		var configData interface{}

		configFile, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		ext := filepath.Ext(path)
		switch ext {
		case ".yaml", ".yml":
			err = yaml.Unmarshal(configFile, &configData)
			if err != nil {
				return err
			}

		case ".json":
			err = json.Unmarshal(configFile, &configData)
			if err != nil {
				return err
			}

		default:
			return errors.ErrConfigInvalidType
		}

		var exploded map[string]interface{}
		switch t := configData.(type) {
		case map[string]interface{}:
			exploded, err = fp.parseMap(t, "")
			if err != nil {
				return err
			}
		case []interface{}:
			exploded, err = fp.parseSlice(t, "")
			if err != nil {
				return err
			}
		default:
			return errors.ErrConfigFileDataTypeNotSupported
		}

		// merge the data of all the files
		for k, v := range exploded {
			data[k] = v
		}

	}
	return nil
}

func (fp *fileProvider) parseMap(input map[string]interface{}, parent string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for k, i := range input {
		if len(parent) > 0 {
			k = parent + fp.delimiter + k
		}
		switch v := i.(type) {
		case []interface{}:
			out, err := fp.parseSlice(v, k)
			if err != nil {
				return nil, err
			}
			for key, value := range out {
				result[key] = value
			}
			result[k] = v
		case map[string]interface{}:
			out, err := fp.parseMap(v, k)
			if err != nil {
				return nil, err
			}
			for key, value := range out {
				result[key] = value
			}
			result[k] = v
		default:
			result[k] = v
		}
	}
	return result, nil
}

func (fp *fileProvider) parseSlice(input []interface{}, parent string) (map[string]interface{}, error) {
	var key string
	result := make(map[string]interface{})
	if len(input) == 0 {
		key = parent
		result[key] = nil
	}

	for k, i := range input {
		if len(parent) > 0 {
			key = parent + fp.delimiter + strconv.Itoa(k)
		} else {
			key = strconv.Itoa(k)
		}

		switch v := i.(type) {
		case []interface{}:
			out, err := fp.parseSlice(v, key)
			if err != nil {
				return nil, err
			}
			for newkey, value := range out {
				result[newkey] = value
			}
			result[key] = v
		case map[string]interface{}:
			out, err := fp.parseMap(v, key)
			if err != nil {
				return nil, err
			}
			for newkey, value := range out {
				result[newkey] = value
			}
			result[key] = v
		default:
			result[key] = v
		}
	}
	return result, nil
}
