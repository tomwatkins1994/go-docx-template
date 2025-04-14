package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Validation Functions

func TestFunctionNameValid(t *testing.T) {
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

			result := FunctionNameValid(tt.fnName)
			assert.Equal(result, tt.expectedResult)
		})
	}
}

func TestFunctionValid(t *testing.T) {
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

			err := FunctionValid(tt.fn)
			assert.Equal((err != nil), tt.expectError)
		})
	}
}
