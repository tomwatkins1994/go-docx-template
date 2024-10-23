package docxtpl

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"io"
	"os"
)

type ContentTypes struct {
	XMLName   xml.Name   `xml:"Types"`
	Defaults  []Default  `xml:"Default"`
	Overrides []Override `xml:"Override"`
}

type Default struct {
	Extension   string `xml:"Extension,attr"`
	ContentType string `xml:"ContentType,attr"`
}

type Override struct {
	PartName    string `xml:"PartName,attr"`
	ContentType string `xml:"ContentType,attr"`
}

func getContentTypes(reader io.ReaderAt, size int64) (*ContentTypes, error) {
	zipReader, err := zip.NewReader(reader, size)
	if err != nil {
		return nil, err
	}

	for _, f := range zipReader.File {
		if f.Name == "[Content_Types].xml" {
			zf, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer zf.Close()

			dataBuf := make([]byte, f.UncompressedSize64)
			_, err = zf.Read(dataBuf)
			if err != nil {
				return nil, err
			}

			var ct ContentTypes
			// Unmarshal the XML data into the struct
			if err := xml.Unmarshal(dataBuf, &ct); err != nil {
				return nil, err
			}

			return &ct, nil
		}
	}

	return nil, errors.New("no content types found")
}

func (d *DocxTmpl) addContentType(filePath string, ext string, contentType string) error {
	// Read the existing Content Types file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var ct ContentTypes
	if err := xml.Unmarshal(data, &ct); err != nil {
		return err
	}

	// Add new default content type
	ct.Defaults = append(ct.Defaults, Default{Extension: ext, ContentType: contentType})

	// Marshal back to XML
	output, err := xml.MarshalIndent(ct, "", "    ")
	if err != nil {
		return err
	}

	// Write back to file
	return os.WriteFile(filePath, output, 0644)
}

// func main() {
// 	// Update Content Types file
// 	err := addContentType("[Content_Types].xml", "png", "image/png")
// 	if err != nil {
// 		panic(err)
// 	}
// 	// Repeat for other extensions as needed
// }
