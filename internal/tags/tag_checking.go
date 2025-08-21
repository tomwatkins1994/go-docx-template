package tags

import (
	"regexp"
)

var tagRegex = regexp.MustCompile("{{.*?}}")
var incompleteTagRegex = regexp.MustCompile("{{[^}]*?(}$|$)|{$|^}")

func TextContainsTags(text string) bool {
	return tagRegex.MatchString(text)
}

func TextContainsIncompleteTags(text string) bool {
	return incompleteTagRegex.MatchString(text)
}
