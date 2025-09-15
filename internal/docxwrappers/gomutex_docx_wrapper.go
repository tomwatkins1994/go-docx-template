package docxwrappers

import (
	"github.com/gomutex/godocx"
	"github.com/gomutex/godocx/docx"
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
