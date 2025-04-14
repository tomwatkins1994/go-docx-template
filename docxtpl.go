package docxtpl

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"io"
	"maps"
	"os"
	"text/template"

	"github.com/fumiama/go-docx"
	"github.com/tomwatkins1994/go-docx-template/internal/contenttypes"
	"github.com/tomwatkins1994/go-docx-template/internal/functions"
	"github.com/tomwatkins1994/go-docx-template/internal/tags"
)

type DocxTmpl struct {
	*docx.Docx
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
	doc, err := docx.Parse(reader, size)
	if err != nil {
		return nil, err
	}

	contentTypes, err := contenttypes.GetContentTypes(reader, size)
	if err != nil {
		return nil, err
	}

	funcMap := make(template.FuncMap)
	maps.Copy(funcMap, functions.DefaultFuncMap)

	return &DocxTmpl{doc, funcMap, contentTypes}, nil
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
	tags.MergeTags(d.Document.Body.Items)

	// Process the template data
	processedData, err := d.processTemplateData(data)
	if err != nil {
		return err
	}

	// Get the document XML
	documentXmlString, err := d.getDocumentXml()
	if err != nil {
		return err
	}

	// Replace the tags in XML
	documentXmlString, err = tags.ReplaceTagsInXml(documentXmlString, processedData, d.funcMap)
	if err != nil {
		return err
	}

	// Unmarshal the modified XML and replace the document body with it
	decoder := xml.NewDecoder(bytes.NewBufferString(documentXmlString))
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if start, ok := t.(xml.StartElement); ok {
			if start.Name.Local == "Body" {
				clear(d.Document.Body.Items)
				err = d.Document.Body.UnmarshalXML(decoder, start)
				if err != nil {
					return err
				}
				break
			}
		}
	}

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
	_, err := d.WriteTo(&buf)
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

func (d *DocxTmpl) getDocumentXml() (string, error) {
	out, err := xml.Marshal(d.Document.Body)
	if err != nil {
		return "", nil
	}

	return string(out), err
}
