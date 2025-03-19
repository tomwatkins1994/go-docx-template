package docxtpl

import (
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
