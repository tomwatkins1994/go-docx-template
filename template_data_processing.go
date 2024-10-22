package docxtpl

import (
	"fmt"
	"os"
	"path"
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

	if err = handleTagValues(d, &mapData); err != nil {
		return nil, err
	}

	return mapData, nil
}

func handleTagValues(d *DocxTmpl, data *map[string]interface{}) error {
	for key, value := range *data {
		if stringVal, ok := value.(string); ok {
			// Check for files
			if isImage, err := isImageFilePath(stringVal); err != nil {
				return err
			} else {
				if isImage {
					image, err := d.CreateInlineImage(stringVal)
					if err != nil {
						return err
					}
					imageXml, err := image.getImageXml()
					if err != nil {
						return err
					}
					(*data)[key] = imageXml
				}
			}
		} else if sliceValue, ok := value.([]map[string]interface{}); ok {
			for _, val := range sliceValue {
				handleTagValues(d, &val)
			}
		}
	}

	return nil
}

func isFilePath(filepath string) (bool, error) {
	// Check if the path exists
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	// Check if it's a file
	if info, err := os.Stat(filepath); err == nil {
		return !info.IsDir(), nil // Return true if it's a file, false if it's a directory
	}

	return false, nil
}

func isImageFilePath(filepath string) (bool, error) {
	ext := path.Ext(filepath)
	validExts := []string{".png", ".jpg", ".jpeg"}
	isValid := false
	for _, v := range validExts {
		if ext == v {
			isValid = true
			break
		}
	}
	if !isValid {
		return false, nil
	}

	isFile, err := isFilePath(filepath)
	if err != nil {
		return false, err
	}

	return isFile, nil
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
		if value.Kind() == reflect.Slice {
			newMapSlice := make([]map[string]interface{}, value.Len())
			for i := 0; i < value.Len(); i++ {
				newMap, err := convertStructToMap(value.Index(i).Interface())
				if err != nil {
					return nil, err
				}
				newMapSlice[i] = newMap
			}
			result[field.Name] = newMapSlice
		} else {
			result[field.Name] = value.Interface()
		}
	}

	return result, nil
}
