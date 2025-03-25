package docxtpl

import (
	"reflect"
	"testing"
)

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
		if err != nil {
			t.Fatalf("Error converting basic struct: %v", err)
		}
	})

	t.Run("Struct with nested data", func(t *testing.T) {
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
		if err != nil {
			t.Fatalf("Error converting struct with nested data: %v", err)
		}
		for key, value := range mapData {
			val := reflect.ValueOf(value)
			if val.Kind() == reflect.Slice {
				for i := range val.Len() {
					if val.Index(i).Kind() == reflect.Struct {
						t.Fatalf("Found struct in data: %v %v", key, value)
					}
				}
			}
		}
	})
}
