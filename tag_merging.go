package main

import (
	"sync"

	"github.com/fumiama/go-docx"
)

func (d *DocxTmpl) mergeTags() error {
	errCh := make(chan error, len(d.Document.Body.Items))
	var wg sync.WaitGroup

	for _, item := range d.Document.Body.Items {
		wg.Add(1)
		switch i := item.(type) {
		case *docx.Paragraph:
			go mergeTagsInParagraph(i, errCh, &wg)
		case *docx.Table:
			go mergeTagsInTable(i, errCh, &wg)
		default:
			wg.Done()
		}
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

func mergeTagsInParagraph(paragraph *docx.Paragraph, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

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
						errCh <- err
						return
					}
					if containsIncompleteTags {
						text.Text = ""
						inIncompleteTag = true
					} else {
						inIncompleteTag = false
						containsTags, err := textContainsTags(currentText)
						if err != nil {
							errCh <- err
							return
						}
						if containsTags {
							text.Text = currentText
						}
					}
				}
			}
		}
	}

	errCh <- nil
}

func mergeTagsInTable(table *docx.Table, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, row := range table.TableRows {
		for _, cell := range row.TableCells {
			pErrCh := make(chan error, len(cell.Paragraphs))
			var pwg sync.WaitGroup
			for _, paragraph := range cell.Paragraphs {
				pwg.Add(1)
				go mergeTagsInParagraph(paragraph, pErrCh, &pwg)
			}
			go func() {
				pwg.Wait()
				close(pErrCh)
			}()
			for err := range pErrCh {
				if err != nil {
					errCh <- err
					return
				}
			}
		}
	}

	errCh <- nil
}
