package docxtpl

import "testing"

func TestTextContainsTags(t *testing.T) {
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

	t.Run("Incomplete tag", func(t *testing.T) {
		containsTags, err := textContainsTags("This text contains an incomplete tag: {{.Tag")
		if err != nil {
			t.Fatalf("Error checking single tag: %v", err)
		}
		if containsTags == true {
			t.Fatalf("Should be false")
		}
	})
}

func TestTextContainsIncompleteTags(t *testing.T) {
	t.Run("Single partial tag", func(t *testing.T) {
		containsIncompleteTags, err := textContainsIncompleteTags("This text contains a partial tag: {{.Tag")
		if err != nil {
			t.Fatalf("Error checking single incomplete tag: %v", err)
		}
		if containsIncompleteTags == false {
			t.Fatalf("Should be true")
		}
	})

	t.Run("Multiple tags", func(t *testing.T) {
		containsIncompleteTags, err := textContainsIncompleteTags("This text contains multiple incomplete tags: {{.Tag1 and {{.Tag2")
		if err != nil {
			t.Fatalf("Error checking multiple tags: %v", err)
		}
		if containsIncompleteTags == false {
			t.Fatalf("Should be true")
		}
	})

	t.Run("No tag", func(t *testing.T) {
		containsIncompleteTags, err := textContainsIncompleteTags("This text contains no incomplete tags")
		if err != nil {
			t.Fatalf("Error checking single tag: %v", err)
		}
		if containsIncompleteTags == true {
			t.Fatalf("Should be false")
		}
	})

	t.Run("Complete tag", func(t *testing.T) {
		containsIncompleteTags, err := textContainsIncompleteTags("This text contains a tag {{.Tag}}")
		if err != nil {
			t.Fatalf("Error checking single tag: %v", err)
		}
		if containsIncompleteTags == true {
			t.Fatalf("Should be false")
		}
	})
}
