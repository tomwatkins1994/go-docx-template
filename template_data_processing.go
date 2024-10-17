package docxtpl

import (
	"fmt"
	"os"
	"reflect"
)

func (d *DocxTmpl) processData(data interface{}) (map[string]interface{}, error) {
	var err error
	mapData := make(map[string]interface{})
	if m, ok := data.(map[string]interface{}); ok {
		mapData = m
	} else {
		if mapData, err = convertStructToMap(data); err != nil {
			return nil, err
		}
	}

	if err = handleTagValues(&mapData); err != nil {
		return nil, err
	}

	return mapData, nil
}

func handleTagValues(data *map[string]interface{}) error {
	for key, value := range *data {
		if stringVal, ok := value.(string); ok {
			// Check for files
			if isFile, err := isFilePath(stringVal); err != nil {
				return err
			} else {
				if isFile {
					(*data)[key] = "Hello"
				}
			}
		}
	}

	return nil
}

func isFilePath(path string) (bool, error) {
	// Check if the path exists
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	// Check if it's a file
	if info, err := os.Stat(path); err == nil {
		return !info.IsDir(), nil // Return true if it's a file, false if it's a directory
	}

	return false, nil
}

func convertStructToMap(s interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	val := reflect.ValueOf(s)

	// Check if the input is a pointer and dereference it
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Ensure we have a struct type
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %s", val.Kind())
	}

	// Iterate over the struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i)

		// Store the field name and value in the map
		result[field.Name] = value.Interface()
	}

	return result, nil
}
