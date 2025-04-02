package docxtpl

import (
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
		assert := assert.New(t)

		data := struct {
			ProjectNumber string
			Client        string
			Status        string
		}{
			ProjectNumber: "B-00001",
			Client:        "TW Software",
			Status:        "New",
		}
		outputMap, err := convertStructToMap(data)
		assert.Equal(map[string]any{
			"ProjectNumber": "B-00001",
			"Client":        "TW Software",
			"Status":        "New",
		}, outputMap)
		assert.Nil(err)
	})

	t.Run("Struct with nested structs and slices", func(t *testing.T) {
		assert := assert.New(t)

		data := struct {
			ProjectNumber string
			Client        string
			Status        string
			ExtraFields   struct {
				Field1 string
				Field2 string
			}
			People []struct {
				Name   string
				Gender string
				Age    uint8
			}
		}{
			ProjectNumber: "B-00001",
			Client:        "TW Software",
			Status:        "New",
			ExtraFields: struct {
				Field1 string
				Field2 string
			}{
				Field1: "Value 1",
				Field2: "Value 2",
			},
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
		outputMap, err := convertStructToMap(data)
		assert.Equal(map[string]any{
			"ProjectNumber": "B-00001",
			"Client":        "TW Software",
			"Status":        "New",
			"ExtraFields": map[string]any{
				"Field1": "Value 1",
				"Field2": "Value 2",
			},
			"People": []map[string]any{
				{
					"Name":   "Tom Watkins",
					"Gender": "Male",
					"Age":    uint8(30),
				},
				{
					"Name":   "Evie Argyle",
					"Gender": "Female",
					"Age":    uint8(29),
				},
			},
		}, outputMap)
		assert.Nil(err)
	})

	t.Run("Pointer to a struct", func(t *testing.T) {
		assert := assert.New(t)

		data := struct {
			ProjectNumber string
			Client        string
			Status        string
		}{
			ProjectNumber: "B-00001",
			Client:        "TW Software",
			Status:        "New",
		}
		outputMap, err := convertStructToMap(&data)
		assert.Equal(map[string]any{
			"ProjectNumber": "B-00001",
			"Client":        "TW Software",
			"Status":        "New",
		}, outputMap)
		assert.Nil(err)
	})

	t.Run("Passing in a non struct value should return error", func(t *testing.T) {
		assert := assert.New(t)

		outputMap, err := convertStructToMap("string")
		assert.Nil(outputMap)
		assert.NotNil(t, err)
	})

	t.Run("Passing in nil should return an error", func(t *testing.T) {
		assert := assert.New(t)

		outputMap, err := convertStructToMap(nil)
		assert.Nil(outputMap)
		assert.NotNil(err)
	})
}
