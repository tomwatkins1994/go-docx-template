package main

import (
	"os"
)

func main() {
	doc, err := Parse("test.docx")
	if err != nil {
		panic(err)
	}

	type TemplateData struct {
		ProjectNumber string
		Client        string
		Status        string
		CreatedBy     string
		SignedOffBy   string
	}
	templateData := TemplateData{ProjectNumber: "B-00001", Client: "TW Software", Status: "New", CreatedBy: "Tom Watkins", SignedOffBy: "Tom Watkins"}
	err = doc.Render(templateData)
	if err != nil {
		panic(err)
	}

	// Create a new file for the output
	f, err := os.Create("generated.docx")
	if err != nil {
		panic(err)
	}
	_, err = doc.WriteTo(f)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
}
