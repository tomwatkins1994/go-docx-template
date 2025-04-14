package functions

import "testing"

func TestTitle(t *testing.T) {
	titleCase := title("tom watkins")
	if titleCase != "Tom Watkins" {
		t.Fatalf("Should return in title case but returned: %v", titleCase)
	}
}
