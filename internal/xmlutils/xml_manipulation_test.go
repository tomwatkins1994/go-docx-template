package docxtpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceTableRangeRows(t *testing.T) {
	tests := []struct {
		name              string
		inputXml          string
		expectedOutputXml string
	}{
		{
			name:              "Basic range tag",
			inputXml:          "<w:tr>{{range . }}</w:tr>",
			expectedOutputXml: "{{range . }}",
		},
		{
			name:              "Basic end tag",
			inputXml:          "<w:tr>{{end}}</w:tr>",
			expectedOutputXml: "{{end}}",
		},
		{
			name:              "Range tag in other text",
			inputXml:          "<w:tbl><w:tr>{{range . }}</w:tr></w:tbl>",
			expectedOutputXml: "<w:tbl>{{range . }}</w:tbl>",
		},
		{
			name:              "End tag in other text",
			inputXml:          "<w:tbl><w:tr>{{end}}</w:tr></w:tbl>",
			expectedOutputXml: "<w:tbl>{{end}}</w:tbl>",
		},
		{
			name:              "Multiple range tags",
			inputXml:          "<w:tbl><w:tr>{{range . }}</w:tr></w:tbl><w:tbl><w:tr>{{range . }}</w:tr></w:tbl>",
			expectedOutputXml: "<w:tbl>{{range . }}</w:tbl><w:tbl>{{range . }}</w:tbl>",
		},
		{
			name:              "Multiple end tags",
			inputXml:          "<w:tbl><w:tr>{{end}}</w:tr></w:tbl><w:tbl><w:tr>{{end}}</w:tr></w:tbl>",
			expectedOutputXml: "<w:tbl>{{end}}</w:tbl><w:tbl>{{end}}</w:tbl>",
		},
		{
			name:              "Full table",
			inputXml:          "<w:tbl><w:tr>{{range . }}</w:tr><w:tr></w:tr><w:tr>{{end}}</w:tr></w:tbl>",
			expectedOutputXml: "<w:tbl>{{range . }}<w:tr></w:tr>{{end}}</w:tbl>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			outputXml, err := replaceTableRangeRows(tt.inputXml)
			assert.Nil(err)
			assert.Equal(outputXml, tt.expectedOutputXml)
		})
	}
}

func TestFixXmlIssuesPostTagReplacement(t *testing.T) {
	tests := []struct {
		name              string
		inputXml          string
		expectedOutputXml string
	}{
		{
			name:              "Drawing tags",
			inputXml:          "<w:t><w:drawing>...</w:drawing></w:t>",
			expectedOutputXml: "<w:drawing>...</w:drawing>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputXml := fixXmlIssuesPostTagReplacement(tt.inputXml)
			assert.Equal(t, outputXml, tt.expectedOutputXml)
		})
	}
}
