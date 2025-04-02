package docxtpl

import (
	"fmt"
	"reflect"
)

func dataToMap(data any) (map[string]any, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	if m, ok := data.(map[string]any); ok {
		return m, nil
	}

	return convertStructToMap(data)
}

func convertStructToMap(s any) (map[string]any, error) {
	result := make(map[string]any)
	val := reflect.ValueOf(s)

	// Check if the input is a pointer and dereference it
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Ensure we have a struct type
	if val.Kind() != reflect.Struct {
		fmt.Printf("is error\n")
		return nil, fmt.Errorf("expected a struct, got %s", val.Kind())
	}

	// Iterate over the struct fields
	for i := range val.NumField() {
		field := val.Type().Field(i)
		value := val.Field(i)

		// Store the field name and value in the map
		if value.Kind() == reflect.Slice {
			newMapSlice := make([]map[string]any, value.Len())
			for j := range value.Len() {
				sliceValue := value.Index(j)
				if sliceValue.Kind() == reflect.Struct {
					newMap, err := convertStructToMap(sliceValue.Interface())
					if err != nil {
						return nil, err
					}
					newMapSlice[j] = newMap
				}
			}
			result[field.Name] = newMapSlice
		} else if value.Kind() == reflect.Struct {
			newMap, err := convertStructToMap(value.Interface())
			if err != nil {
				return nil, err
			}
			result[field.Name] = newMap
		} else {
			result[field.Name] = value.Interface()
		}
	}

	return result, nil
}
