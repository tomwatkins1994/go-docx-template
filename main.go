package main

import (
	"bytes"
	"log"

	"os"
	"regexp"
	"text/template"

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

type TemplateData struct {
	ProjectNumber string
	Status        string
}

func replaceTagsInParagraph(paragraph *docx.Paragraph) error {
	currentText := ""
	inIncompleteTag := false
	for _, pChild := range paragraph.Children {
		run, ok := pChild.(*docx.Run)
		if ok {
			for _, rChild := range run.Children {
				text, ok := rChild.(*docx.Text)
				if ok {
					if inIncompleteTag {
						currentText += text.Text
					} else {
						currentText = text.Text
					}
					containsIncompleteTags, err := textContainsIncompleteTags(currentText)
					if err != nil {
						return err
					}
					if containsIncompleteTags {
						text.Text = ""
						inIncompleteTag = true
					} else {
						inIncompleteTag = false
						containsTags, err := textContainsTags(currentText)
						if err != nil {
							return err
						}
						if containsTags {
							newText, err := replaceTagsInText(currentText)
							if err != nil {
								return err
							}
							text.Text = newText
						}
					}
				}
			}
		}
	}

	return nil
}

func replaceTagsInText(text string) (string, error) {
	tmpl, err := template.New("").Parse(text)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
		return "", err
	}
	templateData := TemplateData{"B-00001", "New"}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, templateData)
	return buf.String(), err
}

func textContainsTags(text string) (bool, error) {
	r, err := regexp.Compile("{{.+}}")
	if err != nil {
		return false, err
	}
	return r.MatchString(text), err
}

func textContainsIncompleteTags(text string) (bool, error) {
	r, err := regexp.Compile("{{[^}]*(}$|$)|{$|^}")
	if err != nil {
		return false, err
	}
	return r.MatchString(text), err
}
