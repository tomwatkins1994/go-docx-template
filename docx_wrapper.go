package docxtpl

import "io"

type DocxWrapper interface {
	GetDocumentXml() (string, error)
	ReplaceDocumentXml(xml string) error
	MergeTags()
	Write(w io.Writer) error
}
