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
	t.Run("Function names only containing letters should be allowed", func(t *testing.T) {
		if goodName("hello") == false {
			t.Fatal(t.Name())
		}
	})

	t.Run("Function names only containing letters and numbers should be allowed", func(t *testing.T) {
		if goodName("hello1") == false {
			t.Fatal(t.Name())
		}
	})

	t.Run("Function names containing _ should be allowed", func(t *testing.T) {
		if goodName("my_function") == false {
			t.Fatal(t.Name())
		}
	})

	t.Run("Blank function names should not be allowed", func(t *testing.T) {
		if goodName("") {
			t.Fatal(t.Name())
		}
	})

	t.Run("Function names not beginning with a letter should not be allowed", func(t *testing.T) {
		if goodName("1") {
			t.Fatal(t.Name())
		}
	})
}

func TestGoodFunc(t *testing.T) {
	t.Run("Valid function", func(t *testing.T) {
		err := goodFunc(func(text string) string {
			return "hello"
		})
		if err != nil {
			t.Fatalf("%v - should be valid: %v", t.Name(), err)
		}
	})

	t.Run("Valid function with error", func(t *testing.T) {
		err := goodFunc(func(text string) (string, error) {
			return "hello", nil
		})
		if err != nil {
			t.Fatalf("%v - should be valid: %v", t.Name(), err)
		}
	})

	t.Run("Invalid function", func(t *testing.T) {
		err := goodFunc(func(text string) {
			// Do nothing here
		})
		if err == nil {
			t.Fatalf("%v - did not return expected error: %v", t.Name(), err)
		}
	})
}

// Custom functions

func TestTitle(t *testing.T) {
	titleCase := title("tom watkins")
	if titleCase != "Tom Watkins" {
		t.Fatalf("Should return in title case but returned: %v", titleCase)
	}
}
