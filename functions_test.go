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

func TestTitle(t *testing.T) {
	titleCase := title("tom watkins")
	if titleCase != "Tom Watkins" {
		t.Fatalf("Should return in title case but returned: %v", titleCase)
	}
}
