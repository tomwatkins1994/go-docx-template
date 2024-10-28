package docxtpl

import (
	"strings"
	"text/template"
)

var defaultFuncMap = template.FuncMap{
	"upper": strings.ToUpper,
	"lower": strings.ToLower,
}
