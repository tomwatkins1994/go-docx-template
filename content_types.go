package docxtpl

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"io"
)

type ContentTypes struct {
	XMLName   xml.Name   `xml:"Types"`
	Xmlns     string     `xml:"xmlns,attr"`
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

var PNG_CONTENT_TYPE = Default{Extension: "png", ContentType: "image/png"}
var JPG_CONTENT_TYPE = Default{Extension: "jpg", ContentType: "image/jpg"}

func (ct *ContentTypes) addContentType(contentType *Default) {
	for _, v := range ct.Defaults {
		if v == *contentType {
			return
		}
	}
	ct.Defaults = append(ct.Defaults, *contentType)
}

func (ct *ContentTypes) marshalXml() (string, error) {
	output, err := xml.MarshalIndent(ct, "", "  ")
	if err != nil {
		return "", err
	}

	xmlString := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n" + string(output)

	return xmlString, nil
}
