package images

import (
	"image"
	"testing"

	"github.com/bep/imagemeta"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tomwatkins1994/go-docx-template/internal/contenttypes"
)

var testJpgImage string = "../../test_templates/test_image.jpg"
var testPngImage string = "../../test_templates/test_image.png"

func TestCreateInlineImage(t *testing.T) {
	t.Run("Should return an image for a valid filepath", func(t *testing.T) {
		assert := assert.New(t)

		img, err := CreateInlineImage(testJpgImage)
		assert.Nil(err)
		assert.NotNil(img.data)
		assert.Equal(img.Ext, ".jpg")
	})

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

		img := &InlineImage{Ext: ".jpg"}
		format, err := img.getImageFormat()
		assert.Equal(format, imagemeta.JPEG)
		assert.Nil(err)
	})

	t.Run("JPEG should return format", func(t *testing.T) {
		assert := assert.New(t)

		img := &InlineImage{Ext: ".jpeg"}
		format, err := img.getImageFormat()
		assert.Equal(format, imagemeta.JPEG)
		assert.Nil(err)
	})

	t.Run("PNG should return format", func(t *testing.T) {
		assert := assert.New(t)

		img := &InlineImage{Ext: ".png"}
		format, err := img.getImageFormat()
		assert.Equal(format, imagemeta.PNG)
		assert.Nil(err)
	})

	t.Run("Not image extension should return error", func(t *testing.T) {
		assert := assert.New(t)

		img := &InlineImage{Ext: ".txt"}
		format, err := img.getImageFormat()
		assert.Equal(format, imagemeta.ImageFormat(0))
		assert.NotNil(err)
	})
}

func TestGetExifData(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	img, err := CreateInlineImage(testJpgImage)
	require.Nil(err)

	exifData, err := img.GetExifData()
	assert.Nil(err)
	assert.Greater(len(exifData), 0)
}

func TestResize(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	inlineImage, err := CreateInlineImage(testJpgImage)
	require.Nil(err)

	originalWEmu, originalHEmu, err := inlineImage.GetSize()
	require.Nil(err)

	wDpi, hDpi := inlineImage.GetResolution()
	newWidthPx := int(originalWEmu/EMUS_PER_INCH) * wDpi * 2
	newHeightPx := int(originalHEmu/EMUS_PER_INCH) * hDpi * 2

	err = inlineImage.Resize(newWidthPx, newHeightPx)
	assert.Nil(err)

	w, h, err := inlineImage.GetSize()
	assert.Nil(err)
	assert.Equal(w, originalWEmu*2)
	assert.Equal(h, originalHEmu*2)
}

func TestGetImage(t *testing.T) {
	t.Run("Get a jpeg image", func(t *testing.T) {
		assert := assert.New(t)

		inlineImage, err := CreateInlineImage(testJpgImage)
		assert.Nil(err)

		img, err := inlineImage.getImage()
		assert.Nil(err)
		assert.NotNil(img)
	})

	t.Run("Get a png image", func(t *testing.T) {
		assert := assert.New(t)

		inlineImage, err := CreateInlineImage(testPngImage)
		assert.Nil(err)

		img, err := inlineImage.getImage()
		assert.Nil(err)
		assert.NotNil(img)
	})

	t.Run("Return an error for an invalid image", func(t *testing.T) {
		assert := assert.New(t)

		inlineImage := &InlineImage{Ext: ".txt"}
		_, err := inlineImage.getImage()
		assert.NotNil(err)
	})
}

func TestReplaceImage(t *testing.T) {
	t.Run("Replace a jpeg image", func(t *testing.T) {
		assert := assert.New(t)

		inlineImage, err := CreateInlineImage(testJpgImage)
		assert.Nil(err)

		rgba := image.NewRGBA(image.Rect(0, 0, 100, 100))
		var resizedImage image.Image = rgba
		err = inlineImage.replaceImage(&resizedImage)
		assert.Nil(err)
	})

	t.Run("Replace a png image", func(t *testing.T) {
		assert := assert.New(t)

		inlineImage, err := CreateInlineImage(testPngImage)
		assert.Nil(err)

		rgba := image.NewRGBA(image.Rect(0, 0, 100, 100))
		var resizedImage image.Image = rgba
		err = inlineImage.replaceImage(&resizedImage)
		assert.Nil(err)
	})
}

func TestGetSize(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	inlineImage, err := CreateInlineImage(testJpgImage)
	require.Nil(err)

	w, h, err := inlineImage.GetSize()
	assert.Nil(err)
	assert.Greater(w, int64(0))
	assert.Greater(h, int64(0))
}

func TestGetResolution(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	inlineImage, err := CreateInlineImage(testJpgImage)
	require.Nil(err)

	wDpi, hDpi := inlineImage.GetResolution()
	assert.Greater(wDpi, 0)
	assert.Greater(hDpi, 0)
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

func TestContentTypes(t *testing.T) {
	t.Run("Get content types for a jpeg image", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		inlineImage, err := CreateInlineImage(testJpgImage)
		require.Nil(err)

		contentTypes, err := inlineImage.GetContentTypes()
		assert.Nil(err)
		assert.Equal(contentTypes[0], &contenttypes.JPG_CONTENT_TYPE)
		assert.Equal(contentTypes[1], &contenttypes.JPEG_CONTENT_TYPE)
	})

	t.Run("Get content types for a png image", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		inlineImage, err := CreateInlineImage(testPngImage)
		require.Nil(err)

		contentTypes, err := inlineImage.GetContentTypes()
		assert.Nil(err)
		assert.Equal(contentTypes[0], &contenttypes.PNG_CONTENT_TYPE)
	})
}
