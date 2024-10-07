package main

import (
	"os"

	"github.com/fumiama/go-docx"
)

func main() {
	readFile, err := os.Open("test.docx")
	if err != nil {
		panic(err)
	}
	fileinfo, err := readFile.Stat()
	if err != nil {
		panic(err)
	}
	size := fileinfo.Size()
	doc, err := docx.Parse(readFile, size)
	if err != nil {
		panic(err)
	}

	for _, item := range doc.Document.Body.Items {
		switch item.(type) {
		case *docx.Paragraph, *docx.Table: // printable
			paragraph, ok := item.(*docx.Paragraph)
			if ok {
				err = replaceTagsInParagraph(paragraph)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	// Create a new file for the output
	f, err := os.Create("generated.docx")
	if err != nil {
		panic(err)
	}
	_, err = doc.WriteTo(f)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
}
