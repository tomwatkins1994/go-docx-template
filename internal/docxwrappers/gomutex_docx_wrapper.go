package docxwrappers

import (
	"bytes"
	"encoding/xml"
	"io"
	"sync"

	"github.com/gomutex/godocx"
	"github.com/gomutex/godocx/docx"
	"github.com/tomwatkins1994/go-docx-template/internal/tags"
)

type GomutexDocx struct {
	*docx.RootDoc
}

func NewGomutexDocxFromFilename(filename string) (*GomutexDocx, error) {
	docx, err := godocx.OpenDocument(filename)
	if err != nil {
		return nil, err
	}

	return &GomutexDocx{docx}, nil
}

func (d *GomutexDocx) GetDocumentXml() (string, error) {
	out, err := xml.Marshal(d.Document.Body)
	if err != nil {
		return "", nil
	}

	return string(out), err
}

func (d *GomutexDocx) ReplaceDocumentXml(xmlString string) error {
	decoder := xml.NewDecoder(bytes.NewBufferString(xmlString))
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if start, ok := t.(xml.StartElement); ok {
			if start.Name.Local == "body" {
				clear(d.Document.Body.Children)
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

func (d *GomutexDocx) MergeTags() {
	mergeGomutexTags(d.Document.Body.Children)
}

func mergeGomutexTags(items []docx.DocumentChild) {
	var wg sync.WaitGroup

	for _, item := range items {
		wg.Add(1)
		go func(i any) {
			defer wg.Done()

			switch o := i.(type) {
			case *docx.Paragraph:
				mergeGomutexTagsInParagraph(o)
				// case *docx.Table:
				// 	mergeTagsInTable(o)
			}
		}(item)
	}

	wg.Wait()
}

func mergeGomutexTagsInParagraph(paragraph *docx.Paragraph) {
	currentText := ""
	inIncompleteTag := false
	for _, pChild := range paragraph.GetCT().Children {
		run := pChild.Run
		for _, rChild := range run.Children {
			text := rChild.Text
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

// func mergeGomutexTagsInTable(table *docx.Table) {
// 	var wg sync.WaitGroup

// 	for _, row := range table. {
// 		for _, cell := range row.TableCells {
// 			for _, paragraph := range cell.Paragraphs {
// 				wg.Add(1)
// 				go func() {
// 					defer wg.Done()
// 					mergeTagsInParagraph(paragraph)
// 				}()
// 			}
// 		}
// 	}

// 	wg.Wait()
// }
