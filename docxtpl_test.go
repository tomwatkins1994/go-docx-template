package docxtpl

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
