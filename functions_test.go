package docxtpl

import "testing"

func TestRegisterFunctions(t *testing.T) {
	t.Run("Valid function with no error", func(t *testing.T) {
		doc, err := ParseFromFilename("test_templates/test_basic.docx")
		if err != nil {
			t.Fatalf("%v - Parsing error: %v", t.Name(), err)
		}
		err = doc.RegisterFunction("hello", func(text string) string {
			return "hello"
		})
		if err != nil {
			t.Fatalf("%v - Error registering function: %v", t.Name(), err)
		}
	})

	t.Run("Valid function with error", func(t *testing.T) {
		doc, err := ParseFromFilename("test_templates/test_basic.docx")
		if err != nil {
			t.Fatalf("%v - Parsing error: %v", t.Name(), err)
		}
		err = doc.RegisterFunction("hello", func(text string) (string, error) {
			return "hello", nil
		})
		if err != nil {
			t.Fatalf("%v - Error registering function: %v", t.Name(), err)
		}
	})

	t.Run("Invalid function", func(t *testing.T) {
		doc, err := ParseFromFilename("test_templates/test_basic.docx")
		if err != nil {
			t.Fatalf("%v - Parsing error: %v", t.Name(), err)
		}
		err = doc.RegisterFunction("hello", func(text string) {
			// Do nothing here
		})
		if err == nil {
			t.Fatalf("%v - Did not return expected error: %v", t.Name(), err)
		}
	})
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
		if goodName("") == false {
			t.Fatal(t.Name())
		}
	})

	t.Run("Function names not beginning with a letter should not be allowed", func(t *testing.T) {
		if goodName("1") == false {
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
