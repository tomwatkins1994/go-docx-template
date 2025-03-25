package docxtpl

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

	t.Run("Basic document with images", func(t *testing.T) {
		data := struct {
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

	t.Run("Document with tables and images", func(t *testing.T) {
		data := struct {
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
		}
		parseAndRender(t, "test_with_tables_and_images.docx", data)
	})

	t.Run("Document with image structs", func(t *testing.T) {
		assert := assert.New(t)

		doc, err := ParseFromFilename("test_templates/test_with_tables_and_images.docx")
		assert.Nil(err)

		testImage, err := CreateInlineImage("test_templates/test_image.png")
		assert.Nil(err)

		profileImage, err := CreateInlineImage("test_templates/test_image.jpg")
		assert.Nil(err)

		err = profileImage.Resize(100, 100)
		assert.Nil(err)

		data := struct {
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
		}

		err = doc.Render(data)
		assert.Nil(err)

		f, err := os.Create("test_templates/generated_test_with_image_structs.docx")
		assert.Nil(err)

		err = doc.Save(f)
		assert.Nil(err)

		err = f.Close()
		assert.Nil(err)
	})

	t.Run("Document with custom function", func(t *testing.T) {
		assert := assert.New(t)

		doc, err := ParseFromFilename("test_templates/test_with_custom_function.docx")
		assert.Nil(err)

		err = doc.RegisterFunction("hello", func(text string) string {
			return fmt.Sprintf("Hello %v", text)
		})
		assert.Nil(err)

		data := struct {
			ProjectNumber string
			Client        string
			Status        string
		}{
			ProjectNumber: "B-00001",
			Client:        "TW Software",
			Status:        "New",
		}

		err = doc.Render(data)
		assert.Nil(err)

		f, err := os.Create("test_templates/generated_test_with_custom_function.docx")
		assert.Nil(err)

		err = doc.Save(f)
		assert.Nil(err)

		err = f.Close()
		assert.Nil(err)
	})
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
