package docxtpl

import (
	"os"
	"path"
	"slices"
)

func handleTagValues(d *DocxTmpl, data *map[string]any) error {
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
					imageXml, err := image.addToDocument()
					if err != nil {
						return err
					}
					(*data)[key] = imageXml
				}
			}
		} else if sliceValue, ok := value.([]map[string]any); ok {
			for _, val := range sliceValue {
				handleTagValues(d, &val)
			}
		} else if inlineImage, ok := value.(*InlineImage); ok {
			imageXml, err := inlineImage.addToDocument()
			if err != nil {
				return err
			}
			(*data)[key] = imageXml
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
	isValid := slices.Contains(validExts, ext)
	if !isValid {
		return false, nil
	}

	isFile, err := isFilePath(filepath)
	if err != nil {
		return false, err
	}

	return isFile, nil
}
