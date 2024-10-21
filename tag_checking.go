package docxtpl

import "regexp"

func textContainsTags(text string) (bool, error) {
	r, err := regexp.Compile("{{.*?}}")
	if err != nil {
		return false, err
	}
	return r.MatchString(text), err
}

func textContainsIncompleteTags(text string) (bool, error) {
	r, err := regexp.Compile("{{[^}]*(}$|$)|{$|^}")
	if err != nil {
		return false, err
	}
	return r.MatchString(text), err
}
