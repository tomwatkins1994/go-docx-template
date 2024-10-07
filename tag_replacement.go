package main

import (
	"bytes"
	"html/template"
	"log"

	"github.com/fumiama/go-docx"
)

func replaceTagsInParagraph(paragraph *docx.Paragraph, data any) error {
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
							// newText, err := replaceTagsInText(currentText, data)
							// if err != nil {
							// 	return err
							// }
							text.Text = currentText //newText
						}
					}
				}
			}
		}
	}

	return nil
}

func replaceTagsInText(text string, data any) (string, error) {
	tmpl, err := template.New("").Parse(text)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
		return "", err
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)
	return buf.String(), err
}
