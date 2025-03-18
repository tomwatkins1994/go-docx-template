package docxtpl

import (
	"testing"
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
			doc, err := ParseFromFilename("test_templates/test_basic.docx")
			if err != nil {
				t.Fatalf("%v - Parsing error: %v", t.Name(), err)
			}

			err = doc.RegisterFunction(tt.fnName, tt.fn)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}

			funcMap := doc.GetRegisteredFunctions()
			_, exists := (*funcMap)[tt.fnName]
			if !tt.expectError && !exists {
				t.Errorf("%v - function not found in map", tt.name)
			}
			if tt.expectError && exists {
				t.Errorf("%v - function found in map after erroring", tt.name)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := goodName(tt.fnName)
			if result != tt.expectedResult {
				t.Fatalf("%v - should return %v but returned %v", tt.fnName, tt.expectedResult, result)
			}
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
			err := goodFunc(tt.fn)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

// Custom functions

func TestTitle(t *testing.T) {
	titleCase := title("tom watkins")
	if titleCase != "Tom Watkins" {
		t.Fatalf("Should return in title case but returned: %v", titleCase)
	}
}
