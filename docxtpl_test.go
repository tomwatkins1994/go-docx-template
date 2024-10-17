package docxtpl

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestParseAndRender(t *testing.T) {
	t.Run("Basic document", func(t *testing.T) {
		data := struct {
			ProjectNumber string
			Client        string
			Status        string
		}{
			ProjectNumber: "B-00001",
			Client:        "TW Software",
			Status:        "New",
		}
		parseAndRender(t, "test_basic.docx", data)
	})

	t.Run("Document with tables", func(t *testing.T) {
		data := struct {
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
		}{
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
		parseAndRender(t, "test_with_tables.docx", data)
	})
}

func parseAndRender(t *testing.T, filename string, data interface{}) {
	// Parse the document
	start := time.Now()
	doc, err := Parse(filename)
	if err != nil {
		t.Fatalf("%v - Parsing error: %v", t.Name(), err)
	}
	fmt.Printf("%v - Parse: %v\n", t.Name(), time.Since(start))

	// Render the document
	start = time.Now()
	err = doc.Render(data)
	if err != nil {
		t.Fatalf("%v - Rendering error: %v", t.Name(), err)
	}
	fmt.Printf("%v - Render: %v\n", t.Name(), time.Since(start))

	// Create a new file for the output
	start = time.Now()
	f, err := os.Create("generated_" + filename)
	if err != nil {
		t.Fatalf("%v - Error creating document: %v", t.Name(), err)
		panic(err)
	}
	_, err = doc.WriteTo(f)
	if err != nil {
		t.Fatalf("%v - Error writing to document: %v", t.Name(), err)
	}
	err = f.Close()
	if err != nil {
		t.Fatalf("%v - Error closing created document: %v", t.Name(), err)
	}
	fmt.Printf("%v - Save: %v\n", t.Name(), time.Since(start))
}
