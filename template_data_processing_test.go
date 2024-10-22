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

}
