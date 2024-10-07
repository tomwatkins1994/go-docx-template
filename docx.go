package main

import (
	"os"

	"github.com/fumiama/go-docx"
)

type DocxTmpl struct {
	*docx.Docx
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

	return &DocxTmpl{doc}, nil
}

func (d *DocxTmpl) Render(data any) error {
	for _, item := range d.Document.Body.Items {
		switch item.(type) {
		case *docx.Paragraph, *docx.Table: // printable
			paragraph, ok := item.(*docx.Paragraph)
			if ok {
				err := replaceTagsInParagraph(paragraph, data)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
