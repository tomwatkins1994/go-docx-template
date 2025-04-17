package xmlutils

import (
	"bytes"
	"encoding/xml"
)

// Escape special XML characters in a string.
func EscapeXmlString(xmlString string) (string, error) {
	var buf bytes.Buffer
	err := xml.EscapeText(&buf, []byte(xmlString))
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
