package docxtpl

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/fumiama/go-docx"
)

type DocxTmpl struct {
	*docx.Docx
	filename string
}

func Parse(filename string) (*DocxTmpl, error) {
	readFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	fileinfo, err := readFile.Stat()
	if err != nil {
		return nil, err
	}
	size := fileinfo.Size()
	doc, err := docx.Parse(readFile, size)
	if err != nil {
		return nil, err
	}

	return &DocxTmpl{doc, filename}, nil
}

func (d *DocxTmpl) Render(data *interface{}) error {
	// Ensure that there are no 'part tags' in the XML document
	err := d.mergeTags()
	if err != nil {
		return err
	}

	// Process the template data
	d.processData(data)
	fmt.Printf("Processed data: %v", data)

	// Get the document XML
	documentXmlString, err := d.getDocumentXml()
	if err != nil {
		return err
	}

	// Prepare the XML for tag replacement
	documentXmlString, err = prepareXmlForTagReplacement(documentXmlString)
	if err != nil {
		return err
	}

	// Replace the tags in XML
	documentXmlString, err = replaceTagsInText(documentXmlString, data)
	if err != nil {
		return err
	}

	// Unmarshal the modified XML and replace the document body with it
	decoder := xml.NewDecoder(bytes.NewBufferString(documentXmlString))
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if start, ok := t.(xml.StartElement); ok {
			if start.Name.Local == "Body" {
				clear(d.Document.Body.Items)
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

func (d *DocxTmpl) getDocumentXml() (string, error) {
	out, err := xml.Marshal(d.Document.Body)
	if err != nil {
		return "", nil
	}

	return string(out), err
}
