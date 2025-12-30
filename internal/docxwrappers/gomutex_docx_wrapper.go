package docxwrappers

import (
	"bytes"
	"encoding/xml"
	"io"
	"sync"

	"github.com/gomutex/godocx"
	"github.com/gomutex/godocx/common/units"
	"github.com/gomutex/godocx/docx"
	"github.com/gomutex/godocx/wml/ctypes"
	"github.com/tomwatkins1994/go-docx-template/internal/images"
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
		if item.Para != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mergeGomutexTagsInParagraph(item.Para.GetCT())
			}()
		}
		if item.Table != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mergeGomutexTagsInTable(item.Table.GetCT())
			}()
		}
	}

	wg.Wait()
}

func mergeGomutexTagsInParagraph(paragraph *ctypes.Paragraph) {
	currentText := ""
	inIncompleteTag := false
	for _, pChild := range paragraph.Children {
		run := pChild.Run
		for _, rChild := range run.Children {
			text := rChild.Text
			if text == nil {
				continue
			}
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

func mergeGomutexTagsInTable(table *ctypes.Table) {
	var wg sync.WaitGroup

	for _, row := range table.RowContents {
		for _, rowContents := range row.Row.Contents {
			for _, cellContent := range rowContents.Cell.Contents {
				if cellContent.Paragraph != nil {
					wg.Add(1)
					go func() {
						defer wg.Done()
						mergeGomutexTagsInParagraph(cellContent.Paragraph)
					}()
				}
				if cellContent.Table != nil {
					wg.Add(1)
					go func() {
						defer wg.Done()
						mergeGomutexTagsInTable(cellContent.Table)
					}()
				}
			}
		}
	}

	wg.Wait()
}

func (d *GomutexDocx) AddInlineImage(i *images.InlineImage) (xmlString string, err error) {
	// Correctly size the image
	w, h, err := i.GetSizeInches()
	if err != nil {
		return "", err
	}

	// Add the image to the document
	image, err := d.AddPicture(i.Filepath, units.Inch(w), units.Inch(h))
	if err != nil {
		return "", err
	}

	// Get the image XML
	out, err := xml.Marshal(image.Para.GetCT().Children[0].Run.Children[0].Drawing)
	if err != nil {
		return "", err
	}

	// Remove the paragraph from the word doc so we don't get the image twice
	// newItems := []docx.DocumentChild{}
	// for _, item := range d.Document.Body.Children {
	// 	if item.Para == image.Para {
	// 		continue
	// 	}
	// 	newItems = append(newItems, item)
	// }
	// d.Document.Body.Children = newItems

	xmlString = string(out)

	return xmlString, nil
}

func (d *GomutexDocx) Save(w io.Writer) error {
	err := d.Document.Root.Write(w)
	if err != nil {
		return err
	}

	return nil
}
