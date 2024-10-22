package docxtpl

import (
	"testing"
)

func TestIsFilePath(t *testing.T) {
	t.Run("Existing image", func(t *testing.T) {
		filepath := "test_templates/test_image.png"
		exists, err := isFilePath(filepath)
		if err != nil {
			t.Fatalf("Error checking for file path: %v", err)
		}
		if exists == false {
			t.Fatalf("File should exist: %v", filepath)
		}
	})

	t.Run("Non existent image", func(t *testing.T) {
		filepath := "test_templates/not_exists.png"
		exists, err := isFilePath(filepath)
		if err != nil {
			t.Fatalf("Error checking for file path: %v", err)
		}
		if exists {
			t.Fatalf("File should not exist: %v", filepath)
		}
	})
}
