package docxtpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateInlineImage(t *testing.T) {
	t.Run("Should return error if not a valid image filename", func(t *testing.T) {
		assert := assert.New(t)

		image, err := CreateInlineImage("test_image.txt")
		assert.Nil(image)
		assert.IsType(&InlineImageError{Message: "File is not a valid image"}, err)
		assert.Equal(err.Error(), "Image error: File is not a valid image")
	})

	t.Run("Should return error if not a image doesn't exist", func(t *testing.T) {
		assert := assert.New(t)

		image, err := CreateInlineImage("image_not_exists.png")
		assert.Nil(image)
		assert.IsType(&InlineImageError{Message: "File is not a valid image"}, err)
		assert.Equal(err.Error(), "Image error: File is not a valid image")
	})
}
