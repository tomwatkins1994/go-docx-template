package docxtpl

import (
	"testing"

	"github.com/bep/imagemeta"
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

func TestGetImageFormat(t *testing.T) {
	t.Run("JPG should return format", func(t *testing.T) {
		assert := assert.New(t)

		img := &InlineImage{ext: ".jpg"}
		format, err := img.getImageFormat()
		assert.Equal(format, imagemeta.JPEG)
		assert.Nil(err)
	})

	t.Run("JPEG should return format", func(t *testing.T) {
		assert := assert.New(t)

		img := &InlineImage{ext: ".jpeg"}
		format, err := img.getImageFormat()
		assert.Equal(format, imagemeta.JPEG)
		assert.Nil(err)
	})

	t.Run("PNG should return format", func(t *testing.T) {
		assert := assert.New(t)

		img := &InlineImage{ext: ".png"}
		format, err := img.getImageFormat()
		assert.Equal(format, imagemeta.PNG)
		assert.Nil(err)
	})

	t.Run("Not image extension should return error", func(t *testing.T) {
		assert := assert.New(t)

		img := &InlineImage{ext: ".txt"}
		format, err := img.getImageFormat()
		assert.Equal(format, imagemeta.ImageFormat(0))
		assert.NotNil(err)
	})
}
