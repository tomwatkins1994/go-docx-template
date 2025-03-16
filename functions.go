package docxtpl

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var defaultFuncMap = template.FuncMap{
	"upper": strings.ToUpper,
	"lower": strings.ToLower,
	"title": title,
}

func (d *DocxTmpl) RegisterFunction(name string, fn any) error {
	if !goodName(name) {
		return fmt.Errorf("function name %q is not a valid identifier", name)
	}

	// Check the function signature
	err := goodFunc(fn)
	if err != nil {
		return fmt.Errorf("error registering function (%s): %s", name, err.Error())
	}

	// Add to the function map
	(*d.funcMap)[name] = fn

	return nil
}

// Validation functions

func goodName(name string) bool {
	if name == "" {
		return false
	}
	for i, r := range name {
		switch {
		case r == '_':
		case i == 0 && !unicode.IsLetter(r):
			return false
		case !unicode.IsLetter(r) && !unicode.IsDigit(r):
			return false
		}
	}
	return true
}

func goodFunc(fn any) error {
	// Check that fn is a function
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		return fmt.Errorf("not a function")
	}

	// We allow functions with 1 result or 2 results where the second is an error.
	typ := v.Type()
	switch numOut := typ.NumOut(); {
	case numOut == 1:
		return nil
	case numOut == 2 && typ.Out(1) == reflect.TypeFor[error]():
		return nil
	case numOut == 2:
		return fmt.Errorf("invalid function signature - second return value should be error; is %s", typ.Out(1))
	default:
		return fmt.Errorf("function has %d return values; should be 1 or 2", typ.NumOut())
	}
}

// Custom functions

func title(text string) string {
	caser := cases.Title(language.English)
	return caser.String(text)
}
