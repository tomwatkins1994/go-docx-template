package docxwrappers

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/fumiama/go-docx"
	"github.com/tomwatkins1994/go-docx-template/internal/contenttypes"
	"github.com/tomwatkins1994/go-docx-template/internal/images"
	"github.com/tomwatkins1994/go-docx-template/internal/tags"
)

type FumiamaDocx struct {
	*docx.Docx

	contentTypes *contenttypes.ContentTypes
}

func NewFumiamaDocx(reader io.ReaderAt, size int64) (*FumiamaDocx, error) {
	doc, err := docx.Parse(reader, size)
	if err != nil {
		return nil, err
	}

	contentTypes, err := contenttypes.GetContentTypes(reader, size)
	if err != nil {
		return nil, err
	}

	return &FumiamaDocx{doc, contentTypes}, nil
}

func NewFumiamaDocxFromFilename(filename string) (*FumiamaDocx, error) {
	reader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	fileinfo, err := reader.Stat()
	if err != nil {
		return nil, err
	}
	size := fileinfo.Size()

	docxtpl, err := NewFumiamaDocx(reader, size)
	if err != nil {
		return nil, err
	}

	return docxtpl, nil
}

func (d *FumiamaDocx) GetDocumentXml() (string, error) {
	out, err := xml.Marshal(d.Document.Body)
	if err != nil {
		return "", nil
	}

	return string(out), err
}

func (d *FumiamaDocx) ReplaceDocumentXml(xmlString string) error {
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
	mergeFumiamaTags(d.Document.Body.Items)
}

func mergeFumiamaTags(items []any) {
	var wg sync.WaitGroup

	for _, item := range items {
		wg.Add(1)
		go func(i any) {
			defer wg.Done()

			switch o := i.(type) {
			case *docx.Paragraph:
				mergeTagsInParagraph(o)
			case *docx.Table:
				mergeTagsInTable(o)
			}
		}(item)
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

func (d *FumiamaDocx) AddInlineImage(i *images.InlineImage) (xmlString string, err error) {
	// Add the image to the document
	paragraph := d.AddParagraph()
	run, err := paragraph.AddInlineDrawing(*i.GetData())
	if err != nil {
		return "", err
	}

	// Append the content types
	contentTypes, err := i.GetContentTypes()
	if err != nil {
		return "", err
	}
	for _, contentType := range contentTypes {
		d.contentTypes.AddContentType(contentType)
	}

	// Correctly size the image
	w, h, err := i.GetSize()
	if err != nil {
		return "", err
	}
	for _, child := range run.Children {
		if drawing, ok := child.(*docx.Drawing); ok {
			drawing.Inline.Extent.CX = w
			drawing.Inline.Extent.CY = h
			break
		}
	}

	// Get the image XML
	out, err := xml.Marshal(run)
	if err != nil {
		return "", err
	}

	// Remove run tags as the tag should be in a run already
	xmlString = string(out)
	xmlString = strings.Replace(xmlString, "<w:r>", "", 1)
	xmlString = strings.Replace(xmlString, "<w:rPr></w:rPr>", "", 1)
	lastIndex := strings.LastIndex(xmlString, "</w:r")
	if lastIndex > -1 {
		xmlString = xmlString[:lastIndex]
	}

	// Remove the paragraph from the word doc so we don't get the image twice
	var newItems []interface{}
	for _, item := range d.Document.Body.Items {
		switch o := item.(type) {
		case *docx.Paragraph:
			if o == paragraph {
				continue
			}
		}
		newItems = append(newItems, item)
	}
	d.Document.Body.Items = newItems

	return xmlString, nil
}

func (d *FumiamaDocx) Save(w io.Writer) error {
	var buf bytes.Buffer
	_, err := d.WriteTo(&buf)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(buf.Bytes())
	zipReader, err := zip.NewReader(reader, int64(buf.Len()))
	if err != nil {
		return err
	}

	generatedZip := zip.NewWriter(w)

	for _, f := range zipReader.File {
		newFile, err := generatedZip.Create(f.Name)
		if err != nil {
			return err
		}

		// Override content types with out calculated types
		// Copy across all other files
		if f.Name == "[Content_Types].xml" {
			contentTypesXml, err := d.contentTypes.MarshalXml()
			if err != nil {
				return err
			}

			_, err = newFile.Write([]byte(contentTypesXml))
			if err != nil {
				return err
			}
		} else {
			zf, err := f.Open()
			if err != nil {
				return err
			}
			defer zf.Close()

			if _, err := io.Copy(newFile, zf); err != nil {
				return err
			}
		}
	}

	if err := generatedZip.Close(); err != nil {
		return err
	}

	return nil
}
