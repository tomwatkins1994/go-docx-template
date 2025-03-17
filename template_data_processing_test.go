package docxtpl

import (
	"reflect"
	"testing"
)

func TestIsFilePath(t *testing.T) {
	tests := []struct {
		name           string
		filepath       string
		expectedResult bool
	}{
		{
			name:           "Existing file",
			filepath:       "test_templates/test_image.png",
			expectedResult: true,
		},
		{
			name:           "Non existent file",
			filepath:       "test_templates/not_exists.docx",
			expectedResult: false,
		},
		{
			name:           "Exists but is a folder",
			filepath:       "test_templates",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := isFilePath(tt.filepath)
			if err != nil {
				t.Fatalf("Unexpected error checking for file path: %v", err)
			}
			if result != tt.expectedResult {
				t.Fatalf("%v - should return %v but returned %v", tt.filepath, tt.expectedResult, result)
			}
		})
	}
}

func TestIsImageFilePath(t *testing.T) {
	tests := []struct {
		name           string
		filepath       string
		expectedResult bool
	}{
		{
			name:           "Existing image",
			filepath:       "test_templates/test_image.png",
			expectedResult: true,
		},
		{
			name:           "File exists but isn't ab image",
			filepath:       "test_templates/test_basic.docx",
			expectedResult: false,
		},
		{
			name:           "Missing file extension",
			filepath:       "test_templates/test_image",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := isImageFilePath(tt.filepath)
			if err != nil {
				t.Fatalf("Unexpected error checking for image file path: %v", err)
			}
			if result != tt.expectedResult {
				t.Fatalf("%v - should return %v but returned %v", tt.filepath, tt.expectedResult, result)
			}
		})
	}
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
