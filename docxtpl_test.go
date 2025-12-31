package docxtpl

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tomwatkins1994/go-docx-template/internal/docxwrappers"
	"github.com/tomwatkins1994/go-docx-template/internal/images"
)

type test struct {
	name           string
	filename       string
	outputFilename string
	data           any
	dataFn         func() any
	fns            map[string]any
}

type testWrapper struct {
	name             string
	docxFromFilename func(filename string) (docxwrappers.DocxWrapper, error)
}

func getWrappers() []testWrapper {
	var wrappers = []testWrapper{
		{
			name: "Fumiama",
			docxFromFilename: func(filename string) (docxwrappers.DocxWrapper, error) {
				return docxwrappers.NewFumiamaDocxFromFilename(filename)
			},
		},
		{
			name: "Gomutex",
			docxFromFilename: func(filename string) (docxwrappers.DocxWrapper, error) {
				return docxwrappers.NewGomutexDocxFromFilename(filename)
			},
		},
	}

	return wrappers
}

func getTests() ([]test, error) {
	testImage, err := images.CreateInlineImage("test_templates/test_image.png")
	if err != nil {
		return nil, err
	}
	profileImage, err := images.CreateInlineImage("test_templates/test_image.jpg")
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
			name:           "Basic document with map data",
			filename:       "test_basic.docx",
			outputFilename: "test_basic_with_map_data.docx",
			dataFn: func() any {
				return map[string]any{
					"ProjectNumber": "B-00001",
					"Client":        "TW Software",
					"Status":        "New",
				}
			},
		},
		{
			name:     "Basic document with XML escaped",
			filename: "test_basic_with_escaping.docx",
			data: struct {
				ProjectNumber string
				Client        string
				Status        string
			}{
				ProjectNumber: "'Single quoted text' and \"Quoted text\"",
				Client:        "Text & more text",
				Status:        "<tag>New</tag>",
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
			name:           "Basic document with images in map data",
			filename:       "test_basic_with_images.docx",
			outputFilename: "test_basic_with_images_map_data.docx",
			dataFn: func() any {
				return map[string]any{
					"ProjectNumber": "B-00001",
					"Client":        "TW Software",
					"Status":        "New",
					"ImagePng":      "test_templates/test_image.png",
					"ImageJpg":      "test_templates/test_image.png",
				}
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
			name:           "Basic document with tables in map data",
			filename:       "test_with_tables.docx",
			outputFilename: "test_with_tables_map_data.docx",
			dataFn: func() any {
				return map[string]any{
					"ProjectNumber": "B-00001",
					"Client":        "TW Software",
					"Status":        "New",
					"CreatedBy":     "Tom Watkins",
					"SignedOffBy":   "Tom Watkins",
					"People": []map[string]any{
						{
							"Name":   "Tom Watkins1",
							"Gender": "Male",
							"Age":    30,
						},
						{
							"Name":   "Evie Argyle",
							"Gender": "Female",
							"Age":    29,
						},
					},
				}
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
			name:           "Basic document with tables and images in map data",
			filename:       "test_with_tables_and_images.docx",
			outputFilename: "test_with_tables_and_images_map_data.docx",
			dataFn: func() any {
				return map[string]any{
					"ProjectNumber": "B-00001",
					"Client":        "TW Software",
					"Status":        "New",
					"Image":         "test_templates/test_image.png",
					"People": []map[string]any{
						{
							"Name":           "Tom Watkins1",
							"Gender":         "Male",
							"Age":            30,
							"ProfilePicture": "test_templates/test_image.jpg",
						},
						{
							"Name":           "Evie Argyle",
							"Gender":         "Female",
							"Age":            29,
							"ProfilePicture": "test_templates/test_image.jpeg",
						},
					},
				}
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
				Image         *images.InlineImage
				People        []struct {
					Name           string
					Gender         string
					Age            uint8
					ProfilePicture *images.InlineImage
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
					ProfilePicture *images.InlineImage
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
	docxWrappers := getWrappers()

	tests, err := getTests()
	require.Nil(t, err)

	for _, tt := range tests {
		for _, wrapper := range docxWrappers {
			t.Run(wrapper.name+"_"+tt.name, func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				docx, err := wrapper.docxFromFilename("test_templates/" + tt.filename)
				require.Nil(err, "Parsing error")
				docxtpl := newDocxTmpl(docx)

				for fnName, fn := range tt.fns {
					err = docxtpl.RegisterFunction(fnName, fn)
					assert.Nil(err)
				}

				if tt.dataFn != nil {
					err = docxtpl.Render(tt.dataFn())
				} else {
					err = docxtpl.Render(tt.data)
				}

				assert.Nil(err, "Rendering error")

				outputFilename := tt.filename
				if len(tt.outputFilename) > 0 {
					outputFilename = tt.outputFilename
				}
				err = docxtpl.SaveToFile("test_templates/generated_" + strings.ToLower(wrapper.name) + "_" + outputFilename)
				assert.Nil(err, "Error saving document")
			})
		}
	}
}

func TestProcessTemplateData(t *testing.T) {
	docxWrappers := getWrappers()
	tests := []struct {
		name         string
		dataFn       func() any
		expectedData map[string]any
	}{
		{
			name: "Test XML escaping with structs",
			dataFn: func() any {
				return struct {
					ProjectNumber string
					Client        string
					Status        string
					CreatedBy     string
				}{
					ProjectNumber: "'Single quoted text'",
					Client:        "Text & more text",
					Status:        "\"Quoted text\"",
					CreatedBy:     "<tag>Text</tag>",
				}
			},
			expectedData: map[string]any{
				"ProjectNumber": "&#39;Single quoted text&#39;",
				"Client":        "Text &amp; more text",
				"Status":        "&#34;Quoted text&#34;",
				"CreatedBy":     "&lt;tag&gt;Text&lt;/tag&gt;",
			},
		},
		{
			name: "Test XML escaping with maps",
			dataFn: func() any {
				return map[string]any{
					"ProjectNumber": "'Single quoted text'",
					"Client":        "Text & more text",
					"Status":        "\"Quoted text\"",
					"CreatedBy":     "<tag>Text</tag>",
				}
			},
			expectedData: map[string]any{
				"ProjectNumber": "&#39;Single quoted text&#39;",
				"Client":        "Text &amp; more text",
				"Status":        "&#34;Quoted text&#34;",
				"CreatedBy":     "&lt;tag&gt;Text&lt;/tag&gt;",
			},
		},
		{
			name: "Struct with nested struct",
			dataFn: func() any {
				return struct {
					ProjectNumber string
					People        []struct {
						Name string
					}
				}{
					ProjectNumber: "A-0001",
					People: []struct {
						Name string
					}{
						{
							Name: "Tom Watkins",
						},
						{
							Name: "Evie Argyle",
						},
					},
				}
			},
			expectedData: map[string]any{
				"ProjectNumber": "A-0001",
				"People": []map[string]any{
					{
						"Name": "Tom Watkins",
					},
					{
						"Name": "Evie Argyle",
					},
				},
			},
		},
		{
			name: "Struct with nested map",
			dataFn: func() any {
				return struct {
					ProjectNumber string
					People        []map[string]any
				}{
					ProjectNumber: "A-0001",
					People: []map[string]any{
						{
							"Name": "Tom Watkins",
						},
						{
							"Name": "Evie Argyle",
						},
					},
				}
			},
			expectedData: map[string]any{
				"ProjectNumber": "A-0001",
				"People": []map[string]any{
					{
						"Name": "Tom Watkins",
					},
					{
						"Name": "Evie Argyle",
					},
				},
			},
		},
		{
			name: "Map with nested struct",
			dataFn: func() any {
				return map[string]any{
					"ProjectNumber": "A-0001",
					"People": []struct {
						Name string
					}{
						{
							Name: "Tom Watkins",
						},
						{
							Name: "Evie Argyle",
						},
					},
				}
			},
			expectedData: map[string]any{
				"ProjectNumber": "A-0001",
				"People": []map[string]any{
					{
						"Name": "Tom Watkins",
					},
					{
						"Name": "Evie Argyle",
					},
				},
			},
		},
		{
			name: "Map with nested map",
			dataFn: func() any {
				return map[string]any{
					"ProjectNumber": "A-0001",
					"People": []map[string]any{
						{
							"Name": "Tom Watkins",
						},
						{
							"Name": "Evie Argyle",
						},
					},
				}
			},
			expectedData: map[string]any{
				"ProjectNumber": "A-0001",
				"People": []map[string]any{
					{
						"Name": "Tom Watkins",
					},
					{
						"Name": "Evie Argyle",
					},
				},
			},
		},
	}

	for _, wrapper := range docxWrappers {
		for _, tt := range tests {
			t.Run(wrapper.name+" "+tt.name, func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				docx, err := wrapper.docxFromFilename("test_templates/test_basic.docx")
				require.NoError(err, "Parsing error")
				docxtpl := newDocxTmpl(docx)

				processedData, err := docxtpl.processTemplateData(tt.dataFn())
				require.NoError(err)

				assert.Equal(tt.expectedData, processedData)
			})
		}
	}
}

func BenchmarkParseAndRender(b *testing.B) {
	docxWrappers := getWrappers()

	tests, err := getTests()
	require.Nil(b, err)
	b.ResetTimer()

	for _, tt := range tests {
		for _, wrapper := range docxWrappers {
			b.Run(wrapper.name+" "+tt.name, func(b *testing.B) {
				require := require.New(b)
				start := time.Now()

				parseStart := time.Now()
				docx, err := wrapper.docxFromFilename("test_templates/" + tt.filename)
				require.Nil(err, "Parsing error")
				docxtpl := newDocxTmpl(docx)
				b.Logf("Parse: %v\n", time.Since(parseStart))

				functionsStart := time.Now()
				for fnName, fn := range tt.fns {
					err = docxtpl.RegisterFunction(fnName, fn)
					require.Nil(err)
				}
				b.Logf("Register custom functions: %v\n", time.Since(functionsStart))

				renderStart := time.Now()
				err = docxtpl.Render(tt.data)
				require.Nil(err, "Rendering error")
				b.Logf("Render: %v\n", time.Since(renderStart))

				saveStart := time.Now()
				outputFilename := tt.filename
				if len(tt.outputFilename) > 0 {
					outputFilename = tt.outputFilename
				}
				f, err := os.Create("test_templates/generated_" + strings.ToLower(wrapper.name) + "_" + outputFilename)
				require.Nil(err, "Error creating document")
				err = docxtpl.Save(f)
				require.Nil(err, "Error saving document")
				err = f.Close()
				require.Nil(err, "Error closing document")
				b.Logf("Save: %v\n", time.Since(saveStart))

				b.Logf("Total: %v\n", time.Since(start))
			})
		}
	}
}
