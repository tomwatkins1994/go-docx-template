package docxtpl

import (
	"strings"

	"github.com/dlclark/regexp2"
)

func prepareXmlForTagReplacement(xmlString string) (string, error) {
	newXmlString, err := replaceTableRangeRows(xmlString)

	return newXmlString, err
}

func replaceTableRangeRows(xmlString string) (string, error) {
	r, err := regexp2.Compile("<w:tr>(?:(?!<w:tr>).)*?({{range .*?}}|{{ range .*? }}|{{end}}|{{ end }})(?:(?!<w:tr>).)*?</w:tr>", 0)
	if err != nil {
		return "", err
	}
	r.MatchTimeout = 500

	newXmlString := xmlString

	m, err := r.FindStringMatch(xmlString)
	if err != nil {
		return "", err
	}
	for m != nil {
		gps := m.Groups()
		newXmlString = strings.Replace(newXmlString, m.String(), gps[1].Captures[0].String(), 1)
		m, _ = r.FindNextMatch(m)
	}

	return newXmlString, nil
}

func fixXmlIssuesPostTagReplacement(xmlString string) string {
	// Fix issues with drawings in text nodes
	xmlString = strings.ReplaceAll(xmlString, "<w:t><w:drawing>", "<w:drawing>")
	xmlString = strings.ReplaceAll(xmlString, "</w:drawing></w:t>", "</w:drawing>")

	return xmlString
}
