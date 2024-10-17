package docxtpl

import (
	"fmt"
	"reflect"
)

func (d *DocxTmpl) processData(data interface{}) (map[string]interface{}, error) {
	var err error
	mapData := make(map[string]interface{})
	if m, ok := data.(map[string]interface{}); ok {
		mapData = m
	} else {
		mapData, err = convertStructToMap(data)
		if err != nil {
			return nil, err
		}
	}

	handleTagValues(&mapData)

	return mapData, nil
}

func handleTagValues(data *map[string]interface{}) {
	for key := range *data {
		if key == "Image" {
			(*data)[key] = "Hello"
		}
	}
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
