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

func TestGetResolutionFromString(t *testing.T) {
	t.Run("Get resolution from string", func(t *testing.T) {
		assert := assert.New(t)

		resolutionString := "10/2"
		resolution, err := getResolutionFromString(resolutionString)
		assert.Equal(resolution, 5)
		assert.Nil(err)
	})

	t.Run("More than one slash returns error", func(t *testing.T) {
		assert := assert.New(t)

		resolutionString := "10/2/1"
		resolution, err := getResolutionFromString(resolutionString)
		assert.Equal(resolution, 0)
		assert.NotNil(err)
	})

	t.Run("When numerator is not int, return error", func(t *testing.T) {
		assert := assert.New(t)

		resolutionString := "ten/2"
		resolution, err := getResolutionFromString(resolutionString)
		assert.Equal(resolution, 0)
		assert.NotNil(err)
	})

	t.Run("When denominator is not int, return error", func(t *testing.T) {
		assert := assert.New(t)

		resolutionString := "10/two"
		resolution, err := getResolutionFromString(resolutionString)
		assert.Equal(resolution, 0)
		assert.NotNil(err)
	})
}
