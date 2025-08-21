package docxtpl

import (
	"archive/zip"
	"bytes"
	"io"
	"maps"
	"os"
	"text/template"

	"github.com/tomwatkins1994/go-docx-template/internal/contenttypes"
	"github.com/tomwatkins1994/go-docx-template/internal/functions"
	"github.com/tomwatkins1994/go-docx-template/internal/tags"
	"github.com/tomwatkins1994/go-docx-template/internal/templatedata"
	"github.com/tomwatkins1994/go-docx-template/internal/xmlutils"
)

type DocxTmpl struct {
	docx         DocxWrapper
	funcMap      template.FuncMap
	contentTypes *contenttypes.ContentTypes
}

// Parse the document from a reader and store it in memory.
// You can it invoke from a file.
//
//	reader, err := os.Open("path_to_doc.docx")
//	if err != nil {
//		panic(err)
//	}
//	fileinfo, err := reader.Stat()
//	if err != nil {
//		panic(err)
//	}
//	size := fileinfo.Size()
//	doc, err := docxtpl.Parse(reader, int64(size))
func Parse(reader io.ReaderAt, size int64) (*DocxTmpl, error) {
	docx, err := NewFumiamaDocx(reader, size)
	if err != nil {
		return nil, err
	}

	contentTypes, err := contenttypes.GetContentTypes(reader, size)
	if err != nil {
		return nil, err
	}

	funcMap := make(template.FuncMap)
	maps.Copy(funcMap, functions.DefaultFuncMap)

	return &DocxTmpl{docx, funcMap, contentTypes}, nil
}

// Parse the document from a filename and store it in memory.
//
//	doc, err := docxtpl.ParseFromFilename("path_to_doc.docx")
func ParseFromFilename(filename string) (*DocxTmpl, error) {
	reader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	fileinfo, err := reader.Stat()
	if err != nil {
		return nil, err
	}
	size := fileinfo.Size()

	doxtpl, err := Parse(reader, size)
	if err != nil {
		return nil, err
	}

	return doxtpl, nil
}

// Replace the placeholders in the document with passed in data.
// Data can be a struct or map
//
//	data := struct {
//		FirstName     string
//		LastName      string
//		Gender        string
//	}{
//		FirstName: "Tom",
//		LastName:  "Watkins",
//		Gender:    "Male",
//	}
//
// err = doc.Render(data)
//
// # OR
//
//	data := map[string]any{
//		"ProjectNumber": "B-00001",
//		"Client":        "TW Software",
//		"Status":        "New",
//	}
//
// err = doc.Render(data)
func (d *DocxTmpl) Render(data any) error {
	// Ensure that there are no 'part tags' in the XML document
	d.docx.MergeTags()

	// Process the template data
	processedData, err := d.processTemplateData(data)
	if err != nil {
		return err
	}

	// Get the document XML
	documentXmlString, err := d.docx.GetDocumentXml()
	if err != nil {
		return err
	}

	// Replace the tags in XML
	documentXmlString, err = tags.ReplaceTagsInXml(documentXmlString, processedData, d.funcMap)
	if err != nil {
		return err
	}

	d.docx.ReplaceDocumentXml(documentXmlString)

	return nil
}

// Save the document to a writer.
// This could be a new file.
//
//	f, err := os.Create(FILE_PATH)
//	if err != nil {
//		panic(err)
//	}
//	err = doc.Save(f)
//	if err != nil {
//		panic(err)
//	}
//	err = f.Close()
//	if err != nil {
//		panic(err)
//	}
func (d *DocxTmpl) Save(writer io.Writer) error {
	var buf bytes.Buffer
	err := d.docx.Write(&buf)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(buf.Bytes())
	zipReader, err := zip.NewReader(reader, int64(buf.Len()))
	if err != nil {
		return err
	}

	generatedZip := zip.NewWriter(writer)

	for _, f := range zipReader.File {
		newFile, err := generatedZip.Create(f.Name)
		if err != nil {
			return err
		}

		// Override content types with out calculated types
		// Copy across all other files
		if f.Name == "[Content_Types].xml" {
			contentTypesXml, err := d.contentTypes.MarshalXml()
			if err != nil {
				return err
			}

			_, err = newFile.Write([]byte(contentTypesXml))
			if err != nil {
				return err
			}
		} else {
			zf, err := f.Open()
			if err != nil {
				return err
			}
			defer zf.Close()

			if _, err := io.Copy(newFile, zf); err != nil {
				return err
			}
		}
	}

	if err := generatedZip.Close(); err != nil {
		return err
	}

	return nil
}

func (d *DocxTmpl) processTemplateData(data any) (map[string]any, error) {
	convertedData, err := templatedata.DataToMap(data)
	if err != nil {
		return nil, err
	}

	var processTagValues func(data *map[string]any) error
	processTagValues = func(data *map[string]any) error {
		for key, value := range *data {
			if stringVal, ok := value.(string); ok {
				// Check for files
				if isImage, err := templatedata.IsImageFilePath(stringVal); err != nil {
					return err
				} else {
					if isImage {
						image, err := CreateInlineImage(stringVal)
						if err != nil {
							return err
						}
						imageXml, err := d.addInlineImage(image)
						if err != nil {
							return err
						}
						(*data)[key] = imageXml
					} else {
						xmlEscapedText, err := xmlutils.EscapeXmlString(stringVal)
						if err != nil {
							return err
						}
						(*data)[key] = xmlEscapedText
					}
				}
			} else if nestedMap, ok := value.(map[string]any); ok {
				if err := processTagValues(&nestedMap); err != nil {
					return err
				}
				(*data)[key] = nestedMap
			} else if sliceValue, ok := value.([]map[string]any); ok {
				for i := range sliceValue {
					if err := processTagValues(&sliceValue[i]); err != nil {
						return err
					}
				}
			} else if inlineImage, ok := value.(*InlineImage); ok {
				imageXml, err := d.addInlineImage(inlineImage)
				if err != nil {
					return err
				}
				(*data)[key] = imageXml
			}
		}

		return nil
	}

	err = processTagValues(&convertedData)
	if err != nil {
		return nil, err
	}

	return convertedData, nil
}
