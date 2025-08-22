package docxtpl

import "github.com/tomwatkins1994/go-docx-template/internal/images"

func CreateInlineImage(filepath string) (*images.InlineImage, error) {
	image, err := images.CreateInlineImage(filepath)
	if err != nil {
		return nil, err
	}

	return image, nil
}
