package docxtpl

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"image"
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

type InlineImage struct {
	doc  *DocxTmpl
	data *[]byte
	ext  string
}

type InlineImageError struct {
	Message string
}

func (e *InlineImageError) Error() string {
	return fmt.Sprintf("Image error: %v", e.Message)
}

func (d *DocxTmpl) CreateInlineImage(filepath string) (*InlineImage, error) {
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

	return &InlineImage{d, &file, ext}, nil
}

func (i *InlineImage) getImageFormat() (imagemeta.ImageFormat, error) {
	switch i.ext {
	case ".jpg":
		return imagemeta.JPEG, nil
	case ".webp":
		return imagemeta.WebP, nil
	case ".png":
		return imagemeta.PNG, nil
	case ".tif", ".tiff":
		return imagemeta.TIFF, nil
	default:
		return 0, errors.New("Unknown image format: " + i.ext)
	}
}

func (i *InlineImage) getExifData() (map[string]imagemeta.TagInfo, error) {
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

func (i *InlineImage) Resize(width int, height int) error {
	// Decode the current image
	imgReader := bytes.NewReader(*i.data)
	src, err := png.Decode(imgReader)
	if err != nil {
		return err
	}

	// Resize
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	// Encode the new image and replace the image data
	var buf bytes.Buffer
	err = png.Encode(&buf, dst)
	if err != nil {
		return err
	}
	newImageData := buf.Bytes()
	i.data = &newImageData

	return nil
}

func (i *InlineImage) getSize() (int64, int64, error) {
	sz, _, err := imgsz.DecodeSize(bytes.NewReader(*i.data))
	if err != nil {
		return 0, 0, nil
	}

	_EMUS_PER_INCH := 914400

	wDpi, hDpi, err := i.getResolution()
	if err != nil {
		return 0, 0, nil
	}

	w, h := int64(sz.Width), int64(sz.Height)
	w = (w / wDpi) * int64(_EMUS_PER_INCH)
	h = (h / hDpi) * int64(_EMUS_PER_INCH)

	return w, h, nil
}

func (i *InlineImage) getResolution() (int64, int64, error) {
	defaultDpi := 72

	exif, err := i.getExifData()
	if err != nil {
		return 0, 0, nil
	}

	getResolution := func(tagName string) int64 {
		resolutionTag, exists := exif[tagName]
		if exists {
			if value, ok := resolutionTag.Value.(string); ok {
				resolution, err := getResolutionFromString(value)
				if err != nil || resolution == 0 {
					return int64(defaultDpi)
				}
				return int64(resolution)
			}
		}
		return int64(defaultDpi)
	}

	return getResolution("XResolution"), getResolution("YResolution"), nil
}

func getResolutionFromString(resolution string) (int, error) {
	// Split the string by the slash
	parts := strings.Split(resolution, "/")
	if len(parts) != 2 {
		return 0, nil
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

func (i *InlineImage) addToDocument() (string, error) {
	// Add the image to the document
	paragraph := i.doc.AddParagraph()
	run, err := paragraph.AddInlineDrawing(*i.data)
	if err != nil {
		return "", err
	}

	// Correctly size the image
	w, h, err := i.getSize()
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
		return "", nil
	}

	// Remove run tags as the tag should be in a run already
	xmlString := string(out)
	xmlString = strings.Replace(xmlString, "<w:r>", "", 1)
	xmlString = strings.Replace(xmlString, "<w:rPr></w:rPr>", "", 1)
	lastIndex := strings.LastIndex(xmlString, "</w:r")
	if lastIndex > -1 {
		xmlString = xmlString[:lastIndex]
	}

	// Remove the paragraph from the word doc so we don't get the image twice
	var newItems []interface{}
	for _, item := range i.doc.Document.Body.Items {
		switch o := item.(type) {
		case *docx.Paragraph:
			if o == paragraph {
				continue
			}
		}
		newItems = append(newItems, item)
	}
	i.doc.Document.Body.Items = newItems

	return xmlString, err
}
