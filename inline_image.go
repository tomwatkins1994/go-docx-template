package docxtpl

import (
	"encoding/xml"
	"strings"

	"github.com/fumiama/go-docx"
)

type InlineImage struct {
	*DocxTmpl
	paragraph *docx.Paragraph
	run       *docx.Run
}

func (d *DocxTmpl) AddImage(filepath string) (*InlineImage, error) {
	paragraph := d.AddParagraph()
	run, err := paragraph.AddInlineDrawingFrom(filepath)
	if err != nil {
		return nil, err
	}

	return &InlineImage{d, paragraph, run}, nil
}

func (i *InlineImage) getImageXml() (string, error) {
	out, err := xml.Marshal(i.run)
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

	//Remove the paragraph from the word doc so we don't get the image twice
	var newItems []interface{}
	for _, item := range i.Document.Body.Items {
		switch o := item.(type) {
		case *docx.Paragraph:
			if o == i.paragraph {
				continue
			}
		}
		newItems = append(newItems, item)
	}
	i.Document.Body.Items = newItems

	return xmlString, err
}
