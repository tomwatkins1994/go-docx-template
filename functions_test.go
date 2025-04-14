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

// Validation Functions

func TestGoodName(t *testing.T) {
	tests := []struct {
		name           string
		fnName         string
		expectedResult bool
	}{
		{
			name:           "Function names only containing letters should be allowed",
			fnName:         "hello",
			expectedResult: true,
		},
		{
			name:           "Function names only containing letters and numbers should be allowed",
			fnName:         "hello123",
			expectedResult: true,
		},
		{
			name:           "Function names containing _ should be allowed",
			fnName:         "my_function",
			expectedResult: true,
		},
		{
			name:           "Blank function names should not be allowed",
			fnName:         "",
			expectedResult: false,
		},
		{
			name:           "Function names not beginning with a letter should not be allowed",
			fnName:         "1",
			expectedResult: false,
		},
		{
			name:           "Function names containing symbols should not be allowed",
			fnName:         "hello!",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			result := goodName(tt.fnName)
			assert.Equal(result, tt.expectedResult)
		})
	}
}

func TestGoodFunc(t *testing.T) {
	tests := []struct {
		name        string
		fn          any
		expectError bool
	}{
		{
			name: "Valid function",
			fn: func(text string) string {
				return text
			},
			expectError: false,
		},
		{
			name: "Valid function with error",
			fn: func(text string) (string, error) {
				return text, nil
			},
			expectError: false,
		},
		{
			name: "Function with no return type",
			fn: func(text string) {
				// Function does nothing
			},
			expectError: true,
		},
		{
			name: "Function with no input",
			fn: func() string {
				return "hello"
			},
			expectError: false,
		},
		{
			name: "Function where second return value is not an error",
			fn: func() (string, string) {
				return "hello", "world"
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			err := goodFunc(tt.fn)
			assert.Equal((err != nil), tt.expectError)
		})
	}
}
