package docxtpl

import (
	"sync"

	"github.com/fumiama/go-docx"
)

func mergeTags(items []any) {
	var wg sync.WaitGroup

	for _, item := range items {
		item := item
		wg.Add(1)
		go func() {
			defer wg.Done()

			switch i := item.(type) {
			case *docx.Paragraph:
				mergeTagsInParagraph(i)
			case *docx.Table:
				mergeTagsInTable(i)
			}
		}()
	}

	wg.Wait()
}

func mergeTagsInParagraph(paragraph *docx.Paragraph) {
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
					containsIncompleteTags := textContainsIncompleteTags(currentText)
					if containsIncompleteTags {
						text.Text = ""
						inIncompleteTag = true
					} else {
						inIncompleteTag = false
						containsTags := textContainsTags(currentText)
						if containsTags {
							text.Text = currentText
						}
					}
				}
			}
		}
	}
}

func mergeTagsInTable(table *docx.Table) {
	var wg sync.WaitGroup

	for _, row := range table.TableRows {
		for _, cell := range row.TableCells {
			for _, paragraph := range cell.Paragraphs {
				paragraph := paragraph
				wg.Add(1)
				go func() {
					defer wg.Done()
					mergeTagsInParagraph(paragraph)
				}()
			}
		}
	}

	wg.Wait()
}
