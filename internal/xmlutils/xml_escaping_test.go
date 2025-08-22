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
			expected: "&#34;Quoted text&#34;",
		},
		{
			input:    "'Single quoted text'",
			expected: "&#39;Single quoted text&#39;",
		},
	}

	for _, tt := range tests {
		t.Run("Should escape: "+tt.input, func(t *testing.T) {
			result, err := EscapeXmlString(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
