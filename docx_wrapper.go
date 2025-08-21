package docxtpl

import "io"

type DocxWrapper interface {
	GetDocumentXml() (string, error)
	ReplaceDocumentXml(xmlString string) error
	MergeTags()
	AddInlineImage(img *InlineImage) (xmlString string, err error)
	Save(w io.Writer) error
}
