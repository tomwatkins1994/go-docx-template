package docxtpl

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataToMap(t *testing.T) {
	t.Run("Passing in nil should return an error", func(t *testing.T) {
		assert := assert.New(t)

		outputMap, err := dataToMap(nil)
		assert.Nil(outputMap)
		assert.NotNil(err)
	})

	t.Run("Passing in a map should return the map back", func(t *testing.T) {
		assert := assert.New(t)

		inputMap := map[string]any{
			"test": 1,
		}
		outputMap, err := dataToMap(inputMap)
		assert.Equal(outputMap, inputMap)
		assert.Nil(err)
	})
}

func TestConvertStructToMap(t *testing.T) {
	t.Run("Basic struct", func(t *testing.T) {
		data := struct {
			ProjectNumber string
			Client        string
			Status        string
		}{
			ProjectNumber: "B-00001",
			Client:        "TW Software",
			Status:        "New",
		}
		_, err := convertStructToMap(data)
		assert.Nil(t, err)
	})

	t.Run("Struct with nested data", func(t *testing.T) {
		assert := assert.New(t)

		data := struct {
			ProjectNumber string
			Client        string
			Status        string
			People        []struct {
				Name   string
				Gender string
				Age    uint8
			}
		}{
			ProjectNumber: "B-00001",
			Client:        "TW Software",
			Status:        "New",
			People: []struct {
				Name   string
				Gender string
				Age    uint8
			}{
				{
					Name:   "Tom Watkins",
					Gender: "Male",
					Age:    30,
				},
				{
					Name:   "Evie Argyle",
					Gender: "Female",
					Age:    29,
				},
			},
		}
		mapData, err := convertStructToMap(data)
		assert.Nil(err)

		for key, value := range mapData {
			val := reflect.ValueOf(value)
			if val.Kind() == reflect.Slice {
				for i := range val.Len() {
					assert.NotEqual(val.Index(i).Kind(), reflect.Struct, fmt.Sprintf("Found struct in data: %v %v", key, value))
				}
			}
		}
	})
}
