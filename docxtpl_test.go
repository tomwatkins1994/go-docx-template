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

	t.Run("Basic document with image", func(t *testing.T) {
		data := struct {
			ProjectNumber string
			Client        string
			Status        string
			Image         string
		}{
			ProjectNumber: "B-00001",
			Client:        "TW Software",
			Status:        "New",
			Image:         "test_templates/test_image.png",
		}
		parseAndRender(t, "test_basic_with_images.docx", data)
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
	start := time.Now()

	// Parse the document
	parseStart := time.Now()
	doc, err := Parse("test_templates/" + filename)
	if err != nil {
		t.Fatalf("%v - Parsing error: %v", t.Name(), err)
	}
	fmt.Printf("%v - Parse: %v\n", t.Name(), time.Since(parseStart))

	// Render the document
	renderStart := time.Now()
	err = doc.Render(data)
	if err != nil {
		t.Fatalf("%v - Rendering error: %v", t.Name(), err)
	}
	fmt.Printf("%v - Render: %v\n", t.Name(), time.Since(renderStart))

	// Create a new file for the output
	saveStart := time.Now()
	f, err := os.Create("test_templates/generated_" + filename)
	if err != nil {
		t.Fatalf("%v - Error creating document: %v", t.Name(), err)
	}
	_, err = doc.WriteTo(f)
	if err != nil {
		t.Fatalf("%v - Error writing to document: %v", t.Name(), err)
	}
	err = f.Close()
	if err != nil {
		t.Fatalf("%v - Error closing created document: %v", t.Name(), err)
	}
	fmt.Printf("%v - Save: %v\n", t.Name(), time.Since(saveStart))

	// Log the overall time taken
	fmt.Printf("%v - Total: %v\n", t.Name(), time.Since(start))
}
