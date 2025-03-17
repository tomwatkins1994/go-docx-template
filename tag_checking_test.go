package docxtpl

import "testing"

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
			result, err := textContainsTags(tt.text)
			if err != nil {
				t.Fatalf("Unexpected error checking for tags: %v", err)
			}
			if result != tt.expectedResult {
				t.Fatalf("%v - should return %v but returned %v", tt.text, tt.expectedResult, result)
			}
		})
	}
}

func TestTextContainsIncompleteTags(t *testing.T) {
	// 	t.Run("Single partial tag", func(t *testing.T) {
	// 		containsIncompleteTags, err := textContainsIncompleteTags("This text contains a partial tag: {{.Tag")
	// 		if err != nil {
	// 			t.Fatalf("Error checking single incomplete tag: %v", err)
	// 		}
	// 		if containsIncompleteTags == false {
	// 			t.Fatalf("Should be true")
	// 		}
	// 	})

	// 	t.Run("Multiple tags", func(t *testing.T) {
	// 		containsIncompleteTags, err := textContainsIncompleteTags("This text contains multiple incomplete tags: {{.Tag1 and {{.Tag2")
	// 		if err != nil {
	// 			t.Fatalf("Error checking multiple tags: %v", err)
	// 		}
	// 		if containsIncompleteTags == false {
	// 			t.Fatalf("Should be true")
	// 		}
	// 	})

	// 	t.Run("No tag", func(t *testing.T) {
	// 		containsIncompleteTags, err := textContainsIncompleteTags("This text contains no incomplete tags")
	// 		if err != nil {
	// 			t.Fatalf("Error checking single tag: %v", err)
	// 		}
	// 		if containsIncompleteTags == true {
	// 			t.Fatalf("Should be false")
	// 		}
	// 	})

	// 	t.Run("Complete tag", func(t *testing.T) {
	// 		containsIncompleteTags, err := textContainsIncompleteTags("This text contains a tag {{.Tag}}")
	// 		if err != nil {
	// 			t.Fatalf("Error checking single tag: %v", err)
	// 		}
	// 		if containsIncompleteTags == true {
	// 			t.Fatalf("Should be false")
	// 		}
	// 	})

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
			result, err := textContainsIncompleteTags(tt.text)
			if err != nil {
				t.Fatalf("Unexpected error checking for incomplete tags: %v", err)
			}
			if result != tt.expectedResult {
				t.Fatalf("%v - should return %v but returned %v", tt.text, tt.expectedResult, result)
			}
		})
	}
}
