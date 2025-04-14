package templatedata

import (
	"os"
	"path"
	"slices"
)

func IsFilePath(filepath string) (bool, error) {
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

func IsImageFilePath(filepath string) (bool, error) {
	ext := path.Ext(filepath)
	validExts := []string{".png", ".jpg", ".jpeg"}
	isValid := slices.Contains(validExts, ext)
	if !isValid {
		return false, nil
	}

	isFile, err := IsFilePath(filepath)
	if err != nil {
		return false, err
	}

	return isFile, nil
}
