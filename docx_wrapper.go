package docxtpl

type DocxWrapper interface {
	GetDocumentXml() (string, error)
	ReplaceDocumentXml(xml string) error
	MergeTags() error
}
