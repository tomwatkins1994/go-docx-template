package docxtpl

import (
	"fmt"
	"reflect"
)

func dataToMap(data any) (map[string]any, error) {
	if m, ok := data.(map[string]any); ok {
		return m, nil
	} else {
		if mapData, err := convertStructToMap(data); err != nil {
			return nil, err
		} else {
			return mapData, nil
		}
	}

	// if err = handleTagValues(d, &mapData); err != nil {
	// 	return nil, err
	// }
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
				newMap, err := convertStructToMap(value.Index(j).Interface())
				if err != nil {
					return nil, err
				}
				newMapSlice[j] = newMap
			}
			result[field.Name] = newMapSlice
		} else {
			result[field.Name] = value.Interface()
		}
	}

	return result, nil
}
