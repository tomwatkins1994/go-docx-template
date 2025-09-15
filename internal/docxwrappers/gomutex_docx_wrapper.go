package docxwrappers

import (
	"bytes"
	"encoding/xml"
	"io"

	"github.com/gomutex/godocx"
	"github.com/gomutex/godocx/docx"
)

type GomutexDocx struct {
	*docx.RootDoc
}

func NewGomutexDocxFromFilename(filename string) (*GomutexDocx, error) {
	docx, err := godocx.OpenDocument(filename)
	if err != nil {
		return nil, err
	}

	return &GomutexDocx{docx}, nil
}

func (d *GomutexDocx) GetDocumentXml() (string, error) {
	out, err := xml.Marshal(d.Document.Body)
	if err != nil {
		return "", nil
	}

	return string(out), err
}

func (d *GomutexDocx) ReplaceDocumentXml(xmlString string) error {
	decoder := xml.NewDecoder(bytes.NewBufferString(xmlString))
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if start, ok := t.(xml.StartElement); ok {
			if start.Name.Local == "body" {
				clear(d.Document.Body.Children)
				err = d.Document.Body.UnmarshalXML(decoder, start)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	return nil
}
