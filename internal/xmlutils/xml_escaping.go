package xmlutils

import "strings"

// Escape special XML characters in a string.
func EscapeXmlString(xmlString string) string {
	xmlString = strings.ReplaceAll(xmlString, "&", "&amp;")
	xmlString = strings.ReplaceAll(xmlString, "<", "&lt;")
	xmlString = strings.ReplaceAll(xmlString, ">", "&gt;")
	xmlString = strings.ReplaceAll(xmlString, "\"", "&quot;")
	xmlString = strings.ReplaceAll(xmlString, "'", "&apos;")

	return xmlString
}
