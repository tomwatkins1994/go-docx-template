package docxtpl

import (
	"bytes"
	"encoding/xml"
	"io"
	"sync"

	"github.com/fumiama/go-docx"
	"github.com/tomwatkins1994/go-docx-template/internal/tags"
)

type FumiamaDocx struct {
	*docx.Docx
}

func NewFumiamaDocx(reader io.ReaderAt, size int64) (*FumiamaDocx, error) {
	doc, err := docx.Parse(reader, size)
	if err != nil {
		return nil, err
	}

	return &FumiamaDocx{doc}, nil
}

func (d *FumiamaDocx) GetDocumentXml() (string, error) {
	out, err := xml.Marshal(d.Document.Body)
	if err != nil {
		return "", nil
	}

	return string(out), err
}

func (d *FumiamaDocx) ReplaceDocumentXml(xmlContent string) error {
	decoder := xml.NewDecoder(bytes.NewBufferString(xmlContent))
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if start, ok := t.(xml.StartElement); ok {
			if start.Name.Local == "Body" {
				clear(d.Document.Body.Items)
				err = d.Document.Body.UnmarshalXML(decoder, start)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	return nil
}

func (d *FumiamaDocx) MergeTags() {
	var wg sync.WaitGroup

	for _, item := range d.Document.Body.Items {
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
					containsIncompleteTags := tags.TextContainsIncompleteTags(currentText)
					if containsIncompleteTags {
						text.Text = ""
						inIncompleteTag = true
					} else {
						inIncompleteTag = false
						containsTags := tags.TextContainsTags(currentText)
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
