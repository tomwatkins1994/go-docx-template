package docxtpl

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestParseAndRender(t *testing.T) {
	start := time.Now()
	doc, err := Parse("test.docx")
	if err != nil {
		t.Fatalf("Parsing error: %v", err)
	}
	fmt.Printf("Parse: %v\n", time.Since(start))

	type TemplateData struct {
		ProjectNumber string
		Client        string
		Status        string
		CreatedBy     string
		SignedOffBy   string
		People        []struct {
			Name   string
			Gender string
			Age    uint8
		}
	}
	templateData := TemplateData{
		ProjectNumber: "B-00001",
		Client:        "TW Software",
		Status:        "New",
		CreatedBy:     "Tom Watkins",
		SignedOffBy:   "Tom Watkins",
		People: []struct {
			Name   string
			Gender string
			Age    uint8
		}{
			{
				Name:   "Tom Watkins",
				Gender: "Male",
				Age:    30,
			},
			{
				Name:   "Evie Argyle",
				Gender: "Female",
				Age:    29,
			},
		},
	}
	start = time.Now()
	err = doc.Render(templateData)
	if err != nil {
		t.Fatalf("Rendering error: %v", err)
	}
	fmt.Printf("Render: %v\n", time.Since(start))

	// Create a new file for the output
	start = time.Now()
	f, err := os.Create("generated.docx")
	if err != nil {
		t.Fatalf("Error creating document: %v", err)
		panic(err)
	}
	_, err = doc.WriteTo(f)
	if err != nil {
		t.Fatalf("Error writing to document: %v", err)

	}
	err = f.Close()
	if err != nil {
		t.Fatalf("Error closing created document: %v", err)
	}
	fmt.Printf("Save: %v\n", time.Since(start))
}
