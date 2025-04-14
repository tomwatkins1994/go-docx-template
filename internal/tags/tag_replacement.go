package tags

import (
	"bytes"
	"fmt"
	"text/template"
)

func ReplaceTagsInText(text string, data map[string]any, funcMap template.FuncMap) (string, error) {
	tmpl, err := template.New("").Funcs(funcMap).Parse(text)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), err
}
