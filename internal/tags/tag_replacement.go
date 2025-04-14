package tags

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/tomwatkins1994/go-docx-template/internal/xmlutils"
)

func ReplaceTagsInXml(xmlString string, data map[string]any, funcMap template.FuncMap) (string, error) {
	// Prepare the XML for tag replacement
	preparedXmlString, err := xmlutils.PrepareXmlForTagReplacement(xmlString)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("").Funcs(funcMap).Parse(preparedXmlString)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}

	// Fix any issues in the XML
	outputXmlString := xmlutils.FixXmlIssuesPostTagReplacement(buf.String())

	return outputXmlString, err
}
