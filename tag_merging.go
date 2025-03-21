package docxtpl

import (
	"sync"

	"github.com/fumiama/go-docx"
)

func (d *DocxTmpl) mergeTags() error {
	var wg sync.WaitGroup
	errCh := make(chan error)

	for _, item := range d.Document.Body.Items {
		item := item
		wg.Add(1)
		go func() {
			defer wg.Done()

			var err error
			switch i := item.(type) {
			case *docx.Paragraph:
				err = mergeTagsInParagraph(i)
			case *docx.Table:
				err = mergeTagsInTable(i)
			}
			if err != nil {
				errCh <- err
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

	return nil
}

func mergeTagsInTable(table *docx.Table) error {
	var wg sync.WaitGroup
	errCh := make(chan error)

	for _, row := range table.TableRows {
		for _, cell := range row.TableCells {
			for _, paragraph := range cell.Paragraphs {
				paragraph := paragraph
				wg.Add(1)
				go func() {
					defer wg.Done()
					err := mergeTagsInParagraph(paragraph)
					if err != nil {
						errCh <- err
					}
				}()
			}
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
