package main

import (
	"sync"

	"github.com/fumiama/go-docx"
)

func (d *DocxTmpl) mergeTags() error {
	var wg sync.WaitGroup
	errCh := make(chan error)

	for _, item := range d.Document.Body.Items {
		wg.Add(1)
		go func() {
			item := item
			defer wg.Done()
			switch i := item.(type) {
			case *docx.Paragraph:
				errCh <- mergeTagsInParagraph(i)
			case *docx.Table:
				errCh <- mergeTagsInTable(i)
			default:
				errCh <- nil
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return err
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
	paragraphs := make([]*docx.Paragraph, 0, len(table.TableRows)*10)
	for _, row := range table.TableRows {
		for _, cell := range row.TableCells {
			paragraphs = append(paragraphs, cell.Paragraphs...)
		}
	}

	var wg sync.WaitGroup
	errCh := make(chan error)

	for _, paragraph := range paragraphs {
		paragraph := paragraph
		wg.Add(1)
		go func() {
			defer wg.Done()
			errCh <- mergeTagsInParagraph(paragraph)
		}()
	}
	go func() {
		wg.Wait()
		close(errCh)
	}()
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}
