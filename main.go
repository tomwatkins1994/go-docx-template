package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	start := time.Now()
	doc, err := Parse("test.docx")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parse: %v\n", time.Since(start))

	type TemplateData struct {
		ProjectNumber string
		Client        string
		Status        string
		CreatedBy     string
		SignedOffBy   string
	}
	templateData := TemplateData{ProjectNumber: "B-00001", Client: "TW Software", Status: "New", CreatedBy: "Tom Watkins", SignedOffBy: "Tom Watkins"}
	start = time.Now()
	err = doc.Render(templateData)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Render: %v\n", time.Since(start))

	// Create a new file for the output
	start = time.Now()
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
	fmt.Printf("Save: %v\n", time.Since(start))
}
