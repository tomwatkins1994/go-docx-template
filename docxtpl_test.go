package docxtpl

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type test struct {
	name           string
	filename       string
	outputFilename string
	data           any
	fns            map[string]any
}

var myMap = make(map[string]int)

func getTests() ([]test, error) {
	testImage, err := CreateInlineImage("test_templates/test_image.png")
	if err != nil {
		return nil, err
	}
	profileImage, err := CreateInlineImage("test_templates/test_image.jpg")
	if err != nil {
		return nil, err
	}

	var tests = []test{
		{
			name:     "Basic document",
			filename: "test_basic.docx",
			data: struct {
				ProjectNumber string
				Client        string
				Status        string
			}{
				ProjectNumber: "B-00001",
				Client:        "TW Software",
				Status:        "New",
			},
		},
		{
			name:     "Basic document with images",
			filename: "test_basic_with_images.docx",
			data: struct {
				ProjectNumber string
				Client        string
				Status        string
				ImagePng      string
				ImageJpg      string
			}{
				ProjectNumber: "B-00001",
				Client:        "TW Software",
				Status:        "New",
				ImagePng:      "test_templates/test_image.png",
				ImageJpg:      "test_templates/test_image.png",
			},
		},
		{
			name:     "Basic document with tables",
			filename: "test_with_tables.docx",
			data: struct {
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
			},
		},
		{
			name:     "Basic document with tables and images",
			filename: "test_with_tables_and_images.docx",
			data: struct {
				ProjectNumber string
				Client        string
				Status        string
				Image         string
				People        []struct {
					Name           string
					Gender         string
					Age            uint8
					ProfilePicture string
				}
			}{
				ProjectNumber: "B-00001",
				Client:        "TW Software",
				Status:        "New",
				Image:         "test_templates/test_image.png",
				People: []struct {
					Name           string
					Gender         string
					Age            uint8
					ProfilePicture string
				}{
					{
						Name:           "Tom Watkins",
						Gender:         "Male",
						Age:            30,
						ProfilePicture: "test_templates/test_image.jpg",
					},
					{
						Name:           "Evie Argyle",
						Gender:         "Female",
						Age:            29,
						ProfilePicture: "test_templates/test_image.jpeg",
					},
				},
			},
		},
		{
			name:           "Basic document with image structs",
			filename:       "test_with_tables_and_images.docx",
			outputFilename: "test_with_image_structs.docx",
			data: struct {
				ProjectNumber string
				Client        string
				Status        string
				Image         *InlineImage
				People        []struct {
					Name           string
					Gender         string
					Age            uint8
					ProfilePicture *InlineImage
				}
			}{
				ProjectNumber: "B-00001",
				Client:        "TW Software",
				Status:        "New",
				Image:         testImage,
				People: []struct {
					Name           string
					Gender         string
					Age            uint8
					ProfilePicture *InlineImage
				}{
					{
						Name:           "Tom Watkins",
						Gender:         "Male",
						Age:            30,
						ProfilePicture: profileImage,
					},
					{
						Name:           "Evie Argyle",
						Gender:         "Female",
						Age:            29,
						ProfilePicture: profileImage,
					},
				},
			},
		},
		{
			name:     "Basic document with custom functions",
			filename: "test_with_custom_functions.docx",
			data: struct {
				ProjectNumber string
				Client        string
				Status        string
			}{
				ProjectNumber: "B-00001",
				Client:        "TW Software",
				Status:        "New",
			},
			fns: map[string]any{
				"hello": func(text string) string {
					return fmt.Sprintf("Hello %v", text)
				},
			},
		},
	}

	return tests, nil
}

func TestParseAndRender(t *testing.T) {
	tests, err := getTests()
	assert.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			doc, err := ParseFromFilename("test_templates/" + tt.filename)
			assert.Nil(err, "Parsing error")

			for fnName, fn := range tt.fns {
				err = doc.RegisterFunction(fnName, fn)
				assert.Nil(err)
			}

			err = doc.Render(tt.data)
			assert.Nil(err, "Rendering error")

			outputFilename := tt.filename
			if len(tt.outputFilename) > 0 {
				outputFilename = tt.outputFilename
			}
			f, err := os.Create("test_templates/generated_" + outputFilename)
			assert.Nil(err, "Error creating document")

			err = doc.Save(f)
			assert.Nil(err, "Error sacing document")

			err = f.Close()
			assert.Nil(err, "Error closing document")
		})
	}
}

func parseAndRender(t *testing.T, filename string, data any) {
	assert := assert.New(t)

	start := time.Now()

	// Parse the document
	parseStart := time.Now()
	doc, err := ParseFromFilename("test_templates/" + filename)
	assert.Nil(err, "Parsing error")
	fmt.Printf("%v - Parse: %v\n", t.Name(), time.Since(parseStart))

	// Render the document
	renderStart := time.Now()
	err = doc.Render(data)
	assert.Nil(err, "Rendering error")
	fmt.Printf("%v - Render: %v\n", t.Name(), time.Since(renderStart))

	// Create a new file for the output
	saveStart := time.Now()
	f, err := os.Create("test_templates/generated_" + filename)
	assert.Nil(err, "Error creating document")
	err = doc.Save(f)
	assert.Nil(err, "Error sacing document")
	err = f.Close()
	assert.Nil(err, "Error closing document")
	fmt.Printf("%v - Save: %v\n", t.Name(), time.Since(saveStart))

	// Log the overall time taken
	fmt.Printf("%v - Total: %v\n", t.Name(), time.Since(start))
}
