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
	f, err := os.Create("generated.docx")
	if err != nil {
		panic(err)
	}
	fmt.Println("Plain text:")
	for _, item := range doc.Document.Body.Items {
		switch item.(type) {
		case *docx.Paragraph, *docx.Table: // printable
			p, ok := item.(*docx.Paragraph)
			if ok {
				fmt.Println("Found paragraph")
				for _, pChild := range p.Children {
					run, ok := pChild.(*docx.Run)
					if ok {
						for i, runChild := range run.Children {
							text, ok := runChild.(*docx.Text)
							if ok {
								text.Text = "Hello"
								fmt.Printf("Found run text %v: %v", i, text.Text)
							}
						}
					}
				}
			}
			fmt.Println(item)
		}
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
