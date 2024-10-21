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
}