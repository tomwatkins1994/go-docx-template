package docxtpl

import (
	"testing"
)

func TestReplaceTableRangeRows(t *testing.T) {
	t.Run("Basic range tag", func(t *testing.T) {
		xmlString := "<w:tr>{{range . }}</w:tr>"
		newXmlString, err := replaceTableRangeRows(xmlString)
		if err != nil {
			t.Fatalf("Error in basic range tag: %v", err)
		}
		if newXmlString != "{{range . }}" {
			t.Fatalf("Text left is: %v", newXmlString)
		}
	})

	t.Run("Basic end tag", func(t *testing.T) {
		xmlString := "<w:tr>{{end}}</w:tr>"
		newXmlString, err := replaceTableRangeRows(xmlString)
		if err != nil {
			t.Fatalf("Error in basic end tag: %v", err)
		}
		if newXmlString != "{{end}}" {
			t.Fatalf("Text left is: %v", newXmlString)
		}
	})

	t.Run("Range tag in other text", func(t *testing.T) {
		xmlString := "<w:tbl><w:tr>{{range . }}</w:tr></w:tbl>"
		newXmlString, err := replaceTableRangeRows(xmlString)
		if err != nil {
			t.Fatalf("Error in range tag in otther text: %v", err)
		}
		if newXmlString != "<w:tbl>{{range . }}</w:tbl>" {
			t.Fatalf("Text left is: %v", newXmlString)
		}
	})

	t.Run("End tag in other text", func(t *testing.T) {
		xmlString := "<w:tbl><w:tr>{{end}}</w:tr></w:tbl>"
		newXmlString, err := replaceTableRangeRows(xmlString)
		if err != nil {
			t.Fatalf("Error in end tag in other text: %v", err)
		}
		if newXmlString != "<w:tbl>{{end}}</w:tbl>" {
			t.Fatalf("Text left is: %v", newXmlString)
		}
	})

	t.Run("Multiple range tags", func(t *testing.T) {
		xmlString := "<w:tbl><w:tr>{{range . }}</w:tr></w:tbl><w:tbl><w:tr>{{range . }}</w:tr></w:tbl>"
		newXmlString, err := replaceTableRangeRows(xmlString)
		if err != nil {
			t.Fatalf("Error in multiple range tag: %v", err)
		}
		if newXmlString != "<w:tbl>{{range . }}</w:tbl><w:tbl>{{range . }}</w:tbl>" {
			t.Fatalf("Text left is: %v", newXmlString)
		}
	})

	t.Run("Multiple end tags", func(t *testing.T) {
		xmlString := "<w:tbl><w:tr>{{end}}</w:tr></w:tbl><w:tbl><w:tr>{{end}}</w:tr></w:tbl>"
		newXmlString, err := replaceTableRangeRows(xmlString)
		if err != nil {
			t.Fatalf("Error in multiple end tags: %v", err)
		}
		if newXmlString != "<w:tbl>{{end}}</w:tbl><w:tbl>{{end}}</w:tbl>" {
			t.Fatalf("Text left is: %v", newXmlString)
		}
	})

	t.Run("Full table", func(t *testing.T) {
		xmlString := "<w:tbl><w:tr>{{range . }}</w:tr><w:tr></w:tr><w:tr>{{end}}</w:tr></w:tbl>"
		newXmlString, err := replaceTableRangeRows(xmlString)
		if err != nil {
			t.Fatalf("Error in full table: %v", err)
		}
		if newXmlString != "<w:tbl>{{range . }}<w:tr></w:tr>{{end}}</w:tbl>" {
			t.Fatalf("Text left is: %v", newXmlString)
		}
	})
}
