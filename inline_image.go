package docxtpl

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"image"
	"image/png"
	"os"
	"path"
	"strings"

	"github.com/fumiama/go-docx"
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

func (i *InlineImage) addToDocument() (string, error) {
	// Add the image to the document
	paragraph := i.doc.AddParagraph()
	run, err := paragraph.AddInlineDrawing(*i.data)
	if err != nil {
		return "", err
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
