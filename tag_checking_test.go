package docxtpl

import "testing"

func TestTagContainsText(t *testing.T) {
	t.Run("Single tag", func(t *testing.T) {
		containsTags, err := textContainsTags("This text contains a tag: {{.Tag}}.")
		if err != nil {
			t.Fatalf("Error checking single tag: %v", err)
		}
		if containsTags == false {
			t.Fatalf("Should be true")
		}
	})

	t.Run("Multiple tags", func(t *testing.T) {
		containsTags, err := textContainsTags("This text contains multiple tags: {{.Tag1}} and {{.Tag2}}.")
		if err != nil {
			t.Fatalf("Error checking multiple tags: %v", err)
		}
		if containsTags == false {
			t.Fatalf("Should be true")
		}
	})

	t.Run("No tag", func(t *testing.T) {
		containsTags, err := textContainsTags("This text contains no tags")
		if err != nil {
			t.Fatalf("Error checking single tag: %v", err)
		}
		if containsTags == true {
			t.Fatalf("Should be false")
		}
	})

	t.Run("Partial tag", func(t *testing.T) {
		containsTags, err := textContainsTags("This text contains part of a tag {{.Tag")
		if err != nil {
			t.Fatalf("Error checking single tag: %v", err)
		}
		if containsTags == true {
			t.Fatalf("Should be false")
		}
	})
}
