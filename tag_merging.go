package main

import (
	"github.com/fumiama/go-docx"
)

func (d *DocxTmpl) mergeTags() error {
	for _, item := range d.Document.Body.Items {
		switch i := item.(type) {
		case *docx.Paragraph:
			err := mergeTagsInParagraph(i)
			if err != nil {
				return err
			}
		case *docx.Table:
			err := mergeTagsInTable(i)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func mergeTagsInParagraph(paragraph *docx.Paragraph) error {
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
							text.Text = currentText
						}
					}
				}
			}
		}
	}

	return nil
}

func mergeTagsInTable(table *docx.Table) error {
	for _, row := range table.TableRows {
		for _, cell := range row.TableCells {
			for _, paragraph := range cell.Paragraphs {
				err := mergeTagsInParagraph(paragraph)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
