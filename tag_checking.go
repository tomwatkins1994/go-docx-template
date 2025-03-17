package docxtpl

import (
	"regexp"
	"sync"
)

var tagRegex *regexp.Regexp
var tagOnce sync.Once

func textContainsTags(text string) (bool, error) {
	var err error
	tagOnce.Do(func() {
		tagRegex, err = regexp.Compile("{{.*?}}")
	})
	if err != nil {
		return false, err
	}

	return tagRegex.MatchString(text), err
}

var incompleteTagRegex *regexp.Regexp
var incompleteTagOnce sync.Once

func textContainsIncompleteTags(text string) (bool, error) {
	var err error
	incompleteTagOnce.Do(func() {
		incompleteTagRegex, err = regexp.Compile("{{[^}]*?(}$|$)|{$|^}")
	})
	if err != nil {
		return false, err
	}

	return incompleteTagRegex.MatchString(text), err
}
