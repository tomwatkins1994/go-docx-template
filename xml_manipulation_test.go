package docxtpl

import (
	"testing"
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
			outputXml, err := replaceTableRangeRows(tt.inputXml)
			if err != nil {
				t.Fatalf("Unexptected error: %v", err)
			}
			if outputXml != tt.expectedOutputXml {
				t.Fatalf("output XML string `%v` does not match expected `%v`", outputXml, tt.expectedOutputXml)
			}
		})
	}
}
