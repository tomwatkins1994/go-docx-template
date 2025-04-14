package docxtpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextContainsTags(t *testing.T) {
	tests := []struct {
		name           string
		text           string
		expectedResult bool
	}{
		{
			name:           "Single tag",
			text:           "This text contains a tag: {{.Tag}}.",
			expectedResult: true,
		},
		{
			name:           "Multiple tags",
			text:           "This text contains multiple tags: {{.Tag1}} and {{.Tag2}}.",
			expectedResult: true,
		},
		{
			name:           "No tag",
			text:           "This text contains no tags",
			expectedResult: false,
		},
		{
			name:           "Incomplete tag",
			text:           "This text contains an incomplete tag: {{.Tag",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := textContainsTags(tt.text)
			assert.Equal(t, result, tt.expectedResult)
		})
	}
}

func TestTextContainsIncompleteTags(t *testing.T) {
	tests := []struct {
		name           string
		text           string
		expectedResult bool
	}{
		{
			name:           "Single incomplete tag",
			text:           "This text contains a tag: {{.Tag",
			expectedResult: true,
		},
		{
			name:           "Multiple incomplete tags",
			text:           "This text contains multiple tags: {{.Tag1 and {{.Tag2.",
			expectedResult: true,
		},
		{
			name:           "No tag",
			text:           "This text contains no tags",
			expectedResult: false,
		},
		{
			name:           "Complete tag",
			text:           "This text contains an complete tag: {{.Tag}}",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := textContainsIncompleteTags(tt.text)
			assert.Equal(t, result, tt.expectedResult)
		})
	}
}
