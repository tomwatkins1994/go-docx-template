package main

import (
	"bytes"
	"html/template"
	"log"
)

func replaceTagsInText(text string, data any) (string, error) {
	tmpl, err := template.New("").Parse(text)
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
