package docxtpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterFunctions(t *testing.T) {
	tests := []struct {
		name        string
		fnName      string
		fn          any
		expectError bool
	}{
		{
			name:   "Valid function",
			fnName: "validFunction",
			fn: func(text string) string {
				return "Hello"
			},
			expectError: false,
		},
		{
			name:        "Invalid function name",
			fnName:      "",
			fn:          nil,
			expectError: true,
		},
		{
			name:        "Invalid function signature",
			fnName:      "validFunction",
			fn:          "not a valid function",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			doc, err := ParseFromFilename("test_templates/test_basic.docx")
			assert.Nil(err)

			err = doc.RegisterFunction(tt.fnName, tt.fn)
			assert.Equal((err != nil), tt.expectError)

			funcMap := doc.GetRegisteredFunctions()
			_, foundInFunctionMap := (*funcMap)[tt.fnName]
			if !tt.expectError {
				assert.True(foundInFunctionMap)
			}
			if tt.expectError {
				assert.False(foundInFunctionMap)
			}
		})
	}
}
