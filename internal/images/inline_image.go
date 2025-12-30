package images

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/bep/imagemeta"
	"github.com/fumiama/imgsz"
	"github.com/tomwatkins1994/go-docx-template/internal/contenttypes"
	"github.com/tomwatkins1994/go-docx-template/internal/templatedata"
	"golang.org/x/image/draw"
)

const (
	EMUS_PER_INCH = 914400
	DEFAULT_DPI   = 72
)

type InlineImage struct {
	data     *[]byte
	Filepath string
	Ext      string
}

type InlineImageError struct {
	Message string
}

func (e *InlineImageError) Error() string {
	return fmt.Sprintf("Image error: %v", e.Message)
}

// Take a filenane for an image and return a pointer to an InlineImage struct.
// Images can be Jpegs (.jpg or .jpeg) or PNGs
//
//	img, err := CreateInlineImage("example_img.png")
func CreateInlineImage(filepath string) (*InlineImage, error) {
	if isImage, err := templatedata.IsImageFilePath(filepath); err != nil {
		return nil, err
	} else {
		if !isImage {
			return nil, &InlineImageError{"File is not a valid image"}
		}
	}

	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	ext := path.Ext(filepath)

	return &InlineImage{&file, filepath, ext}, nil
}

func (i *InlineImage) getImageFormat() (imagemeta.ImageFormat, error) {
	switch i.Ext {
	case ".jpg", ".jpeg":
		return imagemeta.JPEG, nil
	case ".png":
		return imagemeta.PNG, nil
	default:
		return 0, errors.New("Unknown image format: " + i.Ext)
	}
}

// Return a map of EXIF data from the image.
func (i *InlineImage) GetExifData() (map[string]imagemeta.TagInfo, error) {
	var tags imagemeta.Tags
	handleTag := func(ti imagemeta.TagInfo) error {
		tags.Add(ti)
		return nil
	}

	imageFormat, err := i.getImageFormat()
	if err != nil {
		return nil, err
	}

	shouldHandle := func(ti imagemeta.TagInfo) bool {
		return true
	}

	knownWarnings := []*regexp.Regexp{}

	warnf := func(format string, args ...any) {
		s := fmt.Sprintf(format, args...)
		for _, re := range knownWarnings {
			if re.MatchString(s) {
				return
			}
		}
		panic(errors.New(s))
	}

	sources := imagemeta.EXIF

	err = imagemeta.Decode(imagemeta.Options{R: bytes.NewReader(*i.data), ImageFormat: imageFormat, ShouldHandleTag: shouldHandle, HandleTag: handleTag, Warnf: warnf, Sources: sources})
	if err != nil {
		return nil, err
	}

	return tags.EXIF(), nil
}

// Resize the image. Width and height should be pixel values.
func (i *InlineImage) Resize(width int, height int) error {
	src, err := i.getImage()
	if err != nil {
		return err
	}

	// Resize
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(rgba, rgba.Rect, *src, (*src).Bounds(), draw.Over, nil)
	var resizedImage image.Image = rgba

	err = i.replaceImage(&resizedImage)
	if err != nil {
		return err
	}

	return nil
}

func (i *InlineImage) GetData() *[]byte {
	if i.data == nil {
		return nil
	}
	return i.data
}

func (i *InlineImage) getImage() (*image.Image, error) {
	format, err := i.getImageFormat()
	if err != nil {
		return nil, err
	}

	var img image.Image
	imgReader := bytes.NewReader(*i.GetData())

	switch format {
	case imagemeta.JPEG:
		img, err = jpeg.Decode(imgReader)
	case imagemeta.PNG:
		img, err = png.Decode(imgReader)
	}

	return &img, err
}

func (i *InlineImage) replaceImage(rgba *image.Image) error {
	format, err := i.getImageFormat()
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	switch format {
	case imagemeta.JPEG:
		err = jpeg.Encode(&buf, *rgba, &jpeg.Options{Quality: 100})
	case imagemeta.PNG:
		err = png.Encode(&buf, *rgba)
	}
	if err != nil {
		return err
	}

	newImageData := buf.Bytes()
	i.data = &newImageData

	return nil
}

// Get the size of the image in Inches.
func (i *InlineImage) GetSizeInches() (w int64, h int64, err error) {
	sz, _, err := imgsz.DecodeSize(bytes.NewReader(*i.data))
	if err != nil {
		return 0, 0, nil
	}

	wDpi, hDpi := i.GetResolution()

	w = int64(sz.Width / wDpi)
	h = int64(sz.Height / hDpi)

	return w, h, nil
}

// Get the size of the image in EMUs.
func (i *InlineImage) GetSizeEmus() (w int64, h int64, err error) {
	wInches, hInches, err := i.GetSizeInches()
	if err != nil {
		return 0, 0, err
	}

	w = wInches * int64(EMUS_PER_INCH)
	h = hInches * int64(EMUS_PER_INCH)

	return w, h, nil
}

// Get the resolution (DPI) of the image.
// It gets this from EXIF data and defaults to 72 if not found.
func (i *InlineImage) GetResolution() (wDpi int, hDpi int) {
	exif, err := i.GetExifData()
	if err != nil {
		return DEFAULT_DPI, DEFAULT_DPI
	}

	getResolution := func(tagName string) int {
		resolutionTag, exists := exif[tagName]
		if exists {
			if value, ok := resolutionTag.Value.(string); ok {
				resolution, err := getResolutionFromString(value)
				if err != nil || resolution == 0 {
					return DEFAULT_DPI
				}
				return resolution
			}
		}
		return DEFAULT_DPI
	}

	wDpi, hDpi = getResolution("XResolution"), getResolution("YResolution")

	return wDpi, hDpi
}

func getResolutionFromString(resolution string) (int, error) {
	// Split the string by the slash
	parts := strings.Split(resolution, "/")
	if len(parts) != 2 {
		return 0, errors.New("more than one slash found in image resolution string")
	}

	numerator, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	denominator, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	result := numerator / denominator

	return result, nil
}

func (i *InlineImage) GetContentTypes() ([]*contenttypes.ContentType, error) {
	format, err := i.getImageFormat()
	if err != nil {
		return nil, err
	}

	switch format {
	case imagemeta.JPEG:
		return []*contenttypes.ContentType{&contenttypes.JPG_CONTENT_TYPE, &contenttypes.JPEG_CONTENT_TYPE}, nil
	case imagemeta.PNG:
		return []*contenttypes.ContentType{&contenttypes.PNG_CONTENT_TYPE}, nil
	}

	return []*contenttypes.ContentType{}, nil
}
