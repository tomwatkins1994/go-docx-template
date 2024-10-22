package docxtpl

import (
	"encoding/xml"
	"strings"

	"github.com/fumiama/go-docx"
)

type InlineImage struct {
	doc      *DocxTmpl
	filepath string
}

func (d *DocxTmpl) CreateInlineImage(filepath string) (*InlineImage, error) {
	return &InlineImage{d, filepath}, nil
}

func (i *InlineImage) addToDocument() (string, error) {
	// Add the image to the document
	paragraph := i.doc.AddParagraph()
	run, err := paragraph.AddInlineDrawingFrom(i.filepath)
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
