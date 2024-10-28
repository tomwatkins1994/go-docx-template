package docxtpl

import (
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var defaultFuncMap = template.FuncMap{
	"upper": strings.ToUpper,
	"lower": strings.ToLower,
	"title": title,
}

func title(text string) string {
	caser := cases.Title(language.English)
	return caser.String(text)
}
