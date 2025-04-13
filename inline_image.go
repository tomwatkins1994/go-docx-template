package docxtpl

import (
	"bytes"
	"encoding/xml"
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
	"github.com/fumiama/go-docx"
	"github.com/fumiama/imgsz"
	"golang.org/x/image/draw"
)

const (
	EMUS_PER_INCH = 914400
	DEFAULT_DPI   = 72
)

type InlineImage struct {
	data *[]byte
	Ext  string
}

type InlineImageError struct {
	Message string
}

func (e *InlineImageError) Error() string {
	return fmt.Sprintf("Image error: %v", e.Message)
}

func CreateInlineImage(filepath string) (*InlineImage, error) {
	if isImage, err := isImageFilePath(filepath); err != nil {
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

	return &InlineImage{&file, ext}, nil
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

func (i *InlineImage) getImage() (*image.Image, error) {
	format, err := i.getImageFormat()
	if err != nil {
		return nil, err
	}

	var img image.Image
	imgReader := bytes.NewReader(*i.data)

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

// Get the size of the image in EMUs.
func (i *InlineImage) GetSize() (w int64, h int64, err error) {
	sz, _, err := imgsz.DecodeSize(bytes.NewReader(*i.data))
	if err != nil {
		return 0, 0, nil
	}

	wDpi, hDpi := i.GetResolution()

	w = int64(sz.Width/wDpi) * int64(EMUS_PER_INCH)
	h = int64(sz.Height/hDpi) * int64(EMUS_PER_INCH)

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

func (i *InlineImage) getContentTypes() ([]*ContentType, error) {
	format, err := i.getImageFormat()
	if err != nil {
		return nil, err
	}

	switch format {
	case imagemeta.JPEG:
		return []*ContentType{&JPG_CONTENT_TYPE, &JPEG_CONTENT_TYPE}, nil
	case imagemeta.PNG:
		return []*ContentType{&PNG_CONTENT_TYPE}, nil
	}

	return []*ContentType{}, nil
}

func (d *DocxTmpl) addInlineImage(i *InlineImage) (xmlString string, err error) {
	// Add the image to the document
	paragraph := d.AddParagraph()
	run, err := paragraph.AddInlineDrawing(*i.data)
	if err != nil {
		return "", err
	}

	// Append the content types
	contentTypes, err := i.getContentTypes()
	if err != nil {
		return "", err
	}
	for _, contentType := range contentTypes {
		d.contentTypes.addContentType(contentType)
	}

	// Correctly size the image
	w, h, err := i.GetSize()
	if err != nil {
		return "", err
	}
	for _, child := range run.Children {
		if drawing, ok := child.(*docx.Drawing); ok {
			drawing.Inline.Extent.CX = w
			drawing.Inline.Extent.CY = h
			break
		}
	}

	// Get the image XML
	out, err := xml.Marshal(run)
	if err != nil {
		return "", err
	}

	// Remove run tags as the tag should be in a run already
	xmlString = string(out)
	xmlString = strings.Replace(xmlString, "<w:r>", "", 1)
	xmlString = strings.Replace(xmlString, "<w:rPr></w:rPr>", "", 1)
	lastIndex := strings.LastIndex(xmlString, "</w:r")
	if lastIndex > -1 {
		xmlString = xmlString[:lastIndex]
	}

	// Remove the paragraph from the word doc so we don't get the image twice
	var newItems []interface{}
	for _, item := range d.Document.Body.Items {
		switch o := item.(type) {
		case *docx.Paragraph:
			if o == paragraph {
				continue
			}
		}
		newItems = append(newItems, item)
	}
	d.Document.Body.Items = newItems

	return xmlString, nil
}
