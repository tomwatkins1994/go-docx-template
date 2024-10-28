package docxtpl

import (
	"bytes"
	"log"
	"text/template"
)

func replaceTagsInText(text string, data map[string]interface{}, funcMap *template.FuncMap) (string, error) {
	tmpl, err := template.New("").Funcs(*funcMap).Parse(text)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
		return "", err
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), err
}
