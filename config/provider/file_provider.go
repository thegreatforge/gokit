package provider

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

type fileProvider struct {
	paths     []string
	delimiter string
}

func NewFileProvider(paths []string) Provider {
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
			return errors.New("config file read failure: " + err.Error())
		}

		ext := filepath.Ext(path)
		switch ext {
		case ".yaml", ".yml":
			err = yaml.Unmarshal(configFile, &configData)
			if err != nil {
				return errors.New("config file unmarshal failure: " + err.Error())
			}

		case ".json":
			err = json.Unmarshal(configFile, &configData)
			if err != nil {
				return errors.New("config file unmarshal failure: " + err.Error())
			}

		default:
			return errors.New("config file extension not supported: " + ext)
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
			return errors.New("config file data type not supported: " + ext)
		}

		// merge the data
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
		case nil:
			result[k] = v
		case int:
			result[k] = v
		case float64:
			result[k] = v
		case string:
			result[k] = v
		case bool:
			result[k] = v
		case []interface{}:
			out, err := fp.parseSlice(v, k)
			if err != nil {
				return nil, err
			}
			for key, value := range out {
				result[key] = value
			}
		case map[string]interface{}:
			out, err := fp.parseMap(v, k)
			if err != nil {
				return nil, err
			}
			for key, value := range out {
				result[key] = value
			}
		default:
			//nothing
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
		case nil:
			result[key] = v
		case int:
			result[key] = v
		case float64:
			result[key] = v
		case string:
			result[key] = v
		case bool:
			result[key] = v
		case []interface{}:
			out, err := fp.parseSlice(v, key)
			if err != nil {
				return nil, err
			}
			for newkey, value := range out {
				result[newkey] = value
			}
		case map[string]interface{}:
			out, err := fp.parseMap(v, key)
			if err != nil {
				return nil, err
			}
			for newkey, value := range out {
				result[newkey] = value
			}
		default:
			// do nothing
		}
	}
	return result, nil
}
