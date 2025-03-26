package docxtpl

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type test struct {
	name           string
	filename       string
	outputFilename string
	data           any
	fns            map[string]any
}

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
	require.Nil(t, err)

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
			assert.Nil(err, "Error saving document")
			err = f.Close()
			assert.Nil(err, "Error closing document")
		})
	}
}

func BenchmarkParseAndRender(b *testing.B) {
	tests, err := getTests()
	require.Nil(b, err)
	b.ResetTimer()

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			require := require.New(b)
			start := time.Now()

			parseStart := time.Now()
			doc, err := ParseFromFilename("test_templates/" + tt.filename)
			require.Nil(err, "Parsing error")
			b.Logf("Parse: %v\n", time.Since(parseStart))

			functionsStart := time.Now()
			for fnName, fn := range tt.fns {
				err = doc.RegisterFunction(fnName, fn)
				require.Nil(err)
			}
			b.Logf("Register custom functions: %v\n", time.Since(functionsStart))

			renderStart := time.Now()
			err = doc.Render(tt.data)
			require.Nil(err, "Rendering error")
			b.Logf("Render: %v\n", time.Since(renderStart))

			saveStart := time.Now()
			outputFilename := tt.filename
			if len(tt.outputFilename) > 0 {
				outputFilename = tt.outputFilename
			}
			f, err := os.Create("test_templates/generated_" + outputFilename)
			require.Nil(err, "Error creating document")
			err = doc.Save(f)
			require.Nil(err, "Error saving document")
			err = f.Close()
			require.Nil(err, "Error closing document")
			b.Logf("Save: %v\n", time.Since(saveStart))

			b.Logf("Total: %v\n", time.Since(start))
		})
	}
}
