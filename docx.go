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
