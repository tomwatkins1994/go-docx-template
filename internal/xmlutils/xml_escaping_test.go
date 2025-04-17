package xmlutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeXmlString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "<tag>Text</tag>",
			expected: "&lt;tag&gt;Text&lt;/tag&gt;",
		},
		{
			input:    "Text & more text",
			expected: "Text &amp; more text",
		},
		{
			input:    "\"Quoted text\"",
			expected: "&quot;Quoted text&quot;",
		},
		{
			input:    "'Single quoted text'",
			expected: "&apos;Single quoted text&apos;",
		},
	}

	for _, test := range tests {
		result := EscapeXmlString(test.input)
		if result != test.expected {
			assert.Equal(t, test.expected, result)
		}
	}
}
