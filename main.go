package main

import (
	"fmt"
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
	fmt.Println("Plain text:")
	for _, it := range doc.Document.Body.Items {
		switch it.(type) {
		case *docx.Paragraph, *docx.Table: // printable
			fmt.Println(it)
		}
	}
}
