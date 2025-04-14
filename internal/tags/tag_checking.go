package tags

import (
	"regexp"
)

var tagRegex = regexp.MustCompile("{{.*?}}")
var incompleteTagRegex = regexp.MustCompile("{{[^}]*?(}$|$)|{$|^}")

func textContainsTags(text string) bool {
	return tagRegex.MatchString(text)
}

func textContainsIncompleteTags(text string) bool {
	return incompleteTagRegex.MatchString(text)
}
