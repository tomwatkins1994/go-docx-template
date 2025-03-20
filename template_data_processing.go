package docxtpl

import (
	"os"
	"path"
	"slices"
)

func processData(d *DocxTmpl, data *any) (map[string]any, error) {
	convertedData, err := dataToMap(data)
	if err != nil {
		return nil, err
	}

	var processTagValues func(d *DocxTmpl, data *map[string]any) error
	processTagValues = func(d *DocxTmpl, data *map[string]any) error {
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
			} else if nestedMap, ok := value.(map[string]any); ok {
				if err := processTagValues(d, &nestedMap); err != nil {
					return err
				}
				(*data)[key] = nestedMap
			} else if sliceValue, ok := value.([]map[string]any); ok {
				for i := range sliceValue {
					if err := processTagValues(d, &sliceValue[i]); err != nil {
						return err
					}
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

	err = processTagValues(d, &convertedData)
	if err != nil {
		return nil, err
	}

	return convertedData, nil
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
