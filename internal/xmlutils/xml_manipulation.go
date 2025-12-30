package xmlutils

import (
	"strings"
	"time"

	"github.com/dlclark/regexp2"
)

func PrepareXmlForTagReplacement(xmlString string) (string, error) {
	newXmlString, err := replaceTableRangeRows(xmlString)

	return newXmlString, err
}

var tableRangeRowRegex = regexp2.MustCompile("<w:tr>(?:(?!<w:tr>).)*?({{range .*?}}|{{ range .*? }}|{{end}}|{{ end }})(?:(?!<w:tr>).)*?</w:tr>", 0)

func replaceTableRangeRows(xmlString string) (string, error) {
	tableRangeRowRegex.MatchTimeout = 50 * time.Millisecond

	newXmlString := xmlString

	m, err := tableRangeRowRegex.FindStringMatch(xmlString)
	if err != nil {
		return "", err
	}
	for m != nil {
		gps := m.Groups()
		newXmlString = strings.Replace(newXmlString, m.String(), gps[1].Captures[0].String(), 1)
		m, _ = tableRangeRowRegex.FindNextMatch(m)
	}

	return newXmlString, nil
}

func FixXmlIssuesPostTagReplacement(xmlString string) string {
	// Fix issues with drawings in text nodes
	xmlString = strings.ReplaceAll(xmlString, "<w:t><w:drawing>", "<w:drawing>")
	xmlString = strings.ReplaceAll(xmlString, "</w:drawing></w:t>", "</w:drawing>")

	return xmlString
}
