package main

import (
	"bytes"
	"log"
	"text/template"

	"github.com/fumiama/go-docx"
)

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
