package docx_wrappers

import (
	"io"

	"github.com/tomwatkins1994/go-docx-template/internal/images"
)

type DocxWrapper interface {
	GetDocumentXml() (string, error)
	ReplaceDocumentXml(xmlString string) error
	MergeTags()
	AddInlineImage(img *images.InlineImage) (xmlString string, err error)
	Save(w io.Writer) error
}
